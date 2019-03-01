//+build ignore

// This command generates Go bindings for the Reserve Dollar smart contracts.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func main() {
	// Run sol-compiler.
	{
		cmd := exec.Command("npx", "sol-compiler")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		check(cmd.Run(), "running sol-compiler")
	}

	// Generate bindings from the compiled artifacts.
	check(filepath.Walk("artifacts", func(path string, info os.FileInfo, err error) error {
		check(err, "walking artifacts/ directory")

		if info.IsDir() {
			return nil
		}

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

		return nil
	}), "filepath.Walk")
}

// check exits the program with a formatted error message if err is non-nil.
// If err is nil, check does nothing.
func check(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", fmt.Sprintf(format, args...), err)
		os.Exit(1)
	}
}
