// Package bitutils provides a collection of utilities to deal with bits.
package bitutils

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

var (
	// shift[i] is (1 << i).
	shift [W]Word
	// shiftNot[i] is ^shift[i].
	shiftNot [W]Word
)

func init() {
	for i := 0; i < len(shift); i++ {
		shift[i] = Word(1) << uint(i)
		shiftNot[i] = ^shift[i]
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
	w = ^w
	return w.Count1()
}

// Count returns the number of b[0]'s contained in w.
func (w Word) Count(b int) int {
	w = w ^ (^Word(0) + Word(b))
	return w.Count1()
}

// Get returns w[i].
func (w Word) Get(i int) Word {
	w = w >> uint(i)
	return w & shift[0]
}

// Set1 sets w[i] to 1.
func (w Word) Set1(i int) Word {
	return w | shift[i]
}

// Set0 sets w[i] to 0.
func (w Word) Set0(i int) Word {
	return w & shiftNot[i]
}

// Flip flips w[i].
func (w Word) Flip(i int) Word {
	return w ^ shift[i]
}

// Lsb returns the index of the first 1 in w if w != 0 and -1 otherwise.
func (w Word) Lsb() int {
	if w == 0 {
		return -1
	} else {
		w = (w - 1) ^ w
		return w.Count1() - 1
	}
}

// Rank1 returns the number of ones in w[0]...w[i].
func (w Word) Rank1(i int) int {
	w = w << uint(W-i-1)
	return w.Count1()
}

// Rank0 returns the number of zeros in w[0]...w[i].
func (w Word) Rank0(i int) int {
	w = ^w << uint(W-i-1)
	return w.Count1()
}
