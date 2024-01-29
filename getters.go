package pytricia

import (
	"net"
)

// Get returns value associated with CIDR or IP
func (t *PyTricia) Get(cidr string) interface{} {
	if node := t.GetNode(cidr); node != nil {
		return node.value
	}
	return nil
}

// GetNode returns the node associated with CIDR or IP
func (t *PyTricia) GetNode(cidr string) *PyTricia {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return nil
	}

	currentNode := t
	var currentValue *PyTricia = nil
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			break
		}
		currentNode = currentNode.children[bit]
		if currentNode.value != nil {
			currentValue = currentNode
		}
	}
	return currentValue
}

// Contains returns whether a CIDR or IP is contained within the trie
func (t *PyTricia) Contains(cidr string) bool {
	return t.Get(cidr) != nil
}

// HasKey returns whether a CIDR or IP is a key within the trie
func (t *PyTricia) HasKey(cidr string) bool {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return false
	}

	currentNode := t
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			return false
		}
		currentNode = currentNode.children[bit]
	}
	return currentNode.value != nil
}

// IsRoot returns whether this PyTricia object is the trie's root node
// (Only needed for more manual operations)
func (t *PyTricia) IsRoot() bool {
	return t.parent == nil
}

// CIDR returns the CIDR representation of the current PyTricia node in the trie.
func (node *PyTricia) CIDR() *net.IPNet {
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
