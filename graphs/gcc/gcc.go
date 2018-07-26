package main

import (
    "dbg"
    "readFiletoDs"
    "os"
    "bufio"
    "strings"
    "strconv"
)

type vertex uint
type graph map[vertex][]vertex
type leaderList struct {
    leader vertex
    memCnt int
}

var leader vertex
var timestamp int64
var sccLeader map[vertex]vertex

func findMaxN(sl map[vertex][]vertex, n int) []leaderList {
    totalCount := 0
    maxCount := make([]leaderList, 0, n)

    for i := 0;i < n; i++ {
        totalCount += len(sl)
        maxIdx, max := findMax(sl)
        delete(sl, maxIdx)
        maxCount = append(maxCount, leaderList{maxIdx, max})
    }
    dbg.Dprint(1, "Total vertices in max:", totalCount)

    return maxCount
}

func findMax(sl map[vertex][]vertex) (vertex, int) {
    max := 0
    var maxk vertex
    for k,v := range sl {
        if max < len(v) {
            max = len(v)
            maxk = k
        }
    }

    return maxk, max
}

func main () {
    if len(os.Args) < 3 {
        dbg.ErrOut("Insufficient input")
    }

    dbg.Set(os.Args[2])

    g, grev, vl := readtoDs(os.Args[1])
    /* Initialize the vlist to number of vertices */
    scc := make(graph, 0)
    sccLeader = make(map[vertex]vertex, len(vl))
    timestamp = int64(len(vl))

    dbg.Dprint(3, "graph:", g)
    dbg.Dprint(3, "reversed graph:", grev)
    dbg.Dprint(1, "num Vertex:", len(g), len(grev), len(vl))

    /* Do reverse DFS */
    pl := dfs_loop(grev, vl, nil)

    /* Do forward DFS */
    dfs_loop(g, pl, scc)
    dbg.Dprint(1, "Number of SCC=", len(scc), len(pl))
    dbg.Dprint(2, "SCCs:\n", scc)
    dbg.Dprint(2, "SCCLeaders:\n", sccLeader)

    ret := findMaxN(scc, 5)
    dbg.Cprint("Biggest SCC:", ret)
}

/* Reads an input file to required Data Structures *
* Input:
*   s : Input file
* Output:
*   g    : adjacency list as a map from v to its nbrs
*   grev : reverse adjacency list as a map from v to its prev nbrs
*   vl   : list of vertices as a slice
*/
func readtoDs(s string) (g graph, grev graph, vl []vertex) {
    g = make(graph, 0)
    grev = make(graph, 0)
    vl = make([]vertex, 0)
    vadded := make(map[vertex]bool, 0)

    ls := readFiletoDs.ReadFiletoScanner(s)
    ls.Split(bufio.ScanLines)
    for ls.Scan() {
        ws := bufio.NewScanner(strings.NewReader(ls.Text()))
        ws.Split(bufio.ScanWords)
        /* Will return 2 ints, forming an edge */
        var v,w vertex
        for ws.Scan() {
            i, _ := strconv.Atoi(ws.Text())
            if v == 0 {
                v = vertex(i)
            } else {
                w = vertex(i)
            }
        }

        if vadded[v] == false {
            vl = append(vl, v)
            vadded[v] = true
        }
        if vadded[w] == false {
            vl = append(vl, w)
            vadded[w] = true
        }
        /* Edge is read, add to the graph */
        if g[v] == nil {
            g[v] = make([]vertex, 0)
        }
        g[v] = append(g[v], w)

        /* Preparing reverse graph */
        if grev[w] == nil {
            grev[w] = make([]vertex, 0)
        }
        grev[w] = append(grev[w], v)
    }
    vadded = nil

    return
}

func dfs_loop(g graph, vl []vertex, scc graph) ([]vertex) {
    seen := make(map[vertex]bool, len(vl))
    pl := make([]vertex, len(vl))

    for _, v  := range vl {
        /* Node not yet seen, new dfs starts from here */
        if seen[v] == false {
            if scc != nil {
                leader = v
                if scc[leader] == nil {
                    scc[leader] = make([]vertex, 0)
                    if sccLeader[leader] != 0 {
                        dbg.Dprint(1, "Leader", sccLeader[leader], "already assigned for node", leader, ", trying to assign", leader)
                        dbg.ErrOut("")
                    }
                    scc[leader] = append(scc[leader], leader)
                    sccLeader[leader] = leader
                }
            } else {
                dbg.Dprint(2, "pl reverse leader:", v, timestamp)
            }
            dfs(g, seen, v, scc, pl)
        }
    }
    seen = nil

    return pl
}

func dfs(g graph, seen map[vertex]bool, v vertex, scc graph, pl []vertex) {
    seen[v] = true
    if scc != nil {
    } else {
        dbg.Dprint(3, "Timestamp for node", v, "is", timestamp)
        dbg.Dprint(2, "pl:", len(pl), timestamp)
        if timestamp == 0 {
            dbg.ErrOut("timestamp hit 0 too soon")
        }
    }

    for _, w := range g[v] {
        /* Node not yet seen, go deeper here */
        if seen[w] == false {
            if scc != nil {
                if sccLeader[w] != 0 {
                    dbg.Dprint(1, "Leader", sccLeader[w], "already assigned for node", w, ", trying to assign", leader)
                    dbg.ErrOut("")
                }
                scc[leader] = append(scc[leader], w)
                sccLeader[w] = leader
            }
            dfs(g, seen, w, scc, pl)
        }
    }

    timestamp --
    if scc == nil {
        pl[timestamp] = v
    }
}
