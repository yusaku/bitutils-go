package bitutils

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var r *rand.Rand
var w Word
var ws, wsR string

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < len(runes)/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func testCase() (w Word, ws string) {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	w = Word(r.Uint32()) | (Word(r.Uint32()) << 32)
	ws = w.String()
	wsR = reverse(ws)
	return
}

func TestMain(m *testing.M) {
	w, ws = testCase()
	os.Exit(m.Run())
}

func TestCount1(t *testing.T) {
	got, want := w.Count1(), strings.Count(ws, "1")
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestCount0(t *testing.T) {
	got, want := w.Count0(), strings.Count(ws, "0")
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestCount(t *testing.T) {
	for i := 0; i < 2; i++ {
		got, want := w.Count(i), strings.Count(ws, strconv.Itoa(i))
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestGet(t *testing.T) {
	for i := 0; i < W; i++ {
		n, err := strconv.ParseUint(wsR[i:i+1], 2, 0)
		if err != nil {
			t.Errorf("ParseUint failed")
		}
		got, want := w.Get(i), Word(n)
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestSet1(t *testing.T) {
	for i := 0; i < W; i++ {
		v := w.Set1(i)
		for j := 0; j < W; j++ {
			got, want := v.Get(j), w.Get(j)
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
			got, want := v.Get(j), w.Get(j)
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
			got, want := v.Get(j), w.Get(j)
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

func TestLeast1(t *testing.T) {
	v := w.Least1()
	j := strings.Index(wsR, "1")
	if got := v.Get(j); got != 1 {
		t.Errorf("got %d, want %d", got, 1)
	}
	v = v.Flip(j)
	if v != 0 {
		t.Errorf("got %d, want %d", v, 0)
	}
}

func TestLeastIndex1(t *testing.T) {
	got, want := w.LeastIndex1(), strings.Index(wsR, "1")
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestRank1(t *testing.T) {
	for i := 0; i < W; i++ {
		got, want := w.Rank1(i), strings.Count(wsR[0:i+1], "1")
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestRank0(t *testing.T) {
	for i := 0; i < W; i++ {
		got, want := w.Rank0(i), strings.Count(wsR[0:i+1], "0")
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

func BenchmarkLeast1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Least1()
	}
}

func BenchmarkLeastIndex1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.LeastIndex1()
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
