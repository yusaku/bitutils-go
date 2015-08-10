package bitutils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var r *rand.Rand
var w Word
var ws string

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < len(runes)/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func testCase() (w Word, ws string) {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	w = Word(r.Uint32()) | Word(r.Uint32())<<32
	ws = fmt.Sprintf("%064b", w)
	return
}

func TestMain(m *testing.M) {
	w, ws = testCase()
	os.Exit(m.Run())
}

func TestCount1(t *testing.T) {
	got := w.Count1()
	want := strings.Count(ws, "1")
	if got != want {
		t.Errorf("got %d, want %d for %s", got, want, ws)
	}
}

func TestCount0(t *testing.T) {
	want := w.Count0()
	got := strings.Count(ws, "0")
	if got != want {
		t.Errorf("got %d, want %d for %s", got, want, ws)
	}
}

func TestCount(t *testing.T) {
	for i := 0; i < 2; i++ {
		want := w.Count(i)
		got := strings.Count(ws, strconv.Itoa(i))
		if got != want {
			t.Errorf("got %d, want %d for %s", got, want, ws)
		}
	}
}

func TestGet(t *testing.T) {
	for i := 0; i < W; i++ {
		j := W - i - 1 // corresponding index in ws.

		n, err := strconv.ParseUint(ws[j:j+1], 10, 0)
		if err != nil {
			t.Errorf("ParseUint failed")
		}

		got := w.Get(i)
		want := Word(n)
		if got != want {
			t.Errorf("got %d, want %d for %s", got, want, ws)
		}
	}
}

func TestSet1(t *testing.T) {
	for i := 0; i < W; i++ {
		v := w.Set1(i)
		for j := 0; j < W; j++ {
			got := v.Get(j)
			want := w.Get(j)
			if i == j {
				want = 1
			}
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		}
	}
}

func TestSet0(t *testing.T) {
	for i := 0; i < W; i++ {
		v := w.Set0(i)
		for j := 0; j < W; j++ {
			got := v.Get(j)
			want := w.Get(j)
			if i == j {
				want = 0
			}
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		}
	}
}

func TestFlip(t *testing.T) {
	for i := 0; i < W; i++ {
		v := w.Flip(i)
		for j := 0; j < W; j++ {
			got := v.Get(j)
			want := w.Get(j)
			if i == j {
				if want == 1 {
					want = 0
				} else {
					want = 1
				}
			}
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		}
	}
}

func TestLsb(t *testing.T) {
	got := w.Lsb()
	want := strings.Index(reverse(ws), "1")
	if got != want {
		t.Errorf("got %d, want %d for %s", got, want, ws)
	}
}

func TestRank1(t *testing.T) {
	for i := 0; i < W; i++ {
		got := w.Rank1(i)
		want := strings.Count(ws[len(ws)-i-1:len(ws)], "1")
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestRank0(t *testing.T) {
	for i := 0; i < W; i++ {
		got := w.Rank0(i)
		want := strings.Count(ws[len(ws)-i-1:len(ws)], "0")
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func BenchmarkCount1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Count1()
	}
}

func BenchmarkCount0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Count0()
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Count(1)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Get(i % W)
	}
}

func BenchmarkSet1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Set1(i % W)
	}
}

func BenchmarkSet0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Set0(i % W)
	}
}

func BenchmarkFlip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Flip(i % W)
	}
}


func BenchmarkLsb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Lsb()
	}
}

func BenchmarkRank1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Rank1(i % W)
	}
}

func BenchmarkRank0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Rank0(i % W)
	}
}
