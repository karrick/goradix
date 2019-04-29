package goradix

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
)

// offsetOfMismatch returns the offset of the first bytes that do not match from
// the argument strings.
func offsetOfMismatch(s1, s2 string) int {
	// NOTE: While bytes work for the purposes of splitting keys, when a
	// multi-byte rune is encountered, this algorithm may cause a single rune to
	// be split across two different nodes. This would only be an issue for
	// printing out the nodes of the trie. When searching, deleting, or storing
	// new values, the byte split will not cause any problems as to the
	// integrity of the data structure.
	l := len(s1)
	if l2 := len(s2); l2 < l {
		l = l2
	}
	for i := 0; i < l; i++ {
		if s1[i] != s2[i] {
			return i
		}
	}
	return l
}

type Alpha struct {
	prefix   string
	value    interface{}
	children []*Alpha // sorted here based on their prefixes
}

// searchChildren returns the index of the child node that the key should be
// found inside.
func (tn *Alpha) searchChildren(key string) int {
	const debug = false
	if debug {
		log.Printf("searching for %q at %+v", key, tn)
	}
	i := sort.Search(len(tn.children), func(i int) bool {
		return tn.children[i].prefix >= key
	})
	// NOTE: Index could be to the left, e.g. "samuel" sorts after "sally",
	// but ought to cause a split at "sa".
	if j := i - 1; j >= 0 && offsetOfMismatch(key, tn.children[j].prefix) > 0 {
		if debug {
			log.Printf("moving to the left due to shared prefix")
		}
		return j
	}
	return i
}

// Delete removes the specified key value pair from the Trie.
func (tn *Alpha) Delete(key string) {
	const debug = false

	curr := tn // start at this node
	var previous *Alpha
	var childIndex int

	for {
		if debug {
			log.Printf("delete %q at %+v", key, curr)
		}

		i := offsetOfMismatch(key, curr.prefix)

		if i < len(curr.prefix) {
			if debug {
				log.Printf("fewer characters shared than curr.prefix: %q %q", curr.prefix[:i], curr.prefix)
			}
			return
		}

		if debug {
			log.Printf("i: %d; len(curr.children): %d", i, len(curr.children))
		}

		if i == len(key) && len(curr.children) > 0 && curr.children[0].prefix == "" {
			if debug {
				log.Printf("delete 0th element from children: %+v", curr.children[0])
			}
			copy(curr.children[0:], curr.children[1:])
			curr.children[len(curr.children)-1] = nil
			curr.children = curr.children[:len(curr.children)-1]

			switch len(curr.children) {
			case 0:
				if debug {
					log.Printf("delete %d element from previous' children: %+v", childIndex, previous.children[childIndex])
				}
				copy(previous.children[childIndex:], previous.children[childIndex+1:])
				previous.children[len(previous.children)-1] = nil
				previous.children = previous.children[:len(previous.children)-1]

				if len(previous.children) == 1 {
					if debug {
						log.Printf("merge previous with its only child: %+v; %+v", previous, previous.children[0])
					}
					previous.prefix += previous.children[0].prefix
					previous.children[0].prefix = ""
					previous.children[0].value = previous.children[0].children[0].value
					previous.children[0].children = previous.children[0].children[:0]
				}
			case 1:
				if debug {
					log.Printf("merge current with its only child: %+v; %+v", curr, curr.children[0])
					for _, child := range curr.children {
						log.Printf("    %+v", child)
					}
				}
				curr.prefix += curr.children[0].prefix
				curr.children[0].prefix = ""
				curr.children[0].value = curr.children[0].children[0].value
				curr.children[0].children = curr.children[0].children[:0]
			}

			return
		}

		suffix := key[i:]
		if debug {
			log.Printf("shared %d; prefix: %q; suffix: %q", i, key[:i], suffix)
		}

		childIndex = curr.searchChildren(suffix)
		if childIndex == len(curr.children) {
			// When insertion point is beyond right end of slice, then the
			// key is not found.
			if debug {
				log.Printf("key not found: %q", key)
			}
			return
		}
		// Next visit the child most likely to have this key.
		if debug {
			log.Printf("going to child: %d", childIndex)
		}
		previous = curr
		curr = curr.children[childIndex]
		key = suffix
	}
}

// Load returns the value associated with the specified key, along with a
// boolean which is true when the trie has the specified key.
func (tn *Alpha) Load(key string) (interface{}, bool) {
	const debug = false

	curr := tn // start at this node

	for {
		if debug {
			log.Printf("load %q at %+v", key, curr)
		}

		i := offsetOfMismatch(key, curr.prefix)

		if i < len(curr.prefix) {
			if debug {
				log.Printf("fewer characters shared than curr.prefix: %q %q", curr.prefix[:i], curr.prefix)
			}
			return nil, false
		}

		if debug {
			log.Printf("i: %d; len(curr.children): %d", i, len(curr.children))
		}

		if i == len(key) && len(curr.children) > 0 && curr.children[0].prefix == "" {
			if debug {
				log.Printf("found value: %+v", curr.children[0])
			}
			return curr.children[0].value, true
		}

		suffix := key[i:]
		if debug {
			log.Printf("shared %d; prefix: %q; suffix: %q", i, key[:i], suffix)
		}

		i = curr.searchChildren(suffix)
		if i == len(curr.children) {
			// When insertion point is beyond right end of slice, then the
			// key is not found.
			return nil, false
		}
		// Next visit the child most likely to have this key.
		if debug {
			log.Printf("going to child: %d", i)
		}
		curr = curr.children[i]
		key = suffix
	}
}

// Store stores the specified key and value in the Trie.
func (tn *Alpha) Store(key string, value interface{}) {
	const debug = false

	curr := tn // start at this node
	var previous *Alpha
	var childIndex int

	for {
		if debug {
			log.Printf("store %q at %+v; previous: %+v", key, curr, previous)
		}

		i := offsetOfMismatch(key, curr.prefix)
		if debug {
			log.Printf("shared: %d; %q", i, key[:i])
		}

		if i < len(curr.prefix) {
			// "sally" --> "sam", shared will be "sa"
			shared := key[:i]
			if debug {
				log.Printf("need to split: %+v; shared: %q", curr, shared)
			}
			newParent := &Alpha{
				prefix:   shared,
				children: make([]*Alpha, 2),
			}
			newNode := &Alpha{
				prefix:   key[i:],
				children: []*Alpha{&Alpha{value: value}},
			}
			curr.prefix = curr.prefix[i:]
			if newNode.prefix < curr.prefix {
				newParent.children[0] = newNode
				newParent.children[1] = curr
			} else {
				newParent.children[0] = curr
				newParent.children[1] = newNode
			}
			// re-wire parent to newNode
			if debug {
				log.Printf("childIndex: %d", childIndex)
				log.Printf("previous:\n%s", string(previous.Bytes()))
				log.Printf("newParent:\n%s", string(newParent.Bytes()))
			}
			previous.children[childIndex] = newParent
			return
		}

		if i == len(key) {
			if debug {
				log.Printf("entire key here")
			}
			if len(curr.children) > 0 && curr.children[0].prefix == "" {
				if debug {
					log.Printf("updating value at first child")
				}
				curr.children[0].value = value
				return
			}
			if debug {
				log.Printf("inserting new child to store value")
			}
			curr.insertChildAtIndex(0, &Alpha{value: value})
			return
		}

		suffix := key[i:]
		childIndex = curr.searchChildren(suffix)

		// does new key share any characters with the specified child?
		if debug {
			log.Printf("childIndex: %d", childIndex)
		}
		if childIndex == len(curr.children) || offsetOfMismatch(suffix, curr.children[childIndex].prefix) == 0 {
			if debug {
				log.Printf("inserting new child to store value")
			}
			curr.insertChildAtIndex(childIndex, &Alpha{
				prefix:   suffix,
				children: []*Alpha{&Alpha{value: value}},
			})
			return
		}

		// go visit the child
		if debug {
			log.Printf("going to child: %d", childIndex)
		}
		previous = curr
		curr = curr.children[childIndex]
		key = suffix
	}
}

func (tn *Alpha) insertChildAtIndex(i int, node *Alpha) {
	if i == len(tn.children) {
		// handles nil list, as well as optimizing for tail of list
		tn.children = append(tn.children, node)
		return
	}
	// Without two copies and mandatory allocation, insert string into array at
	// index.
	tn.children = append(tn.children, tn.children[len(tn.children)-1])
	copy(tn.children[i+1:], tn.children[i:len(tn.children)-1])
	tn.children[i] = node
}

// Bytes returns a slice of bytes representing a hierarchical display of the
// Trie.
func (tn *Alpha) Bytes() []byte {
	bb := new(bytes.Buffer)
	tn.bytes(0, true, bb)
	return bb.Bytes()
}

func (tn *Alpha) bytes(level int, last bool, bb *bytes.Buffer) {
	var indention string
	switch level {
	default:
		for i := 0; i < level; i++ {
			indention += "    "
		}
		if last {
			indention += "└── "
		} else {
			indention += "├── "
		}
	}
	bb.Write([]byte(indention))
	switch tn.prefix {
	case "":
		bb.Write([]byte(fmt.Sprintf(". = %v\n", tn.value)))
	default:
		bb.Write(append([]byte(tn.prefix), '\n'))
	}
	for i, child := range tn.children {
		child.bytes(level+1, i == len(tn.children)-1, bb)
	}
}

// Display writes a hierarchical display of the Trie to the specified io.Writer.
func (tn *Alpha) Display(w io.Writer) {
	tn.display(w, 0)
}

func (tn *Alpha) display(w io.Writer, level int) {
	var indentation string
	for i := 0; i < level; i++ {
		indentation += "    "
	}
	if tn.prefix == "" {
		w.Write([]byte(fmt.Sprintf("%s %v\n", indentation, tn.value)))
	} else {
		w.Write([]byte(fmt.Sprintf("%s%q\n", indentation, tn.prefix)))
	}
	for _, child := range tn.children {
		child.display(w, level+1)
	}
}

// Keys returns a list of strings with the specified prefix, to a max of

// Keys returns a slice of keys strings from the Trie with the specified prefix,
// with no more strings than the specified limit. An empty prefix string matches
// all Trie keys. A limit of 0 returns all matching keys, not just the first N
// found.
func (tn *Alpha) Keys(prefix string, limit int) []string {
	const debug = true

	curr, extra := tn.keysFindStartingNode(prefix)
	if curr == nil {
		if debug {
			log.Printf("no starting node found")
		}
		return nil
	}
	if debug {
		log.Printf("starting node: %+v", curr)
	}
	return curr.keysCollectDescendants(prefix, extra, limit)
}

// keysFindStartingNode returns the Trie element that matches the specified
// prefix.
func (tn *Alpha) keysFindStartingNode(prefix string) (*Alpha, int) {
	const debug = false
	curr := tn

	// goal is to eat up the prefix key

	for {
		if debug {
			log.Printf("looking for %q at %+v", prefix, curr)
		}

		i := offsetOfMismatch(prefix, curr.prefix)

		if i == len(prefix) {
			if debug {
				log.Printf("consumed all of the prefix: %+v", curr)
			}
			return curr, len(curr.prefix) - i
		}

		// if more prefix left over, and only portion of curr.prefix matched,
		// then there are no matches for this prefix.
		if i < len(curr.prefix) {
			if debug {
				log.Printf("??? nothing will match this prefix: %q", prefix)
			}
			return nil, 0
		}

		// check children there's more of the prefix to match on
		suffix := prefix[i:]
		i = curr.searchChildren(suffix)
		if debug {
			log.Printf("going to child: %d", i)
		}
		curr = curr.children[i]
		prefix = suffix
	}
}

func (tn *Alpha) keysCollectDescendants(prefix string, extra, limit int) []string {
	const debug = true
	var list []string

	// There may be extra characters on this node's prefix that we want to add
	// to all of the descendants below.
	if extra > 0 { // guard used for performance reason only.
		prefix += tn.prefix[len(tn.prefix)-extra:]
	}
	if debug {
		log.Printf("prefix: %q", prefix)
		log.Printf("collect descendants of %+v", tn)
	}

	for _, child := range tn.children {
		if debug {
			log.Printf("child: %+v", child)
		}
		if child.prefix == "" {
			if debug {
				log.Printf("adding node: %q", prefix)
			}
			list = append(list, "")
			if limit > 0 && len(list) == limit {
				return list
			}
			continue // data node does not have children
		}
		for _, item := range child.keysCollectDescendants(child.prefix, 0, limit) {
			if debug {
				log.Printf("item: %v", item)
			}
			list = append(list, prefix+child.prefix+item)
			if limit > 0 && len(list) >= limit {
				return list[:limit]
			}
		}
		if debug {
			log.Printf("list: %v", list)
		}
	}
	return list
}
