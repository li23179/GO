package main

import (
	"container/list"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var debug *bool

func manager(bank *bank, executorId int, transactionQueue <-chan transaction, readyTranscationQueue chan transaction, done chan bool, cond_Accounts map[*account]*sync.Cond) {
	// TODO: goes through the transaction queue and schedules the transactions in an optimal way.
	// manager should be the only thread doing the locking

	// manager schedule work -> send ready transaction to the internal queue -> executor execute
	for {
		t := <-transactionQueue

		fromAccount := bank.getAccount(t.from)
		toAccount := bank.getAccount(t.to)

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		// schedules logic:
		// one of the account is locked then the other account needs to wait
		// send the transaction to buffered chan until both locks are ready

		// But given a transaction from A to B an executor should only be allowed
		// to lock account A if it can lock account B.

		bank.lockAccount(t.from, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "locked account", from)

		bank.lockAccount(t.to, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "locked account", to)

		// while one of them is locked, release the other lock and wait (suspense the thread)
		// use condition variable

		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			for fromAccount.locked {
				cond_Accounts[toAccount].Wait()
			}
			wg.Done()
		}()

		go func() {
			for toAccount.locked {
				cond_Accounts[fromAccount].Wait()
			}
			wg.Done()
		}()

		wg.Wait()

		readyTranscationQueue <- t

		<-done

		bank.unlockAccount(t.from, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", from)

		bank.unlockAccount(t.to, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", to)
	}
}

// An executor is a type of a worker goroutine that handles the incoming transactions.
func executor(bank *bank, executorId int, readyTranscationQueue <-chan transaction, done chan<- bool) {
	for {
		t := <-readyTranscationQueue

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)
		fmt.Println("Executor\t", executorId, "attempting transaction from", from, "to", to)
		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.

		bank.execute(t, executorId)

		bank.removeCompleted(e, executorId) // Removing this line will break visualisations.
		done <- true
	}
}

func toChar(i int) rune {
	return rune('A' + i)
}

// main creates a bank and executors that will be handling the incoming transactions.
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	debug = flag.Bool("debug", false, "generate DOT graphs of the state of the bank")
	flag.Parse()

	bankSize := 6 // Must be even for correct visualisation.
	transactions := 1000

	accounts := make([]*account, bankSize)
	for i := range accounts {
		accounts[i] = &account{name: string(toChar(i)), balance: 1000}
	}

	bank := bank{
		accounts:               accounts,
		transactionsInProgress: list.New(),
		gen:                    newGenerator(),
	}

	startSum := bank.sum()

	transactionQueue := make(chan transaction, transactions)
	expectedMoneyTransferred := 0
	for i := 0; i < transactions; i++ {
		t := bank.getTransaction()
		expectedMoneyTransferred += t.amount
		transactionQueue <- t
	}

	done := make(chan bool)

	// Initialise condition variables for each account
	cond_Accounts := make(map[*account]*sync.Cond, len(bank.accounts))
	readyTranscationQueue := make(chan transaction, 3)

	for _, account := range bank.accounts {
		cond_Accounts[account] = sync.NewCond(&account.mutex)
	}

	for i := 0; i < bankSize; i++ {
		go manager(&bank, i, transactionQueue, readyTranscationQueue, done, cond_Accounts)
		go executor(&bank, i, readyTranscationQueue, done)
	}

	for total := 0; total < transactions; total++ {
		fmt.Println("Completed transactions\t", total)
		<-done
	}

	fmt.Println()
	fmt.Println("Expected transferred", expectedMoneyTransferred)
	fmt.Println("Actual transferred", bank.moneyTransferred)
	fmt.Println("Expected sum", startSum)
	fmt.Println("Actual sum", bank.sum())
	if bank.sum() != startSum {
		panic("sum of the account balances does not much the starting sum")
	} else if len(transactionQueue) > 0 {
		panic("not all transactions have been executed")
	} else if bank.moneyTransferred != expectedMoneyTransferred {
		panic("incorrect amount of money was transferred")
	} else {
		fmt.Println("The bank works!")
	}
}
