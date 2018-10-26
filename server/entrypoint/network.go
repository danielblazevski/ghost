package entrypoint

const num_chains = 2
const nodes_per_chain = 3

type Node struct {
	Host string
	Port int
}

type Chain struct {
	nodes []Node
}

// global variable with all the chains.  entrypoint server wil randomly pick a chain head for uploads
// Pick a random node for download within a chain?
var chains = [num_chains]Chain{Chain{[]Node{Node{Host: "chain1_head", Port: 8080},
	Node{Host: "chain1_replica1", Port: 8081},
	Node{Host: "chain1_tail", Port: 8082}}},

	Chain{[]Node{Node{Host: "chain2_head", Port: 8083},
		Node{Host: "chain2_replica1", Port: 8084},
		Node{Host: "chain2_tail", Port: 8085}}}}
