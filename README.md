# Go Concurrency Patterns Playground

This repository contains notes, sketches and code examples and experiments based on Go code presented during various conference talks.

## Introduction

What is `concurrency`?

- It's the composition of independently executing computations.
- It's a way to structure software, as a way to write a clean code that interacts with the real world.

Concurrent features of *Go*:

- Easy to understand
- Easy to use
- Easy to reason about
- No need to be an expert to start using them!

## Examples

Boring example:

```go

```

## Goroutines

What is a goroutine? It's an independently executing function

## Synchronization

When the main function executes `<-c`, it will waitin for a value to be sent.
Similarly, when the boring function executes `c <- value`. it waits for a receiver to be ready.

A sender and receiver *must both be ready* to play their part in the communication. Otherwise we wait until they are.

Thus *channles both communicate and synchronize*.

### Buffered Channels

Go channels can also be created with a buffer.

*Buffering removes synchronization!*

Buffering makes channels more like Erlang's mailboxes.

Buffered channels can be important for some problems but they are more subtle to reason about.

## Go approach

*Don't communicate by sharing memory, share memory by communicating.*

## Patterns

### Genertor - function that returns a channel

Channels are first-class values, just liek strings or integers!

#### Channels as a handle on a service

Our boring function returns a channel that lets us communicate with the boring service it provides.

We can have more instances of the service!

```go
func main() {
    joe := boring("Joe")
    mark := boring("Mark")

    for i := 0; i < 4; i++ {
        fmt.Println(<-joe)
        fmt.Println(<-ann)
    }
    fmt.Println("You're both boring. I am leaving")
}
```

Example:

```go
func boringGenerator(msg string) <-chan string { // Returns receive-only channel of strings.
 c := make(chan string)
 go func() { // We launch the goroutine from iside the function.
  for i := 0; ; i++ {
   c <- fmt.Sprintf("%s %d", msg, i)
   time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
  }
 }()
 return c // Return the channel to the caller.
}

func RunMainGenerator() {
    // create a channel by calling boringGenerator func
    c := boringGenerator("Hello from boring generator!")

 // create a loop and take values from the channel and print them
 for i := 0; i < 6; i++ {
  fmt.Printf("You say: %q\n", <-c)
 }

 // finishing func by printing out the message.
 fmt.Println("You are boring; I am leaving for good!")
}
```

Example: (services)

```go
func boringGenerator(msg string) <-chan string { // Returns receive-only channel of strings.
 c := make(chan string)
 go func() { // We launch the goroutine from iside the function.
  for i := 0; ; i++ {
   c <- fmt.Sprintf("%s %d", msg, i)
   time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
  }
 }()
 return c // Return the channel to the caller.
}

func RunMainGeneratorService() {
 joe := boringGenerator("Joe")
 mark := boringGenerator("Mark")
 for i := 0; i < 5; i++ {
  // important note about synchronization here:
  // If joe is not ready yet, mark won't be able to send values.
  // In other words mark needs to wait for joe
  fmt.Println(<-joe)
  fmt.Println(<-mark)
 }
 fmt.Println("You are both boring. I am leaving.")
}
```

Note: `mark` needs to wait for `joe` - synchronization issues! What to do if `mark` is more talkative? How to refactor and address this?

### Multiplexing (fanIn)

We can use `fan-in` pattern (function) to let whoever is ready to talk! No waiting like the previous example.

### Resoring sequencing - send chan inside a chan

todo

### Select

Control structure unique to concurrency.

The reason channels and goroutines are built into the language.

The `select` statement provides another way to handle multiple channels. It's like `switch`, but each case is a communication:

- all channels are evaluated
- selection blocks until one communication can proceed, which then does
