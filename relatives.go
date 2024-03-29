package pytricia

// AllChildren returns all child nodes with non-nil values in a slice.
func (t *PyTricia) Children(cidr string) map[string]interface{} {
	children := make(map[string]interface{})
	node := t.getNode(cidr)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	stack := []*PyTricia{node}
	for len(stack) > 0 {
		currentNode := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if currentNode.value != nil {
			c := currentNode.cidr()
			if c != nil {
				children[c.String()] = currentNode.value
			}
		}
		if currentNode.children[1] != nil {
			stack = append(stack, currentNode.children[1])
		}
		if currentNode.children[0] != nil {
			stack = append(stack, currentNode.children[0])
		}
	}
	return children
}

// returns parent, if any
func (t *PyTricia) Parent(cidr string) (string, interface{}) {
	node := t.getNode(cidr)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Start from the current node and traverse up
	currentNode := node.parent
	for currentNode != nil {
		// Check if the current ancestor node has a non-nil value
		if currentNode.value != nil {
			c := currentNode.cidr()
			if c != nil {
				return c.String(), currentNode.value
			}
		}
		currentNode = currentNode.parent
	}
	// If no ancestor with a non-nil value is found, return nil
	return "", nil
}
