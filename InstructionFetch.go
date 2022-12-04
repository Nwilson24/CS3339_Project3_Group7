package main

import "fmt"

var pc = 96
var breakFetched = false

func InstructionFetch(pI []Instruction, buffer *[4]Instruction) {
	if breakFetched {
		return
	}
	//make sure instructions are moved to the top of the array
	empty := Instruction{}
	for i := 1; i < 4; i++ {
		if buffer[i] != empty && buffer[i-1] == empty {
			buffer[i-1] = buffer[i]
			buffer[i] = empty
		}
	}
	//make sure there's space in the buffer
	if buffer[3] != empty {
		return
	}

	//fetch instruction
	hit, _ := CacheRead(pc)
	if !hit {
		return //wait a cycle for it to be in cache
	}
	//detect branch
	i := (pc - 96) / 4
	if pI[i].opcode >= 160 && pI[i].opcode <= 191 {
		pc += pI[i].im * 4
		return
	}
	//detect conditional branch
	if pI[i].opcode >= 1440 && pI[i].opcode <= 1447 { //detect conditional branch
		reg := Registers[pI[i].conditional]
		if reg == 0 {
			pc += pI[i].im * 4
		}
		return
	}
	if pI[i].opcode >= 1448 && pI[i].opcode <= 1455 {
		reg := Registers[pI[i].conditional]
		if reg != 0 {
			pc += pI[i].im * 4
		}
		return
	}
	//check break
	if pI[i].opcode == 2038 {
		breakFetched = true
		return
	}
	//check for nop then add to buffer and increment pc
	if pI[i].opcode != 0 {
		pushBuffer(pI[(pc-96)/4], buffer)
		pc += 4
	}

	//check that buffer isn't now full
	if buffer[3] != empty {
		return
	}

	//there's probably a more elegant way to do this
	//fetch instruction
	hit, _ = CacheRead(pc)
	if !hit {
		return //wait a cycle for it to be in cache
	}
	//detect branch
	i = (pc - 96) / 4
	if pI[i].opcode >= 160 && pI[i].opcode <= 191 {
		pc += pI[i].im * 4
		return
	}
	//detect conditional branch
	if pI[i].opcode >= 1440 && pI[i].opcode <= 1447 { //detect conditional branch
		reg := Registers[pI[i].conditional]
		if reg == 0 {
			pc += pI[i].im * 4
		}
		return
	}
	if pI[i].opcode >= 1448 && pI[i].opcode <= 1455 {
		reg := Registers[pI[i].conditional]
		if reg != 0 {
			pc += pI[i].im * 4
		}
		return
	}
	//check break
	if pI[i].opcode == 2038 {
		breakFetched = true
		return
	}
	//check for nop then add to buffer and increment pc
	if pI[i].opcode != 0 {
		pushBuffer(pI[(pc-96)/4], buffer)
		pc += 4
	}
}

func pushBuffer(v Instruction, buffer *[4]Instruction) {
	//make sure instructions are moved to the top of the array
	empty := Instruction{}
	for i := 1; i < 4; i++ {
		if buffer[i] != empty && buffer[i-1] == empty {
			buffer[i-1] = buffer[i]
			buffer[i] = empty
		}
	}
	//push in new instruction if there is space
	for i := 0; i < 4; i++ {
		if buffer[i] == empty {
			buffer[i] = v
			return
		}
	}
	//if the program reaches here, something has gone wrong
	fmt.Println("error: push to preIssue buffer failed")
}
