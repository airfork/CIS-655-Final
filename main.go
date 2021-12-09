package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	checkErr(err)

	instructions, err := getInstructions(content)
	checkErr(err)

	// Avoids needing to deep copy
	instructionsCopy, _ := getInstructions(content)

	f := excelize.NewFile()
	sheetName := "No Stalls"
	// creating new sheet for consistency, default sheet has different cell width than created ones
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	sheets := []sheetInfo{
		{
			f:           f,
			name:        sheetName,
			createSheet: false,
			hazardFunc: func() ([]instruction, int) {
				return checkHazards(instructions, false)
			},
		},
		{
			f:           f,
			name:        "No Forwarding",
			createSheet: true,
			hazardFunc: func() ([]instruction, int) {
				return checkHazards(instructions, true)
			},
		},
		{
			f:           f,
			name:        "Forwarding",
			createSheet: true,
			hazardFunc: func() ([]instruction, int) {
				return checkHazardsForwarding(instructionsCopy)
			},
		},
	}

	for _, sheet := range sheets {
		err = createSheet(sheet)
		checkErr(err)
	}

	// Save the spreadsheet by the given path.
	if err := f.SaveAs("output.xlsx"); err != nil {
		fmt.Println(err)
	}
}
