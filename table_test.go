package clt

import (
	"fmt"
	c "github.com/smartystreets/goconvey/convey"
	"reflect"
	"strings"
	"testing"
)

func TestTerminalSizeCheck(t *testing.T) {
	t.Skipf("Terminal check skipped. No TTY.")
	h, w, err := getTerminalSize()
	if err != nil || h == -1 || w == -1 {
		fmt.Printf("Cannot determine terminal size for Table.  McClintock will still work, but will not be able to automagically determine sizes.")
	}
}

func TestCreateTable(t *testing.T) {
	table := NewTable(3)
	if len(table.columns) != 3 {
		t.Errorf("Table should have %d columns, has %d.", 3, table.columns)
	}
	for _, col := range table.columns {
		if !reflect.DeepEqual(col.justify, jLeft) {
			t.Errorf("Table default justification should be %v, got %v.", jLeft, col.justify)
		}
	}
	if table.MaxHeight == 0 || table.MaxWidth == 0 {
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
	if table.title.style != Style(Bold) {
		fmt.Errorf("Default title style should be chalk.Bold.")
	}
}

func TestSetColumnHeaders(t *testing.T) {
	table := NewTable(2)
	table.SetColumnHeaders("Header1", "Header2")
	want := []cell{cell{
		value: "Header1",
		style: Style(Bold, Underline),
		width: len("Header1"),
	},
		cell{
			value: "Header2",
			style: Style(Bold, Underline),
			width: len("Header2"),
		},
	}
	for i, header := range table.headers {
		if header.value != want[i].value {
			fmt.Errorf("Header should be %s, got %s", want[i].value, header.value)
		}
	}
}

func TestRenderHelpers(t *testing.T) {
	n := []int{1, 3, 2}
	table := &table{}
	table.columns = []col{col{
		naturalWidth:  10,
		computedWidth: 12,
	},
		col{
			naturalWidth:  12,
			computedWidth: 14,
		},
	}
	table.pad = 2

	c.Convey("TestRenderHelpers", t, func() {
		c.So(mapAdd(n, 1), c.ShouldResemble, []int{2, 4, 3})
		c.So(sum(n), c.ShouldEqual, 6)
		i, m := max(n)
		c.So(m, c.ShouldEqual, 3)
		c.So(i, c.ShouldEqual, 1)
		c.So(sumWithoutIndex(n, 1), c.ShouldEqual, 3)
		c.So(wrappedWidthOk(51, 100), c.ShouldBeTrue)
		c.So(wrappedWidthOk(49, 100), c.ShouldBeFalse)
		c.So(extractComputedWidth(table), c.ShouldResemble, []int{12, 14})
		c.So(extractNatWidth(table), c.ShouldResemble, []int{10, 12})
		c.So(table.width(), c.ShouldEqual, 12+14+8)
	})

}

// test helper to get string of length n
func s(n int) string {
	return strings.Repeat("x", n)
}

func TestSimpleStrategy(t *testing.T) {
	table := NewTable(3)
	table.MaxWidth = 80
	table.SetColumnHeaders(s(4), s(4), s(4))
	table.AddRow(s(10), s(12), s(14))
	table.pad = 2
	table.computeColWidths()
	c.Convey("NaturalColWidths < MaxWidth", t, func() {
		c.So(extractNatWidth(table), c.ShouldResemble, []int{10, 12, 14})
		c.So(extractComputedWidth(table), c.ShouldResemble, []int{10, 12, 14})
	})

	// headers bigger than content
	table.SetColumnHeaders(s(15), s(16), s(17))
	table.computeColWidths()
	c.Convey("Big headers, NaturalWidth < MaxWidth", t, func() {
		c.So(extractNatWidth(table), c.ShouldResemble, []int{15, 16, 17})
		c.So(extractComputedWidth(table), c.ShouldResemble, []int{15, 16, 17})
	})
}

func TestWrapWidest(t *testing.T) {
	table := NewTable(3)
	table.MaxWidth = 60
	table.AddRow(s(10), s(20), s(40))

	c.Convey("Last column wrap, no padding", t, func() {
		table.pad = 0
		table.computeColWidths()
		c.So(extractNatWidth(table), c.ShouldResemble, []int{10, 20, 40})
		c.So(extractComputedWidth(table), c.ShouldResemble, []int{10, 20, 30})
	})

}

func TestOverflow(t *testing.T) {
	table := NewTable(3)
	table.MaxWidth = 10
	table.AddRow(s(10), s(20), s(40))

	c.Convey("Overflow to natural width as last resort", t, func() {
		table.pad = 0
		table.computeColWidths()
		c.So(extractNatWidth(table), c.ShouldResemble, []int{10, 20, 40})
		c.So(extractComputedWidth(table), c.ShouldResemble, []int{10, 20, 40})
	})

}

func TestJustifcation(t *testing.T) {
	s := s(4)
	c.Convey("Center justify text with padding", t, func() {
		width := 14
		pad := 2
		sty := Style(Default)
		want := fmt.Sprintf("       %s       ", sty.ApplyTo(s))
		c.So(justCenter(s, width, pad, sty), c.ShouldEqual, want)
	})
	c.Convey("Center justify offest left on uneven", t, func() {
		width := 13
		pad := 2
		sty := Style(Default)
		want := fmt.Sprintf("      %s       ", sty.ApplyTo(s))
		c.So(justCenter(s, width, pad, sty), c.ShouldEqual, want)
	})
	c.Convey("Left justify text with padding", t, func() {
		width := 8
		pad := 2
		sty := Style(Default)
		want := fmt.Sprintf("  %s      ", sty.ApplyTo(s))
		c.So(justLeft(s, width, pad, sty), c.ShouldEqual, want)
	})
	c.Convey("Right justify text with padding", t, func() {
		width := 8
		pad := 2
		sty := Style(Default)
		want := fmt.Sprintf("      %s  ", sty.ApplyTo(s))
		c.So(justRight(s, width, pad, sty), c.ShouldEqual, want)
	})
	c.Convey("Fallback to string + padding if widths jacked up", t, func() {
		width := 1
		pad := 2
		sty := Style(Default)
		want := fmt.Sprintf("  %s  ", sty.ApplyTo(s))
		c.So(justCenter(s, width, pad, sty), c.ShouldEqual, want)
		c.So(justLeft(s, width, pad, sty), c.ShouldEqual, want)
		c.So(justRight(s, width, pad, sty), c.ShouldEqual, want)
	})

}

func TestRenderTitle(t *testing.T) {
	table := NewTable(2)
	table.AddRow(s(10), s(10))
	table.pad = 0
	table.MaxWidth = 30
	table.SetTitle("Test Title")
	table.computeColWidths()
	want := fmt.Sprintf("     %s     ", Style(Bold).ApplyTo("Test Title"))
	c.Convey("Title should be bold and centered", t, func() {
		c.So(renderTitle(table), c.ShouldEqual, want)
	})
}

func TestRenderCell(t *testing.T) {
	table := NewTable(1)
	table.AddRow(s(10))
	table.AddRow(s(14))
	table.MaxWidth = 30
	table.pad = 2
	table.computeColWidths()
	c.Convey("Cell should be rendered with correct justification", t, func() {
		want := fmt.Sprintf("  %s      ", Style(Default).ApplyTo(s(10)))
		st := renderCell(table.rows[0].cells[0].value, table.columns[0], table.pad)
		c.So(st, c.ShouldResemble, want)
		table.columns[0].justify = jCenter
		want = fmt.Sprintf("    %s    ", Style(Default).ApplyTo(s(10)))
		st = renderCell(table.rows[0].cells[0].value, table.columns[0], table.pad)
		c.So(st, c.ShouldResemble, want)
		table.columns[0].justify = jRight
		want = fmt.Sprintf("      %s  ", Style(Default).ApplyTo(s(10)))
		st = renderCell(table.rows[0].cells[0].value, table.columns[0], table.pad)
		c.So(st, c.ShouldResemble, want)
	})
}

func TestRenderRow(t *testing.T) {
	table := NewTable(2)
	table.AddRow(s(10), s(10))
	table.AddRow(s(10), s(20))
	table.pad = 2
	table.MaxWidth = 28
	table.computeColWidths()
	c10 := Style(Default).ApplyTo(s(10))
	cEmpty := Style(Default).ApplyTo("")
	c.Convey("Non-wrapped row rendered normally", t, func() {

		want := fmt.Sprintf("  %s    %s  \n", c10, c10)
		renderedRow := renderRow(table.rows[0].cells, table.columns, table.pad)
		c.So(renderedRow, c.ShouldResemble, want)
	})
	c.Convey("Wrapped row rendered as multiple lines", t, func() {

		want := fmt.Sprintf("  %s    %s  \n  %s              %s  \n", c10, c10, cEmpty, c10)
		renderedRow := renderRow(table.rows[1].cells, table.columns, table.pad)
		c.So(renderedRow, c.ShouldResemble, want)
	})
}

func TestRenderTable(t *testing.T) {
	table := NewTable(2)
	table.AddRow(s(10), s(10))
	table.AddRow(s(10), s(20))
	table.pad = 2
	table.MaxWidth = 28
	table.SetTitle("Test Table")
	c10 := Style(Default).ApplyTo(s(10))
	cEmpty := Style(Default).ApplyTo("")
	cTitle := Style(Bold).ApplyTo("Test Table")
	c.Convey("Table with wrapped + non-wrapped rows rendered appropriately", t, func() {
		want0 := fmt.Sprintf("         %s         \n\n", cTitle)
		want1 := fmt.Sprintf("  %s    %s  \n", c10, c10)
		want2 := fmt.Sprintf("  %s    %s  \n  %s              %s  \n", c10, c10, cEmpty, c10)
		want := want0 + want1 + want2
		renderedTable := table.renderTableAsString()
		c.So(renderedTable, c.ShouldResemble, want)
	})
}
