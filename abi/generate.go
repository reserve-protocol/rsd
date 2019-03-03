//+build ignore

// This command generates Go bindings for the Reserve Dollar smart contracts.
//
// It is intended to be invoked with `make abi/bindings` at the root of the repo.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func main() {
	// Generate bindings from the compiled artifacts.
	artifactPaths, _ := filepath.Glob(filepath.FromSlash("artifacts/*.json"))
	for _, path := range artifactPaths {
		// Read artifact file.
		b, err := ioutil.ReadFile(path)
		check(err, "reading %q", path)

		// Extract name, ABI, and bytecode.
		var compiled struct {
			ContractName   string
			CompilerOutput struct {
				ABI json.RawMessage
				EVM struct {
					Bytecode struct {
						Object string
					}
				}
			}
		}
		check(json.Unmarshal(b, &compiled), "parsing %v", path)

		// Generate bindings.
		code, err := bind.Bind(
			[]string{compiled.ContractName},
			[]string{string(compiled.CompilerOutput.ABI)},
			[]string{compiled.CompilerOutput.EVM.Bytecode.Object},
			"abi",
			bind.LangGo,
		)
		check(err, "generating bindings for %q", compiled.ContractName)

		// Write to .go file.
		check(
			ioutil.WriteFile(
				filepath.Join("abi", compiled.ContractName+".go"),
				[]byte(code),
				0644,
			),
			"writing bindings to disk",
		)

		// Generate event bindings.
		//
		// func ParseReserveDollarLog(log *types.Log) (interface{}, error) {
		//	...
		// }
		//
		// func (f ReserveDollarFrozen) String() string {
		//	return fmt.Sprintf("ReserveDollarFrozen(%v, %v)", f.Freezer.Hex(), f.Account.Hex())
		// }
		buf := new(bytes.Buffer)
		parsedABI, err := abi.JSON(bytes.NewReader(compiled.CompilerOutput.ABI))
		check(err, "parsing ABI JSON")
		check(template.Must(template.New("").Funcs(
			template.FuncMap{
				"flags": func(inputs abi.Arguments) string {
					result := make([]string, len(inputs))
					for i := range result {
						switch inputs[i].Type.String() {
						case "string":
							result[i] = "%q"
						default:
							result[i] = "%v"
						}
					}
					return strings.Join(result, ", ")
				},
				"format": func(inputs abi.Arguments) string {
					result := make([]string, len(inputs))
					for i := range result {
						arg := "e." + abi.ToCamelCase(inputs[i].Name)
						switch inputs[i].Type.String() {
						case "address":
							arg = arg + ".Hex()"
						}
						result[i] = arg
					}
					return strings.Join(result, ",")
				},
			},
		).Parse(`
			// This file is auto-generated. Do not edit.

			package abi

			import (
				"fmt"

				"github.com/ethereum/go-ethereum/core/types"
			)

			{{$contract := .Contract}}

			func (c *{{$contract}}Filterer) ParseLog(log *types.Log) (fmt.Stringer, error) {
				var event fmt.Stringer
				var eventName string
				switch log.Topics[0].Hex() {
				{{- range .Events}}
				case {{with .Id}}{{printf "%q" .Hex}}{{end}}: // {{.Name}}
					event = new({{$contract}}{{.Name}})
					eventName = "{{.Name}}"
				{{- end}}
				default:
					return nil, fmt.Errorf("no such event hash for {{$contract}}: %v", log.Topics[0])
				}

				err := c.contract.UnpackLog(event, eventName, *log)
				if err != nil {
					return nil, err
				}
				return event, nil
			}

			{{range .Events}}
			func (e {{$contract}}{{.Name}}) String() string {
				return fmt.Sprintf("{{$contract}}.{{.Name}}({{flags .Inputs}})",{{format .Inputs}})
			}
			{{end}}
		`)).Execute(buf, map[string]interface{}{
			"Contract": compiled.ContractName,
			"Events":   parsedABI.Events,
		}), "generating event bindings")
		eventCode, err := format.Source(buf.Bytes())
		check(err, "running gofmt")
		check(
			ioutil.WriteFile(
				filepath.Join("abi", compiled.ContractName+"Events.go"),
				eventCode,
				0644,
			),
			"writing event bindings to disk",
		)
	}
}

// check exits the program with a formatted error message if err is non-nil.
// If err is nil, check does nothing.
func check(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", fmt.Sprintf(format, args...), err)
		os.Exit(1)
	}
}
