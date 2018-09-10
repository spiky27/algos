package floydwarshall

import (
//  "fmt"
    "../../parse/toadjlist"
)

/*
type Node toadjlist.Node
type Cost toadjlist.Cost
type Graph map[Node]map[Node]Cost
*/
const MAXCOST toadjlist.Cost = 0x0FFFFFFF

func FindAllSp(gIn, gRevIn, vListIn interface{}) map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost {
    g := gIn.(toadjlist.Graph)
    vList := vListIn.([]toadjlist.Node)
    numVer := len(vList)

    var sp [2]map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost
    sp[0] = make(map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost, numVer)
    sp[1] = make(map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost, numVer)

    //Initialize the first indexes
    for _, src := range vList {
        sp[0][src] = make(map[toadjlist.Node]toadjlist.Cost, numVer)
        sp[1][src] = make(map[toadjlist.Node]toadjlist.Cost, numVer)
        for _, dst := range vList {
            if src == dst {
                sp[0][src][dst] = 0
            } else if val, ok := g[src][dst]; ok {
                sp[0][src][dst] = val
            } else {
                sp[0][src][dst] = MAXCOST
            }
        }
    }

    for iter := 1;iter < numVer; iter ++ {
        for _, src := range vList {
            for _, dst := range vList {
                sp[iter%2][src][dst] = toadjlist.Cost(min(sp[(iter+1)%2][src][dst], (sp[(iter+1)%2][src][vList[iter]] + sp[(iter+1)%2][vList[iter]][dst])))
            }
        }
    }

    //Check for negative cycle
    for idx := toadjlist.Node(0); int(idx) < numVer; idx ++ {
        if sp[(numVer-1)%2][idx][idx] < toadjlist.Cost(0) {
            return nil
        }
    }

    return sp[(numVer+1)%2]
}

func min(aIn, bIn interface{}) int {
    a, b := aIn.(toadjlist.Cost), bIn.(toadjlist.Cost)
    if a < b {
        return int(a)
    }
    return int(b)
}
