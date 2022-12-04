package main

/*
import (
	"fmt"
	"os"
)

func executionDecision(pI []Instruction, d *DataArray, file *os.File) {
	cycle := 1
	out := ""
	for i := 0; i < len(pI); i++ {
		out = fmt.Sprintf("====================\ncycle:%v\t", cycle)
		offset := 0
		switch true {
		case pI[i].opcode >= 160 && pI[i].opcode <= 191:
			offset = bExecution(pI[i])
			out += bPrint(pI[i])[33:] //slice of previous output excluding the binary
		case pI[i].opcode == 1104, pI[i].opcode == 1112, pI[i].opcode == 1360,
			pI[i].opcode == 1624, pI[i].opcode == 1690, pI[i].opcode == 1691,
			pI[i].opcode == 1692, pI[i].opcode == 1872:
			rExecution(pI[i])
			out += rPrint(pI[i])[36:]
		case pI[i].opcode >= 1160 && pI[i].opcode <= 1161,
			pI[i].opcode >= 1672 && pI[i].opcode <= 1673:
			iExecution(pI[i])
			out += iPrint(pI[i])[35:]
		case pI[i].opcode >= 1440 && pI[i].opcode <= 1447,
			pI[i].opcode >= 1448 && pI[i].opcode <= 1455:
			offset = cbExecution(pI[i])
			out += cbPrint(pI[i])[34:]
		case pI[i].opcode >= 1684 && pI[i].opcode <= 1687,
			pI[i].opcode >= 1940 && pI[i].opcode <= 1943:
			imExecution(pI[i])
			out += imPrint(pI[i])[35:]
		case pI[i].opcode == 1984, pI[i].opcode == 1986:
			dExecution(pI[i], d)
			out += dPrint(pI[i])[36:]
		case pI[i].opcode == 2038:
			out += fmt.Sprintf("\t%v\tBREAK\n", pI[i].programCount)
		default:
			nopExecution()
			out += fmt.Sprintf("\t%v\tNOP\n", pI[i].programCount)
		}
		//print registers and memory
		out += DataRegPrint(d)
		out += "\n"
		//write to the output file
		file.WriteString(out)
		//update cycle and fix i from branch instructions if necessary
		i += offset
		cycle++
		if cycle > 1000000 { //make sure we aren't stuck in an infinite loop
			break
		}
	}
}

func bExecution(pI Instruction) int {
	return pI.im - 1
}
func rExecution(pI Instruction) {
	rn := pI.rn
	rm := pI.rm
	rd := pI.rd

	switch true {
	case pI.opcode == 1104:
		Registers[rd] = Registers[rm] & Registers[rn]
	case pI.opcode == 1112:
		Registers[rd] = Registers[rm] + Registers[rn]
	case pI.opcode == 1360:
		Registers[rd] = Registers[rm] | Registers[rn]
	case pI.opcode == 1624:
		Registers[rd] = Registers[rn] - Registers[rm]
	case pI.opcode == 1690:
		Registers[rd] = int(uint(Registers[rn]) >> pI.shamt)
	case pI.opcode == 1691:
		Registers[rd] = int(uint(Registers[rn]) << pI.shamt)
	case pI.opcode == 1692:
		Registers[rd] = Registers[rn] >> pI.shamt
	case pI.opcode == 1872:
		Registers[rd] = Registers[rm] ^ Registers[rn]
	default:
		fmt.Println("R Instruction ERROR")
	}
}
func iExecution(pI Instruction) {
	rn := pI.rn
	im := pI.im
	rd := pI.rd

	switch true {
	case pI.opcode >= 1160 && pI.opcode <= 1161:
		Registers[rd] = Registers[rn] + im
	case pI.opcode >= 1672 && pI.opcode <= 1673:
		Registers[rd] = Registers[rn] - im
	default:
		fmt.Println("I Instruction ERROR")
	}
}
func cbExecution(pI Instruction) int {
	con := pI.conditional

	switch true {
	case pI.opcode >= 1440 && pI.opcode <= 1447:
		{
			if Registers[con] == 0 {
				return pI.im - 1
			} else {
				return 0
			}
		}
	case pI.opcode >= 1448 && pI.opcode <= 1455:
		{
			if Registers[con] != 0 {
				return pI.im - 1
			} else {
				return 0
			}
		}
	default:
		return 0
	}
}
func imExecution(pI Instruction) {
	rd := pI.rd

	switch true {
	case pI.opcode >= 1684 && pI.opcode <= 1687:
		Registers[rd] = int(pI.field << (pI.shiftCode * 16))
	case pI.opcode >= 1940 && pI.opcode <= 1943:
		temp := uint64(pI.field) << (pI.shiftCode * 16)
		if pI.shiftCode == 0 {
			temp2 := uint64(Registers[rd]) & 0xFFFFFFFFFFFF0000
			Registers[rd] = int(temp2 | temp)
		}
		if pI.shiftCode == 1 {
			temp2 := uint64(Registers[rd]) & 0xFFFFFFFF0000FFFF
			Registers[rd] = int(temp2 | temp)
		}
		if pI.shiftCode == 2 {
			temp2 := uint64(Registers[rd]) & 0xFFFF0000FFFFFFFF
			Registers[rd] = int(temp2 | temp)
		}
		if pI.shiftCode == 3 {
			temp2 := uint64(Registers[rd]) & 0x0000FFFFFFFFFFFF
			Registers[rd] = int(temp2 | temp)
		}
	default:
		fmt.Println("IM Instruction ERROR")
	}
}
func dExecution(pI Instruction, d *DataArray) {
	switch true {
	case pI.opcode == 1984:
		d.WriteData(Registers[pI.rd], int(uint16(Registers[pI.rn])+pI.address))
	case pI.opcode == 1986:
		Registers[pI.rd] = d.ReadData(int(uint16(Registers[pI.rn]) + pI.address))
	default:
		fmt.Println("D Instruction ERROR")
	}
}
func nopExecution() {
	return
}
*/
