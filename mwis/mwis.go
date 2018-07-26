package main

import (
    "fmt"
    "os"
    "strconv"
    "bufio"
    "dbg"
)

type node struct {
    weight, maxw int
    mwisPath bool
}

func main () {
    if len(os.Args) < 4 {
        dbg.ErrOut("Insufficient parameters")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    f, err := os.Open(os.Args[1])
    dbg.AbortIfErr(err)

    scanner := bufio.NewScanner(bufio.NewReader(f))

    cntVertex := 0
    var vList []node
    var vIdx int
    for scanner.Scan () {
        i, _ := strconv.Atoi(scanner.Text())
        if cntVertex == 0 {
            cntVertex = i
            vList = make([]node, cntVertex)
        } else {
            vList[vIdx].weight = i
            if vIdx == 0 {
                vList[vIdx].maxw = vList[vIdx].weight
            } else if vIdx == 1 {
                vList[vIdx].maxw = max(vList[vIdx-1].maxw, vList[vIdx].weight)
            } else {
                vList[vIdx].maxw = max(vList[vIdx-1].maxw, vList[vIdx].weight + vList[vIdx-2].maxw)
            }
            vIdx ++
        }
    }

    for vIdx = len(vList)-1; vIdx > 0; vIdx -- {
        if vList[vIdx].maxw > vList[vIdx-1].maxw {
            vList[vIdx].mwisPath = true
            vIdx --
        }
    }
    if vList[1].mwisPath == false {
        vList[0].mwisPath = true
    }

    dbg.Dprint(4, -1, "Complete mwis", vList)

    bitshow := map[int]uint {
        1:   0x10000000,
        2:   0x01000000,
        3:   0x00100000,
        4:   0x00010000,
        17:  0x00001000,
        117: 0x00000100,
        517: 0x00000010,
        997: 0x00000001,
    }

    var output uint
    for k, v := range bitshow {
        if vList[k-1].mwisPath == true {
            output |= v
        }
    }
    fmt.Printf("Final output %x", output)
}

func max(i, j int) int {
    if i > j {
        return i
    }
    return j
}
