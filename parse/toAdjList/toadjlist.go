package toadjlist

import (
    //"fmt"
    "os"
    "bufio"
    "strconv"
    "dbg"
    "log"
    //"errors"
)

type Cost int
type Node int
type Adj map[Node]Cost
type Graph map[Node]map[Node]Cost


func ParseInput(s string) (Graph, Graph, []Node, int, int) {
    dbg.Dprint(2, -1, "Parse Input from", s)
    var g, gRev Graph
    var vList []Node
    var vSeen map[Node]bool

    f, err := os.Open(s)
    if err != nil {
        log.Fatal(err)
    }

    scanner := bufio.NewScanner(bufio.NewReader(f))
    scanner.Split(scanToken)

    numVertex, numEdge := 0, 0
    var src, dst Node
    var cost Cost

    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        if err != nil {
            log.Fatal(err)
        }

        if numVertex == 0 {
            numVertex = i
            g = make(Graph, numVertex)
            gRev = make(Graph, numVertex)
            vList = make([]Node, 0, numVertex)
            vSeen = make(map[Node]bool, numVertex)
        } else if numEdge == 0 {
            numEdge = i
        } else if src == Node(0) {
            src = Node(i)
            if _, ok := g[src]; !ok {
                g[src] = make(Adj)
            }
            if vSeen[src] == false {
                vSeen[src] = true
            }
        } else if dst == Node(0) {
            dst = Node(i)
            if _, ok := gRev[dst]; !ok {
               gRev[dst] = make(Adj)
            }
            if vSeen[dst] == false {
                vSeen[dst] = true
            }
        } else {
            cost = Cost(i)
            g[src][dst] = cost
            gRev[dst][src] = cost
            src = 0
            dst = 0
        }
    }

    for n, seen := range vSeen {
        if seen == true {
            vList = append(vList, n)
        }
    }
    return g, gRev, vList, numVertex, numEdge
}

func scanToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF == true {
        if len(data) == 0 {
            return 0, nil, nil
        }
        return len(data), data, nil
    }
    for i, v := range data {
        if v == ',' || v == '\n' || v == '\t' || v == ' ' {
            advance = i + 1
            token = data[:i]
            err = nil
            return
        }
    }
    return
}
