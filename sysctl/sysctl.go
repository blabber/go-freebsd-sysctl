// +build freebsd darwin
// +build cgo

// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

// Package sysctl provides a basic read-only interface to the FreeBSD sysctl(3)
// library.
//
// It allows you to retrieve int64 and string values. Tables are not supported.
//
//	GetInt64("hw.ncpu")		// get the number of active CPUs
//	GetString("kern.hostname")	// get the hostname
//
// You might want to use the "syscall" standard package that implements Sysctl*
// functions as well.
//
// On non-FreeBSD or non-Darwin systems this package will not build
package sysctl

import (
	"unsafe"
)

// #include <sys/types.h>
// #include <sys/sysctl.h>
import "C"

// GetInt64 gets a numeric value from sysctl(3) by name
func GetInt64(name string) (value int64, err error) {
	oldlen := C.size_t(8)
	_, err = C.sysctlbyname(C.CString(name), unsafe.Pointer(&value), &oldlen, nil, 0)
	if err != nil {
		return
	}
	return
}

// GetString gets a string value from sysctl(3) by name
func GetString(name string) (value string, err error) {
	oldlen := C.size_t(0)

	// Call C.sysctlbyname once to get the required size of the buffer.
	_, err = C.sysctlbyname(C.CString(name), nil, &oldlen, nil, 0)
	if err != nil {
		return
	}

	oldp := C.CString(string(make([]byte, uint32(oldlen))))
	_, err = C.sysctlbyname(C.CString(name), unsafe.Pointer(oldp), &oldlen, nil, 0)
	if err != nil {
		return
	}
	value = C.GoString(oldp)

	return
}
