package node

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
)

func ArraySource(arr ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range arr {
			out <- v
		}
		close(out)
	}()
	return out
}

func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int)
	go func() {
		buf := make([]byte, 8)
		bytesSize := 0
		for {
			n, err := reader.Read(buf)
			bytesSize += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buf))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesSize >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		writer.Write(buf)
	}
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var arr []int
		for v := range in {
			arr = append(arr, v)
		}
		sort.Ints(arr)
		for _, v := range arr {
			out <- v
		}
		close(out)
	}()
	return out
}

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	mid := len(inputs) / 2
	return Merge(
		MergeN(inputs[:mid]...),
		MergeN(inputs[mid:]...))
}
