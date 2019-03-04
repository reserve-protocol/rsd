package tests

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"

	"github.com/reserve-protocol/reserve-dollar/abi"
)

var delayInSeconds = big.NewInt(12 * 60 * 60)

func TestMintAndBurnAdmin(t *testing.T) {
	suite.Run(t, new(MintAndBurnAdminSuite))
}

type MintAndBurnAdminSuite struct {
	TestSuite

	adminContract        *abi.MintAndBurnAdmin
	adminContractAddress common.Address
	adminAccount         account
	adminSigner          *bind.TransactOpts

	utilContract *bind.BoundContract
}

var (
	// Compile-time check that MintAndBurnAdminSuite implements the interfaces we think it does.
	// If it does not implement these interfaces, then the corresponding setup and teardown
	// functions will not actually run.
	_ suite.BeforeTest    = &MintAndBurnAdminSuite{}
	_ suite.SetupAllSuite = &MintAndBurnAdminSuite{}
)

func (s *MintAndBurnAdminSuite) BlockTime() *big.Int {
	result := new(big.Int)
	s.NoError(s.utilContract.Call(nil, &result, "time"))
	return result
}

func (s *MintAndBurnAdminSuite) SetupSuite() {
	s.setup()

	if coverageEnabled {
		// Print warning that we don't support coverage.
		fmt.Fprintln(os.Stderr, "\nNOTE: Coverage information is not available for MintAndBurnAdmin, because its tests require faking the block timestamp")
		fmt.Fprintln(os.Stderr)
	}

	// Always use the fast node without coverage.
	s.createFastNode()

	// Create a utility contract to get the current block time.
	{
		// bytecode and utilABI are the result of compiling this Solidity file offline:
		/*
			pragma solidity ^0.5.4;
			contract Utility {
				function time() public view returns(uint256) {
					return now;
				}
			}
		*/
		bytecode := "0x6080604052348015600f57600080fd5b5060918061001e6000396000f3fe6080604052348015600f57600080fd5b50600436106044577c0100000000000000000000000000000000000000000000000000000000600035046316ada54781146049575b600080fd5b604f6061565b60408051918252519081900360200190f35b429056fea165627a7a723058205524d6a0c4d80ea5535c2ea64615c2619a21518e242cb929275cbd678b04468f0029"
		utilABI, err := ethabi.JSON(strings.NewReader(`
		[{"constant":true,"inputs":[],"name":"time","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]
		`))
		s.Require().NoError(err)

		code, err := hex.DecodeString(strings.TrimPrefix(bytecode, "0x"))
		s.Require().NoError(err)

		var tx *types.Transaction
		_, tx, s.utilContract, err = bind.DeployContract(s.signer, utilABI, code, s.node)
		s.requireTx(tx, err)( /* assert zero events */ )
	}
}

func (s *MintAndBurnAdminSuite) BeforeTest(suiteName, testName string) {
	// Deploy Reserve Dollar.
	reserveAddress, tx, reserve, err := abi.DeployReserveDollar(s.signer, s.node)
	s.requireTx(tx, err)( /* assert zero events */ )
	s.reserve = reserve
	s.reserveAddress = reserveAddress

	// Use our last test account as the owner of the admin contract.
	s.adminAccount = s.account[len(s.account)-1]
	s.adminSigner = signer(s.adminAccount)

	// Deploy admin contract.
	s.adminContractAddress, tx, s.adminContract, err = abi.DeployMintAndBurnAdmin(s.adminSigner, s.node, reserveAddress)
	s.requireTx(tx, err)( /* assert zero events */ )

	s.logParsers = map[common.Address]logParser{
		s.reserveAddress:       s.reserve,
		s.adminContractAddress: s.adminContract,
	}

	// Give minting power to admin contract.
	s.requireTx(reserve.ChangeMinter(s.signer, s.adminContractAddress))(
		abi.ReserveDollarMinterChanged{NewMinter: s.adminContractAddress},
	)
}

func (s *MintAndBurnAdminSuite) TestAdminContractIsMinter() {
	minter, err := s.reserve.Minter(nil)
	s.NoError(err)
	s.Equal(s.adminContractAddress, minter, "admin contract should have the minter role")
}

func (s *MintAndBurnAdminSuite) TestAdminCanMint() {
	recipient := common.BigToAddress(bigInt(1))
	amount := big.NewInt(100)

	// Propose a new mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))(
		s.proposalCreated(0, recipient, amount, true),
	)

	// Trying to confirm it immediately should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, bigInt(0), recipient, amount, true))

	// Advance time.
	s.Require().NoError(s.node.(backend).AdjustTime(13 * time.Hour))

	// Trying to confirm it should now succeed.
	s.requireTx(s.adminContract.Confirm(s.adminSigner, bigInt(0), recipient, amount, true))(
		s.proposalConfirmed(0, recipient, amount, true),
		mintingTransfer(recipient, amount),
	)

	// The mint should have happened.
	s.assertBalance(recipient, amount)

	// Trying to confirm a second time should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, bigInt(0), recipient, amount, true))
}

func (s *MintAndBurnAdminSuite) TestAdminCanCancelMinting() {
	recipient := common.BigToAddress(bigInt(1))
	amount := big.NewInt(100)

	// Propose a new mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))(
		s.proposalCreated(0, recipient, amount, true),
	)

	// And then cancel that minting.
	s.requireTx(s.adminContract.Cancel(s.adminSigner, bigInt(0), recipient, amount, true))(
		s.proposalCancelled(0, recipient, amount, true),
	)

	// Advance time.
	s.Require().NoError(s.node.(backend).AdjustTime(14 * time.Hour))

	// Trying to confirm it should now fail, even though time has advanced.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, bigInt(0), recipient, amount, true))

	// The mint should not have happened.
	s.assertBalance(recipient, bigInt(0))

	// Trying to confirm a second time should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, bigInt(0), recipient, amount, true))
}

func (s *MintAndBurnAdminSuite) TestConstructor() {
	addr, err := s.adminContract.Admin(nil)
	s.Nil(err)

	s.Equal(s.adminAccount.address(), addr, "admin should be set to the adminContractAddress")

	addr, err = s.adminContract.Reserve(nil)
	s.Nil(err)

	s.Equal(s.reserveAddress, addr, "reserve should be set to the Reserve contract address")
}

func (s *MintAndBurnAdminSuite) TestPropose() {
	recipient := common.BigToAddress(bigInt(1))
	amount := big.NewInt(100)
	futureAmount := big.NewInt(1000)

	// Retrieving the 0th proposal should fail when no proposals have been created yet.
	_, err := s.adminContract.Proposals(nil, bigInt(0))
	s.Error(err, "there should be no proposals on a newly-initialized MintAndBurnAdmin")

	// Trying to propose as someone other than the admin signer should fail.
	s.requireTxFails(s.adminContract.Propose(s.signer, recipient, amount, true))

	// Propose a new mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))(
		s.proposalCreated(0, recipient, amount, true),
	)

	// Retrieving the 0th proposal should now work.
	_, err = s.adminContract.Proposals(nil, bigInt(0))
	s.NoError(err)

	// Advance time by 12 hours.
	s.Require().NoError(s.node.(backend).AdjustTime(12 * time.Hour))

	// Propose a second mint.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, futureAmount, true))(
		s.proposalCreated(1, recipient, futureAmount, true),
	)

	// Proposals should now contain the first mint proposal at index 0
	proposalOne, err := s.adminContract.Proposals(nil, bigInt(0))
	s.Nil(err)
	proposalTwo, err := s.adminContract.Proposals(nil, bigInt(1))
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
	recipient := common.BigToAddress(bigInt(1))
	amount := big.NewInt(100)
	index := bigInt(0)

	// Create a proposal.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))(
		s.proposalCreated(0, recipient, amount, true),
	)

	// Trying to cancel as someone other than the admin should fail.
	s.requireTxFails(s.adminContract.Cancel(s.signer, index, recipient, amount, true))

	// Trying to cancel with mismatching recipient should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, common.BigToAddress(bigInt(2)), amount, true))

	// Trying to cancel with mismatching amount should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, recipient, bigInt(1), true))

	// Trying to cancel with mismatching isMint should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, recipient, amount, false))

	// Trying to cancel with mismatching index should fail.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, bigInt(1), recipient, amount, true))

	// Should be able to cancel proposal when supplied properly.
	s.requireTx(s.adminContract.Cancel(s.adminSigner, index, recipient, amount, true))(
		s.proposalCancelled(0, recipient, amount, true),
	)

	// Should be marked as closed.
	proposal, err := s.adminContract.Proposals(nil, index)
	s.NoError(err)
	s.True(proposal.Closed)

	// Should not be able to cancel a second time.
	s.requireTxFails(s.adminContract.Cancel(s.adminSigner, index, recipient, amount, true))
}

func (s *MintAndBurnAdminSuite) TestConfirm() {
	recipient := common.BigToAddress(bigInt(1))
	amount := big.NewInt(100)
	index := bigInt(0)

	// Create a proposal.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, amount, true))(
		s.proposalCreated(0, recipient, amount, true),
	)

	// Should not be able to confirm until time has passed.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, true))

	// Advance time.
	s.Require().NoError(s.node.(backend).AdjustTime(13 * time.Hour))

	// Trying to confirm as someone other than the admin should fail.
	s.requireTxFails(s.adminContract.Confirm(s.signer, index, recipient, amount, true))

	// Trying to confirm with mismatching recipient should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, common.BigToAddress(bigInt(2)), amount, true))

	// Trying to confirm with mismatching amount should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, bigInt(1), true))

	// Trying to confirm with mismatching isMint should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, false))

	// Trying to confirm with mismatching index should fail.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, bigInt(1), recipient, amount, true))

	// Confirm proposal.
	s.requireTx(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, true))(
		s.proposalConfirmed(0, recipient, amount, true),
		mintingTransfer(recipient, amount),
	)

	// Should be marked as closed.
	proposal, err := s.adminContract.Proposals(nil, index)
	s.NoError(err)
	s.True(proposal.Closed)

	// Should not be able to confirm proposal a second time.
	s.requireTxFails(s.adminContract.Confirm(s.adminSigner, index, recipient, amount, true))

	// Confirm mint went through.
	s.assertBalance(recipient, amount)
}

func (s *MintAndBurnAdminSuite) TestCancelAll() {
	// Create several proposals.
	for i := 0; i < 5; i++ {
		recipient := common.BigToAddress(bigInt(uint32(i)))
		value := bigInt(uint32((i + 1) * 100))
		s.requireTx(s.adminContract.Propose(s.adminSigner, recipient, value, i%2 == 0))(
			s.proposalCreated(i, recipient, value, i%2 == 0),
		)
	}

	// Confirm that CancelAll cannot be called by someone unauthorized.
	s.requireTxFails(s.adminContract.CancelAll(signer(s.account[2])))

	// Cancel all.
	s.requireTx(s.adminContract.CancelAll(s.adminSigner))(
		abi.MintAndBurnAdminAllProposalsCancelled{},
	)

	// Confirm that we can no longer retrieve the 0th proposal.
	_, err := s.adminContract.Proposals(nil, bigInt(0))
	s.Error(err, "should not be able to retrieve proposals after cancelling all")
}

func (s *MintAndBurnAdminSuite) TestBurn() {
	recipient := s.account[2]
	value := bigInt(1525)

	// Mint some tokens.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient.address(), value, true))(
		s.proposalCreated(0, recipient.address(), value, true),
	)
	s.Require().NoError(s.node.(backend).AdjustTime(13 * time.Hour))
	s.requireTx(s.adminContract.Confirm(s.adminSigner, bigInt(0), recipient.address(), value, true))(
		s.proposalConfirmed(0, recipient.address(), value, true),
		mintingTransfer(recipient.address(), value),
	)

	// Confirm the mint.
	s.assertBalance(recipient.address(), value)

	// Approve burning.
	s.requireTx(s.reserve.Approve(signer(recipient), s.adminContractAddress, value))(
		abi.ReserveDollarApproval{Holder: recipient.address(), Spender: s.adminContractAddress, Value: value},
	)

	// Burn the tokens.
	s.requireTx(s.adminContract.Propose(s.adminSigner, recipient.address(), value, false))(
		s.proposalCreated(1, recipient.address(), value, false),
	)
	s.Require().NoError(s.node.(backend).AdjustTime(13 * time.Hour))
	s.requireTx(s.adminContract.Confirm(s.adminSigner, bigInt(1), recipient.address(), value, false))(
		s.proposalConfirmed(1, recipient.address(), value, false),
		burningTransfer(recipient.address(), value),
		abi.ReserveDollarApproval{Holder: recipient.address(), Spender: s.adminContractAddress, Value: bigInt(0)},
	)

}

func (s *MintAndBurnAdminSuite) proposalCreated(
	i int, addr common.Address, value *big.Int, isMint bool,
) abi.MintAndBurnAdminProposalCreated {
	return abi.MintAndBurnAdminProposalCreated{
		Index:      bigInt(uint32(i)),
		Addr:       addr,
		Value:      value,
		IsMint:     isMint,
		DelayUntil: new(big.Int).Add(s.BlockTime(), delayInSeconds),
	}
}

func (s *MintAndBurnAdminSuite) proposalConfirmed(
	i int, addr common.Address, value *big.Int, isMint bool,
) abi.MintAndBurnAdminProposalConfirmed {
	return abi.MintAndBurnAdminProposalConfirmed{
		Index:  bigInt(uint32(i)),
		Addr:   addr,
		Value:  value,
		IsMint: isMint,
	}
}

func (s *MintAndBurnAdminSuite) proposalCancelled(
	i int, addr common.Address, value *big.Int, isMint bool,
) abi.MintAndBurnAdminProposalCancelled {
	return abi.MintAndBurnAdminProposalCancelled{
		Index:  bigInt(uint32(i)),
		Addr:   addr,
		Value:  value,
		IsMint: isMint,
	}
}
