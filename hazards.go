package main

import "fmt"

func checkHazards(instructions []instruction, stall bool) ([]instruction, int) {
	var hazards []string
	cycleCount := len(instructions) + 4

	helper := func(curr, next int) bool {
		length := len(instructions)
		if curr >= length || next >= length {
			return false
		}

		current := instructions[curr]
		nextInstr := instructions[next]
		destination := current.arguments[1]

		if current.in == store {
			destination = registerFromMemory(current.arguments[1])
		}

		var nextSources []string
		if len(nextInstr.arguments) == 4 {
			nextSources = nextInstr.arguments[2:]
		} else {
			if nextInstr.in == load {
				nextSources = append(nextSources, registerFromMemory(nextInstr.arguments[2]))
			} else {
				nextSources = append(nextSources, nextInstr.arguments[2])
			}
		}

		if arrayContains(nextSources, destination) {
			// If the current instruction is a store, skip
			if current.in == store {
				return false
			}
			if !stall {
				outputString := fmt.Sprintf(
					"Dependency between instruction #%d (%s) and instruction #%d (%s)",
					current.inNum+1, current.toString(), nextInstr.inNum+1, nextInstr.toString(),
				)
				hazards = append(hazards, outputString)
			} else {
				cycleCount++
				if next-curr == 1 {
					cycleCount++
				}

				stages := []string{"F", "S"}
				laterOffset := 1
				if next-curr == 1 {
					stages = append(stages, "S")
					laterOffset++
				}
				stages = append(stages, "D", "X", "M", "W")
				instructions[next].stages = stages

				instructions = offsetCycleStart(laterOffset, next+1, instructions)
				return true
			}
		} else {
			if stall && len(instructions[curr].stages) == 0 {
				instructions[curr].stages = []string{"F", "D", "X", "M", "W"}
			}
		}

		return false
	}

	for i := range instructions {
		stalled := helper(i, i+1)

		if !stalled {
			_ = helper(i, i+2)
		}

		if !stall {
			instructions[i].stages = []string{"F", "D", "X", "M", "W"}
		}
	}

	if len(hazards) == 0 && !stall {
		fmt.Println("No hazards")
	} else {
		for _, h := range hazards {
			fmt.Println(h)
		}
	}

	return instructions, cycleCount
}

func checkHazardsForwarding(instructions []instruction) ([]instruction, int) {
	cycleCount := len(instructions) + 4

	for i, current := range instructions {
		next := i + 1

		if len(current.stages) == 0 {
			instructions[i].stages = []string{"F", "D", "X", "M", "W"}
		}

		if next >= len(instructions) {
			break
		}

		nextInstr := instructions[next]
		destination := current.arguments[1]

		if current.in == store {
			destination = registerFromMemory(current.arguments[1])
		}

		var nextSources []string
		if len(nextInstr.arguments) == 4 {
			nextSources = nextInstr.arguments[2:]
		} else {
			if nextInstr.in == load {
				nextSources = append(nextSources, registerFromMemory(nextInstr.arguments[2]))
			} else {
				nextSources = append(nextSources, nextInstr.arguments[2])
			}
		}

		if arrayContains(nextSources, destination) {
			stallNeeded := current.in == load && len(nextInstr.arguments) == 4
			if !stallNeeded {
				continue
			}

			cycleCount++
			instructions[next].stages = []string{"F", "D", "S", "X", "M", "W"}
			next++
			if next >= len(instructions) {
				continue
			}

			instructions[next].stages = []string{"F", "S", "D", "X", "M", "W"}

			next++
			if next >= len(instructions) {
				continue
			}

			instructions[next].startCycle++
		}
	}

	return instructions, cycleCount
}
