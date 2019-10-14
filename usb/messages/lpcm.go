package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/*
All this are just guesses about what this could be from looking at the HexDump.
I guessed that this block could be LPCM (linear pulse code modulation) information although i Dont know.
Seems like 0x 00 00 00 00 00 70 E7 40 6B is some kind of separator or constant, as it is repeated 3 times.
0x6D 63 70 6C renders to "mcpl" (or lpcm in bigendian) in ascii. followed by a few ints.

HexDump:
00 00 00 00 00 70 E7 40 6D 63 70 6C 0C 00 00 00 04 00 00 00 01 00 00 00 04 00 00 00 02 00 00 00 10 00 00 00 00 00 00 00 00 00 00 00 00 70 E7 40 00 00 00 00 00 70 E7 40 6B

*/

//LPCMData contains 7 ints that probably mean something related to linear pulse code audio things
type LPCMData struct {
	Unknown_int1 uint32
	Unknown_int2 uint32
	Unknown_int3 uint32
	Unknown_int4 uint32
	Unknown_int5 uint32
	Unknown_int6 uint32
	Unknown_int7 uint32
}

func (data LPCMData) String() string {
	return fmt.Sprintf("[%d,%d,%d,%d,%d,%d,%d]", data.Unknown_int1, data.Unknown_int2, data.Unknown_int3,
		data.Unknown_int4, data.Unknown_int5, data.Unknown_int6, data.Unknown_int7)
}

const (
	separator uint64 = 0x40E7700000000000
	LpcmMagic uint32 = 0x6C70636D
)

//NewLPCMDataFromBytes reads 7 uint32 and puts them into a LPCMData struct
func NewLPCMDataFromBytes(data []byte) (LPCMData, error) {
	r := bytes.NewReader(data)
	var lpcmData LPCMData
	err := binary.Read(r, binary.LittleEndian, &lpcmData)
	if err != nil {
		return lpcmData, err
	}
	return lpcmData, nil
}

func createLpcmInfo() []byte {
	lpcmBytes := make([]byte, 56)
	binary.LittleEndian.PutUint64(lpcmBytes, separator)
	var index = 8
	binary.LittleEndian.PutUint32(lpcmBytes[index:], LpcmMagic)
	index += 4

	binary.LittleEndian.PutUint32(lpcmBytes[index:], 12)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 4)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 1)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 4)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 2)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 16)
	index += 4
	binary.LittleEndian.PutUint32(lpcmBytes[index:], 0)
	index += 4

	binary.LittleEndian.PutUint64(lpcmBytes[index:], separator)
	index += 8
	binary.LittleEndian.PutUint64(lpcmBytes[index:], separator)
	return lpcmBytes
}
