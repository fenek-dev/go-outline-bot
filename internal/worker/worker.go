package worker

import (
	"context"
	"os"
)

type Service interface {
	CheckProlongations(ctx context.Context) (err error)
	CheckExpireSubscriptions(ctx context.Context) (err error)
	UpdateBandwidths(ctx context.Context) (err error)
	CheckBandwidthLimits(ctx context.Context) (err error)
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

func (w *Worker) Start() {
	w.RunCheckExpireSubscriptions()
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
			}
		}
	}()
}

func (w *Worker) RunCheckProlongations() {
	go func() {
		for {
			select {
			case <-w.stopSignal:
				return
			default:
				w.service.CheckProlongations(context.Background())
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
			}
		}
	}()
}
