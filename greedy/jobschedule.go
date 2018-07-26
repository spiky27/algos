package main

import (
    "dbg"
    "readFiletoDs"
    "os"
    "strconv"
    "bufio"
    "sort"
    "strings"
)

type jobNode struct {
    weight, length int
}

type custJList struct {
    jList []jobNode
    less func(x, y jobNode) bool
}

func (job custJList) Len() int {
    return len(job.jList)
}

func (job custJList) Less(i, j int) bool {
    return job.less(job.jList[j], job.jList[i])
}

func (job custJList) Swap(i, j int) {
    job.jList[i], job.jList[j] = job.jList[j], job.jList[i]
}

func lessDiff(x, y jobNode) bool {
    diff := (x.weight-x.length) - (y.weight-y.length)

    if diff == 0 {
        return (x.weight <= y.weight)
    }
    return (diff < 0)
}

func lessRatio(x, y jobNode) bool {
    diff := (float64(x.weight)/float64(x.length)) - (float64(y.weight)/float64(y.length))

    return (diff <= 0)
}


func main () {
    if len(os.Args) < 4 {
        dbg.ErrOut("Insufficient inputs")
    }

    dbg.SetLevel(os.Args[2])
    dbg.SetSession(os.Args[3])

    lscanner := readFiletoDs.ReadFiletoScanner(os.Args[1])
    lscanner.Split(bufio.ScanLines)

    cntJobs := 0
    var jobList []jobNode
    for lscanner.Scan() {
        wscanner := bufio.NewScanner(strings.NewReader(lscanner.Text()))
        wscanner.Split(bufio.ScanWords)

        w := 0
        for wscanner.Scan() {
            i, _ := strconv.Atoi(wscanner.Text())
    dbg.Dprint(4, -1, "Input st:", i)
            if cntJobs == 0 {
                cntJobs = i
                jobList = make([]jobNode, 0, cntJobs)
                dbg.Dprint(3, -1, "Job List Init:", jobList, "num jobs:", cntJobs)
            } else if w == 0 {
                w = i
            } else {
                jobList = append(jobList, jobNode{w, i})
                dbg.Dprint(3, -1, "Job List append:", jobList, "job added:", w, i)
            }
        }
    }
    if cntJobs == 0 {
        dbg.ErrOut("No jobs")
    }

    dbg.Dprint(2, -1, "Job List:", jobList)

    var wsum, currlen int64
    sort.Sort(custJList{jobList, lessDiff})
    dbg.Dprint(2, -1, "Job List after sorting:", jobList)

    for _, v := range jobList {
        dbg.Dprint(2, -1, "Job:", v.weight, v.length)
        currlen += int64(v.length)
        wsum += int64(v.weight) * int64(currlen)
    }
    dbg.Cprint("Weighted Sum:", wsum)

    wsum, currlen = 0, 0
    sort.Sort(custJList{jobList, lessRatio})
    dbg.Dprint(2, -1, "Job List after sorting:", jobList)

    for _, v := range jobList {
        dbg.Dprint(2, -1, "Job:", v.weight, v.length)
        currlen += int64(v.length)
        wsum += int64(v.weight) * int64(currlen)
    }
    dbg.Cprint("Weighted Sum:", wsum)
}
