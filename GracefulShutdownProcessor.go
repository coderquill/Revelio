package main

import (
	"fmt"
	"github.com/opencareer/interview-excercises/ex_02/lib"
	"log"
	"sync"
	"sync/atomic"
)

// Part 3
// GracefulShutdownProcessor handles task processing with the capability to gracefully shutdown.
type GracefulShutdownProcessor struct{}

// ProcessTasks processes tasks with a mechanism to handle cancellation signals for graceful shutdown.
func (gsp GracefulShutdownProcessor) ProcessTasks() (int, uint64, error) {
	// Establish a new connection with the task queue.
	taskQueueConn, err := lib.NewTaskQueueConn()
	if err != nil {
		return 0, 0, err
	}

	// Start listening for tasks from the queue.
	tasks := taskQueueConn.Listen()

	// Listen for a cancellation signal to gracefully shutdown.
	cancellation := listenCancellation()

	// Create a channel for tasks to be distributed to workers.
	taskChannel := make(chan lib.Task, workerCount)
	var waitGroup sync.WaitGroup
	var processedTasks int32

	// Initialize a fixed number of workers to process tasks concurrently.
	for i := 0; i < workerCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for task := range taskChannel {
				// Attempt to process the task.
				if err := task.Do(); err != nil {
					log.Printf("Failed to process task: %v", err)
				} else {
					// Successfully processed the task.
					atomic.AddInt32(&processedTasks, 1)
				}
			}
		}()
	}

	// Loop to distribute tasks to workers, listening for the cancellation signal.
DistributeTasks:
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				// If the task channel is closed, indicating no more tasks, exit the loop.
				break DistributeTasks
			}
			// Attempt to send the task to a worker for processing.
			select {
			case taskChannel <- task:
				// Task successfully sent to a worker.
			default:
				// The task queue for workers is full, unable to queue the task.
				log.Println("Task queue is full, unable to queue task for worker.")
			}
		case <-cancellation:
			// Received a cancellation signal, initiating shutdown procedure.
			fmt.Println("Cancellation signal received. Initiating shutdown...")
			taskQueueConn.Shutdown() // Close the task queue connection to stop fetching new tasks.
			break DistributeTasks
		}
	}

	close(taskChannel) // Close the task channel, signaling all workers to finish up.
	waitGroup.Wait()   // Wait for all workers to complete their processing.

	// Fetch the count of dropped tasks during the processing.
	droppedTasks := taskQueueConn.Dropped()

	return int(processedTasks), droppedTasks, nil
}
