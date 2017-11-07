package main

type Channeler interface {
	Channel(chan int) chan int
	ChannelReceive(<-chan int) <-chan int
	ChannelSend(chan<- int) chan<- int
	ChannelPointer(*chan int) *chan int
	ChannelInterface(chan interface{}) chan interface{}
}
