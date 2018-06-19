package main

import (
	"fmt"
	"sync"
)

const (
	defaultPrefix = "G"
)

var (
	defaultSymbolGenerator = symbolGenerator{Prefix: defaultPrefix}
	defaultMutex           sync.Mutex
)

func gensym() string {
	defaultMutex.Lock()
	defer defaultMutex.Unlock()
	return defaultSymbolGenerator.next()
}

type symbolGenerator struct {
	Prefix string
	Suffix string
	count  uint64
}

func (s *symbolGenerator) next() string {
	s.count++
	return fmt.Sprintf("%s%d%s", s.Prefix, s.count, s.Suffix)
}

func (s *symbolGenerator) reset() {
	s.count = 0
}
