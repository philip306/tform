package main

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/assert"
)

func TestParseMain(t *testing.T) {
	got := ParseMain(".\\testfiles\\proj1\\")
	want := "~> 1.0.0"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestParseMain2(t *testing.T) {
	got := ParseMain(".\\testfiles\\proj2\\")
	want := ">= 0.12"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestParseMain3(t *testing.T) {
	got := ParseMain(".\\testfiles\\proj3\\")
	want := "0.13"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestCheckInstVersion(t *testing.T) {

	insVer, err := CheckInstVersion(".\\testfiles\\lib1\\")
	if err != nil {
		t.Error(err)
	}
	var got []string
	for _, ver := range insVer {
		got = append(got, ver.String())
	}

	want := []string{"0.12.24", "0.12.26", "0.14.0", "0.14.11"}

	if reflect.DeepEqual(got, want) != true {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestCheckInstVersion2(t *testing.T) {

	insVer, err := CheckInstVersion(".\\testfiles\\lib2\\")
	if err != nil {
		t.Error(err)
	}
	got := []string{}
	for _, ver := range insVer {
		got = append(got, ver.String())
	}

	want := []string{}

	if reflect.DeepEqual(got, want) != true {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestCheckInstVersion3(t *testing.T) {
	var got string
	insVer, err := CheckInstVersion(".\\foo\\bar\\")
	if err != nil {
		got = err.Error()
	} else {
		t.Errorf("got %q, wanted error", insVer)
	}

	want := "open .\\foo\\bar\\: The system cannot find the path specified."

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestCheckInstVersion4(t *testing.T) {
	var got string
	insVer, err := CheckInstVersion(".\\testfiles\\lib3\\")
	if err != nil {
		got = err.Error()
	} else {
		t.Errorf("got %q, wanted an error", insVer)
	}

	want := "Malformed version: not.a.dir"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestExists(t *testing.T) {
	var got bool
	check := Exists(".\\testfiles\\lib1")
	if check {
		got = check
	} else {
		t.Errorf("got %t, wanted true", got)
	}

	want := true
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
}

func TestExists2(t *testing.T) {
	var got bool
	check := Exists(".\\foo\\bar")
	if check {
		t.Errorf("got %t, wanted false", got)
	} else {
		got = check
	}

	want := false
	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
}

func TestChooseVer(t *testing.T) {
	var got string
	v, err := ChooseVer(".\\testfiles\\proj1\\", ".\\testfiles\\lib4")
	if err != nil {
		t.Errorf("got %q, wanted 1.0.5", err)
	} else {
		got = v
	}

	want := "1.0.5"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestChooseVer2(t *testing.T) {
	var got string
	v, err := ChooseVer(".\\testfiles\\proj4\\", ".\\testfiles\\lib4\\")
	if err != nil {
		got = err.Error()
	} else {
		t.Errorf("got %q, wanted error", v)
	}

	want := "Malformed constraint: ~> 1.foo.bar"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestChooseVer3(t *testing.T) {
	var got string
	v, err := ChooseVer(".\\testfiles\\proj1\\", ".\\testfiles\\lib1\\")
	if err != nil {
		got = err.Error()
	} else {
		t.Errorf("got %q, wanted error", v)
	}

	want := "No version satisfying constraint: ~> 1.0.0 found on system.  Check your chocolately lib or try 'choco install terraform --version [VERSION] -my' and try again."

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestChooseVer4(t *testing.T) {
	var got string
	v, err := ChooseVer(".\\testfiles\\proj1\\", ".\\foo\\bar\\")
	if err != nil {
		got = err.Error()
	} else {
		t.Errorf("got %q, wanted error", v)
	}

	want := "open .\\foo\\bar\\: The system cannot find the path specified."

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestChooseVer5(t *testing.T) {
	var got string
	v, err := ChooseVer(".\\foo\\bar\\", ".\\testfiles\\proj1\\")
	if err != nil {
		got = err.Error()
	} else {
		t.Errorf("got %q, wanted error", v)
	}

	want := "File path: .\\foo\\bar\\ not found on system"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestExecuteTform(t *testing.T) {
	var got error
	e := ExecuteTform(".\\testfiles\\proj1", ".\\testfiles\\lib4", true)
	if e != nil {
		t.Errorf("got %q, wanted nil", e.Error())
	} else {
		got = e
	}

	want := error(nil)

	if !errors.Is(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestExecuteTform2(t *testing.T) {
	var got string
	e := ExecuteTform(".\\foo\\bar", ".\\testfiles\\lib4\\", true)
	if assert.Errorf(t, e, "error") {
		got = e.Error()
	} else {
		t.Errorf("got %q, wanted an error", e)
	}

	want := "File path: .\\foo\\bar not found on system"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestExecuteTform3(t *testing.T) {
	var got error
	e := ExecuteTform(".\\testfiles\\proj1", ".\\testfiles\\lib4", false)
	if e != nil {
		t.Errorf("got %q, wanted nil", e.Error())
	} else {
		got = e
	}

	want := error(nil)

	if !errors.Is(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestMain(t *testing.T) {
	err := os.Chdir(".\\testfiles\\proj1")
	if err != nil {
		panic(err)
	}
	main()
}

func TestMain2(t *testing.T) {
	err := os.Chdir("..\\..\\testfiles\\proj4")
	if err != nil {
		panic(err)
	}

	if os.Getenv("BE_CRASHER") == "1" {
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMain2")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
