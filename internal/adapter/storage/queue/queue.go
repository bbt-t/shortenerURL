package queue

import (
	"context"
	"log"
	"sync"
)

type Queue struct {
	name   string
	jobs   chan Job
	ctx    context.Context
	cancel context.CancelFunc
}

type Job struct {
	Name   string
	Action func() error
}

type Worker struct {
	Queue *Queue
}

func NewQueue(name string) *Queue {
	ctx, cancel := context.WithCancel(context.Background())

	return &Queue{
		jobs:   make(chan Job),
		name:   name,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (q *Queue) AddJobs(jobs []Job) {
	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(job Job) {
			q.AddJob(job)
			wg.Done()
		}(job)
	}

	go func() {
		wg.Wait()
		q.cancel()
	}()

}

func (q *Queue) AddJob(job Job) {
	q.jobs <- job
}

func (j Job) Run() error {
	log.Printf("Job running: %s", j.Name)

	err := j.Action()
	if err != nil {
		return err
	}
	return nil
}

func NewWorker(queue *Queue) *Worker {
	return &Worker{
		Queue: queue,
	}
}

func (w *Worker) DoWork() bool {
	for {
		select {
		// if context was canceled.
		case <-w.Queue.ctx.Done():
			log.Printf("Work done in queue %s: %s!", w.Queue.name, w.Queue.ctx.Err())
			return true
		// if job received.
		case job := <-w.Queue.jobs:
			err := job.Run()
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
}
