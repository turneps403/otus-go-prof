package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 0 {
		m = 0
	}
	// no buffers
	reqCh, succCh, errCh, graceCh := make(chan Task), make(chan bool), make(chan bool), make(chan bool)
	defer close(graceCh)
	defer close(errCh)
	defer close(succCh)
	defer close(reqCh)

	// run workers
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(reqCh, succCh, errCh, graceCh)
		}()
	}

	for i, done, send := 0, 0, 0; ; {
		// run tasks and read result
		select {
		case reqCh <- tasks[i]:
			send++
			if i == len(tasks)-1 {
				// no body will listen this new dummy-channel, and written is impossible
				reqCh = make(chan Task)
				defer close(reqCh)
			} else {
				i++
			}
		case <-succCh:
			done++
		case <-errCh:
			m--
			done++
		}

		if m < 0 || done == len(tasks) {
			// we reach end of tasks or collect limit errors
			// await last answers
			for done < send {
				select {
				case <-succCh:
					done++
				case <-errCh:
					m--
					done++
				}
			}

			// inform go-routine to shut itself down
			for i := 0; i < n; i++ {
				graceCh <- true
			}

			break
		}
	}
	wg.Wait()

	if m < 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(req <-chan Task, succ chan<- bool, err chan<- bool, graceCh <-chan bool) {
	for {
		select {
		case t := <-req:
			if t() == nil {
				succ <- true
			} else {
				err <- true
			}
		case <-graceCh:
			return
		}
	}
}
