// Copyright 2018 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package disasm_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/go-interpreter/wagon/disasm"
	"github.com/go-interpreter/wagon/wasm"
)

var testPaths = []string{
	"../wasm/testdata",
	"../exec/testdata",
	"../exec/testdata/spec",
}

func TestAssemble(t *testing.T) {
	for _, dir := range testPaths {
		fnames, err := filepath.Glob(filepath.Join(dir, "*.wasm"))
		if err != nil {
			t.Fatal(err)
		}
		for _, fname := range fnames {
			name := fname
			t.Run(filepath.Base(name), func(t *testing.T) {
				raw, err := ioutil.ReadFile(name)
				if err != nil {
					t.Fatal(err)
				}

				r := bytes.NewReader(raw)
				m, err := wasm.DecodeModule(r)
				if err != nil {
					t.Fatalf("error reading module %v", err)
				}
				if m.Code == nil {
					t.SkipNow()
				}
				for _, f := range m.Code.Bodies {
					d, err := disasm.Disassemble(f.Code)
					if err != nil {
						t.Fatalf("disassemble failed: %v", err)
					}
					code, err := disasm.Assemble(d)
					if err != nil {
						t.Fatalf("assemble failed: %v", err)
					}
					if !bytes.Equal(f.Code, code) {
						t.Fatal("code is different")
					}
				}
			})
		}
	}
}
