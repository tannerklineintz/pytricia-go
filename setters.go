package pytricia

import (
	"errors"
)

// Insert inserts an IP or CIDR and its value into the trie. This
// overwrites the value if already present
func (t *PyTricia) Insert(cidr string, value interface{}) error {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return err
	}

	t.mutex.RLock()
	currentNode := t
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			t.mutex.RUnlock()
			t.mutex.Lock()
			currentNode.children[bit] = &PyTricia{
				parent:   currentNode,
				children: [2]*PyTricia{nil, nil},
				value:    nil,
			}
			t.mutex.Unlock()
			t.mutex.RLock()
		}
		currentNode = currentNode.children[bit]
	}
	t.mutex.RUnlock()

	t.mutex.Lock()
	currentNode.value = value
	currentNode.ipType = typeIP(cidr)
	t.mutex.Unlock()

	return nil
}

// Sets value of IP or CIDR, only if the value already exists;
// returns error if CIDR not already inserted.
func (t *PyTricia) Set(cidr string, value interface{}) error {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return err
	}

	t.mutex.RLock()
	currentNode := t
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			t.mutex.RUnlock()
			return errors.New("CIDR not present")
		}
		currentNode = currentNode.children[bit]
	}
	if currentNode.value == nil {
		t.mutex.RUnlock()
		return errors.New("CIDR not present")
	}
	t.mutex.RUnlock()

	t.mutex.Lock()
	currentNode.value = value
	currentNode.ipType = typeIP(cidr)
	t.mutex.Unlock()

	return nil
}

// Sets value of IP or CIDR, only if the value doesn't already exist;
// returns error if CIDR is already inserted.
func (t *PyTricia) Add(cidr string, value interface{}) error {
	ip, ones, err := parseCIDR(cidr)
	if err != nil {
		return err
	}

	t.mutex.RLock()
	currentNode := t
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			t.mutex.RUnlock()
			t.mutex.Lock()
			currentNode.children[bit] = &PyTricia{
				parent:   currentNode,
				children: [2]*PyTricia{nil, nil},
				value:    nil,
			}
			t.mutex.Unlock()
			t.mutex.RLock()
		}
		currentNode = currentNode.children[bit]
	}
	if currentNode.value != nil {
		t.mutex.RUnlock()
		return errors.New("CIDR already present")
	}
	t.mutex.RUnlock()

	t.mutex.Lock()
	currentNode.value = value
	currentNode.ipType = typeIP(cidr)
	t.mutex.Unlock()

	return nil
}
