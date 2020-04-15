package cleaningpipe

import (
	"io/ioutil"
	"strings"
	"testing"
)



var marker error
var markerbuffer []byte
func benchmarkCleanBodyPipe_Read(len int, b *testing.B) {
	test := getTestString(len)
	var err error
	var buffer []byte
	benchy := NewCleaningPipe(demoCleaner, strings.NewReader(test))
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buffer, err = ioutil.ReadAll(benchy)
	}
	markerbuffer = buffer
	marker = err
}

func benchmarkBaseLineReadAllAndReplace_Read(len int, b *testing.B) {
	test := getTestString(len)
	var err error
	var buffer []byte
	benchy := strings.NewReader(test)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buffer, err = ioutil.ReadAll(benchy)
		buffer = demoCleaner(buffer)
	}
	markerbuffer = buffer
	marker = err}

func BenchmarkCleanBodyPipe_Read100(b *testing.B)   { benchmarkCleanBodyPipe_Read(100, b) }
func BenchmarkCleanBodyPipe_Read1000(b *testing.B)   { benchmarkCleanBodyPipe_Read(1000, b) }
func BenchmarkCleanBodyPipe_Read100000(b *testing.B) { benchmarkCleanBodyPipe_Read(100000, b) }

func BenchmarkBaseLineReadAllAndReplace_Read100(b *testing.B)   { benchmarkBaseLineReadAllAndReplace_Read(100, b) }
func BenchmarkBaseLineReadAllAndReplace_Read1000(b *testing.B)   { benchmarkBaseLineReadAllAndReplace_Read(1000, b) }
func BenchmarkBaseLineReadAllAndReplace_Read100000(b *testing.B) { benchmarkBaseLineReadAllAndReplace_Read(100000, b) }
