package main

import (
	"fmt"
)

var _ Channeler = &FakeChanneler{}

func testChannel() {
	hookCalled := false
	expected := 9

	f := &FakeChanneler{
		ChannelHook: func(c chan int) chan int {
			hookCalled = true

			go func() {
				c <- expected
			}()

			return c
		},
	}

	input := make(chan int)
	defer close(input)

	output := f.Channel(input)
	found := <-output

	if expected != found {
		panic(fmt.Sprintf("Channel returned unexpected value: %s", found))
	}

	if len(f.ChannelCalls) != 1 {
		panic(fmt.Sprintf("ChannelCalls: %d", len(f.ChannelCalls)))
	}
	if !hookCalled {
		panic("ChannelHook not called")
	}
	if !f.ChannelCalled() {
		panic("ChannelCalled: Channel not called")
	}
	if !f.ChannelCalledOnce() {
		panic("ChannelCalledOnce: Channel not called once")
	}
	if f.ChannelNotCalled() {
		panic("ChannelNotCalled: Channel not called")
	}
	if !f.ChannelCalledN(1) {
		panic("ChannelCalledN: Channel not called once")
	}
}

func testChannelReceive() {
	var found int
	hookCalled := false
	expected := 9

	f := &FakeChanneler{
		ChannelReceiveHook: func(c <-chan int) <-chan int {
			hookCalled = true
			go func() {
				found = <-c
			}()

			return c
		},
	}

	input := make(chan int)
	defer close(input)

	output := f.ChannelReceive(input)

	input <- expected

	if found != expected {
		panic(fmt.Sprintf("ChannelReceive returned unexpected value: %s", found))
	}
	if output != input {
		panic(fmt.Sprintf("ChannelReceive returned unexpected channel: %s", output))
	}

	if len(f.ChannelReceiveCalls) != 1 {
		panic(fmt.Sprintf("ChannelReceiveCalls: %d", len(f.ChannelReceiveCalls)))
	}
	if !hookCalled {
		panic("ChannelReceiveHook not called")
	}
	if !f.ChannelReceiveCalled() {
		panic("ChannelReceiveCalled: ChannelReceive not called")
	}
	if !f.ChannelReceiveCalledOnce() {
		panic("ChannelReceiveCalledOnce: ChannelReceive not called once")
	}
	if f.ChannelReceiveNotCalled() {
		panic("ChannelReceiveNotCalled: ChannelReceive not called")
	}
	if !f.ChannelReceiveCalledN(1) {
		panic("ChannelReceiveCalledN: ChannelReceive not called once")
	}
}

func testChannelSend() {
	hookCalled := false
	expected := 9

	f := &FakeChanneler{
		ChannelSendHook: func(c chan<- int) chan<- int {
			hookCalled = true

			go func() {
				c <- expected
				defer close(c)
			}()

			return c
		},
	}

	input := make(chan int)

	output := f.ChannelSend(input)

	if output != input {
		panic(fmt.Sprintf("ChannelSend returned unexpected value: %s", output))
	}

	if len(f.ChannelSendCalls) != 1 {
		panic(fmt.Sprintf("ChannelSendCalls: %d", len(f.ChannelSendCalls)))
	}
	if !hookCalled {
		panic("ChannelSendHook not called")
	}
	if !f.ChannelSendCalled() {
		panic("ChannelSendCalled: ChannelSend not called")
	}
	if !f.ChannelSendCalledOnce() {
		panic("ChannelSendCalledOnce: ChannelSend not called once")
	}
	if f.ChannelSendNotCalled() {
		panic("ChannelSendNotCalled: ChannelSend not called")
	}
	if !f.ChannelSendCalledN(1) {
		panic("ChannelSendCalledN: ChannelSend not called once")
	}
}

func testChannelPointer() {
	hookCalled := false
	expected := 9

	f := &FakeChanneler{
		ChannelPointerHook: func(cp *chan int) *chan int {
			hookCalled = true
			c := *cp

			go func() {
				c <- expected
				defer close(c)
			}()

			return &c
		},
	}

	input := make(chan int)

	output := f.ChannelPointer(&input)

	if *output != input {
		panic(fmt.Sprintf("ChannelPointer returned unexpected value: %s", output))
	}

	if len(f.ChannelPointerCalls) != 1 {
		panic(fmt.Sprintf("ChannelPointerCalls: %d", len(f.ChannelPointerCalls)))
	}
	if !hookCalled {
		panic("ChannelPointerHook not called")
	}
	if !f.ChannelPointerCalled() {
		panic("ChannelPointerCalled: ChannelPointer not called")
	}
	if !f.ChannelPointerCalledOnce() {
		panic("ChannelPointerCalledOnce: ChannelPointer not called once")
	}
	if f.ChannelPointerNotCalled() {
		panic("ChannelPointerNotCalled: ChannelPointer not called")
	}
	if !f.ChannelPointerCalledN(1) {
		panic("ChannelPointerCalledN: ChannelPointer not called once")
	}
}

func testChannelInterface() {
	hookCalled := false
	expected := &FakeChanneler{}

	f := &FakeChanneler{
		ChannelInterfaceHook: func(c chan interface{}) chan interface{} {
			hookCalled = true

			go func() {
				c <- expected
			}()

			return c
		},
	}

	input := make(chan interface{})
	defer close(input)

	output := f.ChannelInterface(input)
	found := <-output

	if expected != found {
		panic(fmt.Sprintf("ChannelInterface returned unexpected value: %s", found))
	}

	if len(f.ChannelInterfaceCalls) != 1 {
		panic(fmt.Sprintf("ChannelInterfaceCalls: %d", len(f.ChannelInterfaceCalls)))
	}
	if !hookCalled {
		panic("ChannelInterfaceHook not called")
	}
	if !f.ChannelInterfaceCalled() {
		panic("ChannelInterfaceCalled: ChannelInterface not called")
	}
	if !f.ChannelInterfaceCalledOnce() {
		panic("ChannelInterfaceCalledOnce: ChannelInterface not called once")
	}
	if f.ChannelInterfaceNotCalled() {
		panic("ChannelInterfaceNotCalled: ChannelInterface not called")
	}
	if !f.ChannelInterfaceCalledN(1) {
		panic("ChannelInterfaceCalledN: Channel not called once")
	}
}

func main() {
	testChannel()
	testChannelReceive()
	testChannelSend()
	testChannelPointer()
	testChannelInterface()
}
