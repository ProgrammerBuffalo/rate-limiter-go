package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ProgrammerBuffalo/rate-limiter-go/limiter"
	"github.com/ProgrammerBuffalo/rate-limiter-go/worker"
)

func main() {
	wg := sync.WaitGroup{}
	limiter := limiter.NewLimiter(time.Duration(1)*time.Second, 5)

	ch := make(chan *worker.Job, 100)

	go func() {
		for i := range 100 {
			ch <- &worker.Job{
				Id: i,
				Do: doRequest,
			}
		}
		close(ch)
	}()

	for i := range 5 {
		wg.Add(1)
		w := worker.NewWorker(i, limiter, ch)
		go w.Run(&wg)
	}

	wg.Wait()
}

func doRequest() {
	sec := rand.Intn(2)
	time.Sleep(time.Duration(sec) * time.Second)
}
