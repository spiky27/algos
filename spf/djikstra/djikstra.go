package djikstra

import (
    "dbg"
    "../../parse/toadjlist"
    "container/heap"
)

type vertex struct {
    node toadjlist.Node
    weight toadjlist.Cost
}
type candHeap []vertex

func (ch *candHeap) Len () int {
    return len(*ch)
}

func (ch *candHeap) Less (i, j int) bool {
    return (*ch)[i].weight < (*ch)[j].weight
}

func (ch *candHeap) Swap (i, j int) {
    h := *ch
    h[i].weight, h[j].weight = h[j].weight, h[i].weight
    h[i].node, h[j].node = h[j].node, h[i].node
}

func (ch *candHeap) Push(x interface{}) {
    *ch = append(*ch, x.(vertex))
}

func (ch *candHeap) Pop() interface{} {
    old := *ch
    n := len(old)
    x := old[n-1]
    *ch = old[:n-1]
    return x
}

func FindSp(srcIn, gIn, gRevIn, vListIn interface{}) map[toadjlist.Node]toadjlist.Cost {
    src := srcIn.(toadjlist.Node)
    g := gIn.(toadjlist.Graph)
    vList := vListIn.([]toadjlist.Node)
    dbg.Dprint(4, -1, "Find Sp using Djikstra", src)

    sp := make(map[toadjlist.Node]toadjlist.Cost, len(vList))
    candidate := make(candHeap, 0)

    heap.Push(&candidate, vertex{src, 0})
    dbg.Dprint(4, -1, "candidate List:", candidate)

    for candidate.Len() != 0 {
        curr := heap.Pop(&candidate).(vertex)
    dbg.Dprint(4, -1, "candidate List:", candidate)
        if _, ok := sp[curr.node]; !ok {
            sp[curr.node] = curr.weight
            dbg.Dprint(4, -1, "SP updated:", sp)
        }

        for nbr, cost := range g[curr.node] {
            if _, ok := sp[nbr]; !ok {
                nbrCost := sp[curr.node] + cost
                heap.Push(&candidate, vertex{nbr, nbrCost})
    dbg.Dprint(4, -1, "candidate List:", candidate)
            }
        }
    }

    return sp
}
