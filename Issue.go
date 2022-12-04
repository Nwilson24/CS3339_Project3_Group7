package main

var hazRegs [32]int

func Issue(preIssue *[4]Instruction, preALU *[2]Instruction, preMem *[2]Instruction, postMem *Instruction) {
	//take instructions, check for RAW Hazards, issue instruction to mam and alu buffers
	empty := Instruction{}
	//update array of unavailable registers
	for i := 0; i < len(hazRegs); i++ {
		if hazRegs[i] != 0 {
			hazRegs[i]--
		}
		if i == int(postMem.rd) {
			hazRegs[i] = 1
		}
	}

	var inst1 Instruction
	var index int
	for i := 0; i < 4; i++ {
		if preIssue[i] != empty && hazRegs[preIssue[i].rd] == 0 {
			inst1 = preIssue[i]
			index = i
			break
		}
	}
	//just in case no instructions can be passed on
	if inst1 == empty {
		return
	}

	if (inst1.opcode == 1984 || inst1.opcode == 1986) && (preMem[0] == empty || preMem[1] == empty) {
		//check and account for raw hazards
		if inst1.opcode == 1986 {
			hazRegs[inst1.rd] = -1 //restricted until detected in post-mem buffer
		}
		//sort mem buffer if necessary
		if preMem[0] == empty && preMem[1] != empty {
			preMem[0] = preMem[1]
			preMem[1] = empty
		}
		//send instruction to mem buffer
		preIssue[index] = empty
		if preMem[0] == empty {
			preMem[0] = inst1
		} else {
			preMem[1] = inst1
		}
	} else if preALU[0] == empty || preALU[1] == empty {
		//check and account for raw hazards
		hazRegs[inst1.rd] = 2
		//sort ALU buffer if necessary
		if preALU[0] == empty && preALU[1] != empty {
			preALU[0] = preALU[1]
			preALU[1] = empty
		}
		//send instruction to ALU buffer
		preIssue[index] = empty
		if preALU[0] == empty {
			preALU[0] = inst1
		} else {
			preALU[1] = inst1
		}
	}

	//again, there probably is a more elegant way to do this
	//do same thing again with second instruction
	var inst2 Instruction
	for i := index; i < 4; i++ {
		if preIssue[i] != empty && hazRegs[preIssue[i].rd] == 0 {
			inst2 = preIssue[i]
			index = i
			break
		}
	}
	//make sure we don't waste time passing along an empty instruction
	if inst2 == empty {
		return
	}

	if (inst2.opcode == 1984 || inst2.opcode == 1986) && (preMem[0] == empty || preMem[1] == empty) {
		//check and account for raw hazards
		if inst2.opcode == 1986 {
			hazRegs[inst2.rd] = -1 //restricted until detected in post-mem buffer
		}
		//sort mem buffer if necessary
		if preMem[0] == empty && preMem[1] != empty {
			preMem[0] = preMem[1]
			preMem[1] = empty
		}
		//send instruction to mem buffer
		preIssue[index] = empty
		if preMem[0] == empty {
			preMem[0] = inst2
		} else {
			preMem[1] = inst2
		}
	} else if preALU[0] == empty || preALU[1] == empty {
		//check and account for raw hazards
		hazRegs[inst2.rd] = 2
		//sort ALU buffer if necessary
		if preALU[0] == empty && preALU[1] != empty {
			preALU[0] = preALU[1]
			preALU[1] = empty
		}
		//send instruction to ALU buffer
		preIssue[index] = empty
		if preALU[0] == empty {
			preALU[0] = inst2
		} else {
			preALU[1] = inst2
		}
	}
}

