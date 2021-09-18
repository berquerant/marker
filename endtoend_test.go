package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	b := newMarker(t)
	defer b.close()

	for _, tc := range []*struct {
		title       string
		fileName    string
		typeNames   []string
		methodNames []string
	}{
		{
			title:       "generate a method",
			fileName:    "method1.go",
			typeNames:   []string{"StatementExpr"},
			methodNames: []string{"IsExpr"},
		},
		{
			title:       "generate 2 methods",
			fileName:    "method2.go",
			typeNames:   []string{"AddExpr"},
			methodNames: []string{"IsExpr", "IsOperator"},
		},
		{
			title:       "generate 2x2 methods",
			fileName:    "method2x2.go",
			typeNames:   []string{"FlexExpr", "FixedExpr"},
			methodNames: []string{"IsExpr", "IsNode"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			b.compileAndRun(t, tc.fileName, tc.typeNames, tc.methodNames)
		})
	}
}

type marker struct {
	dir, marker string
}

func newMarker(t *testing.T) *marker {
	t.Helper()
	s := &marker{}
	s.init(t)
	return s
}

func (s *marker) init(t *testing.T) {
	t.Helper()
	dir, err := ioutil.TempDir("", "marker")
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dir, "marker")
	// build marker
	if err := run("go", "build", "-o", marker); err != nil {
		t.Fatal(err)
	}
	s.dir = dir
	s.marker = marker
}

func (s *marker) close() {
	os.RemoveAll(s.dir)
}

func (s *marker) compileAndRun(t *testing.T, fileName string, typeNames, methodNames []string) {
	t.Helper()
	src := filepath.Join(s.dir, fileName)
	if err := copyFile(src, filepath.Join("testdata", fileName)); err != nil {
		t.Fatal(err)
	}
	markerSrc := filepath.Join(s.dir, fmt.Sprintf("%s_marker.go", typeNames[0]))
	// run marker
	if err := run(s.marker, "-type", strings.Join(typeNames, ","), "-method", strings.Join(methodNames, ","), "-output", markerSrc, src); err != nil {
		t.Fatal(err)
	}
	// run testfile with generated file
	if err := run("go", "run", markerSrc, src); err != nil {
		t.Fatal(err)
	}
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(to, from string) error {
	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	_, err = io.Copy(toFile, fromFile)
	return err
}
