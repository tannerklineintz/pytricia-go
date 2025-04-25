package pytricia

// ---------------------------------------------------------------------------
// ToMap: snapshot of every <CIDR,value> in the trie.
// ---------------------------------------------------------------------------
func (t *PyTricia) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	if t == nil {
		return out
	}

	// Pin the root briefly.
	t.mutex.RLock()
	start := t
	t.mutex.RUnlock()

	stack := []*PyTricia{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if v := n.value; v != nil {
			if c := n.cidr(); c != nil {
				out[c.String()] = v
			}
		}
		if r := n.children[1]; r != nil {
			stack = append(stack, r)
		}
		if l := n.children[0]; l != nil {
			stack = append(stack, l)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Keys: every CIDR stored in the trie.
// ---------------------------------------------------------------------------
func (t *PyTricia) Keys() []string {
	keys := []string{}
	if t == nil {
		return keys
	}

	t.mutex.RLock()
	start := t
	t.mutex.RUnlock()

	stack := []*PyTricia{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n.value != nil {
			if c := n.cidr(); c != nil {
				keys = append(keys, c.String())
			}
		}
		if r := n.children[1]; r != nil {
			stack = append(stack, r)
		}
		if l := n.children[0]; l != nil {
			stack = append(stack, l)
		}
	}
	return keys
}

// ---------------------------------------------------------------------------
// Values: every stored value (order parallels Keys()).
// ---------------------------------------------------------------------------
func (t *PyTricia) Values() []interface{} {
	vals := []interface{}{}
	if t == nil {
		return vals
	}

	t.mutex.RLock()
	start := t
	t.mutex.RUnlock()

	stack := []*PyTricia{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if v := n.value; v != nil {
			vals = append(vals, v)
		}
		if r := n.children[1]; r != nil {
			stack = append(stack, r)
		}
		if l := n.children[0]; l != nil {
			stack = append(stack, l)
		}
	}
	return vals
}
