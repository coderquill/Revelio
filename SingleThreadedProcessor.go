package main

import (
	"github.com/opencareer/interview-excercises/ex_02/lib"
	"log"
)

// Part 1

type SingleThreadedProcessor struct{}

func (stp SingleThreadedProcessor) ProcessTasks() (int, uint64, error) {
	return startProcessing()
}

// Single-threaded design will process each task one by one as they come in
func startProcessing() (processedTasks int, droppedTasks uint64, err error) {
	// Establish a new connection to the task queue.
	taskQueueConn, err := lib.NewTaskQueueConn()
	if err != nil {
		return 0, 0, err
	}

	// Listen to the task queue.
	tasks := taskQueueConn.Listen()

	// Process tasks
	for task := range tasks {
		//If err occurred:
		if err := task.Do(); err != nil {
			log.Printf("Failed to process a task: %s\n", err)
			//If tasks file had a function to retrieve a task id for a particular task,
			//we can log that this particular task failed
		} else {
			log.Printf("Processed a task")
			processedTasks++
		}
	}

	// Retrieves the number of dropped tasks using lib function.
	droppedTasks = taskQueueConn.Dropped()

	// If no err occurred, return all the counts.
	return processedTasks, droppedTasks, nil
}
