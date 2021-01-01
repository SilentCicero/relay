package server

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"net"
	"net/http"
	"strings"
	"time"
)

// Server version string
var Version string

// Server instance struct
type Server struct {
	net.Listener
	stopChan  chan struct{}
	errChan   chan error
	outdaters []repository.Outdater
	waiters   []Waiter
}

type Waiter interface {
	Wait() <-chan struct{}
}

// NewServer creates new `HTTP Relay` server instance and returns it.
func NewServer(listener net.Listener) (server *Server) {
	server = &Server{
		stopChan: make(chan struct{}),
		errChan:  make(chan error, 1),
	}

	server.Listener = listener

	syncRep := repository.NewSyncRep(server.stopChan)
	syncCtrl := controller.NewSyncCtrl(syncRep, server.stopChan)
	http.HandleFunc("/sync/", corsHandler(syncCtrl.Conduct, []string{}))

	linkRep := repository.NewLinkRep(server.stopChan)
	linkCtrl := controller.NewLinkCtrl(linkRep, server.stopChan)
	http.HandleFunc("/link/", corsHandler(linkCtrl.Conduct, []string{}))

	mcastRep := repository.NewMcastRep(server.stopChan)
	mcastCtrl := controller.NewMcastCtrl(mcastRep, server.stopChan)
	http.HandleFunc("/mcast/", corsHandler(mcastCtrl.Conduct, []string{"HttpRelay-SeqId"}))

	proxyRep := repository.NewProxyRep()
	proxyCtrl := controller.NewProxyCtrl(proxyRep, server.stopChan)
	http.HandleFunc("/proxy/", wildcardCorsHandler(proxyCtrl.Conduct))

	server.outdaters = []repository.Outdater{linkRep, mcastRep, proxyRep}
	server.waiters = []Waiter{syncCtrl, linkCtrl, mcastCtrl, proxyCtrl}

	return
}

func corsHandler(h http.HandlerFunc, expose []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		cors(w, r, expose)
		if r.Method != "OPTIONS" {
			h(w, r)
		}
	}
}

func cors(w http.ResponseWriter, r *http.Request, expose []string) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	//w.Header().Set("HttpRelay-Version", Version)

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS, SERVE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	} else {
		expose = append(expose, "Content-Length, X-Real-IP, X-Real-Port, Httprelay-Version, Httprelay-Time, Httprelay-Your-Time, Httprelay-Method, Httprelay-Query")
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(expose, ", "))
	}
}

func wildcardCorsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		wildcardCors(w, r)
		if r.Method != "OPTIONS" {
			h(w, r)
		}
	}
}

func wildcardCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	//w.Header().Set("HttpRelay-Version", Version)

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS, SERVE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
	} else {
		w.Header().Set("Access-Control-Expose-Headers", "*")
	}
}

// Starts server
func (s *Server) Start() <-chan error {
	go repository.Outdate(s.outdaters, time.Minute, s.stopChan)

	go func() {
		if err := http.Serve(s, nil); err != nil && s.Active() {
			s.Stop(time.Second)
			s.errChan <- err
		}
	}()
	return s.errChan
}

// Stops server
func (s *Server) Stop(timeout time.Duration) {
	close(s.stopChan)
	s.waitAll(timeout)
	s.Close()
}

// Returns true if server is active
func (s *Server) Active() bool {
	select {
	case <-s.stopChan:
		return false
	default:
		return true
	}
}

func (s *Server) waitAll(timeout time.Duration) {
	t := time.NewTimer(timeout)
	for _, w := range s.waiters {
		select {
		case <-w.Wait():
		case <-t.C:
		}
	}
	t.Stop()
}
