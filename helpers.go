package pytricia

import (
	"fmt"
	"net"
	"strconv"
)

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
func binaryToCIDR(path []byte, ipType int) *net.IPNet {
	if ipType != 4 && ipType != 6 {
		return nil
	}

	// Initialize variables
	var ipStr string
	maskLen := len(path)
	totalBits := 32 // Default for IPv4

	if ipType == 6 {
		totalBits = 128 // For IPv6
	}

	// Ensure path is the correct length by appending zeros if necessary
	for len(path) < totalBits {
		path = append(path, 0)
	}

	// Convert binary to IP address
	if ipType == 4 {
		// IPv4 - every 8 bits (1 byte) is a part of the IP
		for i := 0; i < totalBits; i += 8 {
			if i > 0 {
				ipStr += "."
			}
			byteVal := binarySliceToByte(path[i:min(i+8, totalBits)])
			ipStr += strconv.Itoa(int(byteVal))
		}
	} else {
		// IPv6 - every 16 bits (2 bytes or 1 hextet) is a part of the IP
		for i := 0; i < totalBits; i += 16 {
			if i > 0 {
				ipStr += ":"
			}
			hextet := binarySliceToUint16(path[i:min(i+16, totalBits)])
			ipStr += fmt.Sprintf("%04x", hextet)
		}
	}

	// Parse the IP and mask
	ip, ipNet, err := net.ParseCIDR(ipStr + "/" + strconv.Itoa(maskLen))
	if err != nil {
		fmt.Println("Error parsing CIDR:", err)
		return nil
	}
	ipNet.IP = ip // Ensure that the IPNet struct has the correct IP

	return ipNet
}

// Helper function to convert a slice of binary to a byte
func binarySliceToByte(bits []byte) byte {
	var num byte
	for _, bit := range bits {
		num = (num << 1) + byte(bit)
	}
	return num
}

// Helper function to convert a slice of binary to a uint16 for IPv6
func binarySliceToUint16(bits []byte) uint16 {
	var num uint16
	for _, bit := range bits {
		num = (num << 1) + uint16(bit)
	}
	return num
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// typeIP returns type of IP address / CIDR
// returns -1 if not a valid ip
func typeIP(cidr string) int {
	for i := 0; i < 6; i++ {
		if cidr[i] == byte('.') {
			return 4
		} else if cidr[i] == byte(':') {
			return 6
		}
	}
	return -1
}
