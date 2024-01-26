package pytricia

import "net"

// Insert inserts a CIDR and its value into the trie.
func (t *PyTricia) Insert(cidr string, value interface{}) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}
	ones, _ := ipnet.Mask.Size()
	currentNode := t
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			currentNode.children[bit] = &PyTricia{
				parent:   currentNode,
				children: [2]*PyTricia{nil, nil},
				value:    nil,
			}
		}
		currentNode = currentNode.children[bit]
	}
	currentNode.value = value
	currentNode.ipType = typeIP(cidr)
}
