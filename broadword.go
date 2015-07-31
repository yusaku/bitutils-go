// Package broadword provides a collection of broadword algorithms exploiting
// word-level parallelism.
package broadword

import (
	"fmt"
)

const (
	// W is the length of a machine word.
	W = 64
)

// Magic constants.
const (
	lowers2 = 0x5555555555555555
	lowers4 = 0x3333333333333333
	lowers8 = 0x0f0f0f0f0f0f0f0f
	lowest8 = 0x0101010101010101
)

// Word represents a 64-bit binary string.
type Word uint64

var (
	// convex[i] sets only its ith bit to 1.
	convex [W]Word
	// concave[i] sets only its ith bit to 0.
	concave [W]Word
)

func init() {
	for i := 0; i < W; i++ {
		w := (Word(1) << uint(i))
		convex[i] = w
		concave[i] = ^w
	}
}

// String returns binary string w[0]w[1]...w[63].
func (w Word) String() string {
	return fmt.Sprintf("%064b", w)
}

// Count1 returns the number of ones contained in w.
func (w Word) Count1() int {
	w -= (w >> 1) & lowers2
	w = (w & lowers4) + ((w >> 2) & lowers4)
	w = (w + (w >> 4)) & lowers8
	return int((w * lowest8) >> 56)
}

// Count0 returns the number of zeros contained in w.
func (w Word) Count0() int {
	return W - w.Count1()
}

// Count returns the number of b[0]'s contained in w.
func (w Word) Count(b int) int {
	m := Word(0) - (Word(b^1) & Word(1))
	return (w ^ m).Count1()
}

// Get returns w[i].
func (w Word) Get(i int) Word {
	return (w >> uint(i)) & convex[0]
}

// Set sets w[i] to 1.
func (w Word) Set(i int) Word {
	return w | convex[i]
}

// Unset sets w[i] to 0.
func (w Word) Unset(i int) Word {
	return w & concave[i]
}

// Flip flips w[i].
func (w Word) Flip(i int) Word {
	return w ^ convex[i]
}
