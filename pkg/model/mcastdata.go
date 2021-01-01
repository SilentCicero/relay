package model

type McastData struct {
	changeChan chan struct{}
	data       *TeeData
}

func NewMcastData() *McastData {
	return &McastData{
		changeChan: make(chan struct{}),
	}
}

func (md *McastData) Close() {
	select {
	case <-md.changeChan:
		return
	default:
		close(md.changeChan) // Releasing read waiters
	}
}

func (md *McastData) Write(data *TeeData) {
	md.data = data
	close(md.changeChan)
}

func (md *McastData) Size() int {
	select {
	case <-md.changeChan:
		return md.data.Content.Size()
	default:
		return 0
	}
}

//func (this *McastData) IsEmpty() bool {
//	return this.meta == nil
//}

func (md *McastData) Read(closeChan <-chan struct{}) (data *TeeData, ok bool) {
	select {
	case <-md.changeChan:
		ok = true
		data = md.data
	case <-closeChan:
	}
	return
}
