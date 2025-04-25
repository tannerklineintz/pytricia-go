package pytricia

// Children returns every descendant whose value is non-nil.
// Snapshot semantics: may miss nodes added *after* the initial lock.
func (t *PyTricia) Children(cidr string) map[string]interface{} {
	out := make(map[string]interface{})

	// 1) Locate the subtree root under a short read-lock.
	t.mutex.RLock()
	start := t.getNode(cidr)
	t.mutex.RUnlock()
	if start == nil {
		return out
	}

	// 2) Depth-first scan without holding the global lock.
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

// Parent returns the first ancestor that carries a non-nil value.
// Snapshot semantics: if another goroutine inserts a closer ancestor
// concurrently you may still see the older one, but never an invalid ptr.
func (t *PyTricia) Parent(cidr string) (string, interface{}) {
	// 1) Pin the start node quickly under read-lock.
	t.mutex.RLock()
	n := t.getNode(cidr)
	t.mutex.RUnlock()
	if n == nil {
		return "", nil
	}

	// 2) Walk upward lock-free.
	for p := n.parent; p != nil; p = p.parent {
		if v := p.value; v != nil {
			if c := p.cidr(); c != nil {
				return c.String(), v
			}
		}
	}
	return "", nil
}
