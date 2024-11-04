package worker

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Service interface {
	CheckExpireSubscriptions(ctx context.Context) (err error)
	CheckBandwidthLimits(ctx context.Context) (err error)
	UpdateBandwidths(ctx context.Context) (err error)
}

type Worker struct {
	service    Service
	stopSignal chan os.Signal
}

func New(service Service, stopSignal chan os.Signal) *Worker {
	return &Worker{
		service:    service,
		stopSignal: stopSignal,
	}
}

func (w *Worker) Run() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			w.Run()
		}
	}()

	w.RunCheckExpireSubscriptions()
	w.RunUpdateBandwidths()
}

func (w *Worker) Stop() {
	w.stopSignal <- os.Interrupt
}

func (w *Worker) RunCheckExpireSubscriptions() {
	go func() {
		for {
			select {
			case <-w.stopSignal:
				return
			default:
				w.service.CheckExpireSubscriptions(context.Background())
				time.Sleep(5 * time.Minute)
			}
		}
	}()
}

func (w *Worker) RunUpdateBandwidths() {
	go func() {
		for {
			select {
			case <-w.stopSignal:
				return
			default:
				w.service.UpdateBandwidths(context.Background())
				time.Sleep(5 * time.Minute)
			}
		}
	}()
}

func (w *Worker) RunCheckBandwidthLimits() {
	go func() {
		for {
			select {
			case <-w.stopSignal:
				return
			default:
				w.service.CheckBandwidthLimits(context.Background())
				time.Sleep(5 * time.Minute)
			}
		}
	}()
}
