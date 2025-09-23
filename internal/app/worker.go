package app

import (
	"fmt"
	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/K1la/delayed-notifier/internal/storage"
	"github.com/wb-go/wbf/zlog"
	"time"
)

type Worker struct {
	storage  *storage.MemoryStorage
	interval time.Duration
	done     chan struct{}
}

func NewWorker(storage *storage.MemoryStorage, interval time.Duration) *Worker {
	return &Worker{
		storage:  storage,
		interval: interval,
		done:     make(chan struct{}),
	}
}

func (w *Worker) Start() {
	ticker := time.NewTicker(w.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				w.process()
			case <-w.done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.done)
}

func (w *Worker) process() {
	all := w.storage.GetAll()
	now := time.Now()

	for _, n := range all {
		if n.Status == models.StatusPending && n.SendAt.Before(now) {
			// TODO: сделать реальную реализацию отправки
			fmt.Printf("[WORKER] Sending notification %s via %s to %s: %s\n",
				n.ID, n.Channel, n.To, n.Message)

			n.Status = models.StatusSent
			n.UpdatedAt = time.Now()
			w.storage.Save(n)
			zlog.Logger.Info().
				Str("id", n.ID).
				Str("channel", string(n.Channel)).
				Str("to", n.To).
				Msg("notification sent")
		}
	}
}
