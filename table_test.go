package clt

import (
	"fmt"
	"github.com/ttacon/chalk"
	"reflect"
	"testing"
)

func TestTerminalSizeCheck(t *testing.T) {
	h, w, err := getTerminalSize()
	if err != nil || h == -1 || w == -1 {
		fmt.Printf("Cannot determine terminal size for Table.  McClintock will still work, but will not be able to automagically determine sizes.")
	}
}

func TestCreateTable(t *testing.T) {
	table := NewTable(3)
	if table.columns != 3 {
		t.Errorf("Table should have %d columns, has %d.", 3, table.columns)
	}
	justExpect := []string{"l", "l", "l"}
	if !reflect.DeepEqual(table.justify, justExpect) {
		t.Errorf("Table default justification should be %v, got %v.", justExpect, table.justify)
	}
	if table.MaxHeight == 0 || table.MaxWidth == 0 {
		t.Error("Table should have a width/height.")
	}
}

func TestBasicAddRow(t *testing.T) {
	want := []string{"test1", "test2"}
	table := NewTable(len(want))
	table.AddRow(want)
	if len(table.rows) != 1 {
		t.Errorf("Table should have %d rows, has %d.", len(want), len(table.rows))
	}
	for idx, _ := range want {
		gotX := table.rows[0].cells[idx].value
		if gotX != want[idx] {
			t.Errorf("Table cell 1 should be %s, got %s.", want[idx], gotX)
		}
	}
}

func TestShortAddRow(t *testing.T) {
	want := []string{"test1", "test2"}
	table := NewTable(len(want) + 1)
	table.AddRow(want)
	gotEmpty := table.rows[0].cells[2].value
	if gotEmpty != "" {
		t.Errorf("Table cell 2 should be empty, got %s.", gotEmpty)
	}
}

func TestCellLength(t *testing.T) {
	want := []string{"---7---", "---8----"}
	table := NewTable(len(want))
	table.AddRow(want)
	for idx, _ := range want {
		lenX := table.rows[0].cells[idx].width
		if lenX != len(want[idx]) {
			t.Errorf("Table cell length should be %v, got %v.", len(want[idx]), lenX)
		}
	}
}

func TestTitle(t *testing.T) {
	want := "This is the title"
	table := NewTable(1)
	table.SetTitle(want)
	if table.title.value != want {
		fmt.Errorf("Title should be %s, got %s", want, table.title.value)
	}
	if table.title.width != len(want) {
		fmt.Errorf("Title length should be %v, got %v.", len(want), table.title.width)
	}
	if table.title.style != chalk.Bold {
		fmt.Errorf("Default title style should be chalk.Bold.")
	}
}
