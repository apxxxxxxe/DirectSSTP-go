package main

/*
#cgo LDFLAGS: -L. -lopener
#include <stdio.h>
#include <stdlib.h>
#include "opener.h"
#include <opener.h>
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var rep = regexp.MustCompile(`ssp_fmo_header_[a-z0-9]{8}_[a-z0-9]{8}\.(.*)` + string(0x1) + "(.*)")

type copyDataStruct struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

func decodeSjis(str string) (string, error) {
	iostr := strings.NewReader(str)
	rio := transform.NewReader(iostr, japanese.ShiftJIS.NewDecoder())
	ret, err := ioutil.ReadAll(rio)
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func getHwndList(fmo string) [][]string {
	summery := [][]string{}
	data := []string{}
	for _, line := range strings.Split(fmo, "\r\n") {
		if line != "" {
			match := rep.FindStringSubmatch(line)
			if match[1] == "path" {
				data = []string{"", ""}
			} else if match[1] == "hwnd" {
				data[1] = match[2]
			} else if match[1] == "fullname" {
				data[0] = match[2]
				summery = append(summery, data)
			}
		}
	}
	return summery
}

func main() {
	fmo, err := decodeSjis(C.GoString(C.getFMO()))
	if err != nil {
		panic(err)
	}

	summery := getHwndList(fmo)

	for i, d := range summery {
		fmt.Println(i, d)
	}

	hwndInt, err := strconv.Atoi(summery[0][1])
	if err != nil {
		panic(err)
	}

	hwnd := win.HWND(hwndInt)

	msgAry := []string{
		"SEND SSTP/1.4",
		"Sender: おくりぬし",
		`Script: \h\s0あーあーあー。\e`,
		"HWnd: 1024",
		"Charset: UTF-8",
	}

	msg := ""
	for _, m := range msgAry {
		msg += m + "\n"
	}
	msgBytes := []byte(msg)

	cds := copyDataStruct{
		dwData: 9801,
		cbData: uint32(((*reflect.StringHeader)(unsafe.Pointer(&msgBytes))).Len),
		lpData: ((*reflect.StringHeader)(unsafe.Pointer(&msgBytes))).Data,
	}

	ret := win.SendMessage(hwnd, win.WM_COPYDATA, 0, uintptr(unsafe.Pointer(&cds)))
	if ret == 0 {
		fmt.Println("WM_COPYDATA failed")
	}

}
