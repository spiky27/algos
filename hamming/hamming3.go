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
type gData map[uint]*vertex

type vertex struct {
    rank int
    key, leader uint
//    leader *vertex
//    dList []uint
}

func (g gData) String() string {
    var buf bytes.Buffer

    buf.WriteString("\n")
    for idx, v := range g {
        fmt.Fprintf(&buf, "idx: %d\n", idx)
        fmt.Fprintf(&buf, "  { Key: %b(%d), Leader: %b(%d) }\n", v.key, v.key, v.leader, v.leader)
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

    findCntClusters(input)

    cnt := findLeaders(input)
    dbg.Cprint("Cluster Count:", cnt, dupCnt)
}

func findLeaders(graph gData) (cnt int) {
    for _, v := range graph {
        if v.key == v.leader {
            cnt ++
        }
    }
    return
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

    numObj, numBits, bitIdx, numObjSan := 0, 0, 1, 0
    var entry uint

    for scanner.Scan() {
        i, err := strconv.Atoi(scanner.Text())
        if err != nil {
            continue
        }
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
                numObjSan ++
                //update the graph
                if _, ok := graph[entry]; !ok {
                    graph[entry] = &vertex{0, entry, entry}
                } else {
                    dupCnt ++
                }

                //prepare of the next number
                bitIdx = 1
                entry = 0
            }
        }
    }

    dbg.Dprint(1, -1, "Number of entries given vs actual", numObj, numObjSan, len(graph))
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

func setPathLeader(graph gData, newLdr uint, subE *vertex) {
    // Scan all the children of the leader of this cluster
    if subE.key != subE.leader {
        setPathLeader (graph, newLdr, graph[subE.leader])
    }
    subE.leader = newLdr

    return
}

func findLeader(entry *vertex, graph gData) uint {
    for entry.key != entry.leader {
        entry = graph[entry.leader]
    }
    return entry.leader
}

func updateLeaders(graph gData, entry0, entry1 *vertex) {
    l0, l1 := findLeader(entry0, graph), findLeader(entry1, graph)
    dbg.Dprint(2, -1, "Leaders", entry0.key, ":", l0, ";", entry1.key, ":", l1)
    if l0 == l1 {
        return
    }

    var newLdr uint
    var subE *vertex
    // Need to merge the clusters
    // Choose the lesser of the evils
    if graph[l0].rank >= graph[l1].rank {
        newLdr, subE = l0, entry1
        if graph[l0].rank == graph[l1].rank {
            graph[l0].rank ++
        }
    } else if graph[l0].rank < graph[l1].rank {
        newLdr, subE = l1, entry0
    }
    // Update the leaders of the smaller cluster
    setPathLeader(graph, newLdr, subE)
    dbg.Dprint(3, -1, "New Graph", graph)
}

func findCntClusters(graph gData) {
    dbg.Dprint(3, -1, "Input to findCnt:", graph)
    entriesProc := 0

    for key0, entry0 := range graph {
        entriesProc ++
//        dbg.Dprint(1, -1, "Under scan:", gKey, "numEntries:", len(graph[gKey]))
        for i := uint(0); i < 24; i ++ {
            var eVar1, eVar2 uint
            if (key0 & (1 << i)) == 0 {
                eVar1 = key0 | (uint(1) << i)
            } else {
                eVar1 = key0 & ^(uint(1) << i)
            }
            if _, ok := graph[eVar1]; ok {
                updateLeaders (graph, entry0, graph[eVar1])
            }

            for j := i+1; j < 24; j ++ {
                if (eVar1 & (1 << j)) == 0 {
                    eVar2 = eVar1 | (1 << j)
                } else {
                    eVar2 = eVar1 & ^(1 << j)
                }
                if _, ok := graph[eVar2]; ok {
                    updateLeaders (graph, entry0, graph[eVar2])
                }
            }
        }
        dbg.Dprint(1, -1, "Entries Processed", entriesProc)
    }
    dbg.Dprint(3, -1, "Output of findCnt:", graph)
}
