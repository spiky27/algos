package johnson

import (
    "dbg"
    "../../spf/djikstra"
    "../../spf/bellmanford"
    "../../parse/toadjlist"
)

func FindAllSp(gIn, gRevIn, vListIn interface{}) map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost {
    dbg.Dprint(3, -1, "Find all SP by Johnson")
    g := gIn.(toadjlist.Graph)
    gRev := gRevIn.(toadjlist.Graph)
    vList := vListIn.([]toadjlist.Node)
    numVer := len(vList)

    //Add the dummy src node to gRev as Bellman cares only of gRev
    const DUMMYSRC toadjlist.Node = 0x0EFFFFFF
    for _, n := range vList {
        if gRev[n] == nil {
            gRev[n] = make(toadjlist.Adj)
        }
        gRev[n][DUMMYSRC] = 0
    }
    //Add dummy src to vList
    vList = append(vList, DUMMYSRC)

    sp := bellmanford.FindSp(DUMMYSRC, g, gRev, vList)
    if sp == nil {
        return nil
    }

    delete(sp, DUMMYSRC)
    dbg.Dprint(4, -1, "Bellman result", sp)

    //Remove dummy src from vList, g and gRev
    vList = vList[:len(vList)-1]
    delete(g, DUMMYSRC)
    for _, n := range vList {
        delete(gRev[n], DUMMYSRC)
    }

    //Adjust costs
    for src, adj := range g {
        for dst, _ := range adj {
            g[src][dst] += (sp[src] - sp[dst])
        }
    }

    dbg.Dprint(4, -1, "Adjusted cost map", g)

    allSp := make(map[toadjlist.Node]map[toadjlist.Node]toadjlist.Cost, numVer)
    //Call djikstra
    for _, src := range vList {
        allSp[src] = djikstra.FindSp(src, g, gRev, vList)
        //Adjust sp costs
        for dst, _ := range allSp[src] {
            allSp[src][dst] += sp[dst] - sp[src]
        }
    }

    return allSp
}
