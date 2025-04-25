package pytricia

import "errors"

// Insert: overwrite or create
func (t *PyTricia) Insert(cidr string, value interface{}) error {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return err
	}

	node := t
	i := 0

	// 1) Walk under read-lock until we hit a nil edge
	t.mutex.RLock()
	for ; i < ones; i++ {
		b := bit(ip, i)
		if next := node.children[b]; next != nil {
			node = next
			continue
		}
		break
	}
	t.mutex.RUnlock()

	// 2) If we stopped early, grab the write-lock once and finish path
	if i < ones {
		t.mutex.Lock()
		for ; i < ones; i++ {
			b := bit(ip, i)
			if node.children[b] == nil { // **double-check after lock**
				node.children[b] = &PyTricia{
					parent:   node,
					children: [2]*PyTricia{},
					value:    nil,
					ipType:   typeIP(cidr),
				}
			}
			node = node.children[b]
		}
		// still holding write-lock → set value
		node.value = value
		node.ipType = typeIP(cidr)
		t.mutex.Unlock()
		return nil
	}

	// 3) Path existed; just update value (very short write-lock)
	t.mutex.Lock()
	node.value = value
	node.ipType = typeIP(cidr)
	t.mutex.Unlock()
	return nil
}

// Set: overwrite only if CIDR already present
func (t *PyTricia) Set(cidr string, value interface{}) error {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return err
	}

	node := t
	t.mutex.RLock()
	for i := 0; i < ones; i++ {
		b := bit(ip, i)
		if node = node.children[b]; node == nil {
			t.mutex.RUnlock()
			return errors.New("CIDR not present")
		}
	}
	if node.value == nil {
		t.mutex.RUnlock()
		return errors.New("CIDR not present")
	}
	t.mutex.RUnlock()

	// value exists → acquire write-lock just to mutate
	t.mutex.Lock()
	node.value = value
	node.ipType = typeIP(cidr)
	t.mutex.Unlock()
	return nil
}

// Add: insert only if CIDR *not* already present
func (t *PyTricia) Add(cidr string, value interface{}) error {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return err
	}

	node := t
	i := 0

	// 1) Read-only walk until gap or end
	t.mutex.RLock()
	for ; i < ones; i++ {
		b := bit(ip, i)
		if next := node.children[b]; next != nil {
			node = next
			continue
		}
		break
	}
	alreadyExists := (i == ones && node.value != nil)
	t.mutex.RUnlock()

	if alreadyExists {
		return errors.New("CIDR already present")
	}

	// 2) Need to create nodes or set value → single write-lock
	t.mutex.Lock()
	// (re-do the walk from the point we left off; node is still correct)
	for ; i < ones; i++ {
		b := bit(ip, i)
		if node.children[b] == nil {
			node.children[b] = &PyTricia{
				parent:   node,
				children: [2]*PyTricia{},
				value:    nil,
				ipType:   typeIP(cidr),
			}
		}
		node = node.children[b]
	}

	if node.value != nil { // in case another writer beat us
		t.mutex.Unlock()
		return errors.New("CIDR already present")
	}
	node.value = value
	node.ipType = typeIP(cidr)
	t.mutex.Unlock()
	return nil
}
