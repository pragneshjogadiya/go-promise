package main

import (
	"fmt"
	"sync"
	"time"
)

//Global Wait Group
var wgGlobal sync.WaitGroup

//Constants pending, fulfilled and rejected which represents state of promise
const (
	pending = iota
	fulfilled
	rejected
)

//Struct Promise that will have value, error and state(State can be pending, fulfilled or rejected)
type promise struct {
	value string
	err   error
	state int
}

//Promise Interface with then, catch and finally methods
type promiseInterface interface {
	then(func(string), func(error))
	catch(func(error))
	finally(func(int))
}

//Function to create new promise pointer
func newPromise(f func() (string, error)) *promise {
	p := &promise{}
	wgGlobal.Add(1)
	go func() {
		p.value, p.err = f()
		p.state = pending
		wgGlobal.Done()
	}()
	return p
}

//then method on struct promise pointer. Takes functions onFulfilled and onRejected as arguments.
func (p *promise) then(onFulfilled func(string), onRejected func(error)) {
	func() {
		wgGlobal.Wait()
		if p.err != nil {
			onRejected(p.err)
			p.state = rejected
			return
		}
		onFulfilled(p.value)
		p.state = fulfilled
	}()
}

//catch method on struct promise pointer. Takes function onRejected as argument.
func (p *promise) catch(onRejected func(error)) {
	func() {
		wgGlobal.Wait()
		if p.err != nil {
			onRejected(p.err)
			p.state = rejected
			return
		}
	}()
}

//finally method on struct promise pointer. Takes function onFinally as argument.
func (p *promise) finally(onFinally func(int)) {
	func() {
		wgGlobal.Wait()
		onFinally(p.state)
	}()
}

//You can implement your own function to return (string,error) pair.
func exampleFunction() (string, error) {
	<-time.Tick(time.Second * 1)

	//Return response with nil error
	return "Hello There", nil

	//Return response with error. Uncomment below line and comment above line to return error. Don't forget to import errors package.
	//return "", errors.New("Error in example function")
}

func main() {
	//Creating Promise
	var promiseCreated = newPromise(exampleFunction)

	//Creating Promise Interface
	var p promiseInterface = promiseCreated

	p.then(
		func(result string) {
			fmt.Println("Inside then method")
			fmt.Println(result)
		},
		func(err error) {
			fmt.Println("Inside then method")
			fmt.Println("Error Occured")
			fmt.Println(err)
		})

	p.catch(func(err error) {
		fmt.Println("Inside catch method")
		fmt.Println("Error Occured")
		fmt.Println(err)
	})

	p.finally(func(result int) {
		if result == 1 {
			fmt.Println("Promise Fulfilled")
		} else if result == 2 {
			fmt.Println("Promise Rejected")
		} else {
			fmt.Println("Error Occured")
		}
	})

}
