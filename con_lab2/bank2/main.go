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

// An executor is a type of a worker goroutine that handles the incoming transactions.
func executor(bank *bank, executorId int, transactionQueue <-chan transaction, done chan<- bool) {
	for {
		t := <-transactionQueue

		fmt.Println("Executor\t", executorId, "attempting transaction from", t.from, "to", t.to)

		// Find order of account and attempt lock accordingly
		if t.from > t.to {
			bank.lockAccount(t.from, strconv.Itoa(executorId))
			fmt.Println("Executor\t", executorId, "locked account", t.from)
			bank.lockAccount(t.to, strconv.Itoa(executorId))
			fmt.Println("Executor\t", executorId, "locked account", t.to)
		} else {
			bank.lockAccount(t.to, strconv.Itoa(executorId))
			fmt.Println("Executor\t", executorId, "locked account", t.to)
			bank.lockAccount(t.from, strconv.Itoa(executorId))
			fmt.Println("Executor\t", executorId, "locked account", t.from)
		}

		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.

		bank.execute(t, executorId)

		bank.unlockAccount(t.from, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", t.from)
		bank.unlockAccount(t.to, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", t.to)

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

	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, transactionQueue, done)
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
