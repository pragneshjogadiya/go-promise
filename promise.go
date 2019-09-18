package main

import (
	"fmt"
	"sync"
	//"sync"
	//"time"
)

var wgGlobal sync.WaitGroup

const (
	pending = iota
	fulfilled
	rejected
)

type promise struct {
	//wg    sync.WaitGroup
	value string
	err   error
	state int
}

type ggPromise interface {
	then(func(string), func(error))
	catch(func(error))
	finally(func(string))
}

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

func (p *promise) then(onFulfilled func(string), onRejected func(error)) {
	//go func() {
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

func (p *promise) catch(onRejected func(error)) {
	//go func() {
	func() {
		wgGlobal.Wait()
		if p.err != nil {
			onRejected(p.err)
			p.state = rejected
			return
		}
		//r(p.res)
	}()
}

func (p *promise) finally(onFinally func(string)) {
	//go func() {
	//func() {
	wgGlobal.Wait()
	//if p.err != nil {
	//	e(p.err)
	//	return
	//}
	onFinally(p.value)
	p.state = fulfilled
	//}()
}

func exampleTicker() (string, error) {
	//<-time.Tick(time.Second * 1)
	return "hi", nil
}

func main() {
	//doneChan := make(chan int)
	var p = newPromise(exampleTicker)
	// fmt.Println("Hello")
	// p.Then(func(result string) { fmt.Println(result);  doneChan <- 1}, func(err error) { fmt.Println(err) })
	// <-doneChan

	// var p = promise{
	// 	value: "Hello",
	// 	err:   nil,
	// 	state: pending,
	// }

	var gg ggPromise = p

	//wgGlobal.Add(1)
	gg.then(func(result string) { fmt.Println(result) }, func(err error) { fmt.Println(err) })

	//wgGlobal.Add(1)
	gg.catch(func(err error) { fmt.Println(err) })

	//wgGlobal.Add(1)
	gg.finally(func(result string) { fmt.Println(result) })

	//wgGlobal.Wait()

}
