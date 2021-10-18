package main

/*
#cgo CFLAGS: -I./c
#cgo LDFLAGS: -L./c -ltexture2dstudio -lstdc++
#include "texture2d.h"
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"
)

/**
需要自己导入 https://github.com/smalls0098/texture2dStudio
使用c++的库，得在LDFLAGS中加入-lstdc++
如果使用dll，LDFLAGS需要加入全名除扩展名libtexturedecoder，如果使用a、dylib、so等只需要texturedecoder前面的lib不需要引入(这里我测试了win是需要lib，使用的是.a)
*/

func main() {
	concurrent()
	//etc2()
	//astc()
	//etc2WithFile()
	//astcWithFile()
	//
	//etc2png()
	//etc2pngToFile()
	//astc2png()
	//astc2pngToFile()
}

func concurrent() {
	for i := 0; i < 1000; i++ {
		fmt.Printf("etc2png:%d\n", i)
		etc2png()
	}
	for i := 0; i < 1000; i++ {
		fmt.Printf("etc2:%d\n", i)
		etc2()
	}
}

//==============================使用文件路径 测试文件压缩

func etc2WithFile() {
	png := C.CString("tests/img/img.png")
	ktx := C.CString("tests/compress/etcfile.ktx")
	C.CompressEtc2RGBWithFile(png, ktx, 0, 0, 1)
}

func astcWithFile() {
	png := C.CString("tests/img/img.png")
	astc := C.CString("tests/compress/astcfile.astc")
	C.CompressAstcWithFile(png, astc, 0, 8, 8, 1)
}

//==============================字节流 测试文件压缩

func etc2() {
	input := "tests/img/img.png"
	output := "tests/compress/ktx.ktx"
	pngBy, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	pngPtr := (*C.uchar)(unsafe.Pointer(&pngBy[0]))
	var pngSize C.ulong = C.ulong(len(pngBy))
	var ktxPtr *C.uchar
	var ktxSize C.ulong
	C.CompressEtc2RGB(pngPtr, pngSize, 0, 0, 5, 1, &ktxPtr, &ktxSize)
	if ktxPtr == nil {
		panic("ktx ptr is nil")
	}
	ktx := UintptrToBytesBySize(uintptr(unsafe.Pointer(ktxPtr)), int(ktxSize))
	fmt.Println(len(ktx))

	if err := PathExists(output); err == nil {
		err := os.Remove(output)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(ktx)
	if err != nil {
		panic(err)
	}
}

func astc() {
	input := "tests/img/img.png"
	output := "tests/compress/astc.astc"
	pngBy, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}
	pngPtr := (*C.uchar)(unsafe.Pointer(&pngBy[0]))
	var pngSize C.ulong = C.ulong(len(pngBy))
	var astcPtr *C.uchar
	var astcSize C.ulong
	C.CompressAstc(pngPtr, pngSize, 0, 8, 8, 1, 1, &astcPtr, &astcSize)
	if astcPtr == nil {
		panic("astc ptr is nil")
	}
	astc := UintptrToBytesBySize(uintptr(unsafe.Pointer(astcPtr)), int(astcSize))
	fmt.Println(len(astc))

	if err := PathExists(output); err == nil {
		err := os.Remove(output)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(astc)
	if err != nil {
		panic(err)
	}
}

//==============================字节流 测试文件解压

func etc2png() {
	input := "tests/bin/ktx.ktx"
	output := "tests/decompress/ktx.png"
	ktxBy, _ := ioutil.ReadFile(input)
	ktxPtr := (*C.uchar)(unsafe.Pointer(&ktxBy[0]))
	var png *C.uchar
	var pngSize C.ulong
	C.DecompressEtc2(ktxPtr, 1024, 1024, &png, &pngSize)
	etc := UintptrToBytesBySize(uintptr(unsafe.Pointer(png)), int(pngSize))
	fmt.Println(len(etc))

	if err := PathExists(output); err == nil {
		err := os.Remove(output)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(etc)
	if err != nil {
		panic(err)
	}
}

//==============================文件路径 测试文件解压

func etc2pngToFile() {
	input := "tests/bin/ktx.ktx"
	output := "tests/decompress/ktxfile.png"
	ktxBy, _ := ioutil.ReadFile(input)
	ktxPtr := (*C.uchar)(unsafe.Pointer(&ktxBy[0]))
	C.DecompressEtc2ToFile(ktxPtr, 1024, 1024, C.CString(output))
}

func astc2png() {
	input := "tests/bin/astc.astc"
	output := "tests/decompress/astc.png"
	ktxBy, _ := ioutil.ReadFile(input)
	ktxPtr := (*C.uchar)(unsafe.Pointer(&ktxBy[0]))
	var png *C.uchar
	var pngSize C.ulong
	C.DecompressAstc(ktxPtr, 1024, 1024, 8, 8, &png, &pngSize)
	etc := UintptrToBytesBySize(uintptr(unsafe.Pointer(png)), int(pngSize))
	fmt.Println(len(etc))

	if err := PathExists(output); err == nil {
		err := os.Remove(output)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(etc)
	if err != nil {
		panic(err)
	}
}

func astc2pngToFile() {
	input := "tests/bin/astc.astc"
	output := "tests/decompress/astcfile.png"
	ktxBy, _ := ioutil.ReadFile(input)
	ktxPtr := (*C.uchar)(unsafe.Pointer(&ktxBy[0]))
	C.DecompressAstcToFile(ktxPtr, 1024, 1024, 8, 8, C.CString(output))
}

func UintptrToBytes(u uintptr) []byte {
	// 获取C返回的指针。
	// 注意C返回的r为char*，对应的Go类型为*byte
	p := (*byte)(unsafe.Pointer(u))
	// 定义一个[]byte切片，用来存储C返回的字符串
	data := make([]byte, 0)
	// 遍历C返回的char指针，直到 '\0' 为止
	for *p != 0 {
		data = append(data, *p)        // 将得到的byte追加到末尾
		u += unsafe.Sizeof(byte(0))    // 移动指针，指向下一个char
		p = (*byte)(unsafe.Pointer(u)) // 获取指针的值，此时指针已经指向下一个char
	}
	return data
}

func UintptrToBytesBySize(u uintptr, size int) []byte {
	p := (*byte)(unsafe.Pointer(u))
	data := make([]byte, size)
	// 遍历C返回的char指针，直到 '\0' 为止
	for i := 0; i < size; i++ {
		data[i] = *p                   // 将得到的byte追加到末尾
		u += unsafe.Sizeof(byte(0))    // 移动指针，指向下一个char
		p = (*byte)(unsafe.Pointer(u)) // 获取指针的值，此时指针已经指向下一个char
	}
	return data
}

func PathExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return err
	}
	return nil
}
