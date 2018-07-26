package main

import (
    "dbg"
    "os"
    "strconv"
    "bufio"
    "sort"
)

func main() {
    defer dbg.TraceTime("huffman")()

    if len(os.Args) < 4 {
        dbg.ErrOut("Insufficient Input")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    input := parseInput(os.Args[1])
    maxLen, minLen := findEncDepth(input)

    dbg.Cprint("Max, Min bits for encoding is", maxLen, minLen)
}

type node struct {
    symbol uint
    weight int
}

type symList []node

func (sl *symList) Len() int {
    symbols := *sl
    return len(symbols)
}

func (sl *symList) Less(i, j int) bool {
    symbols := *sl
    return (symbols[i].weight <= symbols[j].weight)
}

func (sl *symList) Swap(i, j int) {
    symbols := *sl
    symbols[i], symbols[j] = symbols[j], symbols[i]
}

func (sl *symList) push(n node) {
    *sl = append(*sl, n)
}

func (sl *symList) pop() {
    symbols := *sl
    if symbols == nil || len(symbols) == 0 {
        return
    }

    if len(symbols) == 1 {
        symbols = symList{}
    } else {
        symbols = symbols[1:]
    }
    *sl = symbols
}

func (sl *symList) peekW(idx int) (w int) {
    symbols := *sl
    if len(symbols) < idx {
        return 0xEFFFFFFF
    }

    w = symbols[idx-1].weight
    return
}

func parseInput(s string) symList {
    f, err := os.Open(s)
    if err != nil {
        dbg.AbortIfErr(err)
    }

    scanner := bufio.NewScanner(bufio.NewReader(f))

    symCnt, sym := 0, uint(0)
    var inputSym symList
    for scanner.Scan() {
        i, _ := strconv.Atoi(scanner.Text())
        if symCnt == 0 {
            symCnt = i
            inputSym = make(symList, 0, symCnt)
        } else {
            inputSym = append(inputSym, node{sym, i})
            sym ++
        }
    }

    dbg.Dprint(1, -1, "read inputs", symCnt, len(inputSym))
    return inputSym
}

func findNewSym(n1, n2 node, symType uint) uint {
    switch symType {
        case 1:
            //both are from symbols
            if n1.weight < n2.weight {
                return (n2.symbol << 16) | (n1.symbol)
            }
            return (n1.symbol << 16) | (n2.symbol)

        case 2:
            //one of them is from symbols
            break

        case 3:
            //none of them is from symbols
            break

        default:
            break
    }
        return 0
}

func findEncDepth(q0 symList) (int, int) {
    sort.Sort(&q0)
    dbg.Dprint(3, -1, "Input Sorted", q0)
    minSym := q0[0].symbol
    maxSym := q0[len(q0)-1].symbol
    dbg.Dprint(1, -1, "Min,Max sym is", minSym, maxSym)
    q1, q2 := symList{}, symList{}
    minDepth, maxDepth := 0, 0
    newSym := uint(0)

    for len(q0) + len(q1) + len(q2) > 1 {
        tables := findLeast(q0, q1, q2)
        switch tables {
            case 0x1:
                newWeight := q0[0].weight + q0[1].weight
                if q0[0].symbol == minSym {
                    newSym = minSym
                }
                if (q0[0].symbol == maxSym) || (q0[1].symbol == maxSym) {
                    newSym |= (maxSym << 16)
                }
//                newSym = (q0[0].symbol) | (q0[1].symbol << 16)
//                newSym = findNewSym(q0[0], q0[1], minSym)
                dbg.Dprint(3, int(newSym), "Merging q0", q0[0], "q0", q0[1])
                q0.pop()
                q0.pop()
                dbg.Dprint(4, int(newSym), "q0 elements left", len(q0))
                q1.push(node{newSym, newWeight})
                dbg.Dprint(3, int(newSym), "Added node q1", q1)
                dbg.Dprint(4, int(newSym), "q1 elements left", len(q1))
                break

            case 0x2:
                newWeight := q1[0].weight + q1[1].weight
                if ((q1[0].symbol & 0xFFFF) == minSym) || ((q1[1].symbol & 0xFFFF) == minSym) {
                    newSym |= (minSym)
                }
                if ((q1[0].symbol & 0xFFFF0000) == (maxSym << 16)) || ((q1[1].symbol & 0xFFFF0000) == (maxSym << 16)) {
                    newSym |= (maxSym << 16)
                }
//                newSym = findNewSym(q1[0], q1[1], minSym)
//                newSym = (q1[0].symbol & 0xFFFF) | (q1[1].symbol & 0xFFFF0000)
                dbg.Dprint(3, int(newSym), "Merging q1", q1[0], "q1", q1[1])
                q1.pop()
                q1.pop()
                dbg.Dprint(4, int(newSym), "q1 elements left", len(q1))
                q2.push(node{newSym, newWeight})
                dbg.Dprint(3, int(newSym), "Added node q2", q2)
                dbg.Dprint(4, int(newSym), "q2 elements left", len(q2))
                break

            case 0x3:
                newWeight := q0[0].weight + q1[0].weight
                if (q1[0].symbol & 0xFFFF) == minSym {
                    newSym = minSym
                }
                if (q0[0].symbol == maxSym) || (q1[0].symbol & 0xFFFF0000) == (maxSym << 16) {
                    newSym |= (maxSym << 16)
                }
//                newSym = findNewSym(q0[0], q1[0], minSym)
//                newSym = (q1[0].symbol & 0xFFFF) | (q0[0].symbol << 16)
                dbg.Dprint(3, int(newSym), "Merging q0", q0[0], "q1", q1[0])
                q0.pop()
                dbg.Dprint(4, int(newSym), "q0 elements left", len(q0))
                q1.pop()
                dbg.Dprint(4, int(newSym), "q1 elements left", len(q1))
                q2.push(node{newSym, newWeight})
                dbg.Dprint(3, int(newSym), "Added node q2", q2)
                dbg.Dprint(4, int(newSym), "q2 elements left", len(q2))
                break

            case 0x4:
                newWeight := q2[0].weight + q2[1].weight
                if ((q2[0].symbol & 0xFFFF) == minSym) || ((q2[1].symbol & 0xFFFF) == minSym) {
                    newSym = minSym
                }
                if ((q2[0].symbol & 0xFFFF0000) == (maxSym << 16)) || ((q2[1].symbol & 0xFFFF0000) == (maxSym << 16)) {
                    newSym |= (maxSym << 16)
                }
//                newSym = findNewSym(q2[0], q2[1], minSym)
//                newSym = (q2[0].symbol & 0xFFFF) | (q2[1].symbol & 0xFFFF0000)
                dbg.Dprint(3, int(newSym), "Merging q2", q2[0], "q2", q2[1])
                q2.pop()
                q2.pop()
                q2.push(node{newSym, newWeight})
                dbg.Dprint(3, int(newSym), "Added node q2", q2)
                dbg.Dprint(4, int(newSym), "q2 elements left", len(q2))
                break

            case 0x5:
                newWeight := q0[0].weight + q2[0].weight
                if (q2[0].symbol & 0xFFFF) == (minSym) {
                    newSym = (minSym)
                }
                if (q0[0].symbol == maxSym) || ((q2[0].symbol & 0xFFFF0000) == (maxSym << 16)) {
                    newSym |= (maxSym << 16)
                }
//                newSym = findNewSym(q0[0], q2[0], minSym)
//                newSym = (q0[0].symbol) | (q2[0].symbol & 0xFFFF0000)
                dbg.Dprint(3, int(newSym), "Merging q0", q0[0], "q2", q2[0])
                q0.pop()
                dbg.Dprint(4, int(newSym), "q0 elements left", len(q0))
                q2.pop()
                q2.push(node{newSym, newWeight})
                dbg.Dprint(3, int(newSym), "Added node q2", q2)
                dbg.Dprint(4, int(newSym), "q2 elements left", len(q2))
                break

            case 0x6:
                newWeight := q1[0].weight + q2[0].weight
                if ((q1[0].symbol & 0xFFFF) == minSym) || ((q2[0].symbol & 0xFFFF) == minSym) {
                    newSym = (minSym)
                }
                if ((q1[0].symbol & 0xFFFF0000) == (maxSym << 16)) || ((q2[0].symbol & 0xFFFF0000) == (maxSym << 16)) {
                    newSym |= (maxSym << 16)
                }
//                newSym = findNewSym(q1[0], q2[0], minSym)
//                if q1[0].weight < q2[0].weight {
//                    newSym = (q1[0].symbol & 0xFFFF) | (q2[0].symbol & 0xFFFF0000)
//                } else {
//                    newSym = (q2[0].symbol & 0xFFFF) | (q1[0].symbol & 0xFFFF0000)
//                }
                dbg.Dprint(3, int(newSym), "Merging q1", q1[0], "q2", q2[0])
                q1.pop()
                dbg.Dprint(4, int(newSym), "q1 elements left", len(q1))
                q2.pop()
                q2.push(node{newSym, newWeight})
                dbg.Dprint(3, int(newSym), "Added node q2", q2)
                dbg.Dprint(4, int(newSym), "q2 elements left", len(q2))
                break

            default:
                break
        }

        if (newSym & 0xFFFF) == minSym {
            minDepth ++
            dbg.Dprint(2, int(newSym), "Updating minDepth", minDepth)
        }
        if (newSym & 0xFFFF0000) == (maxSym << 16) {
            maxDepth ++
            dbg.Dprint(2, int(newSym), "Updating maxDepth", maxDepth)
        }
        newSym = 0
    }

    return minDepth, maxDepth
}

func findLeast(list, q1, q2 symList) (mask uint) {
    l1, l2 := list.peekW(1), list.peekW(2)
    q11, q12 := q1.peekW(1), q1.peekW(2)
    q21, q22 := q2.peekW(1), q2.peekW(2)

    if l1 < q11 {
        if q11 < q21 {
            //least=l1, q11
            if l2 < q11 {
                //least=l1, l2
                return 0x1
            } else if q12 < l1 {
                //least=q11,q12
                return 0x2
            } else {
                return 0x3
            }
        } else {
            //least=l1, q21
            if l2 < q21 {
                //least=l1, l2
                return 0x1
            } else if q22 < l1 {
                //least=q21,q22
                return 0x4
            } else {
                return 0x5
            }
        }
    } else {
        if l1 < q21 {
            //least=l1, q11
            if l2 < q11 {
                //least=l1, l2
                return 0x1
            } else if q12 < l1 {
                //least=q11,q12
                return 0x2
            } else {
                return 0x3
            }
        } else {
            //least=q11, q21
            if q22 < q11 {
                //least=q21, q22
                return 0x4
            } else if q12 < q21 {
                //least=q11,q12
                return 0x2
            } else {
                return 0x6
            }
        }
    }
}
