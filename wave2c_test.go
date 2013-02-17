package main

import(
	"bufio"
	// to open files
	"os"
	// import testing package
	"testing"
)

// known good test
func Test_AreEqualSlices(t *testing.T) {
     // create two equal slices
     a := []byte{ 0xDE, 0xAD, 0xBE, 0xEF }
     b := []byte{ 0xDE, 0xAD, 0xBE, 0xEF }
     // make test -> log failure but continue testing
     if !AreEqualSlices(a,b) { t.Error("equal slices determined unequal") }
}

// one slice shorter than the other
func Test_AreEqualSlices_one_shorter(t *testing.T) {
     // create two equal slices
     a := []byte{ 0xDE, 0xAD, 0xBE, 0xEF }
     b := []byte{ 0xDE, 0xAD, 0xBE }
     //	make test -> log failure but continue testing
     if AreEqualSlices(a,b) { t.Error("different length slices determined equal") }
}

// two same length slices with different content
func Test_AreEqualSlices_unequal(t *testing.T) {
     // create two equal slices
     a := []byte{ 0xDE, 0xAD, 0xBE, 0xEF }
     b := []byte{ 0xCA, 0xFE, 0xBA, 0xBE }
     // make test -> log failure but continue testing
     if AreEqualSlices(a,b) { t.Error("unequal slices determined equal") }
}

// known good test 11k
func Test_IsWavefileHeaderOk_11k(t *testing.T) {
     wavefile, err := os.Open("11k16bitpcm.wav")
     // if error -> panic
     if err != nil { panic(err) }
     // close file when main returns
     defer wavefile.Close()
     // create reader
     filereader := bufio.NewReader(wavefile)
     if !IsWavefileHeaderOk(filereader) { t.Error("valid Wavefile header misinterpreted") }
}

// known good
func Test_IsWavefileHeaderOk_8bit(t *testing.T) {
     wavefile, err := os.Open("validinput.wav")
     // if error -> panic
     if err != nil { panic(err) }
     // close file when main returns
     defer wavefile.Close()
     // create new reader from file
     reader := bufio.NewReader(wavefile)
     if !IsWavefileHeaderOk(reader) { t.Error("valid Wavefile header misinterpreted") }
}

// open random data
func Test_IsWavefileHeaderOk_rnddata(t *testing.T) {
     wavefile, err := os.Open("rnddata.wav")
     // if error -> panic
     if err != nil { panic(err) }
     // close file when main returns
     defer wavefile.Close()
     filereader := bufio.NewReader(wavefile)
     if IsWavefileHeaderOk(filereader) { t.Error("random data mistaken for valid Wavefile") }
}

// known good
func Test_IsWavefileHeaderOk(t *testing.T) {
     wavefile, err := os.Open("validinput.wav")
     // if error -> panic
     if err != nil { panic(err) }
     // close file when main returns
     defer wavefile.Close()
     // create new reader from file
     reader := bufio.NewReader(wavefile)
     if !IsWavefileHeaderOk(reader) { t.Error("valid Wavefile header misinterpreted") }
}

// get 0
func Test_Byte2IntOk_0(t *testing.T) {
     testbyte := []byte{0x00}
     if 0 != Byte2Int(testbyte) {
     	t.Error("invalid conversation from byte 0x00")
     }
}

// check conversation of 128
func Test_Byte2IntOk_128(t *testing.T) {
     testbyte := []byte{0x80}
     if 128 != Byte2Int(testbyte) {
        t.Error("invalid conversation from byte 0xff, returned ",Byte2Int(testbyte))
     }
}

// check conversation of 255
func Test_Byte2IntOk_255(t *testing.T) {
     testbyte := []byte{0xff}
     if 255 != Byte2Int(testbyte) {
        t.Error("invalid conversation from byte 0xff, returned ", Byte2Int(testbyte))
     }
}

/*/ convert too long slice
func Test_Byte2Int_too_long(t *testing.T) {
     testbyte := []byte{0x00, 0xff}
     Byte2Int(testbyte) {
        t.Error("invalid conversation from byte to int")
     }
}
*/

// check conversation
func Test_Bytes2uint32_ok_00000000(t *testing.T) {
     testbytes := []byte{0x00, 0x00, 0x00, 0x00}
     if 0 != Bytes2Uint32(testbytes) {
        t.Error("invalid byte conversation")
     }
}

// check conversation
func Test_Bytes2uint32_ok_ff00ff00(t *testing.T) {
     testbytes := []byte{0xff, 0x00, 0xff, 0x00}
     if 16711935 != Bytes2Uint32(testbytes) {
        t.Error("invalid byte conversation")
     }
}