package pytricia

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// ipToBinary converts an IP address to a binary representation.
func ipToBinary(ip net.IP) []int {
	// Determine the length based on IP type (IPv4 or IPv6)
	var totalBits int
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
		totalBits = 32 // IPv4
	} else {
		totalBits = 128 // IPv6
	}
	bits := make([]int, totalBits)

	// Process each byte of the IP address
	for i, b := range ip {
		baseIndex := i * 8
		bits[baseIndex] = int((b >> 7) & 1)
		bits[baseIndex+1] = int((b >> 6) & 1)
		bits[baseIndex+2] = int((b >> 5) & 1)
		bits[baseIndex+3] = int((b >> 4) & 1)
		bits[baseIndex+4] = int((b >> 3) & 1)
		bits[baseIndex+5] = int((b >> 2) & 1)
		bits[baseIndex+6] = int((b >> 1) & 1)
		bits[baseIndex+7] = int(b & 1)
	}
	return bits
}

// binaryToCIDR converts a binary path to CIDR notation for both IPv4 and IPv6.
func binaryToCIDR(path []byte, ipType int) *net.IPNet {
	if ipType != 4 && ipType != 6 {
		return nil
	}

	totalBits, increment := 32, 8
	if ipType == 6 {
		totalBits, increment = 128, 16
	}

	// Preallocate slice to required size
	pathLen := len(path)
	if len(path) < totalBits {
		extendedPath := make([]byte, totalBits)
		copy(extendedPath, path)
		path = extendedPath
	}

	// Convert binary to IP address string
	var ipStrBuilder strings.Builder
	for i := 0; i < totalBits; i += increment {
		if ipType == 4 && i > 0 {
			ipStrBuilder.WriteByte('.')
		}
		if ipType == 6 && i > 0 {
			ipStrBuilder.WriteByte(':')
		}

		if ipType == 4 {
			// IPv4: Process each byte
			byteVal := binarySliceToByte(path[i:min(i+8, totalBits)])
			ipStrBuilder.WriteString(strconv.Itoa(int(byteVal)))
		} else {
			// IPv6: Process each hextet
			hextet := binarySliceToUint16(path[i:min(i+16, totalBits)])
			ipStrBuilder.WriteString(fmt.Sprintf("%04x", hextet))
		}
	}

	// Parse the IP and mask
	ip, ipNet, err := net.ParseCIDR(
		ipStrBuilder.String() + "/" + strconv.Itoa(pathLen))
	if err != nil {
		fmt.Println("Error parsing CIDR:", err)
		return nil
	}
	ipNet.IP = ip // Ensure that the IPNet struct has the correct IP

	return ipNet
}

// binarySliceToByte converts a slice of binary to a byte
func binarySliceToByte(bits []byte) byte {
	var num byte
	for _, bit := range bits {
		num = (num << 1) + byte(bit)
	}
	return num
}

// binarySliceToUint16 converts a slice of binary to a uint16 for IPv6
func binarySliceToUint16(bits []byte) uint16 {
	var num uint16
	for _, bit := range bits {
		num = (num << 1) + uint16(bit)
	}
	return num
}

// min is a helper function for min
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

// parseCIDR parses either a bare IP string ("8.8.8.8") or a CIDR
// ("8.8.8.0/24") and returns:
//
//   - ip   – the address as a byte slice (IPv4 always 4 bytes, IPv6 16 bytes)
//   - ones – the prefix length in bits (32 for a lone IPv4 address, 128 for IPv6)
//   - err  – non-nil only if the input isn’t a valid IP/CIDR
func parseCIDR(cidr string) (net.IP, int, error) {
	var (
		ip    net.IP
		ipnet *net.IPNet
		ones  int
	)

	//----------------------------------------------------------------------
	// Case 1: the string is a bare IP address.
	//----------------------------------------------------------------------
	if ip = net.ParseIP(cidr); ip != nil {
		switch typeIP(cidr) {
		case 4:
			ip = ip.To4() // <-- ensure 4-byte representation
			ones = 32
		case 6:
			ones = 128
		default:
			return nil, 0, errors.New("invalid IP/CIDR")
		}
		return ip, ones, nil
	}

	//----------------------------------------------------------------------
	// Case 2: must be CIDR notation (ip/mask).
	//----------------------------------------------------------------------
	var err error
	ip, ipnet, err = net.ParseCIDR(cidr)
	if err != nil {
		return nil, 0, errors.New("invalid IP/CIDR")
	}

	ones, _ = ipnet.Mask.Size()

	// If the CIDR is IPv4 (mask ≤ 32) make sure we return a 4-byte IP.
	if ip4 := ip.To4(); ip4 != nil && ones <= 32 {
		ip = ip4
	}

	return ip, ones, nil
}

// isCIDR returns whether an string is a CIDR
// or simply an IP address
func isCIDR(cidr string) bool {
	for i := len(cidr); i > len(cidr)-4; i-- {
		if cidr[i] == byte('/') {
			return true
		}
	}
	return false
}

// randomIPv4 returns random ipv4 address
func randomIPv4() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

// randomIPv4CIDR returns random ipv4 CIDR
func randomIPv4CIDR() string {
	ip := randomIPv4()
	mask := rand.Intn(32) + 1 // 1 to 32
	return fmt.Sprintf("%s/%d", ip, mask)
}

// randomHexDigit returns random hex digit
func randomHexDigit() string {
	digits := "0123456789abcdef"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return string(digits[rand.Intn(len(digits))])
}

// randomIPv6Group returns random ipv6 group
func randomIPv6Group() string {
	return fmt.Sprintf("%s%s%s%s", randomHexDigit(), randomHexDigit(), randomHexDigit(), randomHexDigit())
}

// randomIPv6Group returns random ipv6 address
func randomIPv6() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
		randomIPv6Group(), randomIPv6Group(), randomIPv6Group(), randomIPv6Group(),
		randomIPv6Group(), randomIPv6Group(), randomIPv6Group(), randomIPv6Group())
}

// randomIPv6Group returns random ipv6 CIDR
func randomIPv6CIDR() string {
	ip := randomIPv6()
	mask := rand.Intn(128) + 1 // 1 to 128
	return fmt.Sprintf("%s/%d", ip, mask)
}

// bit returns the i-th bit (0-based) of the IP address.
func bit(ip []byte, i int) int {
	return int((ip[i/8] >> (7 - uint(i%8))) & 1)
}
