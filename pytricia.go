package pytricia

import (
	"net"
	"strconv"
)

// NewPyTricia initializes pytricia object
func NewPyTricia() *PyTricia {
	return &PyTricia{
		children: [2]*PyTricia{nil, nil},
		parent:   nil,
		value:    nil,
	}
}

// PyTricia represents a node in the PyTricia trie.
type PyTricia struct {
	children [2]*PyTricia
	parent   *PyTricia
	value    interface{}
}

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
}

// IsRoot returns whether this PyTricia object is the trie's root node
// (Only needed for more manual operations)
func (t *PyTricia) IsRoot() bool {
	return t.parent == nil
}

// Contains returns whether a cidr is contained within the trie
func (t *PyTricia) Contains(cidr string) bool {
	return t.Get(cidr) != nil
}

// HasKey returns whether a cidr is a key within the trie
func (t *PyTricia) HasKey(cidr string) bool {
	if net.ParseIP(cidr) != nil {
		return false
	} else {
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return false
		}
		ones, _ := ipnet.Mask.Size()
		currentNode := t
		for i, bit := range ipToBinary(ip) {
			if i >= ones {
				break
			}
			if currentNode.children[bit] == nil {
				return false
			}
			currentNode = currentNode.children[bit]
		}
		return currentNode.value != nil
	}
}

// Get returns value associated with CIDR or IP address
func (t *PyTricia) Get(cidr string) interface{} {
	var ip net.IP
	var ipnet *net.IPNet
	var ones int
	var err error

	if ip = net.ParseIP(cidr); ip != nil {
		ones = 32
	} else {
		ip, ipnet, err = net.ParseCIDR(cidr)
		if err != nil {
			return nil
		}
		ones, _ = ipnet.Mask.Size()
	}

	currentNode := t
	var currentValue interface{} = nil
	for i, bit := range ipToBinary(ip) {
		if i >= ones {
			break
		}
		if currentNode.children[bit] == nil {
			break
		}
		currentNode = currentNode.children[bit]
		if currentNode.value != nil {
			currentValue = currentNode.value
		}
	}
	return currentValue
}

// returns children, if any
func (node *PyTricia) Children() [2]*PyTricia {
	return node.children
}

// returns parent, if any
func (node *PyTricia) Parent() *PyTricia {
	return node.parent
}

// ipToBinary converts an IP address to a binary representation.
func ipToBinary(ip net.IP) []int {
	bits := make([]int, 0)

	// Ensure the IP is in 16-byte format
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}

	for _, b := range ip {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}
	return bits
}

// binaryToCIDR converts a binary path to CIDR notation for both IPv4 and IPv6.
func binaryToCIDR(path []int) *net.IPNet {
	var ip net.IP
	length := len(path)
	for i := 0; i < length; i += 8 {
		byteVal := 0
		for j := 0; j < 8 && i+j < length; j++ {
			byteVal = byteVal*2 + path[i+j]
		}
		ip = append(ip, byte(byteVal))
	}

	// Determine if the address is IPv4 or IPv6 based on the length
	var cidr string
	if length <= 32 {
		cidr = ip.To4().String() + "/" + strconv.Itoa(length)
	} else {
		cidr = ip.To16().String() + "/" + strconv.Itoa(length)
	}

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	return ipnet
}
