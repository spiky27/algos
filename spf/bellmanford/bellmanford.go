package bellmanford

import(
//    "fmt"
    "../../parse/toadjlist"
)
type Node  toadjlist.Node
type Cost  toadjlist.Cost
type Graph toadjlist.Graph

const MAXCOST toadjlist.Cost = 0x0FFFFFFF

func FindSp (src toadjlist.Node, gIn, gInRev interface{}, vList []toadjlist.Node) map[toadjlist.Node]toadjlist.Cost {
    gRev := gInRev.(toadjlist.Graph)
    var sp [2]map[toadjlist.Node]toadjlist.Cost
    sp[0] = make(map[toadjlist.Node]toadjlist.Cost, len(vList))
    sp[1] = make(map[toadjlist.Node]toadjlist.Cost, len(vList))

    /* Init all path costs to MAX */
    for _, node := range vList {
        if node == src {
            sp[0][node] = 0
        } else {
            sp[0][node] = MAXCOST
        }
    }

    currBkt, prevBkt := 0, 1
    var iter int
    var diffFound bool
    for ; iter <= len(vList); iter ++ {
        diffFound = false
        for _, node := range vList {
            newCost := MAXCOST
            //Find min over all edges
            for pHop, edgeCost := range gRev[node] {
                if newCost > sp[prevBkt][pHop] + edgeCost {
                    newCost = sp[prevBkt][pHop] + edgeCost
                }
            }
            if newCost > sp[prevBkt][node] {
                newCost = sp[prevBkt][node]
            }
            sp[currBkt][node] = newCost
            if diffFound == false && newCost != sp[prevBkt][node] {
                diffFound = true
            }
        }
        //Break if none of paths are updated
        if diffFound == false {
            break
        }
        currBkt, prevBkt = prevBkt, currBkt
    }

    if iter == len(vList)+1 && diffFound == true {
        return nil
    }
    return sp[prevBkt]
}
