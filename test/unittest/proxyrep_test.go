package unittest

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"testing"
)

func TestNewProxyRepGetSer(t *testing.T) {
	const aSerId = "myserverA"
	const bSerId = "myserverB"

	pr := repository.NewProxyRep()

	serA1 := pr.GetSer(aSerId)
	serA2 := pr.GetSer(aSerId)

	if serA1 != serA2 {
		t.Fail()
	}

	serB := pr.GetSer(bSerId)

	if serA1 == serB {
		t.Fail()
	}
}
