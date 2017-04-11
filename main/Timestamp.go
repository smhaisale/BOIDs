package main

type VectorTimestamp struct {
    processTime map[string]int
}

func (v VectorTimestamp) increment(name string) {
    v.processTime[name] = v.processTime[name] + 1;
}

func (v VectorTimestamp) update(timestamp VectorTimestamp) {
    for key, value := range timestamp.processTime {
        v.processTime[key] = int(max(v.processTime[key], value))
    }
}

func (v VectorTimestamp) compare(timestamp VectorTimestamp) int {
    allGreater, allSmaller := true, true

    for key, value := range v.processTime {
        if timestamp.processTime[key] > value { allGreater = false }
        if timestamp.processTime[key] < value { allSmaller = false }
    }

    switch {
        case allGreater: return 1
        case allSmaller: return -1
        default: return 0
    }
}

func max(list ...int) int {
    max := list[0]
    for _, i := range list[1:] {
        if i > max {
            max = i
        }
    }
    return max
}

func min(list ...int) int {
    min := list[0]
    for _, i := range list[1:] {
        if i < min {
            min = i
        }
    }
    return min
}
