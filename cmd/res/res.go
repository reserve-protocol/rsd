package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reserve-protocol/reserve-dollar/abi"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultKeys = []string{
	"f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d",
	"5d862464fe9303452126c8bc94274b8c5f9874cbd219789b3eb2128075a76f72",
	"df02719c4df8b9b8ac7f551fcb5d9ef48fa27eef7a66453879f4d8fdc6e78fb1",
	"ff12e391b79415e941a94de3bf3a9aee577aed0731e297d5cfa0b8a1e02fa1d0",
	"752dd9cf65e68cfaba7d60225cbdbc1f4729dd5e5507def72815ed0d8abc6249",
	"efb595a0178eb79a8df953f87c5148402a224cdf725e88c0146727c6aceadccd",
}

func main() {
	root := cobra.Command{
		Use:   "res",
		Short: "A command-line interface to interact with the Reserve Dollar smart contract",
	}
	viper.SetEnvPrefix("rsvd")
	viper.AutomaticEnv()
	root.PersistentFlags().StringP(
		"from",
		"F",
		defaultKeys[0],
		"Hex-encoded private key to sign transactions with. Defaults to the 0th address in the 0x mnemonic.",
	)
	root.PersistentFlags().StringP(
		"address",
		"A",
		"",
		"Address of a deployed copy of the Reserve Dollar contract.",
	)
	//viper.BindPFlags(root.PersistentFlags())

	root.AddCommand(
		deployCmd,

		//balanceOfCmd,
		//allowanceCmd,
		totalSupplyCmd,

		//transferCmd,
		//approveCmd,
		//transferFromCmd,
	)

	//cobra.OnInitialize(func() {
	//viper.ReadInConfig()
	//})

	check(root.Execute())
}

// TODO: show events

func getNode() *ethclient.Client {
	client, err := ethclient.Dial("http://localhost:8545")
	checkMsg(err, "Failed to connect to Ethereum node (is there a node running at localhost:8545?)")
	return client
}

func getSigner() *bind.TransactOpts {
	keyBytes, err1 := hex.DecodeString(viper.GetString("from"))
	key, err2 := crypto.ToECDSA(keyBytes)
	if err1 != nil || err2 != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse a private key to use for transaction signing.")
		fmt.Fprintln(os.Stderr, "To specify a private key, set the --from flag or RSVD_FROM environment variable to a hex-encoded private key of an account that already has enough ETH to cover gas costs.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Or leave that flag and env variable unset to use the default key.")
		os.Exit(1)
	}
	return bind.NewKeyedTransactor(key)
}

func getReserveDollar() *abi.ReserveDollar {
	address := viper.GetString("address")
	if address == "" {
		fmt.Fprintln(os.Stderr, "No address specified for the Reserve Dollar contract.")
		fmt.Fprintln(os.Stderr, "To specify an address, set the --address flag or the RSVD_ADDRESS environment variable.")
		fmt.Fprintln(os.Stderr, "To deploy a new contract and set the RSVD_ADDRESS in your current shell in a single step, run:")
		fmt.Fprintln(os.Stderr, "\t$(res deploy)")
		os.Exit(1)
	}
	rsvd, err := abi.NewReserveDollar(common.HexToAddress(address), getNode())
	checkMsg(err, "Couldn't bind Reserve Dollar contract")
	return rsvd
}

func toDisplay(i *big.Int) string {
	return decimal.NewFromBigInt(i, -18).String()
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new copy of the Reserve Dollar",
	Long: `Deploy a new copy of the Reserve Dollar.

This command also outputs the newly-deployed address in the format:

	export RSVD_ADDRESS=<new address>

Which enables using the following pattern to conveniently deploy a new
contract and use that contract address in subsequent commands:

	$ $(res deploy)
	$ res balanceOf @0
`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, _, _, err := abi.DeployReserveDollar(getSigner(), getNode())
		checkMsg(err, "Failed to deploy Reserve Dollar")
		fmt.Println("export RSVD_ADDRESS=" + addr.Hex())
	},
}

var totalSupplyCmd = &cobra.Command{
	Use:   "totalSupply",
	Short: "Call the totalSupply() function on the current Reserve Dollar",
	Run: func(cmd *cobra.Command, args []string) {
		totalSupply, err := getReserveDollar().TotalSupply(nil)
		checkMsg(err, "totalSupply() call failed")
		fmt.Println(toDisplay(totalSupply))
	},
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func checkMsg(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg+":", err)
		os.Exit(1)
	}
}
