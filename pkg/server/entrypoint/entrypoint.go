package entrypoint

import (
	"math"
)

func some_super_cool_hash(s string) int {
	return int(math.Mod(float64(len(s)), float64(num_chains)))
}

type Node struct {
	Host string
	Port int
}

type NodeStat struct {
	N    Node
	Stat int
}

type Chain struct {
	nodes []Node
}

const num_chains = 2
const nodes_per_chain = 3

// global variable with all the chains.  entrypoint server wil randomly pick a chain head for uploads
// Pick a random node for download within a chain?
var chains = [num_chains]Chain{Chain{[]Node{Node{Host: "chain1_head", Port: 8080},
	Node{Host: "chain1_replica1", Port: 8081},
	Node{Host: "chain1_tail", Port: 8082}}},

	Chain{[]Node{Node{Host: "chain2_head", Port: 8083},
		Node{Host: "chain2_replica1", Port: 8084},
		Node{Host: "chain2_tail", Port: 8085}}}}
