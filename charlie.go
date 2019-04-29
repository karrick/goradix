package goradix

// R-ary existence data structure

type cnode struct {
	children []*cnode
}

type Charlie struct {
	head             *cnode
	mask             uint64
	childCount       uint64
	Count            uint64
	bitInit, bitStep uint8 // these values computed once at init and used in most methods
}

// var isPower2 = function(x) { return (x > 0 && !(x & (x-1))); };

func isPower2(x uint8) bool {
	if x == 0 {
		return false
	}
	y := x - 1
	return x&y != 0
}

// roundToPowerOfTwo returns the argument, possibly rounded down to the nearest
// whole number power of 2.
func roundToPowerOfTwo(a uint8) uint8 {
	switch {
	case a >= 32:
		return 32
	case a >= 16:
		return 16
	case a >= 8:
		return 8
	case a >= 4:
		return 4
	case a >= 2:
		return 2
	default:
		return 1
	}
}

// NewCharlie returns a new tree to hold existence information for 64-bit
// keys. The specified bits must be a power of two greater than 0, such as 1, 2,
// 4, 8, etc.
func NewCharlie(bits uint8) *Charlie {
	bits = roundToPowerOfTwo(bits)
	childCount := uint64(1 << bits)
	return &Charlie{
		head:       &cnode{children: make([]*cnode, childCount)},
		mask:       childCount - 1,
		childCount: childCount,
		bitInit:    64 - bits,
		bitStep:    bits,
	}
}

// find returns the node prior to a nil pointer, followed by one less from the
// number of remaining bits to be shifted so upstream can determine whether key
// was located. Bits will be 255 when the specified key was found. If the key
// mismatched on the final bit, bits will be 0.
func (tree *Charlie) find(key uint64) (*cnode, uint8) {
	var prev *cnode
	curr := tree.head
	bits := tree.bitInit
	mask := tree.mask    // store in local variable so optimizer can see it never changes
	step := tree.bitStep // store in local variable so optimizer can see it never changes

	// Need to execute loop from start to 0 inclusive; therefore, terminate when
	// rolls over down from 0 back up to 255.
	for ; bits < 64; bits -= step {
		prev = curr
		curr = curr.children[(key>>bits)&mask]
		if curr == nil {
			return prev, bits
		}
	}

	return prev, bits
}

// TODO: When heavy additions and removals, cleaning tree for every Delete
// results in needless memory churn. Instead, provide a Compact method to clean
// up dead branches independent of Delete method.

// Delete removes the specified 64-bit key.
func (tree *Charlie) Delete(key uint64) {
	node, bits := tree.find(key)
	if bits < 64 {
		return // key not present
	}
	// Check the final bit to determine whether to remove right or left branch
	// from node.
	node.children[key&tree.mask] = nil // remove branch
	tree.Count--
}

// Load returns whether or not the specified 64-bit key is present.
func (tree *Charlie) Load(key uint64) bool {
	_, bits := tree.find(key)
	return bits >= 64 // key present when bits >= 64
}

// Store stores the existence of the specified 64-bit key.
func (tree *Charlie) Store(key uint64) {
	// walk existing tree branches as much as possible
	node, bits := tree.find(key)
	if bits >= 64 {
		return // key already present
	}

	// create needed branches
	tree.Count += uint64(bits + 1) // below loop adds this many nodes to tree

	childCount := tree.childCount // store in local variable so optimizer can see it never changes
	mask := tree.mask             // store in local variable so optimizer can see it never changes
	step := tree.bitStep          // store in local variable so optimizer can see it never changes

	// Need to execute loop from start to 0 inclusive; therefore, terminate
	// when rolls over down from 0 back up to 255.
	for ; bits < 64; bits -= step {
		newNode := &cnode{children: make([]*cnode, childCount)}
		node.children[(key>>bits)&mask] = newNode
		node = newNode
	}
}
