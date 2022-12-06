package main

import "fmt"

func ALU(preBuffer *[2]Instruction, postBuffer *Instruction) {
	empty := Instruction{}
	if preBuffer[0] == empty {
		return
	}
	pI := preBuffer[0]
	preBuffer[0] = empty
	if preBuffer[1] != empty && preBuffer[0] == empty {
		preBuffer[0] = preBuffer[1]
		preBuffer[1] = empty
	}

	switch true {
	case pI.opcode == 1104, pI.opcode == 1112, pI.opcode == 1360,
		pI.opcode == 1624, pI.opcode == 1690, pI.opcode == 1691,
		pI.opcode == 1692, pI.opcode == 1872:
		rExecution(&pI)
	case pI.opcode >= 1160 && pI.opcode <= 1161,
		pI.opcode >= 1672 && pI.opcode <= 1673:
		iExecution(&pI)
	case pI.opcode >= 1684 && pI.opcode <= 1687,
		pI.opcode >= 1940 && pI.opcode <= 1943:
		imExecution(&pI)
	case pI.opcode == 2038:
	default:
		nopExecution()
	}

	*postBuffer = pI
}

func rExecution(pI *Instruction) { // Can't directly send data to destination register, that will be done in WB
	rn := pI.rn
	rm := pI.rm

	switch true {
	case pI.opcode == 1104:
		pI.aluResult = Registers[rm] & Registers[rn]
	case pI.opcode == 1112:
		pI.aluResult = Registers[rm] + Registers[rn]
	case pI.opcode == 1360:
		pI.aluResult = Registers[rm] | Registers[rn]
	case pI.opcode == 1624:
		pI.aluResult = Registers[rn] - Registers[rm]
	case pI.opcode == 1690:
		pI.aluResult = int(uint(Registers[rn]) >> pI.shamt)
	case pI.opcode == 1691:
		pI.aluResult = int(uint(Registers[rn]) << pI.shamt)
	case pI.opcode == 1692:
		pI.aluResult = Registers[rn] >> pI.shamt
	case pI.opcode == 1872:
		pI.aluResult = Registers[rm] ^ Registers[rn]
	default:
		fmt.Println("R Instruction ERROR")
	}
}
func iExecution(pI *Instruction) {
	rn := pI.rn
	im := pI.im

	switch true {
	case pI.opcode >= 1160 && pI.opcode <= 1161:
		pI.aluResult = Registers[rn] + im
	case pI.opcode >= 1672 && pI.opcode <= 1673:
		pI.aluResult = Registers[rn] - im
	default:
		fmt.Println("I Instruction ERROR")
	}
}
func imExecution(pI *Instruction) {

	switch true {
	case pI.opcode >= 1684 && pI.opcode <= 1687:
		pI.aluResult = int(pI.field << (pI.shiftCode * 16))
	case pI.opcode >= 1940 && pI.opcode <= 1943:
		temp := uint64(pI.field) << (pI.shiftCode * 16)
		if pI.shiftCode == 0 {
			temp2 := uint64(Registers[pI.rd]) & 0xFFFFFFFFFFFF0000
			pI.aluResult = int(temp2 | temp)
		}
		if pI.shiftCode == 1 {
			temp2 := uint64(Registers[pI.rd]) & 0xFFFFFFFF0000FFFF
			pI.aluResult = int(temp2 | temp)
		}
		if pI.shiftCode == 2 {
			temp2 := uint64(Registers[pI.rd]) & 0xFFFF0000FFFFFFFF
			pI.aluResult = int(temp2 | temp)
		}
		if pI.shiftCode == 3 {
			temp2 := uint64(Registers[pI.rd]) & 0x0000FFFFFFFFFFFF
			pI.aluResult = int(temp2 | temp)
		}
	default:
		fmt.Println("IM Instruction ERROR")
	}
}
func nopExecution() {
	return
}
