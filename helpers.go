package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	add   string = "add"
	sub          = "sub"
	load         = "lw"
	store        = "sw"
)

var validInstructions = []string{add, sub, store, load}
var memoryPattern = `\d+\(\$[0-9A-z]+\)`

func arrayContains(arr []string, s string) bool {
	for _, val := range arr {
		if val == s {
			return true
		}
	}

	return false
}

func isRegister(str string) bool {
	return strings.HasPrefix(str, "$")
}

func validateThreeReg(args []string) error {
	for _, reg := range args {
		if !isRegister(reg) {
			return errors.New("'" + reg + "' is not a register")
		}
	}
	return nil
}

func validateMemory(args []string, command string) error {
	register := args[0]
	memory := args[1]

	if command == store {
		memory = args[0]
		register = args[1]
	}

	if !isRegister(register) {
		return errors.New("'" + register + "' is not a register")
	}

	ok, err := regexp.Match(memoryPattern, []byte(memory))
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("'" + memory + "' is not a valid memory access")
	}
	return nil
}

func registerFromMemory(m string) string {
	ok, err := regexp.Match(memoryPattern, []byte(m))
	if err != nil || !ok {
		return ""
	}

	re := regexp.MustCompile(`\$[0-9A-z]+`)
	return string(re.Find([]byte(m)))
}

func offsetCycleStart(offset, start int, instructions []instruction) []instruction {
	for i := start; i < len(instructions); i++ {
		instructions[i].startCycle += offset
	}

	return instructions
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func appendNumber(str string, num int) string {
	return fmt.Sprintf("%s%d", str, num)
}

func parseLine(line string) ([]string, error) {
	split := strings.Split(strings.TrimSpace(line), " ")
	if len(split) == 0 {
		return nil, nil
	}

	command := split[0]
	if !arrayContains(validInstructions, command) {
		return nil, errors.New(fmt.Sprintf("'%s' is not a valid command", command))
	}

	expectedArgLength := 3
	if command == validInstructions[2] || command == validInstructions[3] {
		expectedArgLength = 2
	}

	arguments := strings.Split(strings.Join(split[1:], ""), ",")
	if len(arguments) != expectedArgLength {
		return nil, errors.New(fmt.Sprintf("invalid number of arguments for '%s'. Expected: %d, got: %d", command, expectedArgLength, len(arguments)))
	}

	err := func() error {
		if len(arguments) == 3 {
			return validateThreeReg(arguments)
		}
		return validateMemory(arguments, command)
	}()

	if err != nil {
		return nil, err
	}

	arguments = append([]string{command}, arguments...)

	return arguments, nil
}

func getInstructions(content []byte) ([]instruction, error) {
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	var instructions []instruction

	for i, line := range lines {
		args, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		in := instruction{
			in:         args[0],
			startCycle: i,
			inNum:      i,
			arguments:  args,
			stages:     []string{},
		}
		instructions = append(instructions, in)
	}

	return instructions, nil
}
