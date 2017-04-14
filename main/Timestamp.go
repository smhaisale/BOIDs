package main

type VectorTime struct {
    processTime map[string]int
}

func (v VectorTime) increment(name string) {
    v.processTime[name]  = v.processTime[name] + 1;
}

func (v VectorTime) update(timestamp VectorTime) {
    for key, value := range timestamp.processTime {
        v.processTime[key] = int(max(v.processTime[key], value))
    }
}

func (v VectorTime) compare(timestamp VectorTime) int {
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

//TODO: Put some in
var sampleTimestamp = VectorTime { map[string]int {}}

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
