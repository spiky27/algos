package main

import (
    "fmt"
    "os"
    "./floydwarshall"
    "./johnson"
    "../parse/toadjlist"
    "dbg"
)

type Cost toadjlist.Cost
type Node toadjlist.Node

func main () {
    defer dbg.TraceTime("apsp") ()

    if len(os.Args) < 5 {
        dbg.ErrOut("Insufficient Input")
    }

    dbg.SetLevel(os.Args[3])
    dbg.SetSession(os.Args[4])

    input, inputRev, vList, numVertex, numEdge := toadjlist.ParseInput(os.Args[1])
    dbg.Dprint(1, -1, "Read input:", numVertex, len(vList), numEdge, len(input), len(inputRev))
    dbg.Dprint(4, -1, "vList:", vList)
    dbg.Dprint(4, -1, "graph:", input)

    var sp map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost
    if os.Args[2] == "f" {
        sp = floydwarshall.FindAllSp(input, inputRev, vList)
    } else {
        sp = johnson.FindAllSp(input, inputRev, vList)
    }

    if sp == nil {
        dbg.Cprint("Negative cycle found")
        return
    }
    dbg.Dprint(4, -1, "All Sp", sp)
    //Find min cost
    min := toadjlist.Cost(0x0FFFFFFF)
    for _, n := range vList {
        node := toadjlist.Node(n)
        for _, c := range sp[node] {
            if min > c {
                min = c
            }
        }
    }

    fmt.Println("Min Cost:", min)
}
