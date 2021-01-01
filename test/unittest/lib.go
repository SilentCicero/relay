package unittest

//func prepare() (proxyRep *repository.ProxyRep, proxyCtrl *controller.ProxyCtrl, stopChan *closechan.CloseChan) {
//	stopChan = closechan.NewCloseChan()
//	proxyRep = repository.NewProxyRep()
//	proxyCtrl = controller.NewProxyCtrl(proxyRep, stopChan.C())
//	return
//}

func newData(s string, count int) (b []byte) {
	b = make([]byte, len(s)*count)
	bp := copy(b, s)
	for bp < len(b) {
		copy(b[bp:], b[:bp])
		bp *= 2
	}
	return
}
