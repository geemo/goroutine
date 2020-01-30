package main

import (
	"bufio"
	"fmt"
	"os"
	"pipeline/node"
	"strconv"
)

func main() {
	externalSort()
}

func externalSort() {
	// mergeDemo()
	// pipelineDemo()
	// p := createPipeline("small.in", 512, 4)
	// writeToFile("small.out", p)
	// printFile("small.out")

	p := createNetworkPipeline("small.in", 512, 4)
	writeToFile("small.out", p)
	printFile("small.out")
}

func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunSize := fileSize / chunkCount
	var sortRes []<-chan int
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunSize), 0)
		source := node.ReaderSource(bufio.NewReader(file), chunSize)
		// 内部排序完成，可以通过网络发送到下游节点
		sortRes = append(sortRes, node.InMemSort(source))
	}
	// 使用 ReaderSource 封装读取网络发送过来的排序好的数据
	// 将数据源进行归并
	return node.MergeN(sortRes...)
}

func createNetworkPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunSize := fileSize / chunkCount
	var sortAddr []string
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunSize), 0)
		source := node.ReaderSource(bufio.NewReader(file), chunSize)
		// 内部排序完成，可以通过网络发送到下游节点
		addr := ":" + strconv.Itoa(7000 + i)
		node.NetworkSink(addr, node.InMemSort(source))
		sortAddr = append(sortAddr, addr)
	}
	// 读取网络发送过来的排序好的数据
	var sortRes []<-chan int
	for _, addr := range sortAddr {
		sortRes = append(sortRes, node.NetworkSource(addr))
	}
	// 将数据源进行归并
	return node.MergeN(sortRes...)
}

func writeToFile(filename string, in <-chan int) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	node.WriterSink(file, in)
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := node.ReaderSource(file, -1)
	for v := range p {
		fmt.Println(v)
	}
}

func pipelineDemo() {
	const filename = "small.in"
	const n = 64
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	p := node.RandomSource(n)
	node.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = node.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}

func mergeDemo() {
	p := node.Merge(
		node.InMemSort(node.ArraySource(3, 2, 1, 4, 5, 7, 6)),
		node.InMemSort(node.ArraySource(5, 4, 3, 1, 2, 8, 100, 65)))
	for v := range p {
		fmt.Println(v)
	}
}
