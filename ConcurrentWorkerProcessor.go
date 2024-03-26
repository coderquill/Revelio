package main

import (
	"github.com/opencareer/interview-excercises/ex_02/lib" // Importing the task library
	"log"
	"sync"
	"sync/atomic"
)

// Part 2

type ConcurrentWorkerProcessor struct{}

// ProcessTasks implements the TaskProcessor interface for concurrent processing
func (cwp ConcurrentWorkerProcessor) ProcessTasks() (int, uint64, error) {
	return startProcessingWithWorkers() // Delegates task processing to a dedicated function
}

const workerCount = 10 // Number of workers to process tasks concurrently

// Handles task processing using multiple workers to allow concurrent task execution
func startProcessingWithWorkers() (int, uint64, error) {
	taskQueueConn, err := lib.NewTaskQueueConn() // Establishes a new connection to the task queue
	if err != nil {
		return 0, 0, err // Returns early on connection error
	}

	tasks := taskQueueConn.Listen()      // Starts listening for tasks from the queue
	cancellation := listenCancellation() // Listens for cancellation signals to gracefully terminate

	taskChannel := make(chan lib.Task, workerCount) // Creates a channel to distribute tasks to workers
	var waitGroup sync.WaitGroup                    // WaitGroup to synchronize the completion of all workers
	var processedTasks int32                        // Counter for the number of successfully processed tasks

	// Initializes worker goroutines
	for i := 0; i < workerCount; i++ {
		waitGroup.Add(1)
		go worker(&waitGroup, taskChannel, &processedTasks) // Each worker processes tasks from the channel
	}

	// Distributes tasks to workers and handles cancellation
DistributeTasks:
	for {
		select {
		case task, ok := <-tasks: // Receives tasks from the queue
			if !ok {
				break DistributeTasks // Exits if the task channel is closed
			}
			select {
			case taskChannel <- task: // Attempts to send the task to a worker
			default:
				log.Println("Temporarily unable to queue task for worker (worker queue full)")
			}
		case <-cancellation: // Reacts to cancellation signals
			break DistributeTasks
		}
	}

	close(taskChannel) // Closes the task channel to signal workers to stop
	waitGroup.Wait()   // Waits for all workers to finish processing

	droppedTasks := taskQueueConn.Dropped() // Retrieves the count of dropped tasks

	return int(processedTasks), droppedTasks, nil // Returns the processing results
}

// Represents a worker that processes tasks
func worker(waitGroup *sync.WaitGroup, taskChannel <-chan lib.Task, processedTasks *int32) {
	defer waitGroup.Done()          // Signals completion upon returning
	for task := range taskChannel { // Processes each task received from the channel
		if err := task.Do(); err != nil {
			log.Printf("Failed to process task: %v", err)
		} else {
			log.Printf("Processed a task ")
			atomic.AddInt32(processedTasks, 1) // Increments the task counter
		}
	}
}
