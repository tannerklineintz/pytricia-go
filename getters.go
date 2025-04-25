package pytricia

// Get: longest-prefix match – returns the stored value (or nil)
func (t *PyTricia) Get(cidr string) interface{} {
	if n := t.getNode(cidr); n != nil {
		return n.value
	}
	return nil
}

// GetKey: returns the CIDR string that actually stored the value
func (t *PyTricia) GetKey(cidr string) string {
	if n := t.getNode(cidr); n != nil {
		if c := n.cidr(); c != nil {
			return c.String()
		}
	}
	return ""
}

// GetKV: key + value in one call (avoids 2× parseCIDR)
func (t *PyTricia) GetKV(cidr string) (string, interface{}) {
	if n := t.getNode(cidr); n != nil {
		if c := n.cidr(); c != nil {
			return c.String(), n.value
		}
	}
	return "", nil
}

// Contains: does a prefix (or IP) resolve to *anything*?
func (t *PyTricia) Contains(cidr string) bool { return t.Get(cidr) != nil }

// HasKey: exact-match test (node must *store* a value at that prefix)
func (t *PyTricia) HasKey(cidr string) bool { return t.keyNode(cidr) != nil }

// keyNode: exact-match node (no LPM) – read-lock held only for traversal
func (t *PyTricia) keyNode(cidr string) *PyTricia {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return nil
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	n := t
	for i := 0; i < ones; i++ {
		n = n.children[bit(ip, i)]
		if n == nil {
			return nil
		}
	}
	return n
}

// getNode: longest-prefix match (LPM)
func (t *PyTricia) getNode(cidr string) *PyTricia {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return nil
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	n, best := t, (*PyTricia)(nil)
	for i := 0; i < ones; i++ {
		n = n.children[bit(ip, i)]
		if n == nil {
			break
		}
		if n.value != nil {
			best = n
		}
	}
	return best
}
