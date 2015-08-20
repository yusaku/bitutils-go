// Package bitutils provides a collection of utilities to deal with bits.
package bitutils

import (
	"fmt"
	"strconv"
)

// W is the length of a machine word.
const W = 64

// Magic constants.
const (
	Lsh2 = 0x5555555555555555
	Lsh4 = 0x3333333333333333
	Lsh8 = 0x0f0f0f0f0f0f0f0f
	Lsb8 = 0x0101010101010101
)

// Word represents a 64-bit binary string.
type Word uint64

var (
	Pos  [W + 1]Word // Pos[i] has a 1 only at i.
	PosC [W + 1]Word // PosC[i] has a 0 only at i.
	Lsh  [W + 1]Word // Lsh[i] has 1s in its i LSBs.
	Msh  [W + 1]Word // Msh[i] has 1s in its i MSBs.
)

func init() {
	for i := 0; i < len(Pos); i++ {
		Pos[i] = Word(1) << uint(i)
		PosC[i] = ^Pos[i]
		Lsh[i] = Pos[i] - 1
		Msh[i] = Lsh[i] << uint(W-i)
	}
}

// ParseWord returns a Word from a string.
func ParseWord(s string) (Word, error) {
	w, err := strconv.ParseUint(s, 2, 64)
	return Word(w), err
}

// String returns binary string w[0]w[1]...w[63].
func (w Word) String() string {
	return fmt.Sprintf("%064b", w)
}

// Count1 returns the number of ones contained in w.
func (w Word) Count1() int {
	w -= (w >> 1) & Lsh2
	w = (w & Lsh4) + ((w >> 2) & Lsh4)
	w = (w + (w >> 4)) & Lsh8
	return int((w * Lsb8) >> 56)
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
	return w & Pos[0]
}

// Set1 sets w[i] to 1.
func (w Word) Set1(i int) Word {
	return w | Pos[i]
}

// Set0 sets w[i] to 0.
func (w Word) Set0(i int) Word {
	return w & PosC[i]
}

// Flip flips w[i].
func (w Word) Flip(i int) Word {
	return w ^ Pos[i]
}

// Least1 returns a word that indicates the least 1 in w.
func (w Word) Least1() Word {
	if w == 0 {
		return 0
	}
	w = ((w - 1) ^ w) & w
	return w
}

// LeastIndex1 returns the index of the least 1 in w if exists and -1
// otherwise.
func (w Word) LeastIndex1() int {
	if w == 0 {
		return -1
	}
	w = (w - 1) ^ w
	return w.Count1() - 1
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
