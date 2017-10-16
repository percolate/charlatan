package main

import (
	"fmt"
	"sync"
)

const (
	defaultPrefix = "G"
)

var (
	defaultSymbolGenerator = SymbolGenerator{Prefix: defaultPrefix}
	defaultMutex           sync.Mutex
)

func Gensym() string {
	defaultMutex.Lock()
	defer defaultMutex.Unlock()
	return defaultSymbolGenerator.Next()
}

type SymbolGenerator struct {
	Prefix string
	Suffix string
	count  uint64
}

func (s *SymbolGenerator) Next() string {
	s.count++
	return fmt.Sprintf("%s%d%s", s.Prefix, s.count, s.Suffix)
}
