package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/xdg"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/reserve-protocol/reserve-dollar/abi"
)

var defaultKeys = []string{
	"f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d",
	"5d862464fe9303452126c8bc94274b8c5f9874cbd219789b3eb2128075a76f72",
	"df02719c4df8b9b8ac7f551fcb5d9ef48fa27eef7a66453879f4d8fdc6e78fb1",
	"ff12e391b79415e941a94de3bf3a9aee577aed0731e297d5cfa0b8a1e02fa1d0",
	"752dd9cf65e68cfaba7d60225cbdbc1f4729dd5e5507def72815ed0d8abc6249",
	"efb595a0178eb79a8df953f87c5148402a224cdf725e88c0146727c6aceadccd",
	"83c6d2cc5ddcf9711a6d59b417dc20eb48afd58d45290099e5987e3d768f328f",
	"bb2d3f7c9583780a7d3904a2f55d792707c345f21de1bacb2d389934d82796b2",
	"b2fd4d29c1390b71b8795ae81196bfd60293adf99f9d32a0aff06288fcdac55f",
	"23cb7121166b9a2f93ae0b7c05bde02eae50d64449b2cbb42bc84e9d38d6cc89",
}

const pokeIntro = `A CLI for interacting with the Reserve Dollar smart contract.

This is designed for testing purposes. The goal is to make it easier to run small experiments
on the Reserve Dollar from the command line, without needing to write any code.

When we deploy the Reserve Dollar for real, we will use similar code, but it will go through
a hardware wallet.

The CLI includes access to all of the public functions on the Reserve Dollar.

The CLI is written assuming that it is being run against a local Ethereum node, available
on http://localhost:8545, with the same-prefunded accounts as the 0xorg/devnet docker image.
To run the 0xorg/devnet docker image, use the command:

    make run-geth

To deploy a new copy of the Reserve Dollar locally, run:

    $(poke deploy)

Running this command inside '$(...)' will cause your shell to execute the output of the
command, which will set the appropriate environment variable for your next calls to 'poke'
to run against the contract you just deployed.

To see the owner of the contract you just deployed, run:

    poke owner

This should show '0x5409ED021D9299bf6814279A6A1411A7e866A631', the 0th pre-funded account
in the 0xorg/devnet image. You can check that with 'poke address':

    poke address @0

Anywhere you need to supply an address or a private key to the poke tool, you can use
the special strings '@0' - '@9' to get the corresponding address or key from the ten
pre-funded accounts in the 0xorg/devnet image.

For paid mutator calls, 'poke' will default to using account '@0'. To override this default
per-command, you can use the '-F' (aka '--from') flag, like so:

    poke --from @1 transfer @2 200.5

You can also switch the default for the remainder of the current terminal session with
'poke account':

    $(poke account @3)
`

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

{{range getAllCmds}}
{{.Name}}:{{range .Commands}}
  {{rpad .Name .NamePadding}} {{.Short}}{{end}}
{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

  Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func main() {
	root := cobra.Command{
		Use:   "poke",
		Short: "A command-line interface to interact with the Reserve Dollar smart contract",
		Long:  pokeIntro,
	}
	type cmdBlock struct {
		Name     string
		Commands []*cobra.Command
	}
	allCmds := []cmdBlock{
		{
			"CLI Utility Commands",
			[]*cobra.Command{
				deployCmd,
				addressCmd,
				accountCmd,
			},
		},
		{
			"Basic Information Commands",
			[]*cobra.Command{
				nameCmd,
				symbolCmd,
				decimalsCmd,
			},
		},
		{
			"Other State Reading Commands",
			[]*cobra.Command{
				pausedCmd,
			},
		},
		{
			"ERC-20 Commands",
			[]*cobra.Command{
				balanceOfCmd,
				allowanceCmd,
				totalSupplyCmd,
				transferCmd,
				approveCmd,
				transferFromCmd,
			},
		},
		{
			"ERC-20-Approval-Bug Mitigation Commands",
			[]*cobra.Command{
				increaseAllowanceCmd,
				decreaseAllowanceCmd,
			},
		},
		{
			"Admin Role Reading Commands",
			[]*cobra.Command{
				ownerCmd,
				minterCmd,
				pauserCmd,
				freezerCmd,
				nominatedOwnerCmd,
			},
		},
		{
			"Admin Role Changing Commands",
			[]*cobra.Command{
				changeMinterCmd,
				changePauserCmd,
				changeFreezerCmd,
				nominateNewOwnerCmd,
				acceptOwnershipCmd,
				renounceOwnershipCmd,
			},
		},
		{
			"Pausing and Freezing Commands",
			[]*cobra.Command{
				pauseCmd,
				unpauseCmd,
				freezeCmd,
				unfreezeCmd,
				wipeCmd,
			},
		},
		{
			"Minting and Burning Commands",
			[]*cobra.Command{
				mintCmd,
				burnFromCmd,
			},
		},
		{
			"Upgrade Commands",
			[]*cobra.Command{
				transferEternalStorageCmd,
				changeNameCmd,
			},
		},
		{
			"Ethereum Utility Commands",
			[]*cobra.Command{
				sendEthCmd,
				deployBytecodeCmd,
			},
		},
	}
	for _, block := range allCmds {
		root.AddCommand(block.Commands...)
	}
	cobra.AddTemplateFunc("getAllCmds", func() []cmdBlock {
		return allCmds
	})
	root.SetUsageTemplate(usageTemplate)
	viper.SetEnvPrefix("rsvd")
	viper.AutomaticEnv()
	root.PersistentFlags().StringP(
		"from",
		"F",
		defaultKeys[0],
		"Hex-encoded private key to sign transactions with. Defaults to the 0th address in the 0x mnemonic.",
	)
	root.PersistentFlags().String(
		"tokenAddress",
		"",
		"Address of a deployed copy of the Reserve Dollar contract.",
	)
	root.PersistentFlags().StringP(
		"node",
		"n",
		"http://localhost:8545",
		"URL of an Ethereum node",
	)
	viper.BindPFlags(root.PersistentFlags())

	defer runExitFuncs()

	// Read config file.
	{
		cfgFile, err := xdg.Paths{XDGSuffix: "reserve"}.ConfigFile("poke.toml")
		if err == nil {
			viper.SetConfigFile(cfgFile)
			err = viper.ReadInConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to read config:", err)
				exit(1)
			}
		}
	}

	err := root.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exit(1)
	}
}

var exitFuncs []func()

func atExit(f func()) {
	exitFuncs = append(exitFuncs, f)
}

func runExitFuncs() {
	for _, f := range exitFuncs {
		f()
	}
}

func exit(code int) {
	runExitFuncs()
	os.Exit(code)
}

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	exit(1)
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	exit(1)
}

var client *ethclient.Client

func getNode() *ethclient.Client {
	if client == nil {
		var err error
		nodeAddr := viper.GetString("node")
		client, err = ethclient.Dial(nodeAddr)
		check(err, fmt.Sprintf("Failed to connect to Ethereum node (is there a node running at %q?)", nodeAddr))
	}
	return client
}

var (
	singletonAccount accounts.Account
	singletonWallet  accounts.Wallet
)

func openHardwareWallet() (accounts.Wallet, accounts.Account) {
	if singletonWallet != nil {
		return singletonWallet, singletonAccount
	}

	// Open hardware wallet.
	{
		// Check for connected Ledgers and Trezors.
		ledgerHub, err := usbwallet.NewLedgerHub()
		check(err, "calling usbwallet.NewLedgerHub()")
		trezorHub, err := usbwallet.NewTrezorHub()
		check(err, "calling usbwallet.NewTrezorHub()")

		// Collect them into a single list.
		wallets := accounts.NewManager(ledgerHub, trezorHub).Wallets()

		// Don't proceed unless there is exactly one hardware wallet available.
		if len(wallets) == 0 {
			fatal("No hardware wallets found. Is a hardware wallet plugged in? If it's a Ledger, is it unlocked?")
		}
		if len(wallets) > 1 {
			fatalf("%v hardware wallets found, I don't know which to use", len(wallets))
		}

		wallet := wallets[0]

		// "Open" the wallet.
		// This exchanges initial handshake messages with the wallet.
		// On a Trezor, this may require PIN entry.
		err = wallet.Open("")
		if err == usbwallet.ErrTrezorPINNeeded {
			fmt.Println("enter PIN (following scramble on the Trezor display)")
			scanner := bufio.NewScanner(os.Stdin)
			ok, pin, pinErr := scanner.Scan(), scanner.Text(), scanner.Err()
			if !ok {
				fatal("user cancelled pin entry")
			}
			check(pinErr, "getting PIN input")
			err = wallet.Open(pin)
		}
		check(err, "opening hardware wallet")

		// Notify the wallet when the program exits.
		atExit(func() {
			wallet.Close()
		})

		singletonWallet = wallet
	}

	// Open account.
	{
		var (
			hardened uint32 = 1 << 31
			err      error
		)
		singletonAccount, err = singletonWallet.Derive(
			// Standard Ethereum derivation path: m/44'/60'/0'/0 , with index 0
			// See eg https://ethereum.stackexchange.com/a/19061 for more information about derivation paths
			// Note that we don't really need to use the standard path here, but it's a fine default.
			[]uint32{
				44 | hardened,
				60 | hardened,
				0 | hardened,
				0,
				0,
			},
			true, // "pin" this account -- needed for the wallet object to recognize it later in "wallet.SignTx"
		)
		if err != nil && err.Error() == "reply lacks public key entry" {
			fatal("Failed to get public key. Is the Ethereum app open on the Ledger?")
		}
		check(err, "deriving account")
	}

	return singletonWallet, singletonAccount
}

func getNetID() *big.Int {
	netID, err := getNode().NetworkID(context.Background())
	check(err, "Failed to get Ethereum network id")
	return netID
}

func getSigner() *bind.TransactOpts {
	from := viper.GetString("from")
	if from != "hardware" {
		return bind.NewKeyedTransactor(parseKey(from))
	}

	wallet, account := openHardwareWallet()
	return &bind.TransactOpts{
		From: account.Address,
		Signer: func(
			protocolSigner types.Signer,
			from common.Address,
			tx *types.Transaction,
		) (*types.Transaction, error) {
			if from != account.Address {
				fatalf(
					"unexpected `from` address. from=%v account=%v",
					from.Hex(),
					account.Address.Hex(),
				)
			}
			fmt.Println("Waiting for you to confirm on the hardware wallet...")
			return wallet.SignTx(account, tx, getNetID())
		},
	}
}

var rsvd *abi.ReserveDollar

func getReserveDollar() *abi.ReserveDollar {
	if rsvd == nil {
		address := viper.GetString("tokenAddress")
		if address == "" {
			fmt.Fprintln(os.Stderr, "No address specified for the Reserve Dollar contract.")
			fmt.Fprintln(os.Stderr, "To specify an address, set the --address flag or the RSVD_TOKENADDRESS environment variable.")
			fmt.Fprintln(os.Stderr, "To deploy a new contract and set the RSVD_TOKENADDRESS in your current shell in a single step, run:")
			fmt.Fprintln(os.Stderr, "\t$(poke deploy)")
			exit(1)
		}
		var err error
		rsvd, err = abi.NewReserveDollar(hexToAddress(address), getNode())
		check(err, "Couldn't bind Reserve Dollar contract")
	}
	return rsvd
}

// TODO (issue #13): It'd be cleaner for the addresses currently named
// "@0" through "@9" on the command line to just be name "0" through
// "9" -- and then this setting (and the later parseKey) need not go
// through os.Getenv.
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
			fmt.Fprintf(
				os.Stderr,
				"To use a shorthand argument like %q, there should be a non-empty corresponding environment variable called %q\n",
				s,
				"RSVD_"+s[1:],
			)
			exit(1)
		}
		s = env
	}
	keyBytes, err1 := hex.DecodeString(s)
	key, err2 := crypto.ToECDSA(keyBytes)
	if err1 != nil || err2 != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse private key:", s)
		if strings.HasPrefix(origS, "@") {
			fmt.Fprintf(
				os.Stderr,
				"(From argument %q, which I expanded using the env var %v)\n",
				origS,
				"RSVD_"+origS[1:],
			)
			exit(1)
		}
	}
	return key
}

// HexToAddress parses s into a common.Address.
// Unlike go-ethereum's common.HexToAddress, this version
// exits if s is not a valid address encoding.
// This is copied and lightly modifed from:
// https://github.com/reserve-protocol/reserve/tree/686b03e/pkg/eth
func hexToAddress(s string) common.Address {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	check(err, fmt.Sprintf("invalid hex string %q", s))
	var address common.Address
	if len(b) != len(address) {
		fatalf("invalid address length: %v", len(b))
	}
	address.SetBytes(b)
	return address
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
	return hexToAddress(s)
}

// parseAttoTokens reads a decimal-formatted number of tokens and returns that number time 1e18.
func parseAttoTokens(s string) *big.Int {
	d, err := decimal.NewFromString(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Expected a decimal number, but got %q instead.\n", s)
		exit(1)
	}
	return truncateDecimal(d.Shift(18))
}

// truncateDecimal truncates d to an integer and returns it as a *big.Int.
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

// toDisplay is the inverse of parseAttoTokens
func toDisplay(i *big.Int) string {
	return decimal.NewFromBigInt(i, -18).String()
}

func toAddress(key *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(key.PublicKey)
}

func keyToHex(key *ecdsa.PrivateKey) string {
	return hex.EncodeToString(crypto.FromECDSA(key))
}

// log logs the result of a mutator txn to stdout, including that txn's events.
func log(name string, tx *types.Transaction, err error) {
	check(err, name+" failed")
	receipt, err := bind.WaitMined(context.Background(), getNode(), tx)
	check(err, "waiting for "+name+" to be mined")
	if len(receipt.Logs) > 0 {
		rsvd := getReserveDollar()
		fmt.Println("Done. Events:")
		for _, log := range receipt.Logs {
			parsed, err := rsvd.ParseLog(log)
			if err == nil {
				fmt.Println("\t" + parsed.String())
			} else {
				fmt.Println("\t" + err.Error())
			}
		}
	} else {
		fmt.Println("< Done. No events generated >")
	}
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new copy of the Reserve Dollar.",
	Long: `Deploy a new copy of the Reserve Dollar.

This command also outputs the newly-deployed address in the format:

	export RSVD_TOKENADDRESS=<new address>

Which enables using the following pattern to conveniently deploy a new
contract and use that contract address in subsequent commands:

	$ $(poke deploy)
	$ poke balanceOf @0
`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, tx, _, err := abi.DeployReserveDollar(getSigner(), getNode())
		check(err, "Failed to deploy Reserve Dollar")
		receipt, err := bind.WaitMined(context.Background(), getNode(), tx)
		check(err, "Waiting for deployment tx to be mined")
		if receipt.Status != types.ReceiptStatusSuccessful {
			fatal("Deployment transaction reverted")
		}
		fmt.Println("export RSVD_TOKENADDRESS=" + addr.Hex())
	},
}

func getAddress() common.Address {
	from := viper.GetString("from")
	if from == "hardware" {
		_, account := openHardwareWallet()
		return account.Address
	}
	return toAddress(parseKey(from))
}

var addressCmd = &cobra.Command{
	Use:     "address",
	Short:   "Get the address corresponding to the from account",
	Example: "  poke address\n  poke address -F @1",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getAddress().Hex())
	},
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Change the current acting account (invoke with `$(poke account)`).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("export RSVD_FROM=" + keyToHex(parseKey(args[0])))
	},
}

var balanceOfCmd = &cobra.Command{
	Use:   "balanceOf <address holder>",
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
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		totalSupply, err := getReserveDollar().TotalSupply(nil)
		check(err, "totalSupply() call failed")
		fmt.Println(toDisplay(totalSupply))
	},
}

var nameCmd = &cobra.Command{
	Use:   "name",
	Short: "Get the name of Reserve Dollars.",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := getReserveDollar().Name(nil)
		check(err, "name() call failed")
		fmt.Println(name)
	},
}

var symbolCmd = &cobra.Command{
	Use:   "symbol",
	Short: "Get the symbol of Reserve Dollars.",
	Run: func(cmd *cobra.Command, args []string) {
		symbol, err := getReserveDollar().Symbol(nil)
		check(err, "symbol() call failed")
		fmt.Println(symbol)
	},
}

var decimalsCmd = &cobra.Command{
	Use:   "decimals",
	Short: "Get the decimals of Reserve Dollars.",
	Run: func(cmd *cobra.Command, args []string) {
		decimals, err := getReserveDollar().Decimals(nil)
		check(err, "decimals() call failed")
		fmt.Println(decimals)
	},
}

var pausedCmd = &cobra.Command{
	Use:   "paused",
	Short: "Get the paused state of Reserve Dollars.",
	Run: func(cmd *cobra.Command, args []string) {
		paused, err := getReserveDollar().Paused(nil)
		check(err, "paused() call failed")
		fmt.Println(paused)
	},
}

var transferCmd = &cobra.Command{
	Use:   "transfer <address to> <uint value>",
	Short: "Transfer Reserve Dollars.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Transfer(getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]))
		log("transfer()", tx, err)
	},
}

var approveCmd = &cobra.Command{
	Use:   "approve <address spender> <uint allowance>",
	Short: "Approve a spender to spend Reserve Dollars.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Approve(getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]))
		log("approve()", tx, err)
	},
}

var transferFromCmd = &cobra.Command{
	Use:   "transferFrom <address from> <address to> <uint value>",
	Short: "Transfer tokens from `from` to `to`. Must be sent by an account approved to spend from `from`.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().TransferFrom(
			getSigner(), parseAddress(args[0]), parseAddress(args[1]), parseAttoTokens(args[2]),
		)
		log("transferFrom()", tx, err)
	},
}

var increaseAllowanceCmd = &cobra.Command{
	Use:   "increaseAllowance <address spender> <uint addedValue>",
	Short: "Increase the allowance of a spender.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().IncreaseAllowance(
			getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]),
		)
		log("increaseAllowance()", tx, err)
	},
}

var decreaseAllowanceCmd = &cobra.Command{
	Use:   "decreaseAllowance <address spender> <uint subtractedValue>",
	Short: "Decrease the allowance of a spender.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().DecreaseAllowance(
			getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]),
		)
		log("decreaseAllowance()", tx, err)
	},
}

var minterCmd = &cobra.Command{
	Use:   "minter",
	Short: "Show the current minter.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		minter, err := getReserveDollar().Minter(nil)
		check(err, "minter() call failed")
		fmt.Println(minter.Hex())
	},
}

var pauserCmd = &cobra.Command{
	Use:   "pauser",
	Short: "Show the current pauser.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		pauser, err := getReserveDollar().Pauser(nil)
		check(err, "pauser() call failed")
		fmt.Println(pauser.Hex())
	},
}

var freezerCmd = &cobra.Command{
	Use:   "freezer",
	Short: "Show the current freezer.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		freezer, err := getReserveDollar().Freezer(nil)
		check(err, "freezer() call failed")
		fmt.Println(freezer.Hex())
	},
}

var ownerCmd = &cobra.Command{
	Use:   "owner",
	Short: "Show the current owner.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		owner, err := getReserveDollar().Owner(nil)
		check(err, "owner() call failed")
		fmt.Println(owner.Hex())
	},
}

var nominatedOwnerCmd = &cobra.Command{
	Use:   "nominatedOwner",
	Short: "Show the current nominatedOwner.",
	Args:  cobra.ExactArgs(0),
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
	poke changeMinter $(poke address @1)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().ChangeMinter(getSigner(), parseAddress(args[0]))
		log("changeMinter()", tx, err)
	},
}

var changePauserCmd = &cobra.Command{
	Use:   "changePauser <newPauser>",
	Short: "Change the pauser role. Must be called by the current pauser or owner.",
	Example: `Change the pauser role to the address corresponding to the first default key:
	poke changePauser$(poke address @1)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().ChangePauser(getSigner(), parseAddress(args[0]))
		log("changePauser()", tx, err)
	},
}

var changeFreezerCmd = &cobra.Command{
	Use:   "changeFreezer <newFreezer>",
	Short: "Change the freezer role. Must be called by the current freezer or owner.",
	Example: `Change the freezer role to the address corresponding to the first default key:
	poke changeFreezer$(poke address @1)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().ChangeFreezer(getSigner(), parseAddress(args[0]))
		log("changeFreezer()", tx, err)
	},
}

var nominateNewOwnerCmd = &cobra.Command{
	Use:   "nominateNewOwner <newOwner>",
	Short: "Nominate a new owner for the Reserve Dollar contract. The new owner must call `acceptOwnership`.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().NominateNewOwner(getSigner(), parseAddress(args[0]))
		log("nominateNewOwner()", tx, err)
	},
}

var acceptOwnershipCmd = &cobra.Command{
	Use:   "acceptOwnership",
	Short: "Accept ownership nomination. Must be called by current nominatedOwner.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().AcceptOwnership(getSigner())
		log("acceptOwnership()", tx, err)
	},
}

var renounceOwnershipCmd = &cobra.Command{
	Use:   "renounceOwnership",
	Short: "Renounce ownership. Must be called by current owner.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().RenounceOwnership(getSigner())
		log("renounceOwnership()", tx, err)
	},
}

var transferEternalStorageCmd = &cobra.Command{
	Use:   "transferEternalStorage <address newOwner>",
	Short: "Transfer ownership of eternal storage contract. Must be called by current owner.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().TransferEternalStorage(getSigner(), parseAddress(args[0]))
		log("transferEternalStorage()", tx, err)
	},
}

var changeNameCmd = &cobra.Command{
	Use:   "changeName <string newName> <string newSymbol>",
	Short: "Change the name and/or symbol of the Reserve Dollar. Must be called by current owner.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().ChangeName(getSigner(), args[0], args[1])
		log("changeName()", tx, err)
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the Reserve Dollar contract.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Pause(getSigner())
		log("pause()", tx, err)
	},
}

var unpauseCmd = &cobra.Command{
	Use:   "unpause",
	Short: "Unpause the Reserve Dollar contract.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Unpause(getSigner())
		log("unpause()", tx, err)
	},
}

var freezeCmd = &cobra.Command{
	Use:   "freeze <address who>",
	Short: "Freeze an account.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Freeze(getSigner(), parseAddress(args[1]))
		log("freeze()", tx, err)
	},
}

var unfreezeCmd = &cobra.Command{
	Use:   "unfreeze <address who>",
	Short: "Unfreeze an account.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Unfreeze(getSigner(), parseAddress(args[1]))
		log("unfreeze()", tx, err)
	},
}

var wipeCmd = &cobra.Command{
	Use:   "wipe <address who>",
	Short: "Wipe an account.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Wipe(getSigner(), parseAddress(args[1]))
		log("wipe()", tx, err)
	},
}

var mintCmd = &cobra.Command{
	Use:   "mint <address> <value>",
	Short: "Mint tokens to an address. Must be called by the current minter.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().Mint(getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]))
		log("mint()", tx, err)
	},
}

var burnFromCmd = &cobra.Command{
	Use:   "burnFrom <address> <value>",
	Short: "Burn tokens from an address that has approved the minter.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tx, err := getReserveDollar().BurnFrom(getSigner(), parseAddress(args[0]), parseAttoTokens(args[1]))
		log("burnFrom()", tx, err)
	},
}

var sendEthCmd = &cobra.Command{
	Use:   "send-eth <address> <value>",
	Short: "Send ETH to an address.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		nonce, err := getNode().PendingNonceAt(ctx, getAddress())
		check(err, "retrieving nonce")
		gasPrice, err := getNode().SuggestGasPrice(ctx)
		check(err, "retrieving gas price suggestion")
		address := parseAddress(args[0])
		attoTokens := parseAttoTokens(args[1])
		tx, err := getSigner().Signer(
			types.NewEIP155Signer(getNetID()),
			getAddress(),
			types.NewTransaction(
				nonce,
				address,
				attoTokens,
				21000,
				gasPrice,
				nil,
			),
		)
		check(err, "signing transaction")
		check(getNode().SendTransaction(ctx, tx), "sending transaction")
		fmt.Printf("Sent %v atto-ETH to %v.\n", attoTokens, address.Hex())
	},
}

var deployBytecodeCmd = &cobra.Command{
	Use:   "deploy-bytecode <hex-encoded bin>",
	Short: "Deploy arbitrary EVM bytecode.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		nonce, err := getNode().PendingNonceAt(ctx, getAddress())
		check(err, "retrieving nonce")
		gasPrice, err := getNode().SuggestGasPrice(ctx)
		check(err, "retrieving gas price suggestion")

		bytecode, err := hex.DecodeString(args[0])
		check(err, "failed to hex-decode bytecode")

		gasEstimate, err := getNode().EstimateGas(ctx, ethereum.CallMsg{
			From:     getAddress(),
			To:       nil,
			Gas:      0,
			GasPrice: gasPrice,
			Value:    big.NewInt(0),
			Data:     bytecode,
		})
		check(err, "estimating gas")

		tx, err := getSigner().Signer(
			types.NewEIP155Signer(getNetID()),
			getAddress(),
			types.NewContractCreation(
				nonce,
				big.NewInt(0),
				gasEstimate,
				gasPrice,
				bytecode,
			),
		)
		check(err, "signing transaction")
		err = getNode().SendTransaction(ctx, tx)
		log("deploy-bytecode", tx, err)
	},
}

func check(err error, msg string) {
	if err != nil {
		fatal(msg+":", err)
	}
}
