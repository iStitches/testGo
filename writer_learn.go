package io

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Abc struct {
	Id int
}

func (a *Abc) Write(p []byte) (n int, err error) {
	return 1, nil
}

// 写文件
func WriteToFile(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeTemporary)
	if err != nil {
		return "文件打开错误", err
	}
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	num, err := fmt.Fprintln(file, input)
	if err != nil {
		return "文件写入失败", err
	}
	return string(num), nil
}

// 将文件A内容复制到文件B中
func CopyFile(pathA string, pathB string) (string, error) {
	fileSrc, err := os.Open(pathA)
	if err != nil {
		fileSrc.Close()
		return "来源文件打开失败,请检查!", err
	}
	fileTo, err := os.OpenFile(pathB, os.O_CREATE|os.O_APPEND, os.ModeTemporary)
	if err != nil {
		fileTo.Close()
		return "目的文件打开失败,请检查!", err
	}
	writer := bufio.NewWriter(fileTo)
	writer.ReadFrom(fileSrc)
	writer.Flush()
	return "文件复制完毕!", nil
}

func TestSeek() {
	reader := strings.NewReader("Go语言中文网")
	reader.Seek(-6, io.SeekEnd)
	r, _, _ := reader.ReadRune()
	fmt.Println(string(r))
}

func TestBytes() {
	var ch byte
	fmt.Scanf("%c\n", &ch)
	buffer := new(bytes.Buffer)
	err := buffer.WriteByte(ch)
	if err == nil {
		fmt.Println("成功写入一个字节，准备读取该字节！")
		newCh, _ := buffer.ReadByte()
		fmt.Printf("读取到的字节：%c\n", newCh)
	} else {
		fmt.Println("写入错误！")
	}
}

func TestSectionReader() {
	fileReader := strings.NewReader("Go语言中文网")
	reader := io.NewSectionReader(fileReader, 3, 2)
	tmp := make([]byte, 10)
	reader.Read(tmp)
	fmt.Println(string(tmp))
}

func TestLimitedReader() {
	content := "This Is LimitReader Example"
	reader := strings.NewReader(content)
	limitReader := &io.LimitedReader{R: reader, N: 8}
	for limitReader.N > 0 {
		tmp := make([]byte, 8)
		limitReader.Read(tmp)
		fmt.Println(string(tmp))
	}
}

func TestPipe() {
	pipeReader, pipeWriter := io.Pipe()
	go PipeWrite(pipeWriter)
	go PipeRead(pipeReader)
	time.Sleep(30 * time.Second)
}
func PipeWrite(writer *io.PipeWriter) {
	data := []byte("Go语言中文网")
	for i := 0; i < 3; i++ {
		n, err := writer.Write(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("写入字节 %d\n", n)
	}
	writer.CloseWithError(errors.New("写入端关闭..."))
}

func PipeRead(reader *io.PipeReader) {
	buf := make([]byte, 128)
	for {
		fmt.Println("接受端开始阻塞5秒")
		time.Sleep(5 * time.Second)
		fmt.Println("接收端开始接收")
		n, err := reader.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("收到字节：%d\n buf 内容：%s\n", n, buf)
	}
}

func TestCopy() {
	srcFile, err := os.Open("D:\\a.txt")
	dstFile, _ := os.OpenFile("D:\\b.txt", os.O_CREATE|os.O_APPEND, os.ModeTemporary)
	if err == nil {
		srcReader := bufio.NewReader(srcFile)
		dstWriter := bufio.NewWriter(dstFile)
		n, err := io.Copy(dstWriter, srcReader)
		if err != nil {
			fmt.Println("复制文件时出现了错误！")
		}
		fmt.Printf("一共复制了 %d 个字节", n)
	}
}

func TestMultiReader() {
	readers := []io.Reader{
		strings.NewReader("from strings reader"),
		bytes.NewBufferString("from bytes buffer"),
	}
	reader := io.MultiReader(readers...)
	data := make([]byte, 0, 128)
	buf := make([]byte, 10)
	for n, err := reader.Read(buf); err != io.EOF; n, err = reader.Read(buf) {
		if err != nil {
			panic(err)
		}
		data = append(data, buf[:n]...)
	}
	fmt.Printf("%s\n", data)
}

func TestMultiWriter() {
	file, err := os.Create("02.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	writer := io.MultiWriter(writers...)
	writer.Write([]byte("好好学习Golang!"))
}

// 测试 ioutil.ReadDir() 读取目录下的文件
func ListAll(path string, curHier int) {
	file, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range file {
		if info.IsDir() {
			for tmpHier := curHier; tmpHier > 0; tmpHier-- {
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name(), "\\")
			ListAll(path+"/"+info.Name(), curHier+1)
		} else {
			for tmpHier := curHier; tmpHier > 0; tmpHier-- {
				fmt.Printf("|\t")
			}
			fmt.Println(info.Name())
		}
	}
}

// 测试 ioutil.ReadAll() 读取所有内容
func TestReadAll(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := ioutil.ReadAll(file)
	if err == nil {
		fmt.Println(string(res))
	}
}

// 测试 ReadFile()、WriteFile()
func TestReadWriteFile(path string) {
	res, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功读取文件内容")
	err = ioutil.WriteFile("03.txt", res, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功写入文件")
}

// 测试ReadSlice
func TestReadSlice() {
	reader := bufio.NewReader(strings.NewReader("http://studygolang.com. \nIt is the home of gophers"))
	line, _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n", line)
	// 这里可以换上任意的 bufio 的 Read/Write 操作
	n, _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n", line)
	fmt.Println(string(n))
}

// 测试ReadByts
func TestReadBytes() {
	reader := bufio.NewReader(strings.NewReader("http://studygolang.com. \nIt is the home of gophers"))
	line, _ := reader.ReadBytes('\n')
	fmt.Printf("the line:%s\n", line)
	// 这里可以换上任意的 bufio 的 Read/Write 操作
	n, _ := reader.ReadBytes('\n')
	fmt.Printf("the line:%s\n", line)
	fmt.Println(string(n))
}

// 测试Peek
func TestPeek() {
	reader := bufio.NewReaderSize(strings.NewReader("http://studygolang.com.\t It is the home of gophers"), 14)
	go Peek(reader)
	go reader.ReadBytes('\t')
	time.Sleep(1e8)
}

func Peek(reader *bufio.Reader) {
	line, _ := reader.Peek(14)
	fmt.Printf("%s\n", line)
	time.Sleep(1)
	fmt.Printf("%s\n", line)
}

func TestScacner() {
	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//}
	//if err := scanner.Err(); err != nil {
	//	fmt.Fprintln(os.Stderr, "reading standard input:", err)
	//}
	var num int
	var ch string
	fmt.Scanf("%d", &num)
	fmt.Scanf("%s", &ch)
	fmt.Println(num, ch)
}

// 统计一段英文有多少个单词
func TestScannerSplit1() {
	const input = "This is The Golang Standard Library.\nWelcome you!"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(os.Stderr, "reading input:", err)
	}
	fmt.Println(count)
}

// 读取文件中的数据，一次读取一行
func TestScannerSplit2() {
	file, err := os.Create("scanner.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("http://studygolang.com.\nIt is the home of gophers.\nIf you are studying golang, welcome you!")
	// 将文件指针设置到文件开头
	file.Seek(0, os.SEEK_SET)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func TestBufferWriter() {
	file, err := os.Create("buffer_writer.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString("aaaabbbbcccsssddd")
	writer.Flush()
}
