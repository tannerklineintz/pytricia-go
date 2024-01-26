package pytricia

// AllChildren returns all child nodes with non-nil values in a slice.
func (node *PyTricia) Children() []*PyTricia {
	var children []*PyTricia
	node.collectChildren(&children)
	return children
}

// collectChildren is a recursive helper function for AllChildren.
// It traverses the trie and collects nodes with non-nil values.
func (node *PyTricia) collectChildren(children *[]*PyTricia) {
	// Check if the current node has a non-nil value.
	if node.value != nil {
		*children = append(*children, node)
	}

	// Recursively traverse left and right children, if they exist.
	if node.children[0] != nil {
		node.children[0].collectChildren(children)
	}
	if node.children[1] != nil {
		node.children[1].collectChildren(children)
	}
}

// returns parent, if any
func (node *PyTricia) Parent() *PyTricia {
	// Start from the current node and traverse up
	currentNode := node.parent
	for currentNode != nil {
		// Check if the current ancestor node has a non-nil value
		if currentNode.value != nil {
			return currentNode
		}
		currentNode = currentNode.parent
	}
	// If no ancestor with a non-nil value is found, return nil
	return nil
}
