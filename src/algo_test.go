package main

import (
	"fmt"
	"reflect"
	"testing"
)

func rangeAndAssert(indexes []IndexOperation, resultExpected []uint32, t *testing.T) {
	keys := make(map[int]map[int]int)
	result := []uint32{}

	for id := rangeIndexes(&keys, indexes); id > 0; {
		result = append(result, id)
		id = rangeIndexes(&keys, indexes)
	}

	if !reflect.DeepEqual(result, resultExpected) {
		fmt.Println("Expected:", resultExpected, "Actual:", result)
		t.Fail()
	}
}

func TestAlgo1(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "AND", Indexes: [][]uint32{{8, 6}}},
		IndexOperation{Condition: "AND", Indexes: [][]uint32{{10, 8, 6}}},
	}

	rangeAndAssert(indexes, []uint32{8, 6}, t)
}

func TestAlgo2(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "AND", Indexes: [][]uint32{{100, 80, 70, 20}}},
		IndexOperation{Condition: "AND", Indexes: [][]uint32{{90, 80}, {80, 10}}},
	}

	rangeAndAssert(indexes, []uint32{80}, t)
}

func TestAlgo3(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "AND", Indexes: [][]uint32{{100, 80, 70, 20}}},
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{90, 80}, {100, 80, 10}}},
	}

	rangeAndAssert(indexes, []uint32{100, 80}, t)
}

func TestAlgo4(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "AND", Indexes: [][]uint32{{100, 80, 70, 20}}},
	}

	rangeAndAssert(indexes, []uint32{100, 80, 70, 20}, t)
}

func TestAlgo5(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{100, 80, 70, 20}, {50, 45, 32}, {200, 90}}},
	}

	rangeAndAssert(indexes, []uint32{200, 100, 90, 80, 70, 50, 45, 32, 20}, t)
}

func TestAlgo6(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{100, 80, 70, 20}, {50, 45, 32}, {200, 90}}},
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{50, 45, 32}, {201, 200, 90}}},
	}

	rangeAndAssert(indexes, []uint32{200, 90, 50, 45, 32}, t)
}

func TestAlgo7(t *testing.T) {
	indexes := []IndexOperation{
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{100, 80, 70, 20}, {50, 45, 32}, {200, 90}}},
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{50, 45, 32}, {201, 200, 90}}},
		IndexOperation{Condition: "OR", Indexes: [][]uint32{{200, 90}}},
	}

	rangeAndAssert(indexes, []uint32{200, 90}, t)
}
