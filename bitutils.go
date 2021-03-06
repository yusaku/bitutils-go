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

	Msh2 = 0xaaaaaaaaaaaaaaaa
	Msh4 = 0xcccccccccccccccc
	Msh8 = 0xf0f0f0f0f0f0f0f0

	Lsb2 = 0x5555555555555555
	Lsb4 = 0x1111111111111111
	Lsb8 = 0x0101010101010101

	Msb2 = 0xaaaaaaaaaaaaaaaa
	Msb4 = 0x8888888888888888
	Msb8 = 0x8080808080808080

	Idx8 = 0x8040201008040201
)

// Word represents a 64-bit binary string.
type Word uint64

var (
	Idx  [W + 1]Word // Idx[i] has a 1 only at i.
	IdxC [W + 1]Word // IdxC[i] has a 0 only at i.
	LsbN [W + 1]Word // LsbN[i] has 1s in its i LSBs.
	MsbN [W + 1]Word // MsbN[i] has 1s in its i MSBs.
)

func init() {
	for i := 0; i < len(Idx); i++ {
		Idx[i] = Word(1) << uint(i)
		IdxC[i] = ^Idx[i]
		LsbN[i] = Idx[i] - 1
		MsbN[i] = LsbN[i] << uint(W-i)
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
	return (^w).Count1()
}

// Count returns the number of b[0]'s contained in w.
func (w Word) Count(b int) int {
	return ((LsbN[W] + Word(b)) ^ w).Count1()
}

// Get returns w[i].
func (w Word) Get(i int) Word {
	return (w >> uint(i)) & Idx[0]
}

// Set1 sets w[i] to 1.
func (w Word) Set1(i int) Word {
	return w | Idx[i]
}

// Set0 sets w[i] to 0.
func (w Word) Set0(i int) Word {
	return w & IdxC[i]
}

// Flip flips w[i].
func (w Word) Flip(i int) Word {
	return w ^ Idx[i]
}

// Least1 returns a word that indicates the least 1 in w.
func (w Word) Least1() Word {
	if w == 0 {
		return 0
	}
	return ((w - 1) ^ w) & w
}

// LeastIndex1 returns the index of the least 1 in w if exists and -1
// otherwise.
func (w Word) LeastIndex1() int {
	if w == 0 {
		return -1
	}
	return ((w - 1) ^ w).Count1() - 1
}

// Rank1 returns the number of ones in w[0]...w[i].
func (w Word) Rank1(i int) int {
	return (LsbN[i+1] & w).Count1()
}

// Rank0 returns the number of zeros in w[0]...w[i].
func (w Word) Rank0(i int) int {
	return (MsbN[W-i-1] | w).Count0()
}

func (w Word) zcmp8() Word {
	return ((w | ((w | Msb8) - Lsb8)) & Msb8) >> 7
}

func (w Word) leq8(v Word) Word {
	return (((((v | Msb8) - (w & ^Word(Lsb8))) ^ w) ^ v) & Msb8) >> 7
}

// Select1 returns the ith 1 in w.
func (w Word) Select1(i int) int {
	s := w - ((w & Msb2) >> 1)
	s = (s & Lsh4) + ((s >> 2) & Lsh4)
	s = ((s + (s >> 4)) & Lsh8) * Lsb8
	b := ((s.leq8(Word(i)*Lsb8) * Lsb8) >> 53) & ^Word(0x0111)
	l := Word(i) - (((s << 8) >> b) & 0xff)
	s = ((((w >> b) & 0xff) * Lsb8) & Idx8).zcmp8() * Lsb8
	if w = b + ((s.leq8(l*Lsb8) * Lsb8) >> 56); w != 0x48 {
		return int(w)
	} else {
		return -1
	}
}

// Select0 returns the ith 0 in w.
func (w Word) Select0(i int) int {
	return (^w).Select1(i)
}
