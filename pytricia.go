package pytricia

import (
	"net"
	"sync"
)

// NewPyTricia initializes pytricia object
func NewPyTricia() *PyTricia {
	return &PyTricia{
		children: [2]*PyTricia{nil, nil},
		parent:   nil,
		value:    nil,
		mutex:    sync.RWMutex{},
	}
}

// PyTricia represents a node in the PyTricia trie.
type PyTricia struct {
	ipType   int
	children [2]*PyTricia
	parent   *PyTricia
	value    interface{}
	mutex    sync.RWMutex
}

func (n *PyTricia) cidr() *net.IPNet {
	// ─── 1. Find a node that knows the family (4 or 6). ──────────────────
	fam := n
	for fam != nil && fam.ipType == 0 {
		fam = fam.parent
	}
	if fam == nil {
		return nil // should never happen, but guard anyway
	}

	// ─── 2. Build the full bit-path from *root* to the original node. ────
	// We collect bits in reverse, then reverse once at the end because
	// prepending in a loop explodes the allocator.
	var revBits []byte
	for cur := n; cur.parent != nil; cur = cur.parent {
		if cur == cur.parent.children[0] {
			revBits = append(revBits, 0)
		} else {
			revBits = append(revBits, 1)
		}
	}
	// Reverse into forward order.
	bits := make([]byte, len(revBits))
	for i := range revBits {
		bits[len(revBits)-1-i] = revBits[i]
	}

	// ─── 3. Convert the bit slice to *net.IPNet. ─────────────────────────
	return binaryToCIDR(bits, fam.ipType)
}
