package clt

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/BTBurke/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestTerminalSizeCheck(t *testing.T) {
	t.Skipf("Terminal check skipped. No TTY.")
	h, w, err := getTerminalSize()
	if err != nil || h == -1 || w == -1 {
		fmt.Printf("Cannot determine terminal size for Table. This will still work, but will not be able to automagically determine sizes.")
	}
}

func TestCreateTable(t *testing.T) {
	table := NewTable(3)
	if len(table.columns) != 3 {
		t.Errorf("Table should have %d columns, has %d.", 3, len(table.columns))
	}
	for _, col := range table.columns {
		if !reflect.DeepEqual(col.justify, Left) {
			t.Errorf("Table default justification should be %v, got %v.", Left, col.justify)
		}
	}
	if table.maxHeight == 0 || table.maxWidth == 0 {
		t.Error("Table should have a width/height.")
	}

}

func TestBasicAddRow(t *testing.T) {
	want := []string{"test1", "test2"}
	table := NewTable(len(want))
	table.AddRow(want...)
	if len(table.rows) != 1 {
		t.Errorf("Table should have %d rows, has %d.", len(want), len(table.rows))
	}
	for idx := range want {
		gotX := table.rows[0].cells[idx].value
		if gotX != want[idx] {
			t.Errorf("Table cell 1 should be %s, got %s.", want[idx], gotX)
		}
	}
}

func TestShortAddRow(t *testing.T) {
	want := []string{"test1", "test2"}
	table := NewTable(len(want) + 1)
	table.AddRow(want...)
	gotEmpty := table.rows[0].cells[2].value
	if gotEmpty != "" {
		t.Errorf("Table cell 2 should be empty, got %s.", gotEmpty)
	}
}

func TestCellLength(t *testing.T) {
	want := []string{"---7---", "---8----"}
	table := NewTable(len(want))
	table.AddRow(want...)
	for idx := range want {
		lenX := table.rows[0].cells[idx].width
		if lenX != len(want[idx]) {
			t.Errorf("Table cell length should be %v, got %v.", len(want[idx]), lenX)
		}
	}
}

func TestTitle(t *testing.T) {
	want := "This is the title"
	table := NewTable(1)
	table.Title(want)
	if table.title.value != want {
		t.Errorf("Title should be %s, got %s", want, table.title.value)
	}
	if table.title.width != len(want) {
		t.Errorf("Title length should be %v, got %v.", len(want), table.title.width)
	}
	assert.Equal(t, table.title.style, Styled(Bold))
}

func TestSetColumnHeaders(t *testing.T) {
	table := NewTable(2)
	table.ColumnHeaders("Header1", "Header2")
	want := []Cell{Cell{
		value: "Header1",
		style: Styled(Bold, Underline),
		width: len("Header1"),
	},
		Cell{
			value: "Header2",
			style: Styled(Bold, Underline),
			width: len("Header2"),
		},
	}
	for i, header := range table.headers {
		if header.value != want[i].value {
			t.Errorf("Header should be %s, got %s", want[i].value, header.value)
		}
	}
}

func TestRenderHelpers(t *testing.T) {
	n := []int{1, 3, 2}
	table := &Table{}
	table.columns = []Col{Col{
		naturalWidth:  10,
		computedWidth: 12,
	},
		Col{
			naturalWidth:  12,
			computedWidth: 14,
		},
	}
	table.pad = 2

	assert.Equal(t, mapAdd(n, 1), []int{2, 4, 3})
	assert.Equal(t, sum(n), 6)
	i, m := max(n)
	assert.Equal(t, m, 3)
	assert.Equal(t, i, 1)
	assert.Equal(t, sumWithoutIndex(n, 1), 3)
	assert.True(t, wrappedWidthOk(51, 100))
	assert.False(t, wrappedWidthOk(49, 100))
	assert.Equal(t, extractComputedWidth(table), []int{12, 14})
	assert.Equal(t, extractNatWidth(table), []int{10, 12})
	assert.Equal(t, table.width(), 12+14+8)

}

// test helper to get string of length n
func s(n int) string {
	return strings.Repeat("x", n)
}

func TestSimpleStrategy(t *testing.T) {
	table := NewTable(3)
	table.maxWidth = 80
	table.ColumnHeaders(s(4), s(4), s(4))
	table.AddRow(s(10), s(12), s(14))
	table.pad = 2
	table.computeColWidths()
	t.Run("NaturalColWidths < maxWidth", func(t *testing.T) {
		assert.Equal(t, extractNatWidth(table), []int{10, 12, 14})
		assert.Equal(t, extractComputedWidth(table), []int{10, 12, 14})
	})

	// headers bigger than content
	table.ColumnHeaders(s(15), s(16), s(17))
	table.computeColWidths()
	t.Run("Big headers, NaturalWidth < maxWidth", func(t *testing.T) {
		assert.Equal(t, extractNatWidth(table), []int{15, 16, 17})
		assert.Equal(t, extractComputedWidth(table), []int{15, 16, 17})
	})
}

func TestWrapWidest(t *testing.T) {
	table := NewTable(3)
	table.maxWidth = 60
	table.AddRow(s(10), s(20), s(40))

	t.Run("Last column wrap, no padding", func(t *testing.T) {
		table.pad = 0
		table.computeColWidths()
		assert.Equal(t, extractNatWidth(table), []int{10, 20, 40})
		assert.Equal(t, extractComputedWidth(table), []int{10, 20, 30})
	})

}

func TestOverflow(t *testing.T) {
	table := NewTable(3)
	table.maxWidth = 10
	table.AddRow(s(10), s(20), s(40))

	t.Run("Overflow to natural width as last resort", func(t *testing.T) {
		table.pad = 0
		table.computeColWidths()
		assert.Equal(t, extractNatWidth(table), []int{10, 20, 40})
		assert.Equal(t, extractComputedWidth(table), []int{10, 20, 40})
	})

}

func TestJustifcation(t *testing.T) {
	s := s(4)
	t.Run("Center justify text with padding", func(t *testing.T) {
		width := 14
		pad := 2
		sty := Styled(Default)
		want := fmt.Sprintf("       %s       ", sty.ApplyTo(s))
		assert.Equal(t, justCenter(s, width, pad, sty), want)
	})
	t.Run("Center justify offest left on uneven", func(t *testing.T) {
		width := 13
		pad := 2
		sty := Styled(Default)
		want := fmt.Sprintf("      %s       ", sty.ApplyTo(s))
		assert.Equal(t, justCenter(s, width, pad, sty), want)
	})
	t.Run("Left justify text with padding", func(t *testing.T) {
		width := 8
		pad := 2
		sty := Styled(Default)
		want := fmt.Sprintf("  %s      ", sty.ApplyTo(s))
		assert.Equal(t, justLeft(s, width, pad, sty), want)
	})
	t.Run("Right justify text with padding", func(t *testing.T) {
		width := 8
		pad := 2
		sty := Styled(Default)
		want := fmt.Sprintf("      %s  ", sty.ApplyTo(s))
		assert.Equal(t, justRight(s, width, pad, sty), want)
	})
	t.Run("Fallback to string + padding if widths jacked up", func(t *testing.T) {
		width := 1
		pad := 2
		sty := Styled(Default)
		want := fmt.Sprintf("  %s  ", sty.ApplyTo(s))
		assert.Equal(t, justCenter(s, width, pad, sty), want)
		assert.Equal(t, justLeft(s, width, pad, sty), want)
		assert.Equal(t, justRight(s, width, pad, sty), want)
	})

}

func TestRenderTitle(t *testing.T) {
	table := NewTable(2)
	table.AddRow(s(10), s(10))
	table.pad = 0
	table.maxWidth = 30
	table.Title("Test Title")
	table.computeColWidths()
	want := fmt.Sprintf("     %s     ", Styled(Bold).ApplyTo("Test Title"))
	t.Run("Title should be bold and centered", func(t *testing.T) {
		assert.Equal(t, renderTitle(table), want)
	})
}

func TestRenderCell(t *testing.T) {
	table := NewTable(1)
	table.AddRow(s(10))
	table.AddRow(s(14))
	table.maxWidth = 30
	table.pad = 2
	table.computeColWidths()
	t.Run("Cell should be rendered with correct justification", func(t *testing.T) {
		want := fmt.Sprintf("  %s      ", Styled(Default).ApplyTo(s(10)))
		st := renderCell(table.rows[0].cells[0].value, table.columns[0].computedWidth, table.pad, table.columns[0].style, table.columns[0].justify)
		assert.Equal(t, st, want)
		table.columns[0].justify = Center
		want = fmt.Sprintf("    %s    ", Styled(Default).ApplyTo(s(10)))
		st = renderCell(table.rows[0].cells[0].value, table.columns[0].computedWidth, table.pad, table.columns[0].style, table.columns[0].justify)
		assert.Equal(t, st, want)
		table.columns[0].justify = Right
		want = fmt.Sprintf("      %s  ", Styled(Default).ApplyTo(s(10)))
		st = renderCell(table.rows[0].cells[0].value, table.columns[0].computedWidth, table.pad, table.columns[0].style, table.columns[0].justify)
		assert.Equal(t, st, want)
	})
}

func TestRenderRow(t *testing.T) {
	table := NewTable(2)
	table.AddRow(s(10), s(10))
	table.AddRow(s(10), s(20))
	table.pad = 2
	table.maxWidth = 28
	table.computeColWidths()
	c10 := Styled(Default).ApplyTo(s(10))
	cEmpty := Styled(Default).ApplyTo("")
	t.Run("Non-wrapped row rendered normally", func(t *testing.T) {

		want := fmt.Sprintf("  %s    %s  \n", c10, c10)
		renderedRow := renderRow(table.rows[0].cells, table.columns, table.pad, table.spacing)
		assert.Equal(t, renderedRow, want)
	})
	t.Run("Wrapped row rendered as multiple lines", func(t *testing.T) {
		want := fmt.Sprintf("  %s    %s  \n  %s              %s  \n", c10, c10, cEmpty, c10)
		renderedRow := renderRow(table.rows[1].cells, table.columns, table.pad, table.spacing)
		assert.Equal(t, renderedRow, want)
	})
}

func TestRenderTable(t *testing.T) {
	table := NewTable(2)
	table.AddRow(s(10), s(10))
	table.AddRow(s(10), s(20))
	table.pad = 2
	table.maxWidth = 28
	table.Title("Test Table")
	c10 := Styled(Default).ApplyTo(s(10))
	cEmpty := Styled(Default).ApplyTo("")
	cTitle := Styled(Bold).ApplyTo("Test Table")
	t.Run("Table with wrapped + non-wrapped rows rendered appropriately", func(t *testing.T) {
		want0 := fmt.Sprintf("         %s         \n\n", cTitle)
		want1 := fmt.Sprintf("  %s    %s  \n", c10, c10)
		want2 := fmt.Sprintf("  %s    %s  \n  %s              %s  \n", c10, c10, cEmpty, c10)
		want := want0 + want1 + want2
		renderedTable := table.AsString()
		assert.Equal(t, renderedTable, want)
	})
}

func TestHeadersShort(t *testing.T) {
	table := NewTable(2).
		ColumnHeaders("test")
	t.Run("Table with only one column header set", func(t *testing.T) {
		assert.Equal(t, table.headers[0].value, "test")
		assert.Equal(t, table.headers[1].value, "")
		snapshot.Assert(t, []byte(table.AsString()))
	})
}
