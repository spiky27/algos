package main

import (
    "dbg"
    "os"
    "strconv"
    "bufio"
    "container/heap"
)

type vertex int
type cost int
type vHeap []edge
type edge struct {
    nbr vertex
    w cost
}

func main () {
    //Check for sufficient inputs : input file, debug level, debug session
    if len(os.Args) < 4 {
        dbg.ErrOut("Insufficient inputs")
    }

    //Set debug level and session
    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    //read file into data structure
    graph := parseFile(os.Args[1])

    dbg.Dprint(2, -1, "Input Graph:", graph)
    //Calculate mst
    totCost := calcMst(1, graph)

    //Print total cost
    dbg.Cprint("Total Cost of MST:", totCost)
}

func custScan(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF == true {
        if len(data) == 0 {
            return 0, nil, nil
        }
        return len(data), data, nil
    }

    for i,v := range data {
        if v == ' ' || v == '\n' || v == '\t' {
            advance = i+1
            token = data[:i]
            err = nil
            return
        }
    }
    return
}

// Read file to data structure
func parseFile (file string) map[vertex][]edge {
    f, err := os.Open(file)
    dbg.AbortIfErr(err)

    scanner := bufio.NewScanner(f)
    scanner.Split(custScan)

    vcnt, ecnt := 0, 0
    var v1, v2 vertex
    var c cost
    var graph map[vertex][]edge
    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        if vcnt == 0 {
            vcnt = i
            if vcnt > 0 {
                graph = make(map[vertex][]edge, vcnt)
            } else {
                dbg.ErrOut("No vertices")
            }
        } else if ecnt == 0 {
            ecnt = i
        } else if v1 == 0 {
            v1 = vertex(i)
            if _,err := graph[v1]; err == false {
                graph[v1] = make([]edge, 0)
            }
        } else if v2 == 0 {
            v2 = vertex(i)
            if _,err := graph[v2]; err == false {
                graph[v2] = make([]edge, 0)
            }
        } else {
            c = cost(i)
            dbg.Dprint(4, -1, "Adding edge:", v1, v2, c)
            graph[v1] = append(graph[v1], edge{v2, c})
            dbg.Dprint(4, -1, "New Link:", v1, ":", graph[v1])
            graph[v2] = append(graph[v2], edge{v1, c})
            dbg.Dprint(4, -1, "New Link:", v2, ":", graph[v2])
            v1, v2, c = 0, 0, 0
        }
    }
    return graph
}

func (vh *vHeap) Len () int {
    return len(*vh)
}

func (vh *vHeap) Less (i,j int) bool {
    return (*vh)[i].w <= (*vh)[j].w
}

func (vh *vHeap) Swap (i, j int) {
    (*vh)[i], (*vh)[j] = (*vh)[j], (*vh)[i]
}

func (vh *vHeap) Push (x interface{}) {
    *vh = append(*vh, x.(edge))
}

func (vh *vHeap) Pop () interface{} {
    old := *vh
    n := len(old)
    x := old[n-1]
    *vh = old[:n-1]
    return x
}

// Calculate mst
func calcMst (root vertex, graph map[vertex][]edge) (totCost cost) {
    //Sanity test
    if graph == nil {
        dbg.ErrOut("Failed miserably")
    }

    mst := make(map[vertex]bool, len(graph))
    vh := make(vHeap, 0)

    //Add first node to heap setting cost to 0
    heap.Push(&vh, edge{root, 0})
    dbg.Dprint(3, -1, "calcMst-Added root:", root, vh)

    //Loop over the heap
    for len(vh) != 0 {
        n := heap.Pop(&vh)
        node := n.(edge)
        dbg.Dprint(2, -1, "calcMst-Min node:", node)
        if mst[node.nbr] == true {
            continue
        }

        //Add node to mst and update total cost
        mst[node.nbr] = true
        totCost += node.w

        //Loop over all nbrs and add them to the heap if nbr not in mst already
        for _, e := range graph[node.nbr] {
            if mst[e.nbr] == false {
                heap.Push(&vh, e)
                dbg.Dprint(3, -1, "calcMst-Added node:", e, vh)
            }
        }
    }

    return
}
