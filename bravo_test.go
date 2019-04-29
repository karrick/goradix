package goradix

import (
	"fmt"
	"io"
	"testing"
)

func (dst *Bravo) Display(w io.Writer) {
	if dst.root != nil {
		dst.root.display(w, 0)
	}
}

func (n *bnode) display(w io.Writer, level int) {
	var indentation string
	for i := 0; i < level; i++ {
		indentation += " "
	}
	fmt.Fprintf(w, "%snode: %p %+v\n", indentation, n, n)
	if n.left != nil {
		n.left.display(w, level+1)
	}
	if n.right != nil {
		n.right.display(w, level+1)
	}
}

func TestBravoStore(t *testing.T) {
	dst := NewBravo()

	dst.Store(8)
	dst.Store(0xFFFE)

	// dst.Display(os.Stderr)

	if got, want := dst.Load(9), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := dst.Load(8), true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := dst.Load(0xFFFF), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := dst.Load(0xFFFE), true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestBravoDelete(t *testing.T) {
	dst := NewBravo()

	dst.Store(8)
	if got, want := dst.Load(8), true; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	dst.Delete(8)

	if got, want := dst.Load(8), false; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}
