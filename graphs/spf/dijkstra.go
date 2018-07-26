package main

import (
    "readFileToDs"
    "os"
    "dbg"
    "bufio"
    "strings"
    "strconv"
    "container/heap"
    "sort"
)

type vertex int
type cost int

type edge struct {
    node vertex
    weight cost
}

func scanComma(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF == true {
        if len(data) == 0 {
            return 0, nil, nil
        }
        return len(data), data, nil
    }
    for i,v := range data {
        if v == ',' || v == ' ' || v == '\t' {
            advance = i+1
            token = data[:i]
            err = nil
            return
        }
    }
    return
}

func main () {
    if len(os.Args) < 4 {
        dbg.ErrOut("Insufficient input")
    }
    dbg.SetLevel(os.Args[3])
    if len(os.Args) == 5 {
        dbg.SetSession(os.Args[4])
    }

    // adjacency list
    graph := make(map[vertex][]edge, 0)

    lscanner := readFiletoDs.ReadFiletoScanner(os.Args[1])
    lscanner.Split(bufio.ScanLines)

    edgeCnt := 0
    for lscanner.Scan() {
        dbg.Dprint(4, -1, "Input", lscanner.Text())
        wscanner := bufio.NewScanner(strings.NewReader(lscanner.Text()))
        wscanner.Split(scanComma)

        var v, p vertex
        for wscanner.Scan() {
            dbg.Dprint(5, -1, "Input", wscanner.Text())
            i, _ := strconv.Atoi(wscanner.Text())
            if v == 0 {
                // first number is the vertex
                v = vertex(i)
                graph[v] = make([]edge, 0)
            } else if p == 0 {
                // subsequently, first number is the node
                p = vertex(i)
            } else {
                // and second number is its corresponding edge weight
                graph[v] = append(graph[v], edge{p, cost(i)})
                p = 0
                edgeCnt ++
            }
        }
    }

    dbg.Dprint(1, -1, "Input graph len:", len(graph))
    dbg.Dprint(2, -1, "Input edge cnt:", edgeCnt)
    dbg.Dprint(4, -1, "Input graph:", graph)

    root, _ := strconv.Atoi(os.Args[2])
    sp := djikstra(graph, vertex(root))
    if len(os.Args) == 5 {
        end, _ := strconv.Atoi(os.Args[4])
        dbg.Dprint(0, end, "SPF:", sp[vertex(end)])
    } else {
        dbg.Dprint(0, -1, "SPF:", sp)
    }
}

type heapCandList []edge

func (h *heapCandList) Len() int {
    return len(*h)
}

func (clist *heapCandList) Less(i, j int) bool {
    h := *clist
    return h[i].weight < h[j].weight
}

func (clist *heapCandList) Swap(i, j int) {
    h := *clist
    h[i].node, h[j].node = h[j].node, h[i].node
    h[i].weight, h[j].weight = h[j].weight, h[i].weight
}

func (h *heapCandList) Push(x interface{}) {
    *h = append(*h, x.(edge))
}

func (h *heapCandList) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func djikstra(graph map[vertex][]edge, root vertex) (sp map[vertex]cost) {
    sp = make(map[vertex]cost, 0)

    sp[root] = 0
    clist := make(heapCandList, 0, len(graph[root]))
    clist.addToClist(graph, root, sp)
    dbg.Dprint(2, int(root), "addToClist candidate count:", len(clist))

    for len(clist) != 0 {
        m1 := heap.Pop(&clist)
        min := m1.(edge)
        /*min := findRemMin(&clist)
        if nIdx := sort.Search(len(clist), func (i int) bool { return clist[i].node >= min.node }); nIdx < len(clist) && clist[nIdx].node == min.node {
            dbg.Dprint(1, int(min.node), "Node", min.node, "found again at", nIdx, clist[nIdx])
            dbg.ErrOut("")
        }*/
        // Candidate already calculated the spf. skip it 
        if sp[min.node] != 0 {
            continue
        }
        sp[min.node] = min.weight
        clist.addToClist(graph, min.node, sp)
        dbg.Dprint(2, int(min.node), "addToClist candidate count:", len(clist))
    }
    return
}

func (clist *heapCandList) Search(d vertex) int {
    for i,v := range *clist {
        if v.node == d {
            return i
        }
    }
    return len(*clist)
}

func (clist *heapCandList) addToClist1(graph map[vertex][]edge, finishedNode vertex, sp map[vertex]cost) {
    dbg.Dprint(2, int(finishedNode), "addToClist FinishedNode:", finishedNode, "*****")
    for _,e := range graph[finishedNode] {
        dbg.Dprint(2, int(finishedNode), "addToClist Processing edge", e)
        // check if the neighbor is already added to spf, then move to next
        if _,err := sp[e.node]; err == true {
            dbg.Dprint(2, int(finishedNode), "addToClist already computed:", e.node, sp[e.node])
            continue
        }

        // calculate cost of edge.node from finishedNode
        newNodeCost := sp[finishedNode] + e.weight
        // check if the neighbor is already in clist, then update the cost
        //if nIdx := sort.Search(len(*clist), func (i int) bool { return (*clist)[i].node >= e.node }); nIdx < len(*clist) && (*clist)[nIdx].node == e.node {
        if nIdx := clist.Search(e.node); nIdx != len(*clist) {
            dbg.Dprint(4, int(e.node), "addToClist clist:", *clist)
            dbg.Dprint(2, int(e.node), "addToClist found", e.node, "at location", nIdx)
            dbg.Dprint(2, int(e.node), "addToClist found", (*clist)[nIdx])
            if newNodeCost < (*clist)[nIdx].weight {
                dbg.Dprint(1, int(e.node), "addToClist updating cost for", e.node, "to", newNodeCost, "from", (*clist)[nIdx].weight, "prev", finishedNode)
                (*clist)[nIdx].weight = newNodeCost
                dbg.Dprint(4, int(e.node), "addToClist value[update]:", *clist)
                sort.Sort(clist)
            }
        } else {
            // add the node to candidate list with the appropriate cost
            *clist = append(*clist, edge{e.node, newNodeCost})
            sort.Sort(clist)
            dbg.Dprint(1, int(e.node), "addToClist added node", e.node, "with cost", newNodeCost, "coming from", finishedNode)
            dbg.Dprint(4, int(e.node), "addToClist value[create]:", *clist)
        }
    }
}

func findRemMin(clist *heapCandList) edge {
    e := (*clist)[0]

    dbg.Dprint(4, int(e.node), "Clist while popping:", *clist)
    *clist = (*clist)[1:]
    sort.Sort(clist)
    return e
}

func (clist *heapCandList) addToClist(graph map[vertex][]edge, finishedNode vertex, sp map[vertex]cost) {
    dbg.Dprint(2, int(finishedNode), "addToClist FinishedNode:", finishedNode)
    for _,e := range graph[finishedNode] {
        // check if the neighbor is already added to spf, then move to next
        if _,err := sp[e.node]; err == true {
            dbg.Dprint(2, int(finishedNode), "addToClist already computed:", e.node, sp[e.node])
            continue
        }

        // calculate cost of edge.node from finishedNode
        newNodeCost := sp[finishedNode] + e.weight
        // check if the neighbor is already in clist, then update the cost
        // add the node to candidate list with the appropriate cost
        heap.Push(clist, edge{e.node, newNodeCost})
        dbg.Dprint(1, int(e.node), "addToClist added node", e.node, "with cost", newNodeCost, "coming from", finishedNode)
        dbg.Dprint(4, int(e.node), "addToClist value[create]:", *clist)
    }
}
