package model

type Link struct {
	linkChan chan *LinkData
	comm
}

func NewLink() *Link {
	return &Link{
		linkChan: make(chan *LinkData),
		comm:     newComm(),
	}
}

func (l *Link) Chan() chan *LinkData {
	return l.linkChan
}
