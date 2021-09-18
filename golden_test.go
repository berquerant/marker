package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Golden struct {
	name, input, output string
}

const one = `type First struct {}
`

const two = `type First struct {}
type Second struct {}
`

func TestGolden(t *testing.T) {
	dir, err := ioutil.TempDir("", "marker")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range []*struct {
		title   string
		input   string
		output  string
		types   []string
		methods []string
	}{
		{
			title:   "one_first_1",
			input:   one,
			types:   []string{"First"},
			methods: []string{"isFirst"},
			output: `func (*First) isFirst() {}
`,
		},
		{
			title:   "one_first_2",
			input:   one,
			types:   []string{"First"},
			methods: []string{"isFirst", "hasHand"},
			output: `func (*First) isFirst() {}
func (*First) hasHand() {}
`,
		},
		{
			title:   "two_first_1",
			input:   two,
			types:   []string{"First"},
			methods: []string{"isFirst"},
			output: `func (*First) isFirst() {}
`,
		},
		{
			title:   "two_second_2",
			input:   two,
			types:   []string{"First", "Second"},
			methods: []string{"IsRed", "IsBlue"},
			output: `func (*First) IsRed()   {}
func (*First) IsBlue()  {}
func (*Second) IsRed()  {}
func (*Second) IsBlue() {}
`,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			input := fmt.Sprintf("package test\n%s", tc.input)
			file := fmt.Sprintf("%s.go", tc.title)
			absFile := filepath.Join(dir, file)
			if err := ioutil.WriteFile(absFile, []byte(input), 0644); err != nil {
				t.Error(err)
			}
			g := Generator{}
			g.parsePackage([]string{absFile})
			g.generateMulti(tc.types, tc.methods)
			got := string(g.format())
			assert.Equal(t, tc.output, got)
		})
	}
}
