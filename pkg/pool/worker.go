package pool

import "log"

type JobFunc func(id int64, data interface{})

type Job struct {
	Data    interface{}
	JobFunc JobFunc
}

type Worker struct {
	ID            int64
	WorkerChannel chan chan *Job
	Channel       chan *Job
	End           chan struct{}
	jobFinished   chan bool
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				if job != nil {
					job.JobFunc(w.ID, job.Data)
					w.jobFinished <- true
				}
			case <-w.End:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- struct{}{}
}
