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

func manager(bank *bank, transactionQueue <-chan transaction, executors chan<- transaction, done <-chan bool) {
	waitingTransactions := list.New()
	for t := range transactionQueue {
		waitingTransactions.PushBack(t)
	}
	for waitingTransactions.Len() > 0 {
		element := waitingTransactions.Front()
		for element != nil {
			t := element.Value.(transaction)
			if !bank.accounts[t.from].locked && !bank.accounts[t.to].locked {
				bank.lockAccount(t.from, "M")
				fmt.Println("Manager\t", "locked account", t.from)
				bank.lockAccount(t.to, "M")
				fmt.Println("Manager\t", "locked account", t.to)
				executors <- t
				old := element
				element = element.Next()
				waitingTransactions.Remove(old)
			} else {
				element = element.Next()
			}
		}
		<-done
	}
	fmt.Println("Manager finished")
}

// An executor is a type of a worker goroutine that handles the incoming transactions.
func executor(bank *bank, executorId int, transactionQueue <-chan transaction, doneMain chan<- bool, doneManager chan<- bool) {
	for {
		t := <-transactionQueue

		e := bank.addInProgress(t, executorId)

		fmt.Println("Executor\t", executorId, "executing transaction from", t.from, "to", t.to)
		bank.execute(t, executorId)

		fmt.Println("Executor\t", executorId, "trying to unlock account", t.from)
		bank.unlockAccount(t.from, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", t.from)

		fmt.Println("Executor\t", executorId, "trying to unlock account", t.to)
		bank.unlockAccount(t.to, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "unlocked account", t.to)

		bank.removeCompleted(e, executorId)

		doneMain <- true
		doneManager <- true
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

	bankSize := 6
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
	close(transactionQueue)

	internalQueue := make(chan transaction, 10)
	doneManager := make(chan bool, bankSize)
	go manager(&bank, transactionQueue, internalQueue, doneManager)

	doneMain := make(chan bool)

	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, internalQueue, doneMain, doneManager)
	}

	for total := 0; total < transactions; total++ {
		fmt.Println("Completed transactions\t", total)
		<-doneMain
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
