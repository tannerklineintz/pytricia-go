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

// CIDR returns the CIDR representation of the current PyTricia node in the trie.
func (node *PyTricia) cidr() *net.IPNet {
	// Start from the current node and traverse up to the root to construct the path.
	var path []byte
	currentNode := node
	for currentNode.parent != nil {
		if currentNode == currentNode.parent.children[0] {
			path = append([]byte{0}, path...)
		} else {
			path = append([]byte{1}, path...)
		}
		currentNode = currentNode.parent
	}

	// Convert the path to CIDR.
	ipnet := binaryToCIDR(path, node.ipType)
	if ipnet == nil {
		return nil
	}

	return ipnet
}
