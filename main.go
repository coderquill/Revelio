package main

import (
	"github.com/opencareer/interview-excercises/ex_02/lib"
	"log"
	"time"
)

// ================================= EXERCISE 2 =================================

// PART 1
// You are developing a pipeline that fetches tasks from an external queue (think of
// it as a RabbitMQ queue or something equivalent) and process them. You will use
// a library to manage the connection with the queue (`lib` folder). You should read
// that lib just to read the documentation of public items.
//
// Tasks contain details about which computation to perform and how, but for this
// exercise we suppose these details are abstracted away: to execute the task simply
// call the `task.do()` method.
//
// When all tasks for the day are sent, the queue sends a notification and the lib
// closes the channel you use the get tasks (see `.listen()`). If new tasks are
// pushed into the external queue but the queue is full, the queue will overflow and
// drop the task (producer back-pressure). A message is sent via the connection and
// the lib will log the event and keep a count of dropped tasks.
//
// Implement the pipeline in a simple, single-threaded way:
// - fetch tasks and execute them
// - how many tasks have we processed?
// - how many tasks have we lost?
//
// PART 2
// Recently the number of tasks has increased and the back-pressure has increased. We
// want to improve the pipeline to reduce the amount of dropped tasks.
//
// Improve the service:
// - which approaches could you use? what are pros and cons of each?
// - update the pipeline code in order to not lose any task
// - concurrency must be bounded (if present)
//
// PART 3
// We want to react to certain external signals properly (see `listenCancellation()`).
// When a cancellation signal arrives we want to stop fetching tasks from the queue.
// Specifically we want to shut down the connection with the queue, wait the end of
// tasks currently being processed and exit (see `shutdown()` docs).
//
// Improve the service:
// - listen for cancellation signals
// - shutdown the service properly

func main() {
	start := time.Now()
	duration := time.Since(start)
	processedTasks, droppedTasks, err := startProcessing()

	if err != nil {
		log.Fatalf("An error encountered: %s\n", err)
	}

	log.Printf("Tasks processing started at: %s\n", start)
	log.Printf("Processed Tasks: %d\n", processedTasks)
	log.Printf("Dropped Tasks: %d\n", droppedTasks)
	log.Printf("Total duration: %s\n", duration)
}

// Task 1: single-threaded design will process each task one by one as they come in
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

	// Retrieve the number of dropped tasks using lib function.
	droppedTasks = taskQueueConn.Dropped()

	// If no err occurred, return all the counts.
	return processedTasks, droppedTasks, nil
}
