package pytricia

import "errors"

// Delete removes a prefix (or single IP) and prunes now-empty branches.
func (t *PyTricia) Delete(cidr string) error {
	// 1)  Locate the node quickly under read-lock.
	t.mutex.RLock()
	target := t.keyNode(cidr)
	t.mutex.RUnlock()
	if target == nil {
		return errors.New("CIDR not found")
	}

	// 2)  Promote to write-lock once for the actual mutation.
	t.mutex.Lock()
	// Re-validate in case another writer deleted it meanwhile.
	if target.value == nil && target.children[0] == nil && target.children[1] == nil {
		// Either someone already removed it or it never held a value.
		t.mutex.Unlock()
		return errors.New("CIDR not found")
	}

	// Clear the stored value.
	target.value = nil

	// Prune ancestor chain while the branch is empty.
	for n := target; n.parent != nil &&
		n.value == nil &&
		n.children[0] == nil &&
		n.children[1] == nil; {

		p := n.parent
		if p.children[0] == n {
			p.children[0] = nil
		} else if p.children[1] == n {
			p.children[1] = nil
		}
		n = p
	}
	t.mutex.Unlock()
	return nil
}

// Clear wipes the entire trie in O(1) time while holding the write-lock.
func (t *PyTricia) Clear() {
	t.mutex.Lock()
	// Keep the same mutex instance (can’t replace it while locked).
	t.children[0], t.children[1] = nil, nil
	t.value = nil
	t.parent = nil // root has no parent; safe to set
	t.ipType = 0   // if you track type info
	t.mutex.Unlock()
}
