package reserve_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
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

	accounts []*ecdsa.PrivateKey
	signer   *bind.TransactOpts
	node     backend
	reserve  *abi.ReserveDollar
	assert   *assert.Assertions
	require  *require.Assertions
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

// requireTxReverts is like requireTx, but it requires that the transaction reverts.
func (s *ReserveDollarSuite) requireTxReverts(tx *types.Transaction, err error) {
	_requireTxStatus(s, tx, err, types.ReceiptStatusFailed)
}

func _requireTxStatus(s *ReserveDollarSuite, tx *types.Transaction, err error, status uint64) {
	s.require.NoError(err)
	s.require.NotNil(tx)
	receipt, err := bind.WaitMined(context.Background(), s.node, tx)
	s.require.NoError(err)
	s.require.Equal(status, receipt.Status)
}

func (s *ReserveDollarSuite) assertBalance(address common.Address, amount *big.Int) {
	balance, err := s.reserve.BalanceOf(nil, address)
	s.assert.NoError(err)
	s.assert.Equal(amount.String(), balance.String()) // assert.Equal can mis-compare big.Ints, so compare strings instead
}

func (s *ReserveDollarSuite) SetupSuite() {
	//var err error
	//s.node, err = soltools.NewBackend("http://localhost:8545", repo.Path("fiatcoin/artifacts"), repo.Path("fiatcoin"))
	//s.Require().NoError(err)

	keys := []string{
		"f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d",
		"5d862464fe9303452126c8bc94274b8c5f9874cbd219789b3eb2128075a76f72",
	}
	s.accounts = make([]*ecdsa.PrivateKey, len(keys))
	for i, key := range keys {
		b, err := hex.DecodeString(key)
		s.Require().NoError(err)
		s.accounts[i], err = crypto.ToECDSA(b)
		s.Require().NoError(err)
	}
	s.signer = signer(s.accounts[0])

	genesisAlloc := core.GenesisAlloc{}
	for _, account := range s.accounts {
		genesisAlloc[toAddress(account)] = core.GenesisAccount{
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
	s.assert = s.Assert()
	s.require = s.Require()

	_, tx, reserve, err := abi.DeployReserveDollar(s.signer, s.node)
	s.requireTx(tx, err)

	s.reserve = reserve
}

func (s *ReserveDollarSuite) TestDeploy() {}

func (s *ReserveDollarSuite) TestBalanceOf() {
	s.assertBalance(common.Address{}, common.Big0)
}

/*
func (s *ReserveDollarSuite) TestDeploymentAccountIsAMinter() {
	isMinter, err := s.reserveImpl.IsMinter(nil, toAddress(s.deployerKey))
	s.assert.NoError(err)
	s.assert.True(isMinter)
}

func (s *ReserveDollarSuite) TestAllowsMinting() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)

	s.requireTx(s.reserve.Mint(s.signer, recipient, amount))

	s.assertBalance(recipient, amount)
}

func (s *ReserveDollarSuite) TestName() {
	name, err := s.reserve.Name(nil)
	s.assert.NoError(err)
	s.assert.Equal("Reserve Dollar", name)
}

func (s *ReserveDollarSuite) TestSymbol() {
	symbol, err := s.reserve.Symbol(nil)
	s.assert.NoError(err)
	s.assert.Equal("RSVD", symbol)
}

func (s *ReserveDollarSuite) TestDecimals() {
	decimals, err := s.reserve.Decimals(nil)
	s.assert.NoError(err)
	s.assert.Equal(uint8(18), decimals)
}

func (s *ReserveDollarSuite) TestChangeName() {
	const newName, newSymbol = "Flamingo", "MGO"
	s.requireTx(
		s.reserveImpl.ChangeName(s.signer, newName, newSymbol),
	)

	// Check for ChangeName event.
	nameChangeIter, err := s.reserve.FilterChangeName(nil)
	if s.assert.NoError(err) {
		events := 0
		for nameChangeIter.Next() {
			events++
			s.assert.Equal(newName, nameChangeIter.Event.NewName)
			s.assert.Equal(newSymbol, nameChangeIter.Event.NewSymbol)
		}
		s.assert.Equal(1, events, "expected exactly one ChangeName event")
		s.assert.NoError(nameChangeIter.Error())
		s.assert.NoError(nameChangeIter.Close())
	}

	// Check new name.
	name, err := s.reserve.Name(nil)
	s.assert.NoError(err)
	s.assert.Equal(newName, name)

	// Check new symbol.
	symbol, err := s.reserve.Symbol(nil)
	s.assert.NoError(err)
	s.assert.Equal(newSymbol, symbol)
}

func (s *ReserveDollarSuite) TestPausing() {
	// Pause.
	s.requireTx(s.reserveImpl.Pause(s.signer))

	// Minting is allowed while paused.
	amount := big.NewInt(100)
	s.requireTx(s.reserve.Mint(s.signer, toAddress(s.deployerKey), amount))

	// Transfers are not allowed while paused.
	recipient := common.BigToAddress(common.Big1)
	s.requireTxReverts(s.reserve.Transfer(s.signer, recipient, amount))
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

func toAddress(key *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(key.PublicKey)
}

func signer(key *ecdsa.PrivateKey) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key)
}
