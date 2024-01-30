package pytricia

// ToMap converts the PyTricia trie into a map of CIDR strings to their associated values
func (t *PyTricia) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	if t == nil {
		return result
	}

	stack := [][3]interface{}{{t, []byte{}, 0}}

	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		node := item[0].(*PyTricia)
		path := item[1].([]byte)
		depth := item[2].(int)

		if node.value != nil {
			cidr := binaryToCIDR(path[:depth], node.ipType)
			if cidr != nil {
				result[cidr.String()] = node.value
			}
		}

		if node.children[1] != nil {
			newPath := make([]byte, len(path))
			copy(newPath, path)
			stack = append(stack, [3]interface{}{node.children[1], append(newPath, 1), depth + 1})
		}
		if node.children[0] != nil {
			newPath := make([]byte, len(path))
			copy(newPath, path)
			stack = append(stack, [3]interface{}{node.children[0], append(newPath, 0), depth + 1})
		}
	}
	return result
}

// Keys returns all keys in the trie
func (t *PyTricia) Keys() []string {
	result := []string{}
	if t == nil {
		return result
	}

	stack := [][3]interface{}{{t, []byte{}, 0}}

	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		node := item[0].(*PyTricia)
		path := item[1].([]byte)
		depth := item[2].(int)

		if node.value != nil {
			cidr := binaryToCIDR(path[:depth], node.ipType)
			if cidr != nil {
				result = append(result, cidr.String())
			}
		}

		if node.children[1] != nil {
			newPath := make([]byte, len(path))
			copy(newPath, path)
			stack = append(stack, [3]interface{}{node.children[1], append(newPath, 1), depth + 1})
		}
		if node.children[0] != nil {
			newPath := make([]byte, len(path))
			copy(newPath, path)
			stack = append(stack, [3]interface{}{node.children[0], append(newPath, 0), depth + 1})
		}
	}
	return result
}

// Values returns all keys in the trie
func (t *PyTricia) Values() []interface{} {
	result := []interface{}{}
	if t == nil {
		return result
	}

	stack := [][3]interface{}{{t, []byte{}, 0}}

	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	for len(stack) > 0 {
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		node := item[0].(*PyTricia)
		path := item[1].([]byte)
		depth := item[2].(int)

		if node.value != nil {
			result = append(result, node.value)
		}

		if node.children[1] != nil {
			newPath := make([]byte, len(path))
			copy(newPath, path)
			stack = append(stack, [3]interface{}{node.children[1], append(newPath, 1), depth + 1})
		}
		if node.children[0] != nil {
			newPath := make([]byte, len(path))
			copy(newPath, path)
			stack = append(stack, [3]interface{}{node.children[0], append(newPath, 0), depth + 1})
		}
	}
	return result
}
