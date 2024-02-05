package pytricia

import (
	"testing"
)

func TestPytriciaIPv4(t *testing.T) {
	t.Parallel()

	pt := NewPyTricia()

	pt.Insert("8.8.8.0/24", "testing123")
	if val := pt.Get("8.8.8.8"); val != "testing123" {
		t.Errorf("Error on test 1: %v", val)
	}
	if val := pt.Get("8.8.8.10/31"); val != "testing123" {
		t.Errorf("Error on test 2: %v", val)
	}
	if val := pt.Get("8.8.8.0/24"); val != "testing123" {
		t.Errorf("Error on test 3: %v", val)
	}
	if val := pt.Get("8.8.7.0/24"); val != nil {
		t.Errorf("Error on test 4: %v", val)
	}
	if val := pt.Get("8.8.9.0"); val != nil {
		t.Errorf("Error on test 5: %v", val)
	}
	if val := pt.Get("8.8.7.255"); val != nil {
		t.Errorf("Error on test 6: %v", val)
	}

	if val := pt.GetKey("8.8.8.8"); val != "8.8.8.0/24" {
		t.Errorf("Error on test 4: %v", val)
	}

	if !pt.HasKey("8.8.8.0/24") {
		t.Errorf("Error on test 7")
	}
	if pt.HasKey("8.8.8.10/31") {
		t.Errorf("Error on test 8")
	}

	pt.Insert("8.8.8.10/31", "testing456")

	if keys := pt.Keys(); len(keys) != 2 || keys[0] != "8.8.8.0/24" || keys[1] != "8.8.8.10/31" {
		t.Errorf("Error on test 9: %v", keys)
	}

	if values := pt.Values(); len(values) != 2 || values[0] != "testing123" || values[1] != "testing456" {
		t.Errorf("Error on test 9: %v", values)
	}

	children := pt.Children("8.8.8.0/24")
	if children["8.8.8.0/24"] != "testing123" {
		t.Errorf("Error on test 10: %v", children)
	}
	if children["8.8.8.10/31"] != "testing456" {
		t.Errorf("Error on test 10: %v", children)
	}

	expectedMap := map[string]interface{}{
		"8.8.8.0/24":  "testing123",
		"8.8.8.10/31": "testing456",
	}
	mapOutput := pt.ToMap()
	if len(mapOutput) != len(expectedMap) {
		t.Errorf("Error on test 11: %v", mapOutput)
	}
	if mapOutput["8.8.8.0/24"] != expectedMap["8.8.8.0/24"] {
		t.Errorf("Error on test 11: %v", mapOutput)
	}
	if mapOutput["8.8.8.10/31"] != expectedMap["8.8.8.10/31"] {
		t.Errorf("Error on test 11: %v", mapOutput)
	}
}

func TestPytriciaIPv6(t *testing.T) {
	t.Parallel()

	pt := NewPyTricia()

	pt.Insert("2001:250::/38", "testing123")
	if val := pt.Get("2001:250:1:100::"); val != "testing123" {
		t.Errorf("Error on test 1: %v", val)
	}
	if val := pt.Get("2001:250:1::200/119"); val != "testing123" {
		t.Errorf("Error on test 2: %v", val)
	}
	if val := pt.Get("2001:250::/38"); val != "testing123" {
		t.Errorf("Error on test 3: %v", val)
	}
	if val := pt.Get("2001:251::"); val != nil {
		t.Errorf("Error on test 4: %v", val)
	}
	if val := pt.Get("2001:249::"); val != nil {
		t.Errorf("Error on test 5: %v", val)
	}

	if val := pt.GetKey("2001:250::/38"); val != "2001:250::/38" {
		t.Errorf("Error on test 1: %v", val)
	}

	if !pt.HasKey("2001:250::/38") {
		t.Errorf("Error on test 7")
	}
	if pt.HasKey("2001:250:1::200/119") {
		t.Errorf("Error on test 8")
	}

	pt.Insert("2001:250:1::200/119", "testing456")

	if keys := pt.Keys(); len(keys) != 2 || keys[0] != "2001:250::/38" || keys[1] != "2001:250:1::200/119" {
		t.Errorf("Error on test 9: %v", keys)
	}

	if values := pt.Values(); len(values) != 2 || values[0] != "testing123" || values[1] != "testing456" {
		t.Errorf("Error on test 9: %v", values)
	}

	children := pt.Children("2001:250::/38")
	if children["2001:250::/38"] != "testing123" {
		t.Errorf("Error on test 9: %v", children)
	}
	if children["2001:250:1::200/119"] != "testing456" {
		t.Errorf("Error on test 9: %v", children)
	}

	expectedMap := map[string]interface{}{
		"2001:250::/38":       "testing123",
		"2001:250:1::200/119": "testing456",
	}
	mapOutput := pt.ToMap()
	if len(mapOutput) != len(expectedMap) {
		t.Errorf("Error on test 12: %v", mapOutput)
	}
	if mapOutput["2001:250::/38"] != expectedMap["2001:250::/38"] {
		t.Errorf("Error on test 12: %v", mapOutput)
	}
	if mapOutput["2001:250:1::200/119"] != expectedMap["2001:250:1::200/119"] {
		t.Errorf("Error on test 12: %v", mapOutput)
	}
}

func BenchmarkInsertIPv4(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv4CIDR())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Insert(cidrs[i], "test")
	}
}

func BenchmarkSetIPv4(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv4CIDR())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Set(cidrs[i], "test")
	}
}

func BenchmarkAddIPv4(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv4CIDR())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Add(cidrs[i], "test")
	}
}

func BenchmarkGetIPv4(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv4CIDR())
		pt.Insert(randomIPv4CIDR(), "test")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Get(cidrs[i])
	}
}

func BenchmarkHasKeyIPv4(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv4CIDR())
		pt.Insert(randomIPv4CIDR(), "test")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.HasKey(cidrs[i])
	}
}

func BenchmarkInsertIPv6(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv6CIDR())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Insert(cidrs[i], "test")
	}
}

func BenchmarkSetIPv6(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv6CIDR())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Set(cidrs[i], "test")
	}
}

func BenchmarkAddIPv6(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv6CIDR())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Add(cidrs[i], "test")
	}
}

func BenchmarkGetIPv6(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv6CIDR())
		pt.Insert(randomIPv6CIDR(), "test")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.Get(cidrs[i])
	}
}

func BenchmarkHasKeyIPv6(b *testing.B) {
	pt := NewPyTricia()
	cidrs := []string{}
	for i := 0; i < b.N; i++ {
		cidrs = append(cidrs, randomIPv6CIDR())
		pt.Insert(randomIPv6CIDR(), "test")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pt.HasKey(cidrs[i])
	}
}
