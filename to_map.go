package pytricia

// ToMap converts the PyTricia trie into a map of CIDR strings to their associated values
func (t *PyTricia) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	t.toMapHelper(result, []byte{}, 0)
	return result
}

// toMapHelper is a recursive helper function for ToMap
func (t *PyTricia) toMapHelper(result map[string]interface{}, path []byte, depth int) {
	if t == nil {
		return
	}

	// Check if the current node has a value and add it to the map
	if t.value != nil {
		cidr := binaryToCIDR(path[:depth], t.ipType)
		if cidr != nil {
			result[cidr.String()] = t.value
		}
	}

	// Recursively traverse left and right children
	if t.children[0] != nil {
		t.children[0].toMapHelper(result, append(path, 0), depth+1)
	}
	if t.children[1] != nil {
		t.children[1].toMapHelper(result, append(path, 1), depth+1)
	}
}
