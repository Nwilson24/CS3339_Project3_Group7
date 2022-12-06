package main

import (
	"fmt"
	"strconv"
)

var cache [4][2]CacheEntry

type CacheEntry struct {
	valid      int
	dirty      int
	LRU        int
	tag        int
	word1      int
	word2      int
	memAddress int
}

func CacheSetup() {
	//initialize entries
	var temp [4][2]CacheEntry
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			temp[i][j] = CacheEntry{0, 0, 0, 0, 0, 0, 0}
		}
	}
	//fill cache
	cache = temp
}

func CacheRead(address int) (bool, int) {
	//take address and return true, [value] if hit
	//if miss, return false, 0 and bring it into cache
	i := (address / 8) % 4
	for j := 0; j < 2; j++ {
		if cache[i][j].memAddress == address || cache[i][j].memAddress == (address-4) { //remember blocks are 2 words wide
			cache[i][j].LRU = 1
			//fix other LRU in block
			if j == 0 {
				cache[i][1].LRU = 0
			} else {
				cache[i][0].LRU = 0
			}
			//check which word we need to read from
			if cache[i][j].memAddress == address {
				return true, cache[i][j].word1
			} else {
				return true, cache[i][j].word2
			}
		}
	}
	//handle a miss
	MemFetch(address)
	return false, 0
}

func CacheWrite(address int, value int) bool {
	// search cache for value; return if it finds a hit
	i := (address / 8) % 4
	//check for matching address
	for j := 0; j < 2; j++ {
		if cache[i][j].memAddress == address || cache[i][j].memAddress == (address-4) { //remember blocks are 2 words wide
			//set flags and write
			cache[i][j].dirty = 1
			cache[i][j].LRU = 1
			//check which word we need to write to
			if cache[i][j].memAddress == address {
				cache[i][j].word1 = value
			} else {
				cache[i][j].word2 = value
			}
			//fix other LRU in block
			if j == 0 {
				cache[i][1].LRU = 0
			} else {
				cache[i][0].LRU = 0
			}
			return true
		}
	}
	//check for an empty block
	j := -1
	if cache[i][0].valid == 0 {
		j = 0
	}
	if cache[i][1].valid == 0 {
		j = 1
	}
	if j >= 0 {
		//set flags and write
		cache[i][j].valid = 1
		cache[i][j].dirty = 1
		cache[i][j].tag = address / 32
		cache[i][j].LRU = 1
		//check which word we need to write to
		if address%8 == 0 {
			cache[i][j].word1 = value
			cache[i][j].word2 = Data.ReadData(address + 4)
			cache[i][j].memAddress = address
		} else {
			cache[i][j].word1 = Data.ReadData(address - 4)
			cache[i][j].word2 = value
			cache[i][j].memAddress = address - 4
		}
		//fix other LRU in block
		if j == 0 {
			cache[i][1].LRU = 0
		} else {
			cache[i][0].LRU = 0
		}
		return true
	}
	//handle the miss and kick out least recently used block
	if cache[i][0].LRU == 0 {
		if cache[i][0].dirty == 1 {
			Data.WriteData(cache[i][0].word1, cache[i][0].memAddress)
			Data.WriteData(cache[i][0].word2, cache[i][0].memAddress+4)
		}
		cache[i][0] = CacheEntry{0, 0, 0, 0, 0, 0, 0}
		//fix LRU so its not overwritten by fetch
		cache[i][0].LRU = 1
		cache[i][1].LRU = 0
	} else {
		if cache[i][1].dirty == 1 {
			Data.WriteData(cache[i][1].word1, cache[i][1].memAddress)
			Data.WriteData(cache[i][1].word2, cache[i][1].memAddress+4)
		}
		cache[i][1] = CacheEntry{0, 0, 0, 0, 0, 0, 0}
		//fix LRU so its not overwritten by fetch
		cache[i][1].LRU = 1
		cache[i][0].LRU = 0
	}
	return false
}

func MemFetch(address int) { //**************************************************** still need to set tags!!
	set := (address / 8) % 4
	var t int
	//check and reset LRU values
	if cache[set][0].LRU == 0 {
		cache[set][0].LRU = 1
		cache[set][1].LRU = 0
		t = 0
	} else {
		cache[set][0].LRU = 0
		cache[set][1].LRU = 1
		t = 1
	}
	//write back to mem if dirty
	if cache[set][t].dirty == 1 {
		Data.WriteData(cache[set][t].word1, cache[set][t].memAddress)
		Data.WriteData(cache[set][t].word2, cache[set][t].memAddress+4)
	}
	//write in flags
	cache[set][t].valid = 1
	cache[set][t].dirty = 0
	//write in both words from the block
	if address%8 == 0 {
		cache[set][t].word1 = Data.ReadData(address)
		cache[set][t].word2 = Data.ReadData(address + 4)
		cache[set][t].memAddress = address
	} else {
		cache[set][t].word1 = Data.ReadData(address - 4)
		cache[set][t].word2 = Data.ReadData(address)
		cache[set][t].memAddress = address - 4
	}
	//set tag
	cache[set][t].tag = cache[set][t].memAddress / 32
}

func CacheToString() string {
	s := "Cache\n"
	for i := 0; i < 4; i++ {
		s += fmt.Sprintf("Set %v: LRU = %v\n", i, cache[i][1].LRU)
		//format data
		word1 := ""
		word2 := ""
		if cache[i][0].word1 < 0 {
			word1 = strconv.FormatUint(uint64(cache[i][0].word1), 2)[32:]
		} else {
			word1 = strconv.FormatUint(uint64(cache[i][0].word1), 2)
		}
		if cache[i][0].word2 < 0 {
			word2 = strconv.FormatUint(uint64(cache[i][0].word2), 2)[32:]
		} else {
			word2 = strconv.FormatUint(uint64(cache[i][0].word2), 2)
		}
		s += fmt.Sprintf("Entry 0:[(%v, %v, %v)<%v,%v>]\n", cache[i][0].valid, cache[i][0].dirty, cache[i][0].tag, word1, word2)
		if cache[i][1].word1 < 0 {
			word1 = strconv.FormatUint(uint64(cache[i][1].word1), 2)[33:]
		} else {
			word1 = strconv.FormatUint(uint64(cache[i][1].word1), 2)
		}
		if cache[i][1].word2 < 0 {
			word2 = strconv.FormatUint(uint64(cache[i][1].word2), 2)[33:]
		} else {
			word2 = strconv.FormatUint(uint64(cache[i][1].word2), 2)
		}
		s += fmt.Sprintf("Entry 1:[(%v, %v, %v)<%v,%v>]\n", cache[i][1].valid, cache[i][1].dirty, cache[i][1].tag, word1, word2)
	}
	return s
}

func CacheFlush() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if cache[i][j].dirty == 1 {
				Data.WriteData(cache[i][j].word1, cache[i][j].memAddress)
				Data.WriteData(cache[i][j].word2, cache[i][j].memAddress+4)
				cache[i][j].dirty = 0
			}
		}
	}
}
