package main

// TaskProcessor is an interface that outlines the structure for different task processing strategies.
// It's the blueprint for how tasks should be processed, regardless of the specific method used.
//
// The idea here is pretty straightforward but also quite powerful. By defining this interface,
// I'm setting up a contract that any task processing strategy must follow. This means that regardless
// of whether I'm processing tasks one by one, all at once with multiple workers, or even shutting down
// gracefully in response to a signal, I have to adhere to this specific way of doing things.
//
// ProcessTasks is the method that any type implementing the TaskProcessor interface needs to define.
// It's responsible for the heavy lifting of actually getting the tasks done. More specifically,
// it will:
// - Process the tasks it receives in whatever way it sees fit, based on the specific strategy implemented.
// - Return the total number of tasks that were processed successfully.
// - Return the number of tasks that couldn't be processed, perhaps because the queue was too full and they were dropped.
// - And, if anything goes wrong during the processing, it returns an error to let me know what happened.
type TaskProcessor interface {
	ProcessTasks() (processedTasks int, droppedTasks uint64, err error)
}
