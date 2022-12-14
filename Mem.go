package main

import "fmt"

func MEM(preBuffer *[2]Instruction, postBuffer *Instruction) { // calculates address ONLY
	empty := Instruction{}
	if preBuffer[0] == empty {
		return
	}
	pI := preBuffer[0]
	var hit bool

	switch true {
	case pI.opcode == 1984:
		hit = CacheWrite(Registers[pI.rn]+int(pI.address*4), Registers[pI.rd])
		if !hit {
			return
		}
		preBuffer[0] = empty
		if preBuffer[1] != empty && preBuffer[0] == empty {
			preBuffer[0] = preBuffer[1]
			preBuffer[1] = empty
		}

	case pI.opcode == 1986:
		hit, pI.memResult = CacheRead(Registers[pI.rn] + int(pI.address*4))
		if !hit {
			return
		}
		preBuffer[0] = empty
		if preBuffer[1] != empty && preBuffer[0] == empty {
			preBuffer[0] = preBuffer[1]
			preBuffer[1] = empty
		}
		*postBuffer = pI

	default:
		fmt.Println("D Instruction ERROR")
	}
}
