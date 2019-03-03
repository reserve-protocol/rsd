package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

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
	viper.BindPFlags(root.PersistentFlags())

	root.AddCommand(
		deployCmd,
		addressCmd,
		accountCmd,

		balanceOfCmd,
		allowanceCmd,
		totalSupplyCmd,

		transferCmd,
		approveCmd,
		transferFromCmd,

		ownerCmd,
		minterCmd,
		pauserCmd,
		freezerCmd,
		nominatedOwnerCmd,

		changeMinterCmd,
		changePauserCmd,
		changeFreezerCmd,

		//nominateNewOwnerCmd,
		//acceptOwnershipCmd,

		mintCmd,
	)

	//cobra.OnInitialize(func() {
	//viper.ReadInConfig()
	//})

	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// TODO: show events

func getNode() *ethclient.Client {
	client, err := ethclient.Dial("http://localhost:8545")
	check(err, "Failed to connect to Ethereum node (is there a node running at localhost:8545?)")
	return client
}

func getSigner() *bind.TransactOpts {
	/*
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
	*/
	return bind.NewKeyedTransactor(parseKey(viper.GetString("from")))
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
	check(err, "Couldn't bind Reserve Dollar contract")
	return rsvd
}

func init() {
	for i, key := range defaultKeys {
		envVar := "RSVD_" + strconv.Itoa(i)
		if os.Getenv(envVar) == "" {
			os.Setenv(envVar, key)
		}
	}
}

// parseKey parses a hex-encoded private key from s.
// Alternatively, if s begins with "@", parseKey parses
// a hex-encoded private key from the environment variable
// named "RSVD_<s[1:]>".
func parseKey(s string) *ecdsa.PrivateKey {
	origS := s
	if strings.HasPrefix(s, "@") {
		env := os.Getenv("RSVD_" + s[1:])
		if s == "" {
			fmt.Fprintf(os.Stderr, "To use a shorthand argument like %q, there should be a non-empty corresponding environment variable called %q\n", s, "RSVD_"+s[1:])
			os.Exit(1)
		}
		s = env
	}
	keyBytes, err1 := hex.DecodeString(s)
	key, err2 := crypto.ToECDSA(keyBytes)
	if err1 != nil || err2 != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse private key:", s)
		if strings.HasPrefix(origS, "@") {
			fmt.Fprintf(os.Stderr, "(From argument %q, which I expanded using the env var %v)\n", origS, "RSVD_"+origS[1:])
			os.Exit(1)
		}
	}
	return key
}

// parseAddress parses a hex-encoded address from s.
// Alternatively, if s begins with "@", parseAddress parses
// a hex-encoded private key from the environment variable
// named "RSVD_<s[1:]>", then returns the address corresponding
// to that key.
func parseAddress(s string) common.Address {
	if strings.HasPrefix(s, "@") {
		return toAddress(parseKey(s))
	}
	return common.HexToAddress(s)
}

func parseAttoTokens(s string) *big.Int {
	d, err := decimal.NewFromString(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Expected a decimal number, but got %q instead.\n", s)
		os.Exit(1)
	}
	return truncateDecimal(d.Shift(18))
}

// truncateDecimal truncates d to an integer and returns it as an *Int.
func truncateDecimal(d decimal.Decimal) *big.Int {
	coeff := d.Coefficient()
	exp := d.Exponent()
	z := new(big.Int)
	if exp >= 0 {
		// 	coeff * 10 ** exp
		return coeff.Mul(coeff, z.Exp(big.NewInt(10), big.NewInt(int64(exp)), nil))
	}
	// 	coeff / 10 ** -exp
	return coeff.Div(coeff, z.Exp(big.NewInt(10), big.NewInt(int64(-exp)), nil))
}

func toDisplay(i *big.Int) string {
	return decimal.NewFromBigInt(i, -18).String()
}

func toAddress(key *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(key.PublicKey)
}

func keyToHex(key *ecdsa.PrivateKey) string {
	return hex.EncodeToString(crypto.FromECDSA(key))
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new copy of the Reserve Dollar.",
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
		check(err, "Failed to deploy Reserve Dollar")
		fmt.Println("export RSVD_ADDRESS=" + addr.Hex())
	},
}

var addressCmd = &cobra.Command{
	Use:     "address",
	Short:   "Get the address corresponding to a private key. Accepts @-shorthands.",
	Example: "res address @1",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(toAddress(parseKey(args[0])).Hex())
	},
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Change the current acting account (invoke with `$(res account)`).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("export RSVD_FROM=" + keyToHex(parseKey(args[0])))
	},
}

var balanceOfCmd = &cobra.Command{
	Use:   "balanceOf <addressi holder>",
	Short: "Get an account's balance of Reserve Dollars.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		balance, err := getReserveDollar().BalanceOf(nil, parseAddress(args[0]))
		check(err, "balanceOf() call failed")
		fmt.Println(toDisplay(balance))
	},
}

var allowanceCmd = &cobra.Command{
	Use:   "allowance <address holder> <address spender>",
	Short: "Get the allowance that holder is currently granting to spender.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		allowance, err := getReserveDollar().Allowance(nil, parseAddress(args[0]), parseAddress(args[1]))
		check(err, "allowance() call failed")
		fmt.Println(toDisplay(allowance))
	},
}

var totalSupplyCmd = &cobra.Command{
	Use:   "totalSupply",
	Short: "Get the total supply of Reserve Dollars.",
	Run: func(cmd *cobra.Command, args []string) {
		totalSupply, err := getReserveDollar().TotalSupply(nil)
		check(err, "totalSupply() call failed")
		fmt.Println(toDisplay(totalSupply))
	},
}

var transferCmd = &cobra.Command{
	Use:   "transfer <address to> <uint value>",
	Short: "Transfer Reserve Dollars.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getReserveDollar().Transfer(getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]))
		check(err, "transfer() failed")
		fmt.Println("Done.")
	},
}

var approveCmd = &cobra.Command{
	Use:   "approve <address spender> <uint allowance>",
	Short: "Approve a spender to spend Reserve Dollars.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getReserveDollar().Approve(getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]))
		check(err, "approve() failed")
		fmt.Println("Done.")
	},
}

var transferFromCmd = &cobra.Command{
	Use:   "transferFrom <address from> <address to> <uint value>",
	Short: "Transfer tokens from `from` to `to`. Must be sent by an account approved to spend from `from`.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getReserveDollar().TransferFrom(
			getSigner(), parseAddress(args[0]), parseAddress(args[1]), parseAttoTokens(args[2]),
		)
		check(err, "transferFrom() failed")
		fmt.Println("Done.")
	},
}

var minterCmd = &cobra.Command{
	Use:   "minter",
	Short: "Show the current minter.",
	Run: func(cmd *cobra.Command, args []string) {
		minter, err := getReserveDollar().Minter(nil)
		check(err, "minter() call failed")
		fmt.Println(minter.Hex())
	},
}

var pauserCmd = &cobra.Command{
	Use:   "pauser",
	Short: "Show the current pauser.",
	Run: func(cmd *cobra.Command, args []string) {
		pauser, err := getReserveDollar().Pauser(nil)
		check(err, "pauser() call failed")
		fmt.Println(pauser.Hex())
	},
}

var freezerCmd = &cobra.Command{
	Use:   "freezer",
	Short: "Show the current freezer.",
	Run: func(cmd *cobra.Command, args []string) {
		freezer, err := getReserveDollar().Freezer(nil)
		check(err, "freezer() call failed")
		fmt.Println(freezer.Hex())
	},
}

var ownerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Show the current owner.",
	Run: func(cmd *cobra.Command, args []string) {
		owner, err := getReserveDollar().Owner(nil)
		check(err, "owner() call failed")
		fmt.Println(owner.Hex())
	},
}

var nominatedOwnerCmd = &cobra.Command{
	Use:   "nominatedOwner",
	Short: "Show the current nominatedOwner.",
	Run: func(cmd *cobra.Command, args []string) {
		nominatedOwner, err := getReserveDollar().NominatedOwner(nil)
		check(err, "nominatedOwner() call failed")
		fmt.Println(nominatedOwner.Hex())
	},
}

var changeMinterCmd = &cobra.Command{
	Use:   "changeMinter <newMinter>",
	Short: "Change the minter role. Must be called by the current minter or owner.",
	Example: `Change the minter role to the address corresponding to the first default key:
	res changeMinter $(res address @1)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getReserveDollar().ChangeMinter(getSigner(), toAddress(parseKey(args[0])))
		check(err, "changeMinter() failed")
		fmt.Println("Done.")
	},
}

var changePauserCmd = &cobra.Command{
	Use:   "changePauser <newPauser>",
	Short: "Change the pauser role. Must be called by the current pauser or owner.",
	Example: `Change the pauser role to the address corresponding to the first default key:
	res changePauser$(res address @1)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getReserveDollar().ChangePauser(getSigner(), toAddress(parseKey(args[0])))
		check(err, "changePauser() failed")
		fmt.Println("Done.")
	},
}

var changeFreezerCmd = &cobra.Command{
	Use:   "changeFreezer <newFreezer>",
	Short: "Change the freezer role. Must be called by the current freezer or owner.",
	Example: `Change the freezer role to the address corresponding to the first default key:
	res changeFreezer$(res address @1)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := getReserveDollar().ChangeFreezer(getSigner(), toAddress(parseKey(args[0])))
		check(err, "changeFreezer() failed")
		fmt.Println("Done.")
	},
}

var mintCmd = &cobra.Command{
	Use:   "mint <address> <value>",
	Short: "Mint tokens to an address. Must be called by the current minter.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		recipient := parseAddress(args[0])
		value := parseAttoTokens(args[1])
		_, err := getReserveDollar().Mint(getSigner(), recipient, value)
		check(err, "mint() failed")
		fmt.Println("Done.")
	},
}

func check(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg+":", err)
		os.Exit(1)
	}
}
