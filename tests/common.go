package tests

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"os"

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
	"github.com/reserve-protocol/reserve-dollar/soltools"
)

type TestSuite struct {
	suite.Suite

	account []account
	signer  *bind.TransactOpts
	node    interface {
		bind.ContractBackend
		TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	}
	reserve *abi.ReserveDollar
	assert  *assert.Assertions
	require *require.Assertions
}

// requireTx requires that a transaction is successfully mined and does
// not revert. It also takes an extra error argument, and checks that the
// error is nil. This signature allows the function to directly wrap our
// abigen'd mutator calls.
func (s *TestSuite) requireTx(tx *types.Transaction, err error) {
	s._requireTxStatus(tx, err, types.ReceiptStatusSuccessful)
}

// requireTxFails is like requireTx, but it requires that the transaction either
// reverts or is not successfully made in the first place due to gas estimation
// failing.
func (s *TestSuite) requireTxFails(tx *types.Transaction, err error) {
	if err != nil && err.Error() ==
		"failed to estimate gas needed: gas required exceeds allowance or always failing transaction" {
		return
	}

	fmt.Printf("%q\n", err.Error())

	s._requireTxStatus(tx, err, types.ReceiptStatusFailed)
}

func (s *TestSuite) _requireTxStatus(tx *types.Transaction, err error, status uint64) {
	s.Require().NoError(err)
	s.Require().NotNil(tx)
	receipt, err := bind.WaitMined(context.Background(), s.node, tx)
	s.Require().NoError(err)
	s.Require().Equal(status, receipt.Status)
}

func (s *TestSuite) assertBalance(address common.Address, amount *big.Int) {
	balance, err := s.reserve.BalanceOf(nil, address)
	s.NoError(err)
	s.Equal(amount.String(), balance.String()) // assert.Equal can mis-compare big.Ints, so compare strings instead
}

func (s *TestSuite) assertAllowance(owner, spender common.Address, amount *big.Int) {
	allowance, err := s.reserve.Allowance(nil, owner, spender)
	s.NoError(err)
	s.Equal(amount.String(), allowance.String())
}

func (s *TestSuite) assertTotalSupply(amount *big.Int) {
	totalSupply, err := s.reserve.TotalSupply(nil)
	s.NoError(err)
	s.Equal(amount.String(), totalSupply.String())
}

func (s *TestSuite) createSlowCoverageNode() {
	fmt.Fprintln(os.Stderr, "A local geth node must be running for coverage to work.")
	fmt.Fprintln(os.Stderr, "If one is not already running, start one in a new terminal with:")
	fmt.Fprintln(os.Stderr, "\tdocker run -it --rm -p 8545:8501 0xorg/devnet")

	var err error
	s.node, err = soltools.NewBackend("http://localhost:8545")
	s.Require().NoError(err)

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
		types.NewTransaction(0, common.Address{100}, common.Big0, 21000, common.Big1, nil),
		types.HomesteadSigner{},
		s.account[0].key,
	)
	s.node.SendTransaction(context.Background(), tx)
}

func (s *TestSuite) createFastNode() {
	genesisAlloc := core.GenesisAlloc{}
	for _, account := range s.account {
		genesisAlloc[account.address()] = core.GenesisAccount{
			Balance: big.NewInt(math.MaxInt64),
		}
	}
	s.node = backend{
		backends.NewSimulatedBackend(
			genesisAlloc,
			// TODO: the tests fail if this is 4e6. why?
			8e6, // roughly same order of magnitude as mainnet
		),
	}
}

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
