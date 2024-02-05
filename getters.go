package pytricia

// Get returns value associated with CIDR or IP
func (t *PyTricia) Get(cidr string) interface{} {
	if node := t.getNode(cidr); node != nil {
		return node.value
	}
	return nil
}

// GetKey returns key associated with CIDR or IP
func (t *PyTricia) GetKey(cidr string) string {
	if node := t.getNode(cidr); node != nil {
		return node.cidr().String()
	}
	return ""
}

// GetKey returns key value pair associated with CIDR or IP
func (t *PyTricia) GetKV(cidr string) (string, interface{}) {
	if node := t.getNode(cidr); node != nil {
		return node.cidr().String(), node.value
	}
	return "", nil
}

// Contains returns whether a CIDR or IP is contained within the trie
func (t *PyTricia) Contains(cidr string) bool {
	return t.Get(cidr) != nil
}

// HasKey returns whether a CIDR or IP is a key within the trie
func (t *PyTricia) HasKey(cidr string) bool {
	return t.keyNode(cidr) != nil
}

// getKeyNode returns node associated with CIDR or IP
// only if direct key match present
func (t *PyTricia) keyNode(cidr string) *PyTricia {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return nil
	}

	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	currentNode := t
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			return nil
		}
		currentNode = currentNode.children[bit]
	}
	return currentNode
}

// GetNode returns the node associated with CIDR or IP
func (t *PyTricia) getNode(cidr string) *PyTricia {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return nil
	}

	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
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
