package helpers

import (
	"math/rand"
	"time"
)

// IndexOf returns the index of an item in a slice
func IndexOf[T comparable](slice []T, item T) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

// Contains returns true if a slice contains an item
func Contains[T comparable](slice []T, item T) bool {
	return IndexOf(slice, item) > -1
}

// Remove removes an item from a slice
func Remove[T comparable](slice []T, item T) []T {
	index := IndexOf(slice, item)
	if index > -1 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}

// IsUnique returns true if all items in a slice are unique
func IsUnique[T comparable](slice []T) bool {
	seen := make(map[T]bool)
	for _, value := range slice {
		if _, found := seen[value]; found {
			return false
		}
		seen[value] = true
	}
	return true
}

// IsEmpty returns true if a slice is empty
func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}

// IsNotEmpty returns true if a slice is not empty
func IsNotEmpty[T any](slice []T) bool {
	return !IsEmpty(slice)
}

// Choose returns a random item from a slice
func Choose[T any](slice []T) T {
	if len(slice) == 0 {
		panic("Cannot choose from empty slice")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return slice[r.Intn(len(slice))]
}

// EmliminationPool returns the options with the lowest count in a round
func EliminationPool[T comparable](round map[T]int, options []T) []T {
	var lowest []T
	lowestCount := round[options[0]]

	for _, option := range options[1:] {
		if round[option] < lowestCount {
			lowestCount = round[option]
		}
	}
	for _, option := range options {
		if round[option] == lowestCount {
			lowest = append(lowest, option)
		}
	}
	return lowest
}
