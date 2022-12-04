package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type Instruction struct {
	rawInstruction string
	lineValue      uint64
	programCount   int
	opcode         uint64
	rn             uint8
	rm             uint8
	rd             uint8
	op2            uint8
	conditional    uint8
	shiftCode      uint8
	op             uint16
	address        uint16
	field          uint32
	shamt          int
	im             int
	offset         int
}

// global data/cache despite passing in project 2 code
var Data DataArray = DataArray{0, []int{0, 0, 0, 0, 0, 0, 0, 0}}

func main() {

	var inputFile *string
	var outputFile *string

	//set up input/outputs
	inputFile = flag.String("i", "", "get input file name")
	outputFile = flag.String("o", "", "get output file name")

	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	//set up arrays for instructions and data
	parsedInstructions, parsedData := parseInput(txtlines)

	parsedInstructions = parseOpcode(parsedInstructions)

	//create and fill .dis file
	outDis, err := os.Create(*outputFile + "_dis.txt")
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	for _, pI := range parsedInstructions {
		outDis.WriteString(outPrint(pI))
	}
	breakAdr := parsedInstructions[len(parsedInstructions)-1].programCount
	temp := (breakAdr - 92) / 4
	for i := temp; i < len(parsedData); i++ {
		outDis.WriteString(printData(parsedData[i]))
	}

	//Begin Project 2
	//grab address for start of data and make struct

	Data.startAddress = breakAdr + 4
	//load data from input file into array
	for _, d := range parsedData {
		val, err := strconv.ParseInt(d.rawInstruction, 2, 64)
		if err != nil {
			fmt.Println("Error in reading data")
		}

		//check for negative
		if d.rawInstruction[0] == 49 {
			val2 := uint32(val)
			val2 = ^val2 + 1
			val = int64(val2) * -1
		}

		Data.WriteData(int(val), d.programCount)
	}

	outSim, err := os.Create(*outputFile + "_sim.txt")
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	executionDecision(parsedInstructions, &Data, outSim)

	//begin project 3 :) totally not procrastinated
	//set up cache
	CacheSetup()

	preIssueBuffer := [4]Instruction{}
	preMemBuffer := [2]Instruction{}
	preALUBuffer := [2]Instruction{}
	postMemBuffer := Instruction{}
	InstructionFetch(parsedInstructions, &preIssueBuffer)
	Issue(&preIssueBuffer, &preALUBuffer, &preMemBuffer, &postMemBuffer)
	fmt.Println(CacheToString())
}

func parseInput(txt []string) ([]Instruction, []Instruction) {
	var pInst []Instruction
	var pData []Instruction
	var i int
	for i = 0; i < len(txt); i++ {
		if txt[i] != "" {
			raw := txt[i]
			s := raw[0:11] //contains first 11 bits of the raw
			var opcode int = binaryStringToDecimal(s)

			inst := Instruction{rawInstruction: raw, opcode: uint64(opcode), programCount: (i * 4) + 96}
			pInst = append(pInst, inst)

			data := Instruction{rawInstruction: txt[i], programCount: (i * 4) + 96}
			pData = append(pData, data)

			if opcode == 2038 { //check for break instruction
				i++
				break
			}
		}
	}
	for ; i < len(txt); i++ {
		if txt[i] != "" {
			data := Instruction{rawInstruction: txt[i], programCount: (i * 4) + 96}
			pData = append(pData, data)
		}
	}
	return pInst, pData
}

func binaryStringToDecimal(s string) int { //convert binary numbers in a string to decimal value
	var sum float64 = 0
	for i := 0; i < len(s); i++ {
		if s[i] == 49 { //49 is the character code for "1"
			sum += math.Pow(2.0, float64(len(s)-i-1))
		}
	}
	return int(sum)
}

func parseOpcode(pI []Instruction) []Instruction {
	var parsed []Instruction
	var inst Instruction

	for i := 0; i < len(pI); i++ {
		switch true {
		case pI[i].opcode >= 160 && pI[i].opcode <= 191:
			inst = bFormat(pI[i])
		case pI[i].opcode == 1104, pI[i].opcode == 1112, pI[i].opcode == 1360,
			pI[i].opcode == 1624, pI[i].opcode == 1690, pI[i].opcode == 1691,
			pI[i].opcode == 1692, pI[i].opcode == 1872:
			inst = rFormat(pI[i])
		case pI[i].opcode >= 1160 && pI[i].opcode <= 1161,
			pI[i].opcode >= 1672 && pI[i].opcode <= 1673:
			inst = iFormat(pI[i])
		case pI[i].opcode >= 1440 && pI[i].opcode <= 1447,
			pI[i].opcode >= 1448 && pI[i].opcode <= 1455:
			inst = cBFormat(pI[i])
		case pI[i].opcode >= 1684 && pI[i].opcode <= 1687,
			pI[i].opcode >= 1940 && pI[i].opcode <= 1943:
			inst = iMFormat(pI[i])
		case pI[i].opcode == 1984, pI[i].opcode == 1986:
			inst = dFormat(pI[i])
		default:
			inst = nOPFormat(pI[i])
		}

		parsed = append(parsed, inst)
	}
	return parsed
}

func bFormat(pI Instruction) Instruction {
	temp, err := strconv.ParseInt(pI.rawInstruction, 2, 64)
	if err == nil {
	}
	pI.lineValue = uint64(temp)

	temp4 := pI.lineValue & 4292870144 //isolate 11 bit opcode for reference later
	temp4 = temp4 >> 21
	pI.op = uint16(temp4)

	temp2 := pI.lineValue & 67108863 //mask of 11111111111111111111111111
	pI.im = int(temp2)
	if pI.rawInstruction[6] == 49 { //don't ask how this works, I'd need a paragraph
		temp3 := uint32(temp2)
		temp3 = ^temp3 + 1
		temp3 -= 4227858432
		temp5 := int(temp3) * -1
		pI.im = temp5
	}

	return pI
}
func rFormat(pI Instruction) Instruction {
	temp, err := strconv.ParseInt(pI.rawInstruction, 2, 64)
	if err == nil {
	}
	pI.lineValue = uint64(temp)

	temp4 := pI.lineValue & 4292870144 //isolate 11 bit opcode for reference later
	temp4 = temp4 >> 21
	pI.op = uint16(temp4)

	temp2 := pI.lineValue & 2031616 //mask of 111110000000000000000
	temp2 = temp2 >> 16
	pI.rm = uint8(temp2)

	temp2 = pI.lineValue & 64512 //mask of 1111110000000000
	temp2 = temp2 >> 10
	pI.shamt = int(temp2)

	temp2 = pI.lineValue & 992 //mask of 1111100000
	temp2 = temp2 >> 5
	pI.rn = uint8(temp2)

	temp2 = pI.lineValue & 31 //mask of 11111
	pI.rd = uint8(temp2)

	return pI
}
func iFormat(pI Instruction) Instruction {
	temp, err := strconv.ParseInt(pI.rawInstruction, 2, 64)
	if err == nil {
	}
	pI.lineValue = uint64(temp)

	temp4 := pI.lineValue & 4292870144 //isolate 11 bit opcode for reference later
	temp4 = temp4 >> 21
	pI.op = uint16(temp4)

	temp2 := pI.lineValue & 4193280 //mask of 1111111111110000000000
	temp2 = temp2 >> 10
	pI.im = int(temp2)
	if pI.rawInstruction[10] == 49 { //dont ask how this works
		temp3 := uint32(temp2)
		temp3 = ^temp3 + 1
		temp3 -= 4294963200
		temp5 := int(temp3) * -1
		pI.im = temp5
	}

	temp2 = pI.lineValue & 992 //mask of 1111100000
	temp2 = temp2 >> 5
	pI.rn = uint8(temp2)

	temp2 = pI.lineValue & 31 //mask of 11111
	pI.rd = uint8(temp2)

	return pI
}
func cBFormat(pI Instruction) Instruction {
	temp, err := strconv.ParseInt(pI.rawInstruction, 2, 64)
	if err == nil {
	}
	pI.lineValue = uint64(temp)

	temp4 := pI.lineValue & 4292870144 //isolate 11 bit opcode for reference later
	temp4 = temp4 >> 21
	pI.op = uint16(temp4)

	temp2 := pI.lineValue & 16777184 //mask of 111111111111111111100000
	temp2 = temp2 >> 5
	pI.im = int(temp2)
	if pI.rawInstruction[8] == 49 { //dont ask how this works
		temp3 := uint32(temp2)
		temp3 = ^temp3 + 1
		temp3 -= 4294443008
		temp5 := int(temp3) * -1
		pI.im = temp5
	}

	temp2 = pI.lineValue & 31 // mask of 11111
	pI.conditional = uint8(temp2)

	return pI
}
func iMFormat(pI Instruction) Instruction {
	temp, err := strconv.ParseInt(pI.rawInstruction, 2, 64)
	if err == nil {
	}
	pI.lineValue = uint64(temp)

	temp4 := pI.lineValue & 4292870144 //isolate 11 bit opcode for reference later
	temp4 = temp4 >> 21
	pI.op = uint16(temp4)

	temp2 := pI.lineValue & 6291456 //mask of 11000000000000000000000
	temp2 = temp2 >> 21
	pI.shiftCode = uint8(temp2)

	temp2 = pI.lineValue & 2097120 //mask of 111111111111111100000
	temp2 = temp2 >> 5
	pI.field = uint32(temp2)

	temp2 = pI.lineValue & 31 //mask of 11111
	pI.rd = uint8(temp2)

	return pI
}
func dFormat(pI Instruction) Instruction {
	temp, err := strconv.ParseInt(pI.rawInstruction, 2, 64)
	if err == nil {
	}
	pI.lineValue = uint64(temp)

	temp4 := pI.lineValue & 4292870144 //isolate 11 bit opcode for reference later
	temp4 = temp4 >> 21
	pI.op = uint16(temp4)

	temp2 := pI.lineValue & 2093056 //mask of 111111111000000000000
	temp2 = temp2 >> 12
	pI.address = uint16(temp2)

	temp2 = pI.lineValue & 3072 //mask of 110000000000
	temp2 = temp2 >> 10
	pI.op2 = uint8(temp2)

	temp2 = pI.lineValue & 992 //mask of 1111100000
	temp2 = temp2 >> 5
	pI.rn = uint8(temp2)

	temp2 = pI.lineValue & 31 //mask of 11111
	pI.rd = uint8(temp2)

	return pI
}
func nOPFormat(pI Instruction) Instruction {
	return pI
}

func outPrint(pI Instruction) string {

	switch true {
	case pI.opcode >= 160 && pI.opcode <= 191:
		return bPrint(pI)
	case pI.opcode == 1104, pI.opcode == 1112, pI.opcode == 1360,
		pI.opcode == 1624, pI.opcode == 1690, pI.opcode == 1691,
		pI.opcode == 1692, pI.opcode == 1872:
		return rPrint(pI)
	case pI.opcode >= 1160 && pI.opcode <= 1161,
		pI.opcode >= 1672 && pI.opcode <= 1673:
		return iPrint(pI)
	case pI.opcode >= 1440 && pI.opcode <= 1447,
		pI.opcode >= 1448 && pI.opcode <= 1455:
		return cbPrint(pI)
	case pI.opcode >= 1684 && pI.opcode <= 1687,
		pI.opcode >= 1940 && pI.opcode <= 1943:
		return imPrint(pI)
	case pI.opcode == 1984, pI.opcode == 1986:
		return dPrint(pI)
	case pI.opcode == 2038:
		return breakPrint(pI)
	default:
		return nopPrint(pI)
	}

}

func bPrint(pI Instruction) string {
	s := fmt.Sprintf(pI.rawInstruction[0:6]+" "+pI.rawInstruction[6:32]+"\t%v\tB\t%v \n", pI.programCount, pI.im)
	return s
}
func rPrint(pI Instruction) string {
	var iName string
	switch true {
	case pI.opcode == 1104:
		iName = "AND"
	case pI.opcode == 1112:
		iName = "ADD"
	case pI.opcode == 1360:
		iName = "ORR"
	case pI.opcode == 1624:
		iName = "SUB"
	case pI.opcode == 1690:
		iName = "LSR"
	case pI.opcode == 1691:
		iName = "LSL"
	case pI.opcode == 1692:
		iName = "ASR"
	case pI.opcode == 1872:
		iName = "EOR"
	default:
		iName = "ERROR"
	}
	var s string
	if pI.opcode == 1690 || pI.opcode == 1691 || pI.opcode == 1692 {
		s = pI.rawInstruction[0:11] + " " + pI.rawInstruction[11:16] + " " + pI.rawInstruction[16:22] + " " + pI.rawInstruction[22:27] + " " + pI.rawInstruction[27:32]
		s += fmt.Sprintf("\t%v\t%v\tR%v, R%v, #%v \n", pI.programCount, iName, pI.rd, pI.rn, pI.shamt)
	} else {
		s = pI.rawInstruction[0:11] + " " + pI.rawInstruction[11:16] + " " + pI.rawInstruction[16:22] + " " + pI.rawInstruction[22:27] + " " + pI.rawInstruction[27:32]
		s += fmt.Sprintf("\t%v\t%v\tR%v, R%v, R%v \n", pI.programCount, iName, pI.rd, pI.rn, pI.rm)
	}
	return s
}
func iPrint(pI Instruction) string {
	var iName string
	switch true {
	case pI.opcode >= 1160 && pI.opcode <= 1161:
		iName = "ADDI"
	case pI.opcode >= 1672 && pI.opcode <= 1673:
		iName = "SUBI"
	default:
		iName = "ERROR"
	}
	s := pI.rawInstruction[0:10] + " " + pI.rawInstruction[10:22] + " " + pI.rawInstruction[22:27] + " " + pI.rawInstruction[27:32]
	s += fmt.Sprintf("\t%v\t%v\tR%v, R%v, #%v \n", pI.programCount, iName, pI.rd, pI.rn, pI.im)
	return s
}
func cbPrint(pI Instruction) string {
	var iName string
	switch true {
	case pI.opcode >= 1440 && pI.opcode <= 1447:
		iName = "CBZ"
	case pI.opcode >= 1448 && pI.opcode <= 1455:
		iName = "CBNZ"
	default:
		iName = "ERROR"
	}
	s := fmt.Sprintf(pI.rawInstruction[0:8] + " " + pI.rawInstruction[8:27] + " " + pI.rawInstruction[27:32])
	s += fmt.Sprintf("\t%v\t%v\tR%v #%v \n", pI.programCount, iName, pI.conditional, pI.im)
	return s
}
func imPrint(pI Instruction) string {
	var iName string
	switch true {
	case pI.opcode >= 1684 && pI.opcode <= 1687:
		iName = "MOVZ"
	case pI.opcode >= 1940 && pI.opcode <= 1943:
		iName = "MOVK"
	default:
		iName = "ERROR"
	}
	s := fmt.Sprintf(pI.rawInstruction[0:9] + " " + pI.rawInstruction[9:11] + " " + pI.rawInstruction[11:27] + " " + pI.rawInstruction[27:32])
	s += fmt.Sprintf("\t%v\t%v\tR%v, R%v, LSL %v\n", pI.programCount, iName, pI.rd, pI.field, pI.shiftCode*16)
	return s
}
func dPrint(pI Instruction) string {
	var iName string
	switch true {
	case pI.opcode == 1984:
		iName = "STUR"
	case pI.opcode == 1986:
		iName = "LDUR"
	default:
		iName = "ERROR"
	}
	s := pI.rawInstruction[0:11] + " " + pI.rawInstruction[11:20] + " " + pI.rawInstruction[20:22] + " " + pI.rawInstruction[22:27] + " " + pI.rawInstruction[27:32]
	s += fmt.Sprintf("\t%v\t%v\tR%v, [R%v,#%v] \n", pI.programCount, iName, pI.rd, pI.rn, pI.address)
	return s
}
func nopPrint(pI Instruction) string {
	//this instruction is a catch-all so i need to check if it is actually all zeros; otherwise output that we can't recognize the instruction
	if pI.rawInstruction != "00000000000000000000000000000000" {
		return "Instruction not recognized"
	}
	s := fmt.Sprintf(pI.rawInstruction+"\t%v\tNOP\n", pI.programCount)
	return s
}
func breakPrint(pI Instruction) string {
	s := pI.rawInstruction[0:1] + " " + pI.rawInstruction[1:6] + " " + pI.rawInstruction[6:11] + " " + pI.rawInstruction[11:16] + " " + pI.rawInstruction[16:21] + " " + pI.rawInstruction[21:26] + " " + pI.rawInstruction[26:32]
	s += fmt.Sprintf("\t%v\tBREAK\n", pI.programCount)
	return s
}

func printData(d Instruction) string {
	val, err := strconv.ParseInt(d.rawInstruction, 2, 64)
	if err != nil {
		return "ERROR"
	}

	//check for negative
	if d.rawInstruction[0] == 49 {
		val2 := uint32(val)
		val2 = ^val2 + 1
		val = int64(val2) * -1
	}

	s := fmt.Sprintf(d.rawInstruction+"\t%v\t%v \n", d.programCount, val)
	return s
}
