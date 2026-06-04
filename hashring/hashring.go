package hashring

import (
	"crypto/sha1"
	"sort"
)

// Node represents a cache node
type Node struct {
	Address string
}

// HashRing stores nodes, a sorted list of hashes, and a map from hash to node
type HashRing struct {
	nodes       []Node
	sortedKeys  []int
	hashToNode  map[int]Node
}

// New creates a hash ring with given nodes
func New(addresses []string) *HashRing {
	hr := &HashRing{
		nodes:      []Node{},
		hashToNode: make(map[int]Node),
	}

	for _, addr := range addresses {
		hr.AddNode(addr)
	}
	return hr
}

// AddNode inserts a node into the ring
func (hr *HashRing) AddNode(address string) {
	node := Node{Address: address}
	hash := hr.hashKey(address)
	hr.nodes = append(hr.nodes, node)
	hr.hashToNode[hash] = node
	hr.sortedKeys = append(hr.sortedKeys, hash)
	sort.Ints(hr.sortedKeys)
}

// GetNode returns the node responsible for a given key
func (hr *HashRing) GetNode(key string) Node {
	if len(hr.sortedKeys) == 0 {
		return Node{}
	}

	hash := hr.hashKey(key)

	// search for closest hash >= key hash
	for _, k := range hr.sortedKeys {
		if hash <= k {
			return hr.hashToNode[k]
		}
	}
	
	return hr.hashToNode[hr.sortedKeys[0]]
}

func (hr *HashRing) hashKey(key string) int {
	h := sha1.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	var sum int

	for _, b := range bs {
		sum += int(b)
	}

	return sum
}