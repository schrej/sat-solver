package main

import "fmt"

type Formula interface {
	Val([]bool) bool
}

type S int

func (f S) Val(i []bool) bool {
	return i[f]
}

type Not struct {
	A Formula
}

func (f Not) Val(i []bool) bool {
	return !f.A.Val(i)
}

type And struct {
	A Formula
	B Formula
}

func (f And) Val(i []bool) bool {
	return f.A.Val(i) && f.B.Val(i)
}

type Or struct {
	A Formula
	B Formula
}

func (f Or) Val(i []bool) bool {
	return f.A.Val(i) || f.B.Val(i)
}

type Imp struct {
	A Formula
	B Formula
}

func (f Imp) Val(i []bool) bool {
	return !f.A.Val(i) || f.B.Val(i)
}

type BiImp struct {
	A Formula
	B Formula
}

func (f BiImp) Val(i []bool) bool {
	return f.A.Val(i) == f.A.Val(i)
}

type Interpretation struct {
	Vals  []bool
	Fixed int
}

func (i Interpretation) Next() bool {
	for j := i.Fixed; j < len(i.Vals); j++ {
		i.Vals[j] = !i.Vals[j]
		if i.Vals[j] {
			return true
		}
	}
	return false
}

func (i Interpretation) Print(labels []string) {
	for idx, v := range i.Vals {
		fmt.Printf("%s: %v, ", labels[idx], v)
	}
	fmt.Println()
}
