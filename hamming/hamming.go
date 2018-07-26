package hamming

import (
    "dbg"
    "os"
    "strconv"
    "bufio"
    "fmt"
    "math"
    "bytes"
)

var leaderCnt uint
type gData map[uint]*vertex

type vertex struct {
    rank int
    key, leader uint
    dList map[int]vertex
}

func (g gData) String() string {
    var buf bytes.Buffer

    buf.WriteString("\n")
    for _, v := range g {
        fmt.Fprintf(&buf, "{ Rank: %d, Key: %b(%d), Leader: %b(%d), nbr:{ ", v.rank, v.key, v.key, v.leader, v.leader)
        for dk, dv := range v.dList {
            fmt.Fprintf(&buf, "%d:%b(%d) ", dk, dv.key, dv.key)
        }
        fmt.Fprintf(&buf, "} }\n")
    }
    return buf.String()
}

func main () {
    if len(os.Args) < 4 {
        dbg.ErrOut("Bad Input")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    input := parseFile(os.Args[1])

    cnt := findCntClusters(input)

    dbg.Cprint("Cluster Count:", cnt)
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

func parseFile (s string) (graph gData) {
    f, err := os.Open(s)
    dbg.AbortIfErr(err)

    scanner := bufio.NewScanner(bufio.NewReader(f))
    scanner.Split(custScan)

    numObj, numBits, bitIdx := 0, 0, 1
    var entry uint

    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        if numObj == 0 {
            numObj = i
            graph = make(gData, numObj)
        } else if numBits == 0 {
            numBits = i
        } else {
            //form the number
            entry |= uint(i << uint(numBits-bitIdx))
            bitIdx ++

            if bitIdx > numBits {
                //update the graph
                graph[entry] = &vertex{0, entry, entry, nil}

                //prepare of the next number
                bitIdx = 1
                entry = 0
            }
        }
    }
    return
}

func xor (i, j uint) uint {
    return ((i|j)&^(i&j))
}

func findDistance (i, j uint) (int, int) {
    cnt, idx := 0, 32
    x := xor(i, j)
    for x != 0 {
        cnt ++
        if idx == 32 {
            idx = int(math.Log2(float64(x&^(x-1))))
        }
        x &= (x-1)
    }
    return cnt, idx+1
}

func findCntClusters(graph gData) (cnt int) {
    cluster := []uint{}
    matchedCluster := false

    dbg.Dprint(3, -1, "Input to findCnt:", graph)

    for key, _ := range graph {
        dbg.Dprint(2, -1, "Under scan:", graph[key])
        for _, ldr := range cluster {
            dbg.Dprint(2, -1, "Leader candidate:", ldr)
            if ret, _ := checkAndAdd (key, ldr, graph); ret == true {
                matchedCluster = true
                break
            }
        }
        if matchedCluster == false {
            cluster = append(cluster, key)
            dbg.Dprint(2, -1, "Cluster list:", cluster)
        } else {
            matchedCluster = false
        }
    }
    return len(cluster)
}

func max(i, j int) int {
    if i > j {
        return i
    }
    return j
}

func checkAndAdd(key, node uint, graph gData) (bool, int) {
    dbg.Dprint(3, -1, "Check:", key, node, graph[node])
    d, idx := findDistance(key, node)
    dbg.Dprint(2, -1, "Distance:", d, idx)
    if d > max(2, graph[node].rank + 1) {
        dbg.Dprint(2, -1, "return high d", d, max(2, graph[node].rank + 1))
        return false, 0
    }
    if d <= 1 {
        if graph[node].dList == nil {
            graph[node].dList = make(map[int]vertex, 0)
        }
        graph[node].dList[idx] = vertex{0, key, node, make(map[int]vertex, 0)}
        if graph[node].rank == 0 {
            graph[node].rank = 1
        }
        dbg.Dprint(3, -1, "Graph updated:", graph)
        return true, graph[node].rank
    } else {
        if _, ok := graph[node].dList[idx]; !ok {
            if d == 2 {
                if graph[node].dList == nil {
                    graph[node].dList = make(map[int]vertex, 0)
                }
                graph[node].dList[idx] = vertex{0, key, node, nil}
                if graph[node].rank < 2 {
                    graph[node].rank = 2
                }
                dbg.Dprint(3, -1, "Graph updated:", graph)
                return true, graph[node].rank
            }
            dbg.Dprint(3, -1, "return empty list", graph[node], idx)
            return false, 0
        }
        if ret, _ := checkAndAdd(key, graph[node].dList[idx].key, graph); ret == false {
            if d == 2 {
                if graph[node].dList == nil {
                    graph[node].dList = make(map[int]vertex, 0)
                }
                graph[node].dList[idx] = vertex{0, key, node, nil}
                if graph[node].rank < 2 {
                    graph[node].rank = 2
                }
                dbg.Dprint(3, -1, "Graph updated:", graph)
                return true, graph[node].rank
            }
        }
    }
    return false, 0
}
