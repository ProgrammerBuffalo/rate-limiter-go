package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ProgrammerBuffalo/rate-limiter-go/limiter"
)

type Job struct {
	Id int
	Do func()
}

type Worker struct {
	limiter *limiter.Limiter
	ch      <-chan *Job

	id int
}

func NewWorker(id int, limiter *limiter.Limiter, ch <-chan *Job) *Worker {
	return &Worker{
		id:      id,
		limiter: limiter,
		ch:      ch,
	}
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range w.ch {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Worker %d start to process request #%d\n", w.id, job.Id)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		w.runWithTimeout(ctx, cancel, job)
	}
}

func (w *Worker) runWithTimeout(ctx context.Context, cancellation context.CancelFunc, job *Job) {
	defer cancellation()

	fmt.Printf("Woker %d check request #%d\n", w.id, job.Id)

	if w.limiter.Allow() {
		done := make(chan struct{})

		go func() {
			job.Do()
			close(done)
		}()

		select {
		case <-done:
			fmt.Printf("Worker %d finish to process request #%d\n", w.id, job.Id)
		case <-ctx.Done():
			fmt.Printf("Worker %d couldn't finish to process request #%d (time exceeded)\n", w.id, job.Id)
		}

	} else {
		fmt.Printf("Worker %d is not allowed to process request #%d (requests count exceeded)\n", w.id, job.Id)
	}
}
