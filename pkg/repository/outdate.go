package repository

import (
	"time"
)

type Outdater interface {
	removeOutdated()
}

func Outdate(outdaters []Outdater, delay time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			break
		case <-stopChan:
			break
		}

		for _, o := range outdaters {
			o.removeOutdated()
		}
	}
}
