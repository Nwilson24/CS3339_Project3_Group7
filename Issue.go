package main

var hazRegs [32]int

func Issue(preIssue *[4]Instruction, preALU *[2]Instruction, preMem *[2]Instruction, postMem *Instruction, postALU *Instruction) {
	//take instructions, check for RAW Hazards, issue instruction to mam and alu buffers
	empty := Instruction{}
	//update array of unavailable registers; anything not -1 might cause a raw hazard
	for i := 0; i < len(hazRegs); i++ {
		if hazRegs[i] > 0 {
			hazRegs[i]--
		}
		if hazRegs[i] == -1 && i == int(postMem.rd) {
			hazRegs[i] = 1
		}
		if hazRegs[i] == -1 && i == int(postALU.rd) {
			hazRegs[i] = 1
		}
	}

	var inst1 Instruction
	var index int
	for i := 0; i < 4; i++ {
		if preIssue[i] != empty && hazRegs[preIssue[i].rd] == 0 && hazRegs[preIssue[i].rn] == 0 {
			if preIssue[i].opcode == 1104 || preIssue[i].opcode == 1112 || preIssue[i].opcode == 1360 || preIssue[i].opcode == 1624 || preIssue[i].opcode == 1872 {
				if hazRegs[preIssue[i].rm] == 0 {
					inst1 = preIssue[i]
					index = i
					break
				}
			} else {
				inst1 = preIssue[i]
				index = i
				break
			}
		}
		//make sure no instructions dependent on this one are executed
		if preIssue[i] != empty {
			hazRegs[preIssue[i].rd] = -2
		}
	}
	//just in case no instructions can be passed on
	if inst1 == empty {
		//clear temp hazards
		for j := 0; j < len(hazRegs); j++ {
			if hazRegs[j] == -2 {
				hazRegs[j] = 0
			}
		}
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
		hazRegs[inst1.rd] = -1
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
	index++
	for i := index; i < 4; i++ {
		if preIssue[i] != empty && hazRegs[preIssue[i].rd] == 0 && hazRegs[preIssue[i].rn] == 0 {
			if preIssue[i].opcode == 1104 || preIssue[i].opcode == 1112 || preIssue[i].opcode == 1360 || preIssue[i].opcode == 1624 || preIssue[i].opcode == 1872 {
				if hazRegs[preIssue[i].rm] == 0 {
					inst2 = preIssue[i]
					index = i
					break
				}
			} else {
				inst2 = preIssue[i]
				index = i
				break
			}
		}
		if preIssue[i] != empty {
			hazRegs[preIssue[i].rd] = -2
		}
	}
	//make sure we don't waste time passing along an empty instruction
	if inst2 == empty {
		//clear temp hazards
		for j := 0; j < len(hazRegs); j++ {
			if hazRegs[j] == -2 {
				hazRegs[j] = 0
			}
		}
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
		hazRegs[inst2.rd] = -1
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
	//clear temporary hazards (-2)
	for j := 0; j < len(hazRegs); j++ {
		if hazRegs[j] == -2 {
			hazRegs[j] = 0
		}
	}
}
