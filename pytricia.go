package pytricia

import (
	"math/big"
	"net"
	"strconv"
)

// TrieNode represents a node in the trie.
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

func (t *PyTricia) Get(cidr string) interface{} {
	if net.ParseIP(cidr) != nil {
		return t.getIP(cidr)
	} else {
		return t.getCIDR(cidr)
	}
}

// getCIDR finds the key for a CIDR range.
func (t *PyTricia) getCIDR(cidr string) interface{} {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	startIp := ip
	endIp := lastIP(ipnet)

	startVal := t.getIP(startIp.String())
	endVal := t.getIP(endIp.String())

	if startVal == endVal && startVal != nil {
		return startVal
	}
	return nil
}

// getIP finds the key for an IP address.
func (t *PyTricia) getIP(ip string) interface{} {
	netIp := net.ParseIP(ip)
	currentNode := t
	for _, bit := range ipToBinary(netIp) {
		if currentNode.children[bit] == nil {
			return nil
		}
		currentNode = currentNode.children[bit]
	}
	return currentNode.value
}

// returns children, if any
func (node *PyTricia) Children() [2]*PyTricia {
	return node.children
}

// returns parent, if any
func (node *PyTricia) Parent() *PyTricia {
	return node.parent
}

// last ip address in a cidr
func lastIP(ipnet *net.IPNet) net.IP {
	ip := ipnet.IP

	var lastIP big.Int
	lastIP.SetBytes(ip.To16()) // Ensure the IP is in 16 byte format
	var networkSize big.Int
	networkSize.SetBytes(net.IP(ipnet.Mask).To16())

	var ones, bits = ipnet.Mask.Size()
	var totalIPs big.Int
	totalIPs.Lsh(big.NewInt(1), uint(bits-ones))

	lastIP.Add(&lastIP, &totalIPs)
	lastIP.Sub(&lastIP, big.NewInt(1)) // Subtract 1 to get the last address

	ipBytes := lastIP.Bytes()
	if len(ipBytes) == 16 {
		return net.IP(ipBytes)
	}
	return net.IPv4(ipBytes[12], ipBytes[13], ipBytes[14], ipBytes[15])
}

// ipToBinary converts an IP address to a binary representation.
func ipToBinary(ip net.IP) []int {
	bits := make([]int, 0)

	// Ensure the IP is in 16-byte format
	ip = ip.To16()

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
