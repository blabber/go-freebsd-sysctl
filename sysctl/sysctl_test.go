// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package sysctl

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

const (
	sysctlCmd     = "/sbin/sysctl"
	sysctlCmdFlag = "-n"
)

const errorName = "non.existent"

var stringTests = []string{
	"kern.hostname",
	"kern.osrelease",
}

var int64Tests = []string{
	"hw.ncpu",
	"hw.realmem",
}

func execSysctl(name string) (out string, err error) {
	o, err := exec.Command(sysctlCmd, sysctlCmdFlag, name).Output()
	return strings.TrimRight(string(o), "\n"), err
}

func TestStrings(t *testing.T) {
	for _, name := range stringTests {
		t.Logf("-- Testing %q", name)

		expected, err := execSysctl(name)
		if err != nil {
			t.Fatalf("call to sysctl(8) failed: %v\n", err)
		}

		actual, err := GetString(name)
		if err != nil {
			t.Fatalf("call to GetString(%q) failed: %v\n", name, err)
		}

		t.Logf("%v %v %v => %q", sysctlCmd, sysctlCmdFlag, name, expected)
		t.Logf("GetString(%q) => %q", name, actual)

		if actual != expected {
			t.Fatalf("%q != %q", actual, expected)
		}
	}
}

func TestInt64s(t *testing.T) {
	for _, name := range int64Tests {
		t.Logf("-- Testing %q", name)

		expectedString, err := execSysctl(name)
		if err != nil {
			t.Fatalf("call to sysctl(8) failed: %v\n", err)
		}
		expected, err := strconv.ParseInt(expectedString, 10, 64)
		if err != nil {
			t.Fatalf("string conversion failed: %v\n", err)
		}

		actual, err := GetInt64(name)
		if err != nil {
			t.Fatalf("call to GetInt64(%q) failed: %v\n", name, err)
		}

		t.Logf("%v %v %v => %v", sysctlCmd, sysctlCmdFlag, name, expected)
		t.Logf("GetInt64(%q) => %v", name, actual)

		if actual != expected {
			t.Fatalf("%v != %v", actual, expected)
		}
	}
}

func TestErrorInt64(t *testing.T) {
	_, err := GetInt64(errorName)
	if err == nil {
		t.Fatalf("call to GetInt64(%q) succeeded without error", errorName)
	}
}

func TestErrorString(t *testing.T) {
	_, err := GetString(errorName)
	if err == nil {
		t.Fatalf("call to GetString(%q) succeeded without error", errorName)
	}
}
