# Revelio

================================= EXERCISE 2 =================================
### **Design Decisions and Considerations:**

**Interface-Driven Design**
A key architectural decision was to define a TaskProcessor interface,
encapsulating the task processing logic behind a consistent API. I took this decision to enable flexibility in processing strategy selection and simplifying the integration of additional strategies in the future.

**User Interaction**
Recognizing the importance of adaptability in operational environments, I included a console-based user interaction model. Users can select their preferred processing strategy(Also useful to keep track of which part we are testing) at runtime, choosing between single-threaded processing (Part 1), concurrent processing (Part 2), or initiating a graceful shutdown (Part 3) in response to cancellation signals.
This design empowers users to choose the system's behavior to their current needs.

**Concurrency Management**
The concurrent processing model and graceful shutdown mechanism extensively leverages Go's concurrency primitives, such as goroutines, channels, and sync packages. These tools allowed me to manage task distribution effectively, balance workload across workers, and synchronize system shutdown processes, ensuring that all components operated correctly.

### Alternative design considered: The strategy pattern
The strategy pattern could have been a perfect fit for our implementation in managing different task processing strategies. The strategy pattern is a behavioral design pattern that enables selecting an algorithm's behavior at runtime. In the context of our task processing system, this translates to choosing between different task processing strategies (single-threaded, concurrent workers, and graceful shutdown) based on runtime conditions or user input.

However, I chose not to use the strategy pattern explicitly here. Why? Mainly because Go's interface system inherently encourages a form of the strategy pattern without needing to formalize it. By designing our system around the TaskProcessor interface, we essentially allowed for changing the task processing behavior at runtime based on the user's choice, which is at the heart of what the strategy pattern aims to achieve.

In simpler terms, while I didn't set out to apply the strategy pattern by its textbook definition, the approach I took naturally aligns with the pattern's objectives. Go's interfaces and type system made it straightforward to implement this pattern-like behavior implicitly, demonstrating the language's power and flexibility in supporting various design patterns through its idiomatic constructs.


## PART 1
You are developing a pipeline that fetches tasks from an external queue (think of
it as a RabbitMQ queue or something equivalent) and process them. You will use
a library to manage the connection with the queue (`lib` folder). You should read
that lib just to read the documentation of public items.

Tasks contain details about which computation to perform and how, but for this
exercise we suppose these details are abstracted away: to execute the task simply
call the `task.do()` method.

When all tasks for the day are sent, the queue sends a notification and the lib
closes the channel you use the get tasks (see `.listen()`). If new tasks are
pushed into the external queue but the queue is full, the queue will overflow and
drop the task (producer back-pressure). A message is sent via the connection and
the lib will log the event and keep a count of dropped tasks.

Implement the pipeline in a simple, single-threaded way:
- fetch tasks and execute them
- how many tasks have we processed?
- how many tasks have we lost?

### My Approach for Part 1:
I started developing a pipeline to fetch tasks from an external queue, like you'd imagine with RabbitMQ, and to process these tasks using a lib I had at my disposal. The main thing was that the details of each task were sort of a black box - I just had to run task.do() to execute them.

When the day's tasks were all sent, I'd get a notification, and the lib would shut down the task channel. If the queue got too full and couldn't take more tasks, it would just drop them, which wasn't ideal.

So, my job was to make a simple system where:

I fetch tasks and execute them,
Keep track of how many tasks I've managed to process,
And note how many I've lost because the queue was too full.

## PART 2
Recently the number of tasks has increased and the back-pressure has increased. We
want to improve the pipeline to reduce the amount of dropped tasks.

Improve the service:
- which approaches could you use? what are pros and cons of each?
- update the pipeline code in order to not lose any task
- concurrency must be bounded (if present)

For Part 2, I needed to figure out how to handle more tasks without dropping any because of too much pressure on the system. I decided to use Go's awesome ability to do many things at once, called concurrency, with something called goroutines and channels. This basically means setting up a bunch of workers that can work on different tasks at the same time.

### My Approach for Part 2: Concurrent Workers

**What I did**: I set up a bunch of workers, all listening to a single line of communication, or a "channel," waiting for tasks to do. Each worker grabs a task from this channel as soon as they can and starts working on it. This way, many tasks can get done all at once, rather than one by one. I also put a limit on how many workers can be working at the same time and how many tasks can wait in line to be worked on.

**Why I did it**: This method is pretty straightforward but super effective. It makes sure the system can handle a lot more tasks without getting too complicated. Plus, Go is really good at handling this kind of setup, where lots of little tasks are being worked on at the same time.

### Other Ideas I Thought About

1. **Changing the Number of Workers on the Fly**: This would mean adding more workers if things get busy and taking some away when it's not so busy. It sounds smart because it uses resources wisely, but it's pretty complicated to actually do.

2. **Limiting How Fast Tasks Come In**: By slowing down how fast tasks are added, I could prevent the system from getting overwhelmed. It's a simpler fix but might not use all the system's potential power.

3. **Deciding Which Tasks Are Most Important**: Making a system that picks which tasks to do first based on how important they are. This is great if some tasks need to be done right away, but it makes deciding which task to do next more complicated.

### Why I Picked Concurrent Workers

I went with the concurrent workers because it's a good middle ground. It's not too hard to set up, and it really boosts how many tasks can be done at once. It's like hiring more people for a job – more hands on deck means getting more done, without making things too complicated. Plus, it really plays to Go's strengths, which is all about doing lots of things at the same time, in a simple and efficient way.

So, that's the gist of it. I aimed for something simple but effective, without diving into more complex or fancy solutions that might not even be needed. It's all about getting the job done well, with as little fuss as possible.

## PART 3
We want to react to certain external signals properly (see `listenCancellation()`).
When a cancellation signal arrives we want to stop fetching tasks from the queue.
Specifically we want to shut down the connection with the queue, wait the end of
tasks currently being processed and exit (see `shutdown()` docs).

Improve the service:
- listen for cancellation signals
- shutdown the service properly

### My Approach for Part 3:
The new challenge was to properly handle cancellation signals with listenCancellation(). If a signal came through, I needed to stop fetching new tasks, let the current ones finish, and then neatly shut everything down.

### What I did: 
I implemented a GracefulShutdownProcessor that could listen for these cancellation signals. When a signal was received, it would stop fetching tasks, wait for any ongoing tasks to wrap up, and then close the connection to the queue.
