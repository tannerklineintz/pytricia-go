package pytricia

import (
	"errors"
	"sync"
)

// Delete removes the node associated with the given CIDR or IP.
func (t *PyTricia) Delete(cidr string) error {
	node := t.keyNode(cidr)
	if node == nil {
		return errors.New("CIDR not found")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Set the node's value to nil.
	node.value = nil

	// Remove any unnecessary parent nodes.
	for node.parent != nil && node.children[0] == nil && node.children[1] == nil && node.value == nil {
		parent := node.parent
		if parent.children[0] == node {
			parent.children[0] = nil
		} else {
			parent.children[1] = nil
		}
		node = parent
	}

	return nil
}

// Clear deallocates the entire trie.
func (t *PyTricia) Clear() {
	*t = PyTricia{
		children: [2]*PyTricia{nil, nil},
		parent:   nil,
		value:    nil,
		mutex:    sync.RWMutex{},
	}
}
