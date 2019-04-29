package goradix

import (
	"log"
	"math/rand"
	"testing"
)

func TestCharlieStore(t *testing.T) {
	dst := NewCharlie(8)

	dst.Store(8)
	dst.Store(0xFFFE)

	if got, want := dst.Load(8), true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := dst.Load(9), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := dst.Load(0xFFFE), true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := dst.Load(0xFFFF), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestCharlieDelete(t *testing.T) {
	dst := NewCharlie(4)
	dst.Store(8)
	dst.Store(9)

	prevCount := dst.Count

	dst.Delete(8)

	if got, want := dst.Load(8), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := dst.Count, prevCount-1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	prevCount = dst.Count

	dst.Delete(9)

	if got, want := dst.Load(9), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := dst.Count, prevCount-1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestCharlie2Count(t *testing.T) {
	// t.Skip("no need to run memory calculations each time")
	tree := NewCharlie(2)

	const limit = 1000000 // 1M
	var thresh = 1
	for i := 0; i < limit; i++ {
		if i%thresh == 0 {
			log.Printf("n: %12d; count: %20d; bytes: %12d", i, tree.Count, tree.Count*(24+tree.childCount*8))
			thresh *= 10
		}
		tree.Store(rand.Uint64())
	}
	log.Printf("n: %12d; count: %20d; bytes: %12d", limit, tree.Count, tree.Count*(24+tree.childCount*8))
}

func TestCharlie4Count(t *testing.T) {
	// t.Skip("no need to run memory calculations each time")
	tree := NewCharlie(4)

	const limit = 1000000 // 1M
	var thresh = 1
	for i := 0; i < limit; i++ {
		if i%thresh == 0 {
			log.Printf("n: %12d; count: %20d; bytes: %12d", i, tree.Count, tree.Count*(24+tree.childCount*8))
			thresh *= 10
		}
		tree.Store(rand.Uint64())
	}
	log.Printf("n: %12d; count: %20d; bytes: %12d", limit, tree.Count, tree.Count*(24+tree.childCount*8))
}

func TestCharlie8Count(t *testing.T) {
	// t.Skip("no need to run memory calculations each time")
	tree := NewCharlie(8)

	const limit = 1000000 // 1M
	var thresh = 1
	for i := 0; i < limit; i++ {
		if i%thresh == 0 {
			log.Printf("n: %12d; count: %20d; bytes: %12d", i, tree.Count, tree.Count*(24+tree.childCount*8))
			thresh *= 10
		}
		tree.Store(rand.Uint64())
	}
	log.Printf("n: %12d; count: %20d; bytes: %12d", limit, tree.Count, tree.Count*(24+tree.childCount*8))
}
