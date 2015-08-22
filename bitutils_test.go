package bitutils

import (
	"math/rand"
	"os"
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
	w = Word(r.Uint32()) | (Word(r.Uint32()) << 32)
	ws = reverse(w.String())
	return
}

func TestMain(m *testing.M) {
	w, ws = testCase()
	os.Exit(m.Run())
}

func TestParseWord(t *testing.T) {
	got, err := ParseWord(w.String())
	if got == w && err != nil {
		t.Errorf("got %d, want %d", got, w)
	}
}

func TestCount1(t *testing.T) {
	var want int
	for i := 0; i < len(ws); i++ {
		if ws[i] == '1' {
			want += 1
		}
	}

	if got := w.Count1(); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestCount0(t *testing.T) {
	var want int
	for i := 0; i < len(ws); i++ {
		if ws[i] == '0' {
			want += 1
		}
	}

	if got := w.Count0(); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestCount(t *testing.T) {
	for b := 0; b < 2; b++ {
		var want int
		for i := 0; i < len(ws); i++ {
			if b == 1 {
				if ws[i] == '1' {
					want += 1
				}
			} else {
				if ws[i] == '0' {
					want += 1
				}
			}
		}

		if got := w.Count(b); got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestGet(t *testing.T) {
	var wants [W]Word
	for i := 0; i < len(ws); i++ {
		if ws[i] == '1' {
			wants[i] = 1
		} else {
			wants[i] = 0
		}
	}

	for i := 0; i < W; i++ {
		if got, want := w.Get(i), wants[i]; got != want {
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
	var want Word
	for i := 0; i < len(ws); i++ {
		if ws[i] == '1' {
			want = want.Set1(i)
			break
		}
	}

	if got := w.Least1(); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestLeastIndex1(t *testing.T) {
	var want int
	for i := 0; i < len(ws); i++ {
		if ws[i] == '1' {
			want = i
			break
		}
	}

	if got := w.LeastIndex1(); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestRank1(t *testing.T) {
	var wants [W]int
	var n1 int
	for i := 0; i < len(ws); i++ {
		if ws[i] == '1' {
			n1 += 1
		}
		wants[i] = n1
	}

	for i := 0; i < W; i++ {
		if got, want := w.Rank1(i), wants[i]; got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestRank0(t *testing.T) {
	var wants [W]int
	var n0 int
	for i := 0; i < len(ws); i++ {
		if ws[i] == '0' {
			n0 += 1
		}
		wants[i] = n0
	}

	for i := 0; i < W; i++ {
		if got, want := w.Rank0(i), wants[i]; got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}

func TestSelect1(t *testing.T) {
	var wants [W]int
	var n1 int
	for i := 0; i < len(ws); i++ {
		if ws[i] == 1 {
			wants[n1] = i
			n1 += 1
		}
	}

	for i := 0; i < n1; i++ {
		got, want := w.Select1(i), wants[i]
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

func BenchmarkSelect1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := Word(i)
		_ = w.Select1(i % W)
	}
}
