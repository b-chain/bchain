////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: main.go
// @Date: 2018/12/05 14:03:05
//
////////////////////////////////////////////////////////////////////////////////

package main
/*
import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/validate"
	"github.com/go-interpreter/wagon/wasm"
	"io"
	"bchain.io/log"
	"os"
	"time"
	"reflect"
)

var (
	logTag = "core"
	logger log.Logger
)

func init() {
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}

func main() {
	verbose := flag.Bool("v", false, "enable/disable verbose mode")
	verify := flag.Bool("verify-module", false, "run module verification")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	wasm.SetDebugMode(*verbose)

	run(os.Stdout, flag.Arg(0), flag.Arg(1),  *verify)
}

func run(w io.Writer, fname string, entryFunc string, verify bool) {
	f, err := os.Open(fname)
	if err != nil {
		logger.Critical(err)
	}
	defer f.Close()

	//m, err := wasm.ReadModule(f, importer)
	m, err := wasm.ReadModule(f, importer)
	if err != nil {
		logger.Criticalf("could not read module: %v", err)
	}

	if verify {
		err = validate.VerifyModule(m)
		if err != nil {
			logger.Criticalf("could not verify module: %v", err)
		}
	}

	if m.Export == nil {
		logger.Criticalf("module has no export section")
	}
	now := time.Now()
	vm, err := exec.NewVM(m)
	if err != nil {
		logger.Criticalf("could not create VM: %v", err)
	}
	fmt.Println(time.Since(now))

	if xx, ok := m.Export.Entries[entryFunc]; ok {
		args := make([]uint64, 0)
		i := int64(xx.Index)
		fidx := m.Function.Types[int(i)]
		ftype := m.Types.Entries[int(fidx)]
		switch len(ftype.ReturnTypes) {
		case 1:
			fmt.Fprintf(w, "%s() %s \n", entryFunc, ftype.ReturnTypes[0])
		case 0:
			fmt.Fprintf(w, "%s() \n ", entryFunc)
		default:
			logger.Info("running exported functions with more than one return value is not supported")
		}
		lenPara := len(ftype.ParamTypes)
		if lenPara > 0 {
			args = append(args, 64)
			binary.LittleEndian.PutUint32(vm.Memory()[64:], 44)
			binary.LittleEndian.PutUint32(vm.Memory()[68:], 55)
			binary.LittleEndian.PutUint32(vm.Memory()[72:], 76)
			copy(vm.Memory()[76:], []byte("xxyy\n"))
		}
		o, err := vm.ExecCode(i, args...)
		if err != nil {
			fmt.Fprintf(w, "\n")
			logger.Info("err=%v", err)
			return
		}
		if len(ftype.ReturnTypes) == 0 {
			fmt.Fprintf(w, "\n")
			return
		}
		fmt.Fprintf(w, "%[1]v (%[1]T)\n", o)
	}
	fmt.Println(time.Since(now))
}

func printf(proc *exec.Process, msg_ptr int32) int32 {
	mem := proc.Vm().Memory()
	if int32(len(mem)) <= msg_ptr {
		panic("msg pointer is exceed")
	}
	fmt.Printf(string(mem[msg_ptr:msg_ptr+256]))

	return 0
}
func importer(name string) (*wasm.Module, error) {
	switch name {
	case "include":
		f, err := os.Open(name + ".wasm")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		m, err := wasm.ReadModule(f, nil)
		if err != nil {
			return nil, err
		}
		return m, nil
	case "env":
		m := wasm.NewModule()
		m.Types = &wasm.SectionTypes{
			Entries: []wasm.FunctionSig{
				{
					Form:       0, // value for the 'func' type constructor
					ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
					ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
				},
			},
		}
		m.FunctionIndexSpace = []wasm.Function{
			{
				Sig:  &m.Types.Entries[0],
				Host: reflect.ValueOf(printf),
				Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
			},
		}
		m.Export = &wasm.SectionExports{
			Entries: map[string]wasm.ExportEntry{
				"log": {
					FieldStr: "log",
					Kind:     wasm.ExternalFunction,
					Index:    0,
				},
			},
		}
		return m, nil
	}
	return nil, fmt.Errorf("module %q unknown", name)
}
*/