package structs

import "sort"

type Pair[T any] struct {
	Key   T
	Value int
}

type PairList[T any] []Pair[T]

func (p PairList[T]) Len() int {
	return len(p)
}

func (p PairList[T]) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

func (p PairList[T]) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PairList[T]) Sort() {
	sort.Sort(p)
}

func (p PairList[T]) Reverse() {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
}
