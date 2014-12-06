package models

type (
	Node struct {
		Next *Node
	}

	Loop struct {
		Node
	}

	LoopRange struct {
	}

	Cond struct {
	}

	Call struct {
	}
)
