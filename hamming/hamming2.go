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
    key, leader uint
//    leader *vertex
//    dList []uint
}
type cluster map[uint]map[uint]bool
var cl cluster

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
            cl = make(cluster, numObj)
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
                    graph[entry] = &vertex{entry, entry}
                    cl[entry] = make(map[uint]bool, 1)
                    cl[entry][entry] = true
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


func (cl cluster)addClusterElem(key, child uint) {
    if cl[key] == nil {
        cl[key] = make(map[uint]bool, 1)
    }
    cl[key][child] = true //findDefDistance(child)
}

func setPathLeader(graph gData, oldLdr, newLdr uint,cl cluster) {
    // Scan all the children of the leader of this cluster
    for chkey, chPresent := range cl[oldLdr] {
        if chPresent != false {
            // Update the leader of all the children
            graph[chkey].leader = newLdr
            // add this child to the new cluster
            cl.addClusterElem(newLdr, chkey)
        }
    }
    // delete the old cluster
    cl[oldLdr] = nil
    delete(cl, oldLdr)

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

func findLeader(entry *vertex) uint {
    return entry.leader
}

func updateLeaders(graph gData, entry0, entry1 *vertex,cl cluster) {
    var newLdr, oldLdr uint

    l0, l1 := findLeader(entry0), findLeader(entry1)
    dbg.Dprint(2, -1, "Leaders", entry0.key, ":", l0, ";", entry1.key, ":", l1)
    if l0 == l1 {
        return
    }

    // Need to merge the clusters
    // Choose the lesser of the evils
    if len(cl[l0]) >= len(cl[l1]) {
        newLdr, oldLdr = l0, l1
    } else {
        newLdr, oldLdr = l1, l0
    }

    // Update the leaders of the smaller cluster
    setPathLeader(graph, oldLdr, newLdr,cl)
    dbg.Dprint(3, -1, "New Graph", graph)
}

func findCntClusters(graph gData) (cnt int) {
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
                updateLeaders (graph, entry0, graph[eVar1], cl)
            }

            for j := i+1; j < 24; j ++ {
                if (eVar1 & (1 << j)) == 0 {
                    eVar2 = eVar1 | (1 << j)
                } else {
                    eVar2 = eVar1 & ^(1 << j)
                }
                if _, ok := graph[eVar2]; ok {
                    updateLeaders (graph, entry0, graph[eVar2], cl)
                }
            }
        }
        lCnt, mCnt := cl.findMemCnt ()
        dbg.Dprint(1, -1, "Entries Processed", entriesProc, "leaders", lCnt, "members", mCnt)
    }
    dbg.Dprint(3, -1, "Output of findCnt:", graph)
    dbg.Dprint(3, -1, "Cluster", cl)

    return len(cl)
}
