// Package broadword provides a collection of broadword algorithms exploiting
// word-level parallelism.
package broadword

import (
	"fmt"
)

// W is the length of a machine word.
const W = 64

// Magic constants.
const (
	lowers2 = 0x5555555555555555
	lowers4 = 0x3333333333333333
	lowers8 = 0x0f0f0f0f0f0f0f0f
	lowest8 = 0x0101010101010101
)

// Word represents a 64-bit binary string.
type Word uint64

// exp2[i] is 2^i.
var exp2 [W]Word

func init() {
	for i := uint(0); i < W; i++ {
		exp2[i] = Word(1) << i
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
	w = w >> uint(i)
	return w & exp2[0]
}

// Set1 sets w[i] to 1.
func (w Word) Set1(i int) Word {
	return w | exp2[i]
}

// Set0 sets w[i] to 0.
func (w Word) Set0(i int) Word {
	return w & ^exp2[i]
}

// Flip flips w[i].
func (w Word) Flip(i int) Word {
	return w ^ exp2[i]
}

// Rank1 returns the number of ones in w[0]...w[i].
func (w Word) Rank1(i int) int {
	w = (w << uint(W-i-1))
	return w.Count1()
}
