package clt

import "testing"

func TestStyle1(t *testing.T) {
	s := Style(Red)
	expectBefore := "\x1b[31m"
	expectAfter := "\x1b[39m"
	if s.before != expectBefore || s.after != expectAfter {
		t.Errorf("Expected:\nBefore: %s After: %s\nGot:\nBefore: %s After:%s\n", expectBefore, expectAfter, s.before, s.after)
	}
}

func TestStyle2(t *testing.T) {
	s2 := Style(Red, Underline)
	expectBefore2 := "\x1b[31;4m"
	expectAfter2 := "\x1b[39;24m"
	if s2.before != expectBefore2 || s2.after != expectAfter2 {
		t.Errorf("Expected:\nBefore: %s After: %s\nGot:\nBefore: %s After: %s\n", expectBefore2, expectAfter2, s2.before, s2.after)
	}
}

func TestApplyTo(t *testing.T) {
	s := Style(Red)
	testString := "This is a test"
	expect := "\x1b[31mThis is a test\x1b[39m"
	applyResult := s.ApplyTo(testString)
	if applyResult != expect {
		t.Errorf("Expected: %v\nGot: %v\n", expect, applyResult)
	}
}

func TestApplyTo2(t *testing.T) {
	s := Style(Red, Underline)
	testString := "This is a test"
	expect := "\x1b[31;4mThis is a test\x1b[39;24m"
	applyResult := s.ApplyTo(testString)
	if applyResult != expect {
		t.Errorf("Expected: %v\nGot: %v\n", expect, applyResult)
	}
}
