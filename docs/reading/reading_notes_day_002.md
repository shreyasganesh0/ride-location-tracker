# Go Routines and Channels

## Go Routines
- lightweight threads managed by the go runtime
- go f(x, y, z) is the format

## Channels
- a type of buffer used to send and recive values
    - ch <- v; // to send value of v to channel ch
    - v := <-ch // to get value from ch and store it in v
    - x, y := <-ch, <-ch // each call to  <-ch will emit a single value stored inFIFO
- make a channel like maps and slices
    - ch := make(chan int)
- default behaviour is blocking for channels // allow go routines to sync wihtout lock

- Buffered Channels
    - ch := make(chan int, n)
    - creates a buffered channel that can store n values
    - <-ch, <-ch each call gives out one value
        - only blocks if full when sending to and empty when recieving from channel
- Range and Close
    - senders can close a channel to indicate no more values will be sent
        - close(c)
    - recievers can test this with v, ok := <-ch
        - ok will return false if ch is closed
    - dont have to close its only a way to inform recievers

    - we can get values from a a channel using range
        - for i := range c {}
        - if no close happens it blocks
- Select
    - select lets a goroutine wait on multiple communication operations
    select {
        case c<-x:
            x, y = y, x + y
        case <-quit:
            fmt.Println("quitting");
        default:
            
        }
    - chooses one at random  if multiple are ready
    - use default to send and recieve without blocking since it runs if no
      other case is ready
- Mutex
    - sync.Mutex.Lock and Unlock
    - used to lock critical sections of the code

