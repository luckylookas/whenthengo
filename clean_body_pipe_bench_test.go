package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

var marker error

func benchmarkCleanBodyPipe_Read(len int, b *testing.B) {
	sb := strings.Builder{}
	for i:=0; i < len; i++ {
		sb.WriteString("aA\n")
	}
	test := sb.String()
	var err error
	benchy := CleanBodyPipe{strings.NewReader(test)}

	for n := 0; n < b.N; n++ {
		_, err = ioutil.ReadAll(benchy)
	}
	marker = err
}


func benchmarkNakedReader_Read(len int, b *testing.B) {
	sb := strings.Builder{}
	for i:=0; i < len; i++ {
		sb.WriteString("aa")
	}
	test := sb.String()
	var err error
	benchy := strings.NewReader(test)

	for n := 0; n < b.N; n++ {
		_, err = ioutil.ReadAll(benchy)
	}
	marker = err
}

func BenchmarkCleanBodyPipe_Read1000(b *testing.B)  { benchmarkCleanBodyPipe_Read(100, b) }
func BenchmarkCleanBodyPipe_Read100000(b *testing.B)  { benchmarkCleanBodyPipe_Read(100000, b) }

func BenchmarkNakedReader_Read1000(b *testing.B)  { benchmarkNakedReader_Read(100, b) }
func BenchmarkNakedReader_Read100000(b *testing.B)  { benchmarkNakedReader_Read(100000, b) }
