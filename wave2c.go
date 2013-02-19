/* A small utility to convert RIFF wave audio files into C-arrays.
 * written in 2013 by tpltnt
 * license: MIT
 */

package main

import(
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

// check slices for equality
func AreEqualSlices(a,b []byte) bool {
     // check for same length
     if len(a) != len(b) { return false }
     // check for same capacity
     if cap(a) != cap(b) { return false }
     // loop over a-slice & compare values
     for i := 0; i < len(b); i++ {
         if a[i] != b[i] { return false }
     }
     // same length, capacity & values
     return true
}

// Check header for given Wave-file
func IsWavefileHeaderOk(reader *bufio.Reader) bool {
     // create buffer for 4 bytes
     chunkID := make([]byte, 4)
     // try to read first 4 bytes
     bytesread, err := reader.Read(chunkID)
     // panic if no error (and EOF)
     if err != nil && err != io.EOF { return false }
     // panic if less than 4 bytes read
     if bytesread < 4 { return false }
     // create slice with RIFF header bytes
     riffstring := []byte{0x52, 0x49, 0x46, 0x46}
     // compare slices using equality-test function
     if !AreEqualSlices(riffstring,chunkID) { return false }

     // throw away next 4 bytes
     bytesread, err = reader.Read(chunkID)
     // panic if no error (and EOF)
     if err != nil && err != io.EOF { return false }
     // panic if less than 4 bytes read
     if bytesread < 4 { return false }

     // read WAVE-header
     bytesread, err = reader.Read(chunkID)
     // panic if no error (and EOF)
     if err != nil && err != io.EOF { return false }
     // panic if less than 4 bytes read
     if bytesread < 4 { return false }
     // create slice with WAVE header bytes
     wavestring := []byte{0x57, 0x41, 0x56, 0x45}
     // compare slices using equality-test function
     if !AreEqualSlices(wavestring,chunkID) { return false }

     return true
}

// extract file format information
// assumes to be in beginning of (data) chunk
func IsWavefileFormatOk(reader *bufio.Reader) bool {
     // create buffer for 4 bytes
     chunkID := make([]byte, 4)
     //-- try to read first 4 bytes
     bytesread, err := reader.Read(chunkID)
     // panic if no error (and EOF)
     if err != nil && err != io.EOF {
        log.Println("error reading chunk ID")     
        return false
     }
     // panic if less than 4 bytes read
     if bytesread < 4 { return false }
     // create slice with fmt bytes
     fmtstring := []byte{0x66, 0x6d, 0x74, 0x20}
     // compare slices using equality-test function
     if !AreEqualSlices(fmtstring,chunkID) {
        log.Println("wrong chunk ID")
        return false
     }

     fourbytes := make([]byte, 4)
     twobytes := make([]byte, 2)
     //-- read data size
     bytesread, err = reader.Read(fourbytes)
     // panic if no error (and EOF)
     if err != nil && err != io.EOF {
        log.Println("error reading data size")
        return false
     }
     if 4 != bytesread {
        log.Println("not enough bytes read for data size")
     }
     // 16 bytes (little endian encoded)
     datachunksize := []byte{0x10, 0x00, 0x00, 0x00}
     if !AreEqualSlices(fourbytes,datachunksize) {
        log.Println("wrong data chunk size")
	return false
     }

     //-- read wFormatTag
     bytesread, err = reader.Read(twobytes)
     if err != nil && err != io.EOF {
        log.Println("error reading wFormatTag")
        return false
     }
     if 2 != bytesread {
        log.Println("not enough bytes read for wFormatTag")
        return false
     }
     // PCM = 0x01 0x00 (little endian for 1)
     if !(twobytes[0] == 0x01 && twobytes [1] == 0x00) {
        log.Println("no PCM found")
        return false
     }

     //-- read wChannels
     bytesread, err = reader.Read(twobytes)
     if err != nil && err != io.EOF {
        log.Println("error reading wChannels")
        return false
     }
     if 2 != bytesread {
        log.Println("not enough bytes read for wChannels")
        return false
     }
     // PCM = 0x01 0x00 (little endian for 1)
     if !(twobytes[0] == 0x01 && twobytes [1] == 0x00) {
        log.Println("no mono file found")
        return false
     }

     //-- read dwSamplesPerSec 4bytes
     bytesread, err = reader.Read(fourbytes)
     if err != nil && err != io.EOF {
        log.Println("error reading dwSamplesPerSec")
        return false
     }
     if 4 != bytesread {
        log.Println("not enough bytes read for dwSamplesPerSec")
     }
     // need to be 8 khz 
     samplerate := []byte{0x40, 0x1f, 0x00, 0x00}
     if !(AreEqualSlices(fourbytes,samplerate)) {
        log.Println("not 8khz sample rate")
        return false
     }
     
     //-- read dwAvgBytesPerSec 4bytes (ignore)
     bytesread, err = reader.Read(fourbytes)
     if err != nil && err != io.EOF {
        log.Println("error reading dwAvgBytesPerSec")
        return false
     }
     if 4 != bytesread {
        log.Println("not enough bytes read for dwAvgBytesPerSec")
     }

     //-- read wBlockAlign (ignore, since we convert)
     bytesread, err = reader.Read(twobytes)
     if err != nil && err != io.EOF {
        log.Println("error reading wBlockAlign")
        return false
     }
     if 2 != bytesread {
        log.Println("not enough bytes read for wBlockAlign")
        return false
     }

     //-- nBitsPerSample
     bytesread, err = reader.Read(twobytes)
     if err != nil && err != io.EOF {
        log.Println("error reading nBitsPerSample")
        return false
     }
     if 2 != bytesread {
        log.Println("not enough bytes read for nBitsPerSample")
        return false
     }
     // need to be 8bit
     if !(0x08 == twobytes[0] && 0x00 == twobytes[1]) {
        log.Println("not 8 bit sample depth")
        return false
     }

     return true
}

// convert byte to integer
func Byte2Int(data []byte) uint8 {
     if 1 != len(data) { panic("more than one byte given") }
     buf := bytes.NewBuffer(data)
     var value uint8
     binary.Read(buf, binary.LittleEndian, &value)
     return value
}

// convert 4 bytes to uint32
func Bytes2Uint32(data []byte) uint32 {
     if 4 != len(data) { panic("not 4 byte given") }
     buf := bytes.NewBuffer(data)
     var value uint32
     binary.Read(buf, binary.LittleEndian, &value)
     return value
}

// save data in file with given name
func ConvertData(reader *bufio.Reader, filename string) bool {
     // check for data chunk
     fourbytes := make([]byte,4)
     bytesread, err := reader.Read(fourbytes)
     if err != nil && err != io.EOF {
        log.Println("error reading data chunk ID")
        return false
     }
     if bytesread != 4 {
     	log.Println("not enough bytes read for data chunk ID")
	return false
     }
     datachunkid := []byte{0x64, 0x61, 0x74, 0x61}
     if !(AreEqualSlices(fourbytes,datachunkid)) {
     	log.Println("no data chunk ID found")
	return false
     }

     //-- get data size
     bytesread, err = reader.Read(fourbytes)
     if err != nil && err != io.EOF {
        log.Println("error reading data chunk size")
        return false
     }
     if bytesread != 4 {
        log.Println("not enough bytes read for data chunk size")
        return false
     }

     //-- write file
     // C-header foo
     headerfile, err := os.Create(filename)
     defer headerfile.Close()
     if err != nil { panic(err) }
     headerwriter := bufio.NewWriter(headerfile)
     // write array length (fourbytes -> 32bits unsigned)
     headerwriter.WriteString("const long pcm_length=")
     headerwriter.WriteString(Bytes2Uint32(fourbytes))
     headerwriter.WriteString(";\n")
     headerwriter.WriteString("const unsigned char pcm_samples[] PROGMEM ={")

     databyte := make([]byte,1)
     bytesread, err = reader.Read(databyte)

     for err != io.EOF {
         // convert number string to integer
         //currentInt, err = strconv.Atoi(numberstring)
         //if err != nil { log.Fatal(err) }
         // append to dataslice
         //dataslice = append(dataslice, currentInt)
         //line, err = reader.ReadString('\n')
         //numberstring = strings.TrimSpace(line)
     }
     return true
}

func main() {
     // check commandline-arguments
     if 2 != len(os.Args) { panic("wrong number of arguments") }
     // open file, get error
     wavefile, err := os.Open(os.Args[1])
     // if error -> panic
     if err != nil { panic(err) }
     // close file when main returns
     defer wavefile.Close()
     // create new reader from file
     reader := bufio.NewReader(wavefile)
     // check file header 
     if !IsWavefileHeaderOk(reader) { panic("file header b0rked") }
     fmt.Print("file format is ")
     if !IsWavefileFormatOk(reader) {
        fmt.Print("not ")
     }
     fmt.Println("ok.")
}