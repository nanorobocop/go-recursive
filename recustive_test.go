package main

import "testing"

type A struct {
	a  string
	A  int
	AA B
}

type B struct {
	B int
}

func TestRecursiveStruct(t *testing.T) {
	a := A{"asdf", 11, B{}}

	opts := &WalkOpts{WalkFunc: PrintWalkFunc}
	Walk(a, opts)
}

func TestRecursiveInt(t *testing.T) {
	opts := &WalkOpts{WalkFunc: PrintWalkFunc}
	Walk(11, opts)
}

func TestRecursiveMap(t *testing.T) {
	opts := &WalkOpts{WalkFunc: PrintWalkFunc}
	Walk(map[string]int{}, opts)
}

func TestRecursiveSlice(t *testing.T) {
	opts := &WalkOpts{WalkFunc: PrintWalkFunc}
	Walk([]int{}, opts)
}
