package lib

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// TREAT THIS FILE AS A BLACK-BOX LIB

// NewTaskQueueConn establish a new connection with an external queue
// and returns a TaskQueueConn struct representing that connection.
func NewTaskQueueConn() (TaskQueueConn, error) {
	return TaskQueueConn{
		reqs: make(chan Task, 50),
		stop: make(chan bool),
		wg:   sync.WaitGroup{},
		dc:   0,
	}, nil
}

type TaskQueueConn struct {
	reqs chan Task
	stop chan bool
	wg   sync.WaitGroup
	dc   uint64
}

type Task struct {
	Do func() error
}

// Listen starts fetching tasks from the external queue, delivering
// them to the returned channel. See the Shutdown() method for proper
// termination.
func (tqc *TaskQueueConn) Listen() <-chan Task {
	tqc.wg.Add(1)
	go func() {
		defer func() {
			close(tqc.reqs)
			tqc.wg.Done()
		}()

		dur := time.Millisecond * 25
		ticker := time.NewTicker(dur)
		defer ticker.Stop()

		for i := 0; i <= 1000; i++ {
			select {
			case <-ticker.C:
				select {
				case tqc.reqs <- Task{Do: func() error {
					time.Sleep(dur * 15)
					return nil
				}}:
				default:
					log.Println("[task-queue] WARN: filled queue, task discarded")
					atomic.AddUint64(&tqc.dc, 1)
				}
			case <-tqc.stop:
				log.Println("[task-queue] WARN: stopping")
				return
			}
		}
	}()
	return tqc.reqs
}

// Shutdown will close gracefully the connection. When this method returns
// the connection is terminated and the Listen() chan closed, but some
// tasks may still be in that channel.
func (tqc *TaskQueueConn) Shutdown() {
	close(tqc.stop)
	tqc.wg.Wait()
}

// Dropped returns the number of tasks Dropped by the external queue.
func (tqc *TaskQueueConn) Dropped() uint64 {
	return atomic.LoadUint64(&tqc.dc)
}
