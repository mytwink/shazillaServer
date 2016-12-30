package wav

import (
		//"os"
		"log"
		"bytes"
		"bufio"
		"fmt"
		"encoding/binary"
		"mime/multipart"
		"math"
)

const HEADER_LENGTH = 44

type wavSizeError struct {
	size int64
	msg string
}

func (e wavSizeError) Error() string {
	return fmt.Sprintf("%v: %v", e.size, e.msg)
}

func readString(handle *bufio.Reader, length int64) string {
	buffer := make([]byte, length)
	_, err := handle.Read(buffer)
	if err != nil {
		log.Println("Error while reading string from file", err)
	}
	return string(buffer)
}

func readLong(handle *bufio.Reader) int64 {
	var n uint32
	buffer := make([]byte, 4)
	_, err := handle.Read(buffer)
	buf := bytes.NewReader(buffer)
	err = binary.Read(buf, binary.LittleEndian, &n)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}
	return int64(n)
}

func readWord(handle *bufio.Reader) int64 {
	var n uint16
	buffer := make([]byte, 2)
	_, err := handle.Read(buffer)
	buf := bytes.NewReader(buffer)
	err = binary.Read(buf, binary.LittleEndian, &n)
	if err != nil {
		log.Println("binary.Read failed:", err)
	}
	return int64(n)
}

type header struct {
	Chunkid,Format string
	Chunksize int64
}

type subchunk1 struct {
	Id string
	Size,Audioformat,Numchannels,Samplerate,Byterate,Blockalign,Bitspersample int64
}

type subchunk2 struct {
	Size int64
	Id string
	Data []int64
}

type wav struct {
	Header header
	Subchunk1 subchunk1
	Subchunk2 subchunk2
}

type WavParse struct {
	Wav wav
}


func NewWavParse(file multipart.File) (wavParse *WavParse, err error) {
	/*file, err := os.Open(filename)
	if err!=nil {
		return
	}

	finfo, err := os.Stat(filename)
	if err != nil {
		return
	}

	if finfo.Size() < int64(HEADER_LENGTH) {
		err = wavSizeError{
			finfo.Size(),
			"the file size is too small",
		}
	}*/
	handle := bufio.NewReader(file)

	wavParse = &WavParse{}
	wavParse.Wav.Header.Chunkid = readString(handle, 4)
	wavParse.Wav.Header.Chunksize = readLong(handle)
	wavParse.Wav.Header.Format = readString(handle, 4)

	wavParse.Wav.Subchunk1.Id = readString(handle, 4)
	wavParse.Wav.Subchunk1.Size = readLong(handle)
	wavParse.Wav.Subchunk1.Audioformat = readWord(handle)
	wavParse.Wav.Subchunk1.Numchannels = readWord(handle)
	wavParse.Wav.Subchunk1.Samplerate = readLong(handle)
	wavParse.Wav.Subchunk1.Byterate = readLong(handle)
	wavParse.Wav.Subchunk1.Blockalign = readWord(handle)
	wavParse.Wav.Subchunk1.Bitspersample = readWord(handle)

	wavParse.Wav.Subchunk2.Id = readString(handle, 4)
	wavParse.Wav.Subchunk2.Size = readLong(handle)
	data := []int64{}

	peek := wavParse.Wav.Subchunk1.Bitspersample
	bite := peek/8

	skeepingBytesCount := bite*wavParse.Wav.Subchunk1.Numchannels

	buffer := make([]byte, bite)
	n:=1
	for n>0 {
		var val float64
		var j int64
		
		for j=0; j<wavParse.Wav.Subchunk1.Numchannels; j++ {
			n, err = handle.Read(buffer)
			if err!=nil {
				break;
			}
		}

		if n>0 {
			switch bite {
				case 1:
					val+=float64(buffer[0])
				case 2:
					val+=(float64(buffer[0])+float64(buffer[1]))/2.0
			}
		}
		value := int64(math.Floor(val))
		if value == 0 {
			value++
		}
		data = append(data, value)

		if n>0 {
			buffer = make([]byte, skeepingBytesCount)
			_, err = handle.Read(buffer)
		}
		
	}
	file.Close()
	wavParse.Wav.Subchunk2.Data = data
	return
}