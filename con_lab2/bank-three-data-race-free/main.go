// This version completed by Harry Greentree (2023)
package main

import (
	"container/list"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var debug *bool

// A manager is a type of worker goroutine that manages the transaction queue and schedules transactions.
func manager(bank *bank, transactionQueue chan transaction, executorQueue chan transaction) {
	executor := "Manager"
	for {
		// Pop a transaction from the transaction queue
		transaction := <-transactionQueue
		from := transaction.from
		to := transaction.to

		// We lock the bank to avoid a data race where we check if an account is locked
		// and at the same time unlock an account in an executor thread.
		bank.lock(executor)
		if !bank.isAccountLocked(from) && !bank.isAccountLocked(to) {
			bank.unlock()
			bank.lockAccount(from, executor)
			fmt.Println("Manager has locked account:", from)

			bank.lockAccount(to, executor)
			fmt.Println("Manager has locked account:", to)

			// Hand the transaction to the executor queue to be handled.
			executorQueue <- transaction
		} else {
			bank.unlock()
			// Add any transactions that cant be handled to the end of the transaction queue
			// to be scheduled later.
			transactionQueue <- transaction
		}
	}
}

// An executor is a type of a worker goroutine that handles the incoming transactions.
func executor(bank *bank, executorId int, executorQueue <-chan transaction, done chan<- bool) {
	for {
		t := <-executorQueue

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		fmt.Println("Executor\t", executorId, "attempting transaction from", from, "to", to)
		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.
		bank.execute(t, executorId)

		bank.unlockAccount(t.from, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", from)
		bank.unlockAccount(t.to, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", to)

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

	// Initialize a queue of transactions to be executed by the executor threads.
	// We use a buffer size of 3 as we can only have 3 concurrent transactions at a time.
	executorQueue := make(chan transaction, 3)

	go manager(&bank, transactionQueue, executorQueue)

	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, executorQueue, done)
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
