package reserve_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/reserve-protocol/reserve-dollar/abi"
)

func TestReserveDollar(t *testing.T) {
	suite.Run(t, new(ReserveDollarSuite))
}

type ReserveDollarSuite struct {
	suite.Suite

	account []account
	signer  *bind.TransactOpts
	node    backend
	reserve *abi.ReserveDollar
	assert  *assert.Assertions
	require *require.Assertions
}

var (
	// Compile-time check that ReserveDollarSuite implements the interfaces we think it does.
	// If it does not implement these interfaces, then the corresponding setup and teardown
	// functions will not actually run.
	_ suite.BeforeTest       = &ReserveDollarSuite{}
	_ suite.SetupAllSuite    = &ReserveDollarSuite{}
	_ suite.TearDownAllSuite = &ReserveDollarSuite{}
)

// requireTx requires that a transaction is successfully mined and does
// not revert. It also takes an extra error argument, and checks that the
// error is nil. This signature allows the function to directly wrap our
// abigen'd mutator calls.
func (s *ReserveDollarSuite) requireTx(tx *types.Transaction, err error) {
	_requireTxStatus(s, tx, err, types.ReceiptStatusSuccessful)
}

// requireTxFails is like requireTx, but it requires that the transaction either
// reverts or is not successfully made in the first place due to gas estimation
// failing.
func (s *ReserveDollarSuite) requireTxFails(tx *types.Transaction, err error) {
	if err != nil && err.Error() ==
		"failed to estimate gas needed: gas required exceeds allowance or always failing transaction" {
		return
	}

	fmt.Printf("%q\n", err.Error())

	_requireTxStatus(s, tx, err, types.ReceiptStatusFailed)
}

func _requireTxStatus(s *ReserveDollarSuite, tx *types.Transaction, err error, status uint64) {
	s.Require().NoError(err)
	s.Require().NotNil(tx)
	receipt, err := bind.WaitMined(context.Background(), s.node, tx)
	s.Require().NoError(err)
	s.Require().Equal(status, receipt.Status)
}

func (s *ReserveDollarSuite) assertBalance(address common.Address, amount *big.Int) {
	balance, err := s.reserve.BalanceOf(nil, address)
	s.NoError(err)
	s.Equal(amount.String(), balance.String()) // assert.Equal can mis-compare big.Ints, so compare strings instead
}

func (s *ReserveDollarSuite) assertAllowance(owner, spender common.Address, amount *big.Int) {
	allowance, err := s.reserve.Allowance(nil, owner, spender)
	s.NoError(err)
	s.Equal(amount.String(), allowance.String())
}

func (s *ReserveDollarSuite) assertTotalSupply(amount *big.Int) {
	totalSupply, err := s.reserve.TotalSupply(nil)
	s.NoError(err)
	s.Equal(amount.String(), totalSupply.String())
}

func (s *ReserveDollarSuite) SetupSuite() {
	//var err error
	//s.node, err = soltools.NewBackend("http://localhost:8545", repo.Path("fiatcoin/artifacts"), repo.Path("fiatcoin"))
	//s.Require().NoError(err)

	// The first few keys from the following well-known mnemonic used by 0x:
	//	concert load couple harbor equip island argue ramp clarify fence smart topic
	keys := []string{
		"f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d",
		"5d862464fe9303452126c8bc94274b8c5f9874cbd219789b3eb2128075a76f72",
		"df02719c4df8b9b8ac7f551fcb5d9ef48fa27eef7a66453879f4d8fdc6e78fb1",
		"ff12e391b79415e941a94de3bf3a9aee577aed0731e297d5cfa0b8a1e02fa1d0",
		"752dd9cf65e68cfaba7d60225cbdbc1f4729dd5e5507def72815ed0d8abc6249",
		"efb595a0178eb79a8df953f87c5148402a224cdf725e88c0146727c6aceadccd",
	}
	s.account = make([]account, len(keys))
	for i, key := range keys {
		b, err := hex.DecodeString(key)
		s.Require().NoError(err)
		s.account[i].key, err = crypto.ToECDSA(b)
		s.Require().NoError(err)
	}
	s.signer = signer(s.account[0])

	genesisAlloc := core.GenesisAlloc{}
	for _, account := range s.account {
		genesisAlloc[account.address()] = core.GenesisAccount{
			Balance: big.NewInt(math.MaxInt64),
		}
	}
	s.node = backend{
		backends.NewSimulatedBackend(
			genesisAlloc,
			4000000, // roughly same order of magnitude as mainnet
		),
	}

	/*
		{
			// Throwaway initial transaction.
			// The tests fail if running against a newly-initialized 0xorg/devnet container.
			// I (jeremy) suspect that this is because the node is configured to move through
			// the historical Ethereum hard forks over the course of the first few blocks, rather
			// than all at once in the first block. Meaning the first transactions run against different
			// versions of Etherum than the rest of the transactions:
			//
			//   https://github.com/0xProject/0x-monorepo/blob/e909faa3ef9cea5d9b4044b993251e98afdb0d19/packages/devnet/genesis.json#L4-L9
			//
			// To work around this issue, we try to send a throwaway transaction at the beginning with a
			// Homestead-style signature. This will fail if it is not the first transaction on the chain,
			// but that's ok. If it is the first transaction on the chain, it succeeds and causes the chain
			// to advance by one block, upgrading the Ethereum version and allowing the rest of the tests
			// to pass.
			tx, _ := types.SignTx(
				types.NewTransaction(0, common.Address{}, common.Big0, 21000, common.Big1, nil),
				types.HomesteadSigner{},
				s.deployerKey,
			)
			s.node.SendTransaction(context.Background(), tx)
		}
	*/
}

func (s *ReserveDollarSuite) TearDownSuite() {
	/*
		s.Assert().NoError(s.node.WriteCoverage())
		s.Assert().NoError(s.node.Close())

		if out, err := exec.Command("npx", "istanbul", "report", "html").CombinedOutput(); err != nil {
			fmt.Println()
			fmt.Println("I generated coverage information in coverage/coverage.json.")
			fmt.Println("I tried to process it with `istanbul` to turn it into a readable report, but failed.")
			fmt.Println("The error I got when running istanbul was:", err)
			fmt.Println("Istanbul's output was:\n" + string(out))
		}
	*/
}

func (s *ReserveDollarSuite) BeforeTest(suiteName, testName string) {
	_, tx, reserve, err := abi.DeployReserveDollar(s.signer, s.node)
	s.requireTx(tx, err)

	s.reserve = reserve
}

func (s *ReserveDollarSuite) TestDeploy() {}

func (s *ReserveDollarSuite) TestBalanceOf() {
	s.assertBalance(common.Address{}, common.Big0)
}

func (s *ReserveDollarSuite) TestName() {
	name, err := s.reserve.Name(nil)
	s.NoError(err)
	s.Equal("Reserve Dollar", name)
}

func (s *ReserveDollarSuite) TestSymbol() {
	symbol, err := s.reserve.Symbol(nil)
	s.NoError(err)
	s.Equal("RSVD", symbol)
}

func (s *ReserveDollarSuite) TestDecimals() {
	decimals, err := s.reserve.Decimals(nil)
	s.NoError(err)
	s.Equal(uint8(18), decimals)
}

func (s *ReserveDollarSuite) TestChangeName() {
	const newName, newSymbol = "Flamingo", "MGO"
	s.requireTx(
		s.reserve.ChangeName(s.signer, newName, newSymbol),
	)

	// Check for ChangeName event.
	nameChangeIter, err := s.reserve.FilterNameChanged(nil)
	if s.NoError(err) {
		events := 0
		for nameChangeIter.Next() {
			events++
			s.Equal(newName, nameChangeIter.Event.NewName)
			s.Equal(newSymbol, nameChangeIter.Event.NewSymbol)
		}
		s.Equal(1, events, "expected exactly one NameChanged event")
		s.NoError(nameChangeIter.Error())
		s.NoError(nameChangeIter.Close())
	}

	// Check new name.
	name, err := s.reserve.Name(nil)
	s.NoError(err)
	s.Equal(newName, name)

	// Check new symbol.
	symbol, err := s.reserve.Symbol(nil)
	s.NoError(err)
	s.Equal(newSymbol, symbol)
}

func (s *ReserveDollarSuite) TestAllowsMinting() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)

	s.requireTx(s.reserve.Mint(s.signer, recipient, amount))

	s.assertBalance(s.account[0].address(), common.Big0)
	s.assertBalance(recipient, amount)
	s.assertTotalSupply(amount)
}

func (s *ReserveDollarSuite) TestTransfer() {
	sender := s.account[1]
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)

	// Mint to sender.
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), amount))

	// Transfer from sender to recipient.
	s.requireTx(s.reserve.Transfer(signer(sender), recipient, amount))

	// Check that balances are as expected.
	s.assertBalance(sender.address(), common.Big0)
	s.assertBalance(recipient, amount)
	s.assertBalance(s.account[0].address(), common.Big0)
	s.assertTotalSupply(amount)
}

func (s *ReserveDollarSuite) TestTransferExceedsFunds() {
	sender := s.account[1]
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)
	smallAmount := big.NewInt(10) // must be smaller than amount

	// Mint smallAmount to sender.
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), smallAmount))

	// Transfer from sender to recipient should fail.
	s.requireTxFails(s.reserve.Transfer(signer(sender), recipient, amount))

	// Balances should be as we expect.
	s.assertBalance(sender.address(), smallAmount)
	s.assertBalance(recipient, common.Big0)
	s.assertBalance(s.account[0].address(), common.Big0)
	s.assertTotalSupply(smallAmount)
}

func (s *ReserveDollarSuite) TestApprove() {
	owner := s.account[1]
	spender := s.account[2]
	amount := big.NewInt(53)

	// owner approves spender
	s.requireTx(s.reserve.Approve(signer(owner), spender.address(), amount))

	// approval should be reflected in allowance
	s.assertAllowance(owner.address(), spender.address(), amount)

	// shouldn't be symmetric
	s.assertAllowance(spender.address(), owner.address(), common.Big0)

	// balances shouldn't change
	s.assertBalance(owner.address(), common.Big0)
	s.assertBalance(spender.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}

func (s *ReserveDollarSuite) TestIncreaseAllowance() {
	owner := s.account[1]
	spender := s.account[2]
	amount := big.NewInt(2000)

	// owner approves spender through increaseAllowance
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), amount))

	// approval should be reflected in allowance
	s.assertAllowance(owner.address(), spender.address(), amount)

	// shouldn't be symmetric
	s.assertAllowance(spender.address(), owner.address(), common.Big0)

	// balances shouldn't change
	s.assertBalance(owner.address(), common.Big0)
	s.assertBalance(spender.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}

func maxUint256() *big.Int {
	z := big.NewInt(1)
	z = z.Lsh(z, 256)
	z = z.Sub(z, common.Big1)
	return z
}

func (s *ReserveDollarSuite) TestIncreaseAllowanceOverflow() {
	owner := s.account[1]
	spender := s.account[2]
	initialAmount := big.NewInt(10)

	// owner approves spender for initial amount
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// owner should not be able to increase approval high enough to overflow a uint256
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(owner), spender.address(), maxUint256()))
}

func (s *ReserveDollarSuite) TestDecreaseAllowance() {
	owner := s.account[1]
	spender := s.account[2]
	initialAmount := big.NewInt(10)
	decrease := big.NewInt(6)
	final := big.NewInt(4)

	// owner approves spender for initial amount
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// owner decreases allowance
	s.requireTx(s.reserve.DecreaseAllowance(signer(owner), spender.address(), decrease))

	// allowance should be as we expect
	s.assertAllowance(owner.address(), spender.address(), final)

	// balances shouldn't change
	s.assertBalance(owner.address(), common.Big0)
	s.assertBalance(spender.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}

/*

func (s *ReserveDollarSuite) TestPausing() {
	// Pause.
	s.requireTx(s.reserveImpl.Pause(s.signer))

	// Minting is allowed while paused.
	amount := big.NewInt(100)
	s.requireTx(s.reserve.Mint(s.signer, toAddress(s.deployerKey), amount))

	// Transfers are not allowed while paused.
	recipient := common.BigToAddress(common.Big1)
	s.requireTxFails(s.reserve.Transfer(s.signer, recipient, amount))
	s.assertBalance(recipient, common.Big0)

	// Unpause.
	s.requireTx(s.reserveImpl.Unpause(s.signer))

	// Transfers are allowed while unpaused.
	s.requireTx(s.reserve.Transfer(s.signer, recipient, amount))
	s.assertBalance(recipient, amount)
}
*/

type backend struct {
	*backends.SimulatedBackend
}

func (b backend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	defer b.Commit()
	return b.SimulatedBackend.SendTransaction(ctx, tx)
}

func signer(a account) *bind.TransactOpts {
	return bind.NewKeyedTransactor(a.key)
}

type account struct {
	key *ecdsa.PrivateKey
}

func (a account) address() common.Address {
	return crypto.PubkeyToAddress(a.key.PublicKey)
}
