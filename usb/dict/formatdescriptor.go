package dict

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

//Those are the markers found in the hex dumps.
//For convenience I have added the ASCII representation as a comment
//in normal byte order and reverse byteorder (so you can find them in the hex dumps)
// Note: I have just guessed what the names could be from the marker ascii, I could be wrong ;-)
const (
	FormatDescriptorMagic uint32 = 0x66647363 //fdsc - csdf
	MediaTypeVideo        uint32 = 0x76696465 //vide - ediv
	MediaTypeMagic        uint32 = 0x6D646961 //mdia - aidm
	VideoDimensionMagic   uint32 = 0x7664696D //vdim - midv
	CodecMagic            uint32 = 0x636F6463 //codc - cdoc
	CodecAvc1             uint32 = 0x61766331 //avc1 - 1cva
	ExtensionMagic        uint32 = 0x6578746E //extn - ntxe
)

type FormatDescriptor struct {
	MediaType            uint32
	VideoDimensionWidth  uint32
	VideoDimensionHeight uint32
	Codec                uint32
}

func NewFormatDescriptorFromBytes(data []byte) (FormatDescriptor, error) {

	_, remainingBytes, err := parseLengthAndMagic(data, FormatDescriptorMagic)
	if err != nil {
		return FormatDescriptor{}, err
	}
	mediaType, remainingBytes, err := parseMediaType(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	videoDimensionWidth, videoDimensionHeight, remainingBytes, err := parseVideoDimension(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	codec, remainingBytes, err := parseCodec(remainingBytes)
	if err != nil {
		return FormatDescriptor{}, err
	}

	return FormatDescriptor{
		MediaType:            mediaType,
		VideoDimensionHeight: videoDimensionHeight,
		VideoDimensionWidth:  videoDimensionWidth,
		Codec:                codec,
	}, nil
}

func parseCodec(bytes []byte) (uint32, []byte, error) {
	length, _, err := parseLengthAndMagic(bytes, CodecMagic)
	if err != nil {
		return 0, nil, err
	}
	if length != 12 {
		return 0, nil, fmt.Errorf("invalid length for codec: %d", length)
	}
	codec := binary.LittleEndian.Uint32(bytes[8:])
	return codec, bytes[length:], nil
}

func parseVideoDimension(bytes []byte) (uint32, uint32, []byte, error) {
	length, _, err := parseLengthAndMagic(bytes, VideoDimensionMagic)
	if err != nil {
		return 0, 0, nil, err
	}
	if length != 16 {
		return 0, 0, nil, fmt.Errorf("invalid length for video dimension: %d", length)
	}
	width := binary.LittleEndian.Uint32(bytes[8:])
	height := binary.LittleEndian.Uint32(bytes[12:])
	return width, height, bytes[length:], nil
}

func parseMediaType(bytes []byte) (uint32, []byte, error) {
	length, _, err := parseLengthAndMagic(bytes, MediaTypeMagic)
	if err != nil {
		return 0, nil, err
	}
	if length != 12 {
		return 0, nil, fmt.Errorf("invalid length for media type: %d", length)
	}
	mediaType := binary.LittleEndian.Uint32(bytes[8:])
	return mediaType, bytes[length:], nil
}

func parseLengthAndMagic(bytes []byte, exptectedMagic uint32) (int, []byte, error) {
	length := binary.LittleEndian.Uint32(bytes)
	magic := binary.LittleEndian.Uint32(bytes[4:])
	if int(length) > len(bytes) {
		return 0, bytes, fmt.Errorf("invalid length in header: %d but only received: %d bytes", length, len(bytes))
	}
	if magic != exptectedMagic {
		unknownMagic := string(bytes[4:8])
		return 0, nil, fmt.Errorf("unknown magic type:%s (0x%x), cannot parse value %s", unknownMagic, magic, hex.Dump(bytes))
	}
	return int(length), bytes[8:], nil
}

func (fdsc FormatDescriptor) String() string {
	return fmt.Sprintf("FormatDescriptor:\n\t MediaType %s \n\t VideoDimension:(%dx%d) \n\t Codec:%s \n",
		readableMediaType(fdsc.MediaType), fdsc.VideoDimensionWidth, fdsc.VideoDimensionHeight, readableCodec(fdsc.Codec))
}

func readableCodec(codec uint32) string {
	if codec == CodecAvc1 {
		return "AVC-1"
	}
	return fmt.Sprintf("Unknown(%x)", codec)
}

func readableMediaType(mediaType uint32) string {
	if mediaType == MediaTypeVideo {
		return "Video"
	}
	return fmt.Sprintf("Unknown(%x)", mediaType)
}
