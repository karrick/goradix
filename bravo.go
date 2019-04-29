package goradix

// NOTE: This is not really a DST, but more of an existence trie.

type bnode struct {
	left, right *bnode
}

type Bravo struct {
	root  *bnode
	Count uint64
}

func NewBravo() *Bravo {
	return &Bravo{root: new(bnode)}
}

const initialMask = uint64(1 << 63)

// search returns the node prior to a nil pointer, followed by the mask so
// upstream can determine whether node was located.
func (dst *Bravo) search(key uint64) (*bnode, uint64) {
	var node *bnode
	next := dst.root
	mask := initialMask

	for ; mask != 0; mask >>= 1 {
		node = next
		if key&mask != 0 {
			next = node.right
		} else {
			next = node.left
		}
		if next == nil {
			return node, mask
		}
	}

	return node, mask
}

// TODO: When DST has heavy additions and removals, cleaning tree for every
// Delete results in needless memory churn. Instead, provide a Compact method to
// clean up dead branches independent of Delete method.

// Delete removes the specified 64-bit key from the DST.
func (dst *Bravo) Delete(key uint64) {
	node, mask := dst.search(key)
	if mask != 0 {
		return // key not present
	}
	// Check the final bit to determine whether to remove right or left branch
	// from node.
	if key&1 != 0 {
		node.right = nil // remove right branch
	} else {
		node.left = nil // remove left branch
	}
	dst.Count--
}

// Load returns whether or not the specified 64-bit key is present in the DST.
func (dst *Bravo) Load(key uint64) bool {
	_, mask := dst.search(key)
	return mask == 0
}

// Store stores the existence of the specified 64-bit key.
func (dst *Bravo) Store(key uint64) {
	// walk existing tree branches as much as possible
	node, mask := dst.search(key)

	// create whatever branches needed
	for ; mask != 0; mask >>= 1 {
		newNode := new(bnode)
		dst.Count++
		if key&mask != 0 {
			node.right = newNode
		} else {
			node.left = newNode
		}
		node = newNode
	}
}
