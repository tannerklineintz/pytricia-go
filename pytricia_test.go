package pytricia

import "testing"

func TestPytriciaIPv4(t *testing.T) {
	t.Parallel()

	pt := NewPyTricia()

	if !pt.IsRoot() {
		t.Errorf("Error on test 0")
	}

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

	if !pt.HasKey("8.8.8.0/24") {
		t.Errorf("Error on test 7")
	}
	if pt.HasKey("8.8.8.10/31") {
		t.Errorf("Error on test 8")
	}

	pt.Insert("8.8.8.10/31", "testing456")

	node := pt.GetNode("8.8.8.240")
	if val2 := node.CIDR().String(); val2 != "8.8.8.0/24" {
		t.Errorf("Error on test 9: %v", val2)
	}

	children := node.Children()
	if val2 := children[0].CIDR().String(); val2 != "8.8.8.0/24" {
		t.Errorf("Error on test 10: %v", val2)
	}
	if val2 := children[1].CIDR().String(); val2 != "8.8.8.10/31" {
		t.Errorf("Error on test 11: %v", val2)
	}

	expectedMap := map[string]interface{}{
		"8.8.8.0/24":  "testing123",
		"8.8.8.10/31": "testing456",
	}
	mapOutput := pt.ToMap()
	if len(mapOutput) != len(expectedMap) {
		t.Errorf("Error on test 12: %v", mapOutput)
	}
	if mapOutput["8.8.8.0/24"] != expectedMap["8.8.8.0/24"] {
		t.Errorf("Error on test 12: %v", mapOutput)
	}
	if mapOutput["8.8.8.10/31"] != expectedMap["8.8.8.10/31"] {
		t.Errorf("Error on test 12: %v", mapOutput)
	}

	node2 := pt.GetNode("8.8.8.10/31")
	parent := node2.Parent()
	if parent.CIDR().String() != "8.8.8.0/24" {
		t.Errorf("Error on test 13: %v", parent)
	}
}
