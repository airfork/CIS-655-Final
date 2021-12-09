package main

import "github.com/xuri/excelize/v2"

type sheetInfo struct {
	f           *excelize.File
	name        string
	createSheet bool
	hazardFunc  func() ([]instruction, int)
}

func addCenterStyle(f *excelize.File, sheetName, startCol, endCol string) error {
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
	})

	if err != nil {
		return err
	}

	return f.SetCellStyle(sheetName, startCol, endCol, style)
}

func addHeaderStyles(f *excelize.File, sheetName, startCol, endCol string) error {
	style, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
	})

	if err != nil {
		return err
	}

	return f.SetCellStyle(sheetName, startCol, endCol, style)
}

func setHeaders(f *excelize.File, sheetName string, colLength int) error {
	endColumn, err := excelize.ColumnNumberToName(colLength + 1)
	if err != nil {
		return err
	}

	topRow := make([]int, colLength)
	for i := 0; i < colLength; i++ {
		topRow[i] = i + 1
	}

	headerFunctions := []func() error{
		func() error {
			return f.SetCellValue(sheetName, "B1", "Cycles")
		},
		func() error {
			return f.MergeCell(sheetName, "B1", endColumn+"1")
		},
		func() error {
			return addCenterStyle(f, sheetName, "A1", endColumn+"2")
		},
		func() error {
			return addHeaderStyles(f, sheetName, "B1", endColumn+"1")
		},
		func() error {
			return addHeaderStyles(f, sheetName, "A2", "A2")
		},
		func() error {
			return f.SetCellValue(sheetName, "A2", "Instructions")
		},
		func() error {
			return f.SetColWidth(sheetName, "A", "A", 20)
		},
		func() error {
			return f.SetSheetRow(sheetName, "B2", &topRow)
		},
	}

	for _, function := range headerFunctions {
		err = function()
		if err != nil {
			return err
		}
	}
	return nil
}

func addInstruction(f *excelize.File, sheetName string, in instruction) error {
	row := []string{in.toString()}
	cycle := in.startCycle

	for cycle > 0 {
		row = append(row, " ")
		cycle--
	}

	row = append(row, in.stages...)
	return f.SetSheetRow(sheetName, appendNumber("A", in.inNum+3), &row)
}

func createSheet(info sheetInfo) error {
	if info.createSheet {
		_ = info.f.NewSheet(info.name)
	}

	instructions, colLength := info.hazardFunc()
	endColumn, err := excelize.ColumnNumberToName(colLength + 1)
	checkErr(err)

	err = setHeaders(info.f, info.name, colLength)
	checkErr(err)

	for i, in := range instructions {
		err = addInstruction(info.f, info.name, in)
		checkErr(err)

		err = addCenterStyle(info.f, info.name, appendNumber("A", i+3), appendNumber(endColumn, i+3))
		checkErr(err)
	}
	return nil
}
