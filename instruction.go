package main

import (
	"fmt"
	"strings"
)

type instruction struct {
	in         string
	inNum      int
	startCycle int
	arguments  []string
	stages     []string
}

func (i instruction) toString() string {
	command := i.in
	args := i.arguments[1:]
	return fmt.Sprintf("%s %s", command, strings.Join(args, ", "))
}
