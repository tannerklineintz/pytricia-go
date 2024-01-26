package pytricia

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
	ipType   int
	children [2]*PyTricia
	parent   *PyTricia
	value    interface{}
}
