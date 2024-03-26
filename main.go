package main

import (
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
	err := startProcessing()
	log.Printf("Total duration: %s", time.Since(start))
	if err != nil {
		log.Fatalln(err)
	}
}

func startProcessing() error {
	panic("TODO")
}
