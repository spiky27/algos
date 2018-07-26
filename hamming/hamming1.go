package main

import (
    "dbg"
    "os"
    "strconv"
    "bufio"
    "fmt"
//    "math"
    "bytes"
)

var leaderCnt uint
type gData map[uint8]map[uint]*vertex

type vertex struct {
    key uint
    leader *vertex
//    dList []uint
}

func (g gData) String() string {
    var buf bytes.Buffer

    buf.WriteString("\n")
    for idx, entry := range g {
        fmt.Fprintf(&buf, "idx: %d\n", idx)
        for _, v := range entry {
            fmt.Fprintf(&buf, "  { Key: %b(%d), Leader: %b(%d) }\n", v.key, v.key, v.leader.key, v.leader.key)
        }
/*        for dk, dv := range v.dList {
            fmt.Fprintf(&buf, "%d:%b(%d) ", dk, dv.key, dv.key)
        }
        fmt.Fprintf(&buf, "} }\n")*/
    }
    return buf.String()
}

func main () {
    defer dbg.TraceTime("hamming")();
    if len(os.Args) < 4 {
        dbg.ErrOut("Bad Input")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    input := parseFile(os.Args[1])

    cnt := findCntClusters(input)

    dbg.Cprint("Cluster Count:", cnt, dupCnt)
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
var dupCnt int
func parseFile (s string) (graph gData) {
    defer dbg.TraceTime("parsing") ();
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
        } else if numBits == 0 {
            numBits = i
            graph = make(gData, numBits)
        } else {
            //form the number
            entry |= uint(i << uint(numBits-bitIdx))
            bitIdx ++

            if bitIdx > numBits {
                d := findDefDistance (entry)
                //update the graph
                if graph[d] == nil {
                    graph[d] = make(map[uint]*vertex, 1)
                }
                if _, ok := graph[d][entry]; !ok {
                    graph[d][entry] = &vertex{entry, nil}
                    //assign leader to self
                    graph[d][entry].leader = graph[d][entry]
                } else {
                    dupCnt ++
                }

                //prepare of the next number
                bitIdx = 1
                entry = 0
            }
        }
    }

    sanCnt := 0
    for j := uint8(0); j <=24; j ++ {
        sanCnt += len(graph[j])
    }

    dbg.Dprint(1, -1, "Number of entries given vs actual", numObj, sanCnt)
    return
}

func xor (i, j uint) uint {
    return ((i|j)&^(i&j))
}

func findDefDistance (j uint) uint8 {
    i := uint(0x00FFFFFF)
    cnt := uint8(0)//, uint8(32)
    x := xor(i, j)
    for x != 0 {
        cnt ++
/*        if idx == 32 {
            idx = uint8(math.Log2(float64(x&^(x-1))))
        }*/
        x &= (x-1)
    }
    return cnt
}

func checkIfMaxDist2 (i, j uint) bool {
    cnt := 0
    x := xor(i, j)
    for x != 0 {
        if cnt == 2 {
            return false
        }
        cnt ++
        x &= (x-1)
    }
    return true
}

type cluster map[uint]map[uint]uint8

func (cl cluster)addClusterElem(key, child uint) {
    if cl[key] == nil {
        cl[key] = make(map[uint]uint8, 1)
    }
    cl[key][child] = findDefDistance(child)
}

func setPathLeader(graph gData, oldLdr, newLdr *vertex, cl cluster) {
    // Scan all the children of the leader of this cluster
    for chkey, chDefDist := range cl[oldLdr.key] {
        if chDefDist != 0 {
            // Update the leader of all the children
            graph[chDefDist][chkey].leader = newLdr
            // add this child to the new cluster
            cl.addClusterElem(newLdr.key, chkey)
        }
    }
    oldLdr.leader = newLdr
    cl.addClusterElem(newLdr.key, oldLdr.key)
    // delete the old cluster
    cl[oldLdr.key] = nil
    delete(cl, oldLdr.key)

    return
}

func (cl cluster)findMemCnt() (ldrCnt,memCnt int) {
    for _, leader := range cl {
        if len(leader) > 0 {
            ldrCnt ++
            memCnt += len(leader)
        }
    }

    return ldrCnt, memCnt
}

func findLeader(entry *vertex) *vertex {
    if entry.key == entry.leader.key {
        return entry
    }
    return findLeader(entry.leader)
}

func updateLeaders(graph gData, entry0, entry1 *vertex, cl cluster) {
    var newLdr, oldLdr *vertex

    l0, l1 := findLeader(entry0), findLeader(entry1)
    dbg.Dprint(2, -1, "Leaders", entry0.key, ":", l0.key, ";", entry1.key, ":", l1.key)
    if l0.key == l1.key {
        return
    }

    // Need to merge the clusters
    if true == checkIfMaxDist2(entry0.key, entry1.key) {
        // Choose the lesser of the evils
        if len(cl[l0.key]) >= len(cl[l1.key]) {
            newLdr, oldLdr = l0, l1
        } else {
            newLdr, oldLdr = l1, l0
        }

        // Update the leaders of the smaller cluster
        setPathLeader(graph, oldLdr, newLdr, cl)
        dbg.Dprint(3, -1, "New Graph", graph)
    }
}

func findCntClusters(graph gData) (cnt int) {
    cl := make(cluster, 1)

    dbg.Dprint(3, -1, "Input to findCnt:", graph)

    for gKey := uint8(1); gKey <= uint8(24); gKey++ {
        dbg.Dprint(1, -1, "Under scan:", gKey, "numEntries:", len(graph[gKey]))
        lCnt, mCnt := cl.findMemCnt ()
        dbg.Dprint(1, -1, "leaders", lCnt, "members", mCnt)
        for _, entry0 := range graph[gKey] {
            if len(graph[gKey-1]) < 0 {
                for _, entry1 := range graph[gKey-1] {
                    updateLeaders (graph, entry0, entry1, cl)
                }
                if gKey >= 2 {
                    for _, entry2 := range graph[gKey-2] {
                        updateLeaders (graph, entry0, entry2, cl)
                    }
                }
            } else {
                for i := uint(0); i < 24; i ++ {
                    var eVar1, eVar2 uint
                    if (entry0.key & (1 << i)) == 0 {
                        eVar1 = entry0.key | uint(1 << i)
                    } else {
                        eVar1 = entry0.key & ^uint(1 << i)
                    }
                    if _, ok := graph[gKey-1][eVar1]; ok {
                        updateLeaders (graph, entry0, graph[gKey-1][eVar1], cl)
                    }

                    for j := i+1; j < 24; j ++ {
                        if (eVar1 & (1 << j)) == 0 {
                            eVar2 = eVar1 | (1 << j)
                        } else {
                            eVar2 = eVar1 & ^(1 << j)
                        }
                        if _, ok := graph[gKey-2][eVar2]; ok {
                            updateLeaders (graph, entry0, graph[gKey-2][eVar2], cl)
                        }
                    }
                }
            }
        }
    }
    dbg.Dprint(3, -1, "Output of findCnt:", graph)
    dbg.Dprint(3, -1, "Cluster", cl)

    return len(cl)
}
