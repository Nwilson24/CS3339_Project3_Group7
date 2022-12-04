package main

import (
	"math"
)

type DataArray struct {
	startAddress int
	arr          []int
}

var Registers [32]int

func (d *DataArray) WriteData(num int, address int) { //add data into slice and make sure it stays a multiple of 8
	i := (address - 96) / 4
	if len(d.arr) <= i { //no space; we have to extend the slice by a multiple of 8 to fit the new data
		//find multiple of 8 to append to data slice
		diff := float64(i - (len(d.arr) - 1))
		t := 8 * int(math.Ceil(diff/8))
		//extend slice; I couldn't figure out a more efficient way but this works
		for i := 0; i < t; i++ {
			d.arr = append(d.arr, 0)
		}
	}
	d.arr[i] = num
}

func (d *DataArray) ReadData(address int) int {
	index := (address - 96) / 4
	if index < len(d.arr) {
		return d.arr[index]
	}
	return 0
}
