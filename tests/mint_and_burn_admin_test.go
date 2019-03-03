package tests

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/reserve-protocol/reserve-dollar/abi"
)

func TestMintAndBurnAdmin(t *testing.T) {
	suite.Run(t, new(MintAndBurnAdminSuite))
}

type MintAndBurnAdminSuite struct {
	TestSuite

	adminContract        *abi.MintAndBurnAdmin
	adminContractAddress common.Address
	adminAccount         account
	adminSigner          *bind.TransactOpts
}

var (
	// Compile-time check that MintAndBurnAdminSuite implements the interfaces we think it does.
	// If it does not implement these interfaces, then the corresponding setup and teardown
	// functions will not actually run.
	_ suite.BeforeTest    = &MintAndBurnAdminSuite{}
	_ suite.SetupAllSuite = &MintAndBurnAdminSuite{}
)

func (s *MintAndBurnAdminSuite) SetupSuite() {
	s.setup()

	if testing.CoverMode() != "" {
		// Print warning that we don't support coverage.
		fmt.Fprintln(os.Stderr, "\nNOTE: Coverage information is not available for MintAndBurnAdmin, because its tests require faking the block timestamp")
		fmt.Fprintln(os.Stderr)
	}

	// Always use the fast node without coverage.
	s.createFastNode()
}

func (s *MintAndBurnAdminSuite) BeforeTest(suiteName, testName string) {
	// Deploy Reserve Dollar.
	reserveAddress, tx, reserve, err := abi.DeployReserveDollar(s.signer, s.node)
	s.requireTx(tx, err)
	s.reserve = reserve
	s.reserveAddress = reserveAddress

	// Use our last test account as the owner of the admin contract.
	s.adminAccount = s.account[len(s.account)-1]
	s.adminSigner = signer(s.adminAccount)

	// Deploy admin contract.
	s.adminContractAddress, tx, s.adminContract, err = abi.DeployMintAndBurnAdmin(s.adminSigner, s.node, reserveAddress)
	s.requireTx(tx, err)

	// Give minting power to admin contract.
	s.requireTx(reserve.ChangeMinter(s.signer, s.adminContractAddress))
}

func (s *MintAndBurnAdminSuite) TestAdminContractIsMinter() {
	minter, err := s.reserve.Minter(nil)
	s.NoError(err)
	s.Equal(s.adminContractAddress, minter, "admin contract should have the minter role")
}

func (s *MintAndBurnAdminSuite) TestAdminCanMint() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)

	// Propose a new mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))

	// Trying to confirm it immediately should fail.
	_, err := s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true)
	s.Require().Error(err)

	// Advance time.
	s.Require().NoError(s.node.(backend).AdjustTime(13 * time.Hour))

	// Trying to confirm it should now succeed.
	s.requireTx(s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true))

	// The mint should have happened.
	s.assertBalance(recipient, amount)

	// Trying to confirm a second time should fail.
	_, err = s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true)
	s.Error(err)
}

func (s *MintAndBurnAdminSuite) TestAdminCanCancelMinting() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)

	// Propose a new mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))

	// And then cancel that minting.
	s.requireTx(s.adminContract.Cancel(s.adminSigner, common.Big0, recipient, amount, true))

	// Advance time.
	s.Require().NoError(s.node.(backend).AdjustTime(14 * time.Hour))

	// Trying to confirm it should now fail, even though time has advanced.
	_, err := s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true)
	s.Require().Error(err)

	// The mint should not have happened.
	s.assertBalance(recipient, common.Big0)

	// Trying to confirm a second time should fail.
	_, err = s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true)
	s.Error(err)
}

// TODO: test changing admin

// Unit tests

func (s *MintAndBurnAdminSuite) TestConstructor() {
	addr, err := s.adminContract.Admin(nil)
	s.Nil(err)

	s.Equal(s.adminAccount.address(), addr, "admin should be set to the adminContractAddress")

	addr, err = s.adminContract.Reserve(nil)
	s.Nil(err)

	s.Equal(s.reserveAddress, addr, "reserve should be set to the Reserve contract address")
}

func (s *MintAndBurnAdminSuite) TestPropose() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)
	futureAmount := big.NewInt(1000)

	// nextProposal should be 0.
	nextProposal, err := s.adminContract.NextProposal(nil)
	s.Nil(err)
	s.Equal(common.Big0.String(), nextProposal.String(), "nextProposal should be 0")

	// Proposals should be empty at 0 index.
	proposal, err := s.adminContract.Proposals(nil, common.Big0)
	s.Nil(err)
	s.Equal(common.Address{}, proposal.Addr, "0th proposal address should be the zero value")
	s.Equal(common.Big0.String(), proposal.Value.String(), "0th proposal value should be the zero value")
	s.Equal(common.Big0.String(), proposal.Time.String(), "0th proposal should be zero value")
	s.Equal(false, proposal.IsMint, "0th proposal isMint should be zero value")

	// Trying to propose as someone other than the admin signer should fail.
	s.requireTxFails(s.adminContract.Propose(s.signer, recipient, amount, true))

	// Propose a new mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))

	// nextProposal should now be 1.
	nextProposal, err = s.adminContract.NextProposal(nil)
	s.Nil(err)
	s.Equal(common.Big1.String(), nextProposal.String(), "nextProposal should now be 1")

	// Advance time by 12 hours.
	s.Require().NoError(s.node.(backend).AdjustTime(12 * time.Hour))

	// Propose a second mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, futureAmount, true))

	// Proposals should now contain the first mint proposal at index 0
	proposalOne, err := s.adminContract.Proposals(nil, common.Big0)
	s.Nil(err)
	proposalTwo, err := s.adminContract.Proposals(nil, common.Big1)
	s.Nil(err)

	// Check proposal one is correct
	s.Equal(recipient, proposalOne.Addr, "0th proposal address should be recipient")
	s.Equal(amount.String(), proposalOne.Value.String(), "0th proposal value should be amount")
	s.Equal(true, proposalOne.IsMint, "0th proposal isMint should be true")

	// Check proposal two is correct
	s.Equal(recipient, proposalTwo.Addr, "1th proposal address should be recipient")
	s.Equal(futureAmount.String(), proposalTwo.Value.String(), "1th proposal value should be amount")
	s.Equal(true, proposalTwo.IsMint, "1th proposal isMint should be true")

	// Confirm times are separated by 12 hours plus the blocktime (20 s)
	diff := big.NewInt(0).Sub(proposalTwo.Time, proposalOne.Time)
	s.Equal(big.NewInt(12*3600+20).String(), diff.String(), "proposals should be separated by exactly 12 hours")
}

func (s *MintAndBurnAdminSuite) TestCancel() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)
	index := common.Big0

	// Create a proposal.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))

	// Trying to cancel as someone other than the admin should fail.
	s.requireTxFails(s.adminContract.Cancel(s.signer, index, recipient, amount, true))

	// Trying to cancel with mismatching recipient should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, common.BigToAddress(common.Big2), amount, true))

	// Trying to cancel with mismatching amount should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, recipient, common.Big1, true))

	// Trying to cancel with mismatching isMint should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, recipient, amount, false))

	// Trying to cancel with mismatching index should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, common.Big1, recipient, amount, true))

	// Should be able to cancel proposal when supplied properly.
	s.requireTx(s.adminContract.Cancel(s.adminSigner, index, recipient, amount, true))

	// Should be marked as completed.
	completed, err := s.adminContract.Completed(nil, index)
	s.Nil(err)
	s.Equal(true, completed)

	// Should not be able to cancel a second time.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, recipient, amount, true))
}

func (s *MintAndBurnAdminSuite) TestConfirm() {
	recipient := common.BigToAddress(common.Big1)
	amount := big.NewInt(100)
	index := common.Big0

	// Create a proposal.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))

	// Should not be able to confirm until time has passed.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, true))

	// Advance time.
	s.Require().NoError(s.node.(backend).AdjustTime(13 * time.Hour))

	// Trying to confirm as someone other than the admin should fail.
	s.requireTxFails(s.adminContract.Confirm(s.signer, index, recipient, amount, true))

	// Trying to confirm with mismatching recipient should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, common.BigToAddress(common.Big2), amount, true))

	// Trying to confirm with mismatching amount should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, common.Big1, true))

	// Trying to confirm with mismatching isMint should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, false))

	// Trying to confirm with mismatching index should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, common.Big1, recipient, amount, true))

	// Confirm proposal.
	s.requireTx(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, true))

	// Should be marked as completed.
	completed, err := s.adminContract.Completed(nil, index)
	s.Nil(err)
	s.Equal(true, completed)

	// Should not be able to confirm proposal a second time.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, true))

	// Confirm mint went through.
	s.assertBalance(recipient, amount)
}
