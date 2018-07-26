package main

import (
    "dbg"
    "os"
    "strconv"
    "bufio"
    "sort"
)

type vertex int
type cost   int

type nodeProp struct {
    rank int
    leader vertex
}

type edge struct {
    v1, v2 vertex
    w cost
}

type eList []edge

func main () {
    if len(os.Args) < 4 {
        dbg.ErrOut("Bad Input")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    numClusters := 4

    // Read input file
    edgeDb, nodes := parseFile(os.Args[1])

    // sort edgelist
    sort.Sort(edgeDb)

    // run kruskal
    lList, eIdx := kruskal(edgeDb, nodes, numClusters)

    dbg.Cprint("List of leaders:", lList, "\nMax distance:", edgeDb[eIdx].w)
}

func custScan (data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF == true {
        if len(data) == 0 {
            return 0, nil, nil
        }
        return len(data), data, nil
    }

    for i, v := range data {
        if v == ' ' || v == '\n' || v == '\t' {
            return i + 1, data[:i], nil
        }
    }
    return
}

func parseFile(s string) (edgeDb eList, nodes []nodeProp) {
    f, err := os.Open(s)
    dbg.AbortIfErr(err)

    scanner := bufio.NewScanner(bufio.NewReader(f))
    scanner.Split(custScan)

    numNode := 0
    var v1, v2 vertex

    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        if numNode == 0 {
            numNode = i

            //Init the map for node to nodeProp
            nodes = make([]nodeProp, numNode + 1)
            for idx1 := 1; idx1 <= numNode; idx1++ {
                idx := vertex(idx1)
                nodes[idx].rank = 0
                nodes[idx].leader = vertex(idx)
            }

            //Init the edge list
            numEdge := (numNode * (numNode - 1))/2 + 1
            edgeDb = make(eList, 0, numEdge)
        } else if v1 == 0 {
            v1 = vertex(i)
        } else if v2 == 0 {
            v2 = vertex(i)
            if v1 > v2 {
                v1, v2 = v2, v1
            }
        } else {
            edgeDb = append(edgeDb, edge{v1, v2, cost(i)})
            dbg.Dprint(3, -1, "Input edge:", edgeDb[len(edgeDb)-1])
            v1, v2 = 0, 0
        }
    }

    return
}

func (edges eList) Len () int {
    return len(edges)
}

func (edges eList) Less (i, j int) bool {
    return (edges[i].w <= edges[j].w)
}

func (edges eList) Swap (i, j int) {
    edges[i], edges[j] = edges[j], edges[i]
}

func kruskal(edges eList, nodes []nodeProp, numCl int) (map[vertex][]vertex, int) {
    dbg.Dprint(2, -1, "Kruskal input:", edges, numCl)
    lList := make(map[vertex][]vertex, len(nodes) + 1)

    var v vertex
    for v = 1; int(v) <= len(nodes); v++ {
        lList[v] = make([]vertex, 0, 1)
        lList[v] = append(lList[v], v)
    }

    for idx, edge := range edges {
        l1, l2 := findLeader(nodes, edge.v1), findLeader(nodes, edge.v2)
        if l1 == l2 {
            continue
        }

        //merge v1, v2 into a single cluster
        mergeClusters(l1, l2, nodes, lList)

        //return if only numCl clusters are left
        if len(lList) <= numCl {
            dbg.Dprint(1, -1, "Kruskal returning:", lList, idx)
            return lList, idx
        }
    }

    dbg.Dprint(1, -1, "Kruskal bad scan returning:", lList, len(edges))
    return lList, len(edges)
}

func findLeader(nodes []nodeProp, curr vertex) vertex {
    for nodes[curr].leader != curr {
        curr = nodes[curr].leader
    }
    return curr
}

func mergeClusters(l1, l2 vertex, nodes []nodeProp, lList map[vertex][]vertex) {
    dbg.Dprint(2, -1, "Merge input:", l1, lList[l1], ";", l2, lList[l2])
    if nodes[l1].rank >= nodes[l2].rank {
        nodes[l2].leader = l1
        if nodes[l1].rank == nodes[l2].rank {
            nodes[l1].rank ++
        }
        lList[l1] = append(lList[l1], lList[l2]...)
        lList[l2] = nil
        delete(lList, l2)
    } else {
        nodes[l1].leader = l2
        lList[l2] = append(lList[l2], lList[l1]...)
        lList[l1] = nil
        delete(lList, l1)
    }
    dbg.Dprint(2, -1, "Merge return:", lList)
}
