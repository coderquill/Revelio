package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	//Provide choice to users for which processing stategy they want to implement
	fmt.Println("Choose the task processing strategy:")
	fmt.Println("1: Single-threaded (Part 1)")
	fmt.Println("2: Concurrent workers (Part 2)")
	fmt.Println("3: Process tasks with Graceful Shutdown  (Part 3)")
	fmt.Print("Enter choice (1,2 or 3): ")

	choice, _ := reader.ReadString('\n')

	var processor TaskProcessor

	//Based on the user choice, use the corresponding processor strategy
	switch choice {
	case "1\n":
		processor = SingleThreadedProcessor{}
	case "2\n":
		processor = ConcurrentWorkerProcessor{}
	case "3\n":
		processor = GracefulShutdownProcessor{}
	default:
		fmt.Println("Invalid choice. Exiting.")
		return
	}

	start := time.Now()
	processedTasks, droppedTasks, err := processor.ProcessTasks()
	duration := time.Since(start)

	if err != nil {
		log.Fatalf("An error encountered: %s\n", err)
	}

	log.Printf("Tasks processing started at: %s\n", start)
	log.Printf("Processed Tasks: %d\n", processedTasks)
	log.Printf("Dropped Tasks (queue overflow): %d\n", droppedTasks)
	log.Printf("Total duration: %s\n", duration)
}
