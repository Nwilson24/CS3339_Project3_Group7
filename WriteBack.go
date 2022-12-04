package main

// receives destination register and calculated address from post MEM
// only accessed from LDUR instruction
// writes loaded data into destination register

func WriteBack(memPI *Instruction, aluPI *Instruction) { // destination register(post MEM) = data in given address space(post MEM)
	// post ALU is used to do the WB for the one instruction in the buffer, send to destination register
	// can run both a post ALU and post MEM instruction simultaneously if offered
	// need to check if either buffer is empty
	empty := Instruction{}

	if *memPI == empty {
		return
	} else {
		Registers[memPI.rd] = memPI.memResult
		*memPI = empty
	}

	if *aluPI == empty {
		return
	} else {
		Registers[aluPI.rd] = aluPI.aluResult
		*aluPI = empty
	}
}
