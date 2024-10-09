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

// An executor is a type of a worker goroutine that handles the incoming transactions.
func executor(bank *bank, executorId int, transactionQueue <-chan transaction, done chan<- bool, cond_Accounts map[*account]*sync.Cond) {
	for {
		t := <-transactionQueue

		from := bank.getAccountName(t.from)
		to := bank.getAccountName(t.to)

		var accountFrom *account
		var accountTo *account

		for _, account := range bank.accounts{
			if account.name == from {
				accountFrom = account
			} else if account.name == to {
				accountTo = account
			}
		}

		fmt.Println("Executor\t", executorId, "attempting transaction from", from, "to", to)
		e := bank.addInProgress(t, executorId) // Removing this line will break visualisations.

		bank.lockAccount(t.from, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "locked account", from)

		bank.lockAccount(t.to, strconv.Itoa(executorId))
		fmt.Println("Executor\t", executorId, "locked account", to)

		// But given a transaction from A to B an executor should only be allowed
		// to lock account A if it can lock account B.
		// while both are lock, release the mutex lock and wait (suspense the thread)
		// use condition variable
		
		// Find the matching account mutex
		accountFromCond_var := cond_Accounts[accountFrom]
		accountToCond_var := cond_Accounts[accountTo]

		bothLocked := accountFrom.locked && accountTo.locked

		// while both are lock, release the mutex lock and wait (suspense the thread)		
		for ; !bothLocked;{
			accountFromCond_var.Wait()
			accountToCond_var.Wait()
		}

		bank.execute(t, executorId)

		bank.unlockAccount(t.from, strconv.Itoa(executorId))
		accountFromCond_var.Broadcast()
		fmt.Println("Executor\t", executorId, "unlocked account", from)

		bank.unlockAccount(t.to, strconv.Itoa(executorId))
		accountToCond_var.Broadcast()
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

	// Initialise condition variables for each account
	cond_Accounts := make(map[*account]*sync.Cond, len(bank.accounts))

	for _, account := range bank.accounts{
		cond_Accounts[account] = sync.NewCond(&account.mutex)
	}
	
	for i := 0; i < bankSize; i++ {
		go executor(&bank, i, transactionQueue, done, cond_Accounts)
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
