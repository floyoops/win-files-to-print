package print

import (
	"fmt"
	"sort"
)

type PrinterList struct {
	List      []string
	IndexMap  map[int]string
	nextIndex int
}

func (pl *PrinterList) Add(value string) {
	pl.List = append(pl.List, value)
	pl.IndexMap[pl.nextIndex] = value
	pl.nextIndex++
}

func (pl *PrinterList) GetByIndex(index int) (string, bool) {
	value, found := pl.IndexMap[index]
	return value, found
}

func (pl *PrinterList) Render() {
	var keys []int
	for key := range pl.IndexMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("%d - %s\n", key, pl.IndexMap[key])
	}
}
