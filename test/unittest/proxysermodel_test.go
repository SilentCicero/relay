package unittest

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"testing"
)

func TestNewProxySerModel(t *testing.T) {
	psm := model.NewProxySer()

	select {
	case psm.ReqChan <- nil:
		t.Fail()
	default:
	}
}

func TestNewProxySerModelAddTakeJob(t *testing.T) {
	const jobId = "12345678"
	psm := model.NewProxySer()

	data, _, _, _, _, _, _, _ := newProxyCliData()
	psm.AddJob(jobId, data)

	if takenData, ok := psm.TakeJob(jobId); ok {
		if takenData != data {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if _, ok := psm.TakeJob(jobId); ok {
		t.Fail()
	}
}

func TestNewProxySerModelAddRemoveJob(t *testing.T) {
	const jobId = "12345678"
	psm := model.NewProxySer()

	data, _, _, _, _, _, _, _ := newProxyCliData()
	psm.AddJob(jobId, data)

	psm.RemoveJob(data)

	if _, ok := psm.TakeJob(jobId); ok {
		t.Fail()
	}
}
