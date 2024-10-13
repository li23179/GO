# Concurrency Lab 2

## Question 1 - Random sum

Solutions are provided in the zip file:

- `sumAtomic.go`
  - The simplest solution using atomic add operation.
- `sumMutex.go`
  - Also quite simple but now using mutexes. A near identical solution could also be achieved using binary semaphores.
- `sumChan.go`
  - Channel based solution where the delta to add to the sum is sent by each goroutine.
- `sumChan2.go`
  - Channel based solution that uses a parallel tree reduction.

While all of the above solutions are valid, note that this is an example of a problem where channels do not offer the best solution.

## Question 2 - Producer-consumer problem

### Question 2a

4, 994, 5

### Question 2b

See `pc.go`. It uses two semaphores from `github.com/ChrisGora/semaphore` and one mutex from `sync.Mutex`.

## Question 3 - Bank

See the bank directory in the zip. The given solution uses a manager. The manager is the only goroutine that is allowed to do the locking. Other than removing the `lockAccount` method calls, the executors are basically unchanged.

The manager works as follows:

1. Transfer the entire buffered channel of transactions to an internal linked list.
2. It is guaranteed that at least one transaction will be executable at this point (assuming the list is non-empty).
3. Go through the entire list once.
4. If an executable transaction has been found, lock the `from` and `to` accounts and send of the transaction to the executors.
5. Remove the transaction from the linked list.
6. One full iteration through the list ensures that as many transactions have been scheduled as possible.
7. Block on a channel receive and wait for one executor to finish. It is often the case that thanks to one transaction finishing we can execute another one between the two accounts that it freed up - even if two other executors are still operating.
8. Repeat steps 3-7 until the linked list is empty.

`bank2` is an alternative solution. While not the most optimal, it does avoid using a manager which means it's more scalable.
