package main

import (
	"fmt"
	"sync"
	"time"
)

type Dispatcher struct {
	WorkerPool chan chan Job
	maxWokers int
}

func NewDispatcher(maxWorkers int)  *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Job, maxWorkers),
		maxWokers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	for i:=0; i<d.maxWokers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel<-job
			}(job)
		}
	}
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit chan bool
}

func NewWorker(workerPool chan chan Job) *Worker{
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit: make(chan bool),
	}
}

var wg sync.WaitGroup

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool<-w.JobChannel

			select {
			case job := <-w.JobChannel:
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("worker received job %d\n", job.id)
				// wg.Done()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit<-true
	}()
}

type Job struct {
	id int
}

// public
var JobQueue = make(chan Job)

func main() {

	dispatcher := NewDispatcher(100_000)
	dispatcher.Run()

	id := 0

	for {
		input := 0 
		fmt.Scanf("%d", &input)
		if input == 1 {
			fmt.Println("spawning job")
			JobQueue<-Job{id}
			id += 100
		} else if input == 2 {
			now := time.Now()
			var max = 100_000
			for i:=0; i<max; i++ {
				go func(i int) {JobQueue<-Job{i}}(i)
			}
			fmt.Printf("spawned %d jobs\ntook %d", max, time.Since(now))
		} else if input == -1 {
			return
		}
	}
}