package graphs

import (
	"fmt"
	"math"
)

type node struct {
	mark    bool
	rank    uint64
	value   float64
	content interface{}
	parent  *node
	child   *node
	l       *node
	r       *node
}

type FibanaccyHeap struct {
	min   *node
	nodes map[interface{}]*node
}

func New() *FibanaccyHeap {
	return &FibanaccyHeap{
		nodes: make(map[interface{}]*node),
	}
}

func (fh *FibanaccyHeap) Add(v float64, c interface{}) {
	newNode := &node{
		mark:    false,
		rank:    1,
		value:   v,
		content: c,
	}

	fh.nodes[c] = newNode

	// Add the node to the roots
	if fh.min == nil {
		fh.min = newNode
		newNode.l = newNode
		newNode.r = newNode
		fh.min = newNode
	} else {
		fh.min.l.r = newNode
		newNode.l = fh.min.l
		newNode.r = fh.min
		fh.min.l = newNode

		if fh.min.value > v {
			fh.min = newNode
		}
	}
}

func (fh *FibanaccyHeap) moveToRoot(n *node) {
	if n.parent != nil {
		if n.r != n {
			n.parent.rank -= n.rank
			n.parent.child = n.r
		} else {
			n.parent.child = nil
			n.parent.rank = 1
		}
	}

	n.r.l = n.l
	n.l.r = n.r

	n.r = fh.min.r
	n.l = fh.min
	fh.min.r.l = n
	fh.min.r = n
	n.parent = nil

	if n.value < fh.min.value {
		fh.min = n
	}
}

func (fh *FibanaccyHeap) moveMarkedParentsToRoot(n *node) {
	if n.mark {
		fh.moveMarkedParentsToRoot(n.parent)
		n.mark = false
		fh.moveToRoot(n)
	}
}

func (fh *FibanaccyHeap) DecreaseKey(c interface{}, v float64) {
	node := fh.nodes[c]
	node.value = v

	if node.parent != nil {
		if node.value > node.parent.value {
			return
		}

		if node.parent.parent != nil {
			// Only the not root nodes can be marked
			if node.parent.mark {
				fh.moveMarkedParentsToRoot(node.parent)
			} else {
				node.parent.mark = true
				node.parent.child = nil
			}
		}
	}

	fh.moveToRoot(node)
	node.mark = false
}

func (fh *FibanaccyHeap) Min() (value float64, content interface{}) {
	if fh.min == nil {
		return
	}

	value = fh.min.value
	content = fh.min.content

	delete(fh.nodes, content)

	if fh.min.rank > 1 {
		fh.min.child.parent = nil
		if fh.min.l != fh.min {
			fh.min.l.r = fh.min.child
			fh.min.r.l = fh.min.child.l
			fh.min.child.l.r = fh.min.r
			fh.min.child.l = fh.min.l
		} else {
			// We are on the last node of the roots
			fh.min = fh.min.child
		}
	} else {
		if fh.min.l == fh.min {
			// This is the last node
			fh.min = nil

			return
		}

		fh.min.l.r = fh.min.r
		fh.min.r.l = fh.min.l
	}

	init := true
	min := math.Inf(1)
	node := fh.min.r
	origin := node
	for node != origin || init {
		if min > node.value {
			min = node.value
			fh.min = node
		}
		node = node.r
		init = false
	}

	fh.rebalance()

	return
}

func (fh *FibanaccyHeap) mergeNodes(parent *node, child *node) {
	// Remove the child from the current linked list
	child.l.r = child.r
	child.r.l = child.l

	if parent.child != nil {
		child.r = parent.child.r
		parent.child.r = child
		child.r.l = child
		child.l = parent.child
	} else {
		child.r = child
		child.l = child
	}

	parent.child = child
	child.parent = parent
	parent.rank += child.rank
}

func (fh *FibanaccyHeap) print(from *node, deep int) {
	initial := true
	aux := from
	for aux != from || initial {
		fmt.Println(deep, "Value:", aux.value, "Rank:", aux.rank, "Content:", aux.content, "Marked:", aux.mark)
		if aux.child != nil {
			fh.print(aux.child, deep+1)
		}
		aux = aux.r
		initial = false
	}
}

func (fh *FibanaccyHeap) rebalance() {
	balanced := false
balance:
	for !balanced {
		usedRanks := make(map[uint64]*node)

		node := fh.min
		first := true
		for !balanced || first {
			if _, ok := usedRanks[node.rank]; ok {
				if node.value < usedRanks[node.rank].value {
					fh.mergeNodes(node, usedRanks[node.rank])
				} else {
					fh.mergeNodes(usedRanks[node.rank], node)
				}

				continue balance
			}

			usedRanks[node.rank] = node
			node = node.r
			balanced = node == fh.min
			first = false
		}
	}
}
