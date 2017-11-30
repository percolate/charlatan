package main

type Mapper interface {
	MapParameter(map[string]string)
	MapReturn() map[string]string
}
