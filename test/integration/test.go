package integration

//const host = "staging.httprelay.io"
//const ttl = 1*60
//
//var schemes = []string{"http", "https"}
//var methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
//
//type reqData struct {
//	name string
//	scheme string
//	method string
//	mode string
//	id string
//	query string
//	contentType string
//	body *[]byte
//}
//
//type respData struct {
//	name string
//	status int
//	mode string
//	id string
//	err error
//	duration time.Duration
//	url *url.URL
//	contentType string
//	body *[]byte
//	ip string
//	port string
//	time string
//	yourTime string
//	method string
//	query string
//}
//
//func main(){
//	http.DefaultClient.Timeout = time.Duration(ttl+time.Minute)
//	sync()
//}
//
//func sync() {
//	errChArr := []chan error{}
//	for i := 0; i <= ttl; i++ {
//		uuid := uuid()
//		errCh := make(chan error)
//		errChArr = append(errChArr, errCh)
//		go syncPair(errCh, i, uuid)
//	}
//
//	for i, errCh := range errChArr {
//		err := <-errCh
//		fmt.Println(i, err)
//	}
//}
//
//func syncPair(errCh chan error, idx int, uuid string) {
//	delay := time.Second * time.Duration(idx)
//	schemeComb, _ := permRep(idx, 2, len(schemes))
//	methodComb, _ := permRep(idx, 2, len(methods))
//	bodyA := bytes.Repeat([]byte{ 65 }, idx) // 65="A"
//	bodyB := bytes.Repeat([]byte{ 66 }, idx) // 66="B"
//	rA := &reqData{ name: "A", scheme: schemes[schemeComb[0]], method: methods[methodComb[0]], mode: "sync", id: uuid, query: "peer=A", contentType: "content-type-A", body: &bodyA}
//	rB := &reqData{ name: "B", scheme: schemes[schemeComb[1]], method: methods[methodComb[1]], mode: "sync", id: uuid, query: "peer=B", contentType: "content-type-B", body: &bodyB}
//
//	aC := make(chan *respData)
//	go doReq("A", aC, 0, delay+time.Second, rA)
//	bC := make(chan *respData)
//	go doReq("B", bC, delay, time.Second, rB)
//
//	go func() {
//		if err := reqResValidation(rB, <-aC); err != nil {
//			errCh<- err
//		} else {
//			errCh<- reqResValidation(rA, <-bC)
//		}
//	}()
//}
//
//func doReq(name string, c chan *respData, delay time.Duration, timeout time.Duration, r *reqData) {
//	time.Sleep(delay)
//	u := &url.URL{
//		Scheme:   r.scheme,
//		Host:     host,
//		Path:     fmt.Sprintf("/%s/%s", r.mode, r.id),
//		RawQuery: r.query,
//	}
//
//	req := &http.Request{
//		Method: r.method,
//		Header: http.Header{"Content-Type": []string{r.contentType}},
//		URL:    u,
//		Body:   ioutil.NopCloser(bytes.NewReader(*r.body)),
//	}
//
//	startTime := time.Now()
//
//	client := http.Client{Timeout: timeout}
//	resp, err := client.Do(req)
//
//	rr := &respData{
//		name: name,
//		mode: r.mode,
//		id:       r.id,
//		err:      err,
//		duration: time.Since(startTime),
//	}
//
//	if err == nil {
//		respBody, err := ioutil.ReadAll(resp.Body)
//		rr.url = resp.Request.URL
//		rr.status = resp.StatusCode
//		rr.err = err
//		rr.body = &respBody
//		rr.contentType = resp.Header.Get("Content-Type")
//		rr.ip = resp.Header.Get("X-Real-IP")
//		rr.port = resp.Header.Get("X-Real-Port")
//		rr.time = resp.Header.Get("HttpRelay-Time")
//		rr.yourTime = resp.Header.Get("HttpRelay-Your-Time")
//		rr.method = resp.Header.Get("HttpRelay-Method")
//		rr.query = resp.Header.Get("HttpRelay-Query")
//	}
//	c <- rr
//}
//
//func uuid() string {
//	out, _ := exec.Command("uuidgen").Output()
//	return strings.TrimSpace(string(out))
//}
//
//func reqResValidation(req *reqData, resp *respData) (err error) {
//	names := fmt.Sprintf("%s->%s->%s", req.name, req.mode, resp.name)
//	if resp.err != nil { return fmt.Errorf("%s %s", names, resp.err) }
//	if resp.status != http.StatusOK { return fmt.Errorf("%s Response code %d for %s", names, resp.status, resp.url)}
//	if req.mode != resp.mode { return fmt.Errorf("%s Mode mismatch req=%s resp=%s", names, req.mode, resp.mode)}
//	if req.id != resp.id { return fmt.Errorf("%s Id mismatch req=%s resp=%s", names, req.id, resp.id)}
//	if len(*req.body) != len(*resp.body) { return fmt.Errorf("%s Body length mismatch req=%d resp=%d", names, len(*req.body), len(*resp.body))}
//	if req.method != resp.method { return fmt.Errorf("%s Method mismatch req=%s resp=%s", names, req.method, resp.method)}
//	if req.contentType != resp.contentType { return fmt.Errorf("%s Content-type mismatch req=%s resp=%s", names, req.contentType, resp.contentType)}
//	if req.query != resp.query { return fmt.Errorf("%s Query mismatch req=%s resp=%s", names, req.query, resp.query)}
//	return
//}
//
//func permRep(i int, pos int, dic int) (comb []int, err error) {
//	if dic > 36 { return nil, errors.New("dictionary can't ecxeed 36") }
//	str := strings.Repeat("0", pos) + strconv.FormatInt(int64(i), dic)
//	str = str[len(str)-pos:]
//	for p:=pos-1; p>=0; p-- {
//		val, err := strconv.ParseInt(string(str[p]), dic, dic+1)
//		if err != nil { return nil, err }
//		comb = append(comb, int(val))
//	}
//	return
//}
