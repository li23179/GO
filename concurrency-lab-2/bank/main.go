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

func manager(bank *bank, executorId int, transactionQueue <-chan transaction, readyTranscationQueue chan transaction, done chan bool) {
func manager(bank *bank, executorId int, transactionQueue <-chan transaction, readyTranscationQueue chan transaction, done chan bool) {
	// TODO: goes through the transaction queue and schedules the transactions in an optimal way.
	// manager should be the only thread doing the locking

	// manager schedule work -> send ready transaction to the internal queue -> executor execute
	for {
		t := <-transactionQueue

		exe := strconv.Itoa(executorId)

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		fromAcc := bank.getAccount(t.from)
		toAcc := bank.getAccount(t.to)

		for !fromAcc.locked && !toAcc.locked{
			bank.lockAccount(t.from, exe)
			fmt.Println("Executor\t", executorId, "locked account", from)

			bank.lockAccount(t.to, exe)
			fmt.Println("Executor\t", executorId, "locked account", to)

			readyTranscationQueue <- t
			break
		}
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
	readyTranscationQueue := make(chan transaction, transactions)

	expectedMoneyTransferred := 0
	for i := 0; i < transactions; i++ {
		t := bank.getTransaction()
		expectedMoneyTransferred += t.amount
		transactionQueue <- t
	}

	done := make(chan bool)

	for i := 0; i < bankSize; i++ {
		go manager(&bank, i, transactionQueue, readyTranscationQueue, done)
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
