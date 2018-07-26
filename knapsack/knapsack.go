package main

import (
    "fmt"
    "dbg"
    "os"
    "strconv"
    "bufio"
)

func main () {
    if len(os.Args) < 4 {
        panic("Insufficient Inputs")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    input, numItems, maxCap := parseFile(os.Args[1])
    if input == nil {
        panic("Parse Failed")
    }

    fmt.Println("Answer to knapsack problem is", findKnapsack(input, numItems, maxCap))
}

type itemInfo struct {
    value, weight int
}

type itemList []itemInfo

func parseFile(s string) (*itemList, int, int) {
    f, err := os.Open(s)
    if err != nil {
        panic(err)
    }

    scanner := bufio.NewScanner(bufio.NewReader(f))
    scanner.Split(bufio.ScanWords)

    var input itemList
    numItems, maxCap, idx := 0, 0, 0
    itemVal, itemWgt := 0, 0
    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        dbg.Dprint(5, -1, "Input", i)
        if maxCap == 0 {
            maxCap = i
        } else if numItems == 0 {
            numItems = i
            input = make(itemList, numItems + 1)
            dbg.Dprint(3, -1, "Inputs", numItems, maxCap)
        } else if itemVal == 0 {
            itemVal = i
        } else {
            idx ++
            itemWgt = i
            dbg.Dprint(4, -1, "Item", idx, itemVal, itemWgt)
            input[idx] = itemInfo{itemVal, itemWgt}
            itemVal = 0
        }
    }

    return &input, numItems, maxCap
}

func findMax(i, j int) int {
    if i > j {
        return i
    }
    return j
}

func findKnapsack(in *itemList, numItems, maxCap int) (int) {
    dbg.Dprint(1, -1, "Knapsack input", len(*in), numItems, maxCap)
    input := *in
    tmpKsack := make([][]int, numItems + 1)
    for itemIdx := range tmpKsack {
        tmpKsack[itemIdx] = make([]int, maxCap + 1)
    }
    tmpKsack[numItems] = make([]int, maxCap + 1)

    for itemIdx := 1; itemIdx <= numItems; itemIdx ++ {
        for capIdx := 1; capIdx <= maxCap; capIdx ++ {
            var eliValue, addValue int
            eliValue = tmpKsack[itemIdx-1][capIdx]
            if remCap := capIdx - input[itemIdx].weight; remCap >= 0 {
                addValue = tmpKsack[itemIdx-1][remCap] + input[itemIdx].value
            }
            tmpKsack[itemIdx][capIdx] = findMax(eliValue, addValue)
        }
        dbg.Dprint(2, -1, "Item", itemIdx, ":", tmpKsack[itemIdx])
    }
    return tmpKsack[numItems][maxCap]
}
