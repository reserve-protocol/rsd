package tests

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/reserve-protocol/reserve-dollar/abi"
	"github.com/reserve-protocol/reserve-dollar/soltools"
)

func TestReserveDollar(t *testing.T) {
	suite.Run(t, new(ReserveDollarSuite))
}

type ReserveDollarSuite struct {
	TestSuite
}

var (
	// Compile-time check that ReserveDollarSuite implements the interfaces we think it does.
	// If it does not implement these interfaces, then the corresponding setup and teardown
	// functions will not actually run.
	_ suite.BeforeTest       = &ReserveDollarSuite{}
	_ suite.SetupAllSuite    = &ReserveDollarSuite{}
	_ suite.TearDownAllSuite = &ReserveDollarSuite{}
)

// SetupSuite runs once, before all of the tests in the suite.
func (s *ReserveDollarSuite) SetupSuite() {
	s.setup()
	if testing.CoverMode() == "" {
		s.createFastNode()
	} else {
		s.createSlowCoverageNode()
	}
}

// TearDownSuite runs once, after all of the tests in the suite.
func (s *ReserveDollarSuite) TearDownSuite() {
	if testing.CoverMode() != "" {
		// Write coverage profile to disk.
		s.Assert().NoError(s.node.(*soltools.Backend).WriteCoverage())

		// Close the node.js process.
		s.Assert().NoError(s.node.(*soltools.Backend).Close())

		// Process coverage profile into an HTML report.
		if out, err := exec.Command("npx", "istanbul", "report", "html").CombinedOutput(); err != nil {
			fmt.Println()
			fmt.Println("I generated coverage information in coverage/coverage.json.")
			fmt.Println("I tried to process it with `istanbul` to turn it into a readable report, but failed.")
			fmt.Println("The error I got when running istanbul was:", err)
			fmt.Println("Istanbul's output was:\n" + string(out))
		}
	}
}

// Before test runs before each test in the suite.
func (s *ReserveDollarSuite) BeforeTest(suiteName, testName string) {
	reserveAddress, tx, reserve, err := abi.DeployReserveDollar(s.signer, s.node)
	s.requireTx(tx, err)

	s.reserve = reserve
	s.reserveAddress = reserveAddress

	// Make the deployment account a minter, pauser, and freezer.
	s.requireTx(s.reserve.ChangeMinter(s.signer, s.account[0].address()))
	s.requireTx(s.reserve.ChangePauser(s.signer, s.account[0].address()))
	s.requireTx(s.reserve.ChangeFreezer(s.signer, s.account[0].address()))
}

func (s *ReserveDollarSuite) TestDeploy() {}

func (s *ReserveDollarSuite) TestBalanceOf() {
	s.assertBalance(common.Address{}, bigInt(0))
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
	)(
		abi.ReserveDollarNameChanged{
			NewName:   newName,
			NewSymbol: newSymbol,
		},
	)

	// Check new name.
	name, err := s.reserve.Name(nil)
	s.NoError(err)
	s.Equal(newName, name)

	// Check new symbol.
	symbol, err := s.reserve.Symbol(nil)
	s.NoError(err)
	s.Equal(newSymbol, symbol)
}

func (s *ReserveDollarSuite) TestChangeNameFailsForNonOwner() {
	const newName, newSymbol = "Flamingo", "MGO"
	s.requireTxFails(
		s.reserve.ChangeName(signer(s.account[2]), newName, newSymbol),
	)
}

func (s *ReserveDollarSuite) TestAllowsMinting() {
	recipient := common.BigToAddress(bigInt(1))
	amount := bigInt(100)

	// Mint to recipient.
	s.requireTx(s.reserve.Mint(s.signer, recipient, amount))

	// Check that balances are as expected.
	s.assertBalance(s.account[0].address(), bigInt(0))
	s.assertBalance(recipient, amount)
	s.assertTotalSupply(amount)
}

func (s *ReserveDollarSuite) TestTransfer() {
	sender := s.account[1]
	recipient := common.BigToAddress(bigInt(1))
	amount := bigInt(100)

	// Mint to sender.
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), amount))

	// Transfer from sender to recipient.
	s.requireTx(s.reserve.Transfer(signer(sender), recipient, amount))

	// Check that balances are as expected.
	s.assertBalance(sender.address(), bigInt(0))
	s.assertBalance(recipient, amount)
	s.assertBalance(s.account[0].address(), bigInt(0))
	s.assertTotalSupply(amount)

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(sender.address(), amount),
		{From: sender.address(), To: recipient, Value: amount},
	})

	s.assertZeroApprovalEvents()
}

func (s *ReserveDollarSuite) TestTransferExceedsFunds() {
	sender := s.account[1]
	recipient := common.BigToAddress(bigInt(1))
	amount := bigInt(100)
	smallAmount := bigInt(10) // must be smaller than amount

	// Mint smallAmount to sender.
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), smallAmount))

	// Transfer from sender to recipient should fail.
	s.requireTxFails(s.reserve.Transfer(signer(sender), recipient, amount))

	// Balances should be as we expect.
	s.assertBalance(sender.address(), smallAmount)
	s.assertBalance(recipient, bigInt(0))
	s.assertBalance(s.account[0].address(), bigInt(0))
	s.assertTotalSupply(smallAmount)

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(sender.address(), smallAmount),
	})
	s.assertZeroApprovalEvents()
}

// As long as Minting cannot overflow a uint256, then `transferFrom` cannot overflow.
func (s *ReserveDollarSuite) TestMintWouldOverflow() {
	interestingRecipients := []common.Address{
		common.BigToAddress(bigInt(1)),
		common.BigToAddress(bigInt(255)),
		common.BigToAddress(bigInt(256)),
		common.BigToAddress(bigInt(256)),
		common.BigToAddress(maxUint160()),
		common.BigToAddress(minInt160AsUint160()),
	}
	for _, recipient := range interestingRecipients {
		smallAmount := bigInt(10) // must be smaller than amount
		overflowCausingAmount := maxUint256()
		overflowCausingAmount = overflowCausingAmount.Sub(overflowCausingAmount, bigInt(8))

		// Mint smallAmount to recipient.
		s.requireTx(s.reserve.Mint(s.signer, recipient, smallAmount))

		// Mint a quantity large enough to cause overflow in totalSupply i.e.
		// `10 + (uint256::MAX - 8) > uint256::MAX`
		s.requireTxFails(s.reserve.Mint(s.signer, recipient, overflowCausingAmount))
	}
}

func (s *ReserveDollarSuite) TestApprove() {
	owner := s.account[1]
	spender := s.account[2]
	amount := bigInt(53)

	// Owner approves spender.
	s.requireTx(s.reserve.Approve(signer(owner), spender.address(), amount))

	// Approval should be reflected in allowance.
	s.assertAllowance(owner.address(), spender.address(), amount)

	// Shouldn't be symmetric.
	s.assertAllowance(spender.address(), owner.address(), bigInt(0))

	// Balances shouldn't change.
	s.assertBalance(owner.address(), bigInt(0))
	s.assertBalance(spender.address(), bigInt(0))
	s.assertTotalSupply(bigInt(0))
}

func (s *ReserveDollarSuite) TestIncreaseAllowance() {
	owner := s.account[1]
	spender := s.account[2]
	amount := bigInt(2000)

	// Owner approves spender through increaseAllowance.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), amount))

	// Approval should be reflected in allowance.
	s.assertAllowance(owner.address(), spender.address(), amount)

	// Shouldn't be symmetric.
	s.assertAllowance(spender.address(), owner.address(), bigInt(0))

	// Check for Approval event.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: spender.address(), Value: amount},
	})

	// Balances shouldn't change.
	s.assertBalance(owner.address(), bigInt(0))
	s.assertBalance(spender.address(), bigInt(0))
	s.assertTotalSupply(bigInt(0))
}

func (s *ReserveDollarSuite) TestIncreaseAllowanceWouldOverflow() {
	owner := s.account[1]
	spender := s.account[2]
	initialAmount := bigInt(10)

	// Owner approves spender for initial amount.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// Check for Approval event.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: spender.address(), Value: initialAmount},
	})

	// Owner should not be able to increase approval high enough to overflow a uint256.
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(owner), spender.address(), maxUint256()))

	// Check that no additional Approval event was emitted.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: spender.address(), Value: initialAmount},
	})
}

func (s *ReserveDollarSuite) TestDecreaseAllowance() {
	owner := s.account[1]
	spender := s.account[2]
	initialAmount := bigInt(10)
	decrease := bigInt(6)
	final := bigInt(4)

	// Owner approves spender for initial amount.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// Owner decreases allowance.
	s.requireTx(s.reserve.DecreaseAllowance(signer(owner), spender.address(), decrease))

	// Allowance should be as we expect.
	s.assertAllowance(owner.address(), spender.address(), final)

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: spender.address(), Value: initialAmount},
		{From: owner.address(), To: spender.address(), Value: final},
	})

	// Balances shouldn't change.
	s.assertBalance(owner.address(), bigInt(0))
	s.assertBalance(spender.address(), bigInt(0))
	s.assertTotalSupply(bigInt(0))
}

func (s *ReserveDollarSuite) TestDecreaseAllowanceUnderflow() {
	owner := s.account[1]
	spender := s.account[2]
	initialAmount := bigInt(10)
	decrease := bigInt(11)

	// Owner approves spender for initial amount.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// Owner decreases allowance fails because of underflow.
	s.requireTxFails(s.reserve.DecreaseAllowance(signer(owner), spender.address(), decrease))

	// Allowance should be as we expect.
	s.assertAllowance(owner.address(), spender.address(), initialAmount)

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: spender.address(), Value: initialAmount},
	})

	// Balances shouldn't change.
	s.assertBalance(owner.address(), bigInt(0))
	s.assertBalance(spender.address(), bigInt(0))
	s.assertTotalSupply(bigInt(0))

	// Allowances shouldn't change
	s.assertAllowance(owner.address(), spender.address(), initialAmount)
}

func (s *ReserveDollarSuite) TestDecreaseAllowanceSpenderFrozen() {
	spender := s.account[1]
	owner := s.account[2]

	// Owner approves spender for initial amount.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), bigInt(10)))

	// Freeze spender.
	s.requireTx(s.reserve.Freeze(s.signer, spender.address()))

	// Owner decreases allowance fails because of spender is frozen.
	s.requireTxFails(s.reserve.DecreaseAllowance(signer(owner), spender.address(), bigInt(2)))
}

func (s *ReserveDollarSuite) TestPausing() {
	banker := s.account[1]
	amount := bigInt(1000)
	approveAmount := bigInt(1)
	recipient := s.account[2]
	spender := s.account[3]

	// Give banker funds. Minting is allowed while unpaused.
	s.requireTx(s.reserve.Mint(s.signer, banker.address(), amount))
	s.assertBalance(banker.address(), amount)

	// Approve spender to spend bankers funds.
	s.requireTx(s.reserve.Approve(signer(banker), spender.address(), approveAmount))
	s.assertAllowance(banker.address(), spender.address(), approveAmount)

	// Pause.
	s.requireTx(s.reserve.Pause(s.signer))

	// Minting is not allowed while paused.
	s.requireTxFails(s.reserve.Mint(s.signer, recipient.address(), amount))

	// Transfers from are not allowed while paused.
	s.requireTxFails(s.reserve.TransferFrom(s.signer, banker.address(), recipient.address(), amount))
	s.assertBalance(recipient.address(), bigInt(0))
	s.assertBalance(banker.address(), amount)

	// Transfers are not allowed while paused.
	s.requireTxFails(s.reserve.Transfer(signer(banker), recipient.address(), amount))
	s.assertBalance(recipient.address(), bigInt(0))
	s.assertBalance(banker.address(), amount)

	// Approving is not allowed while paused.
	s.requireTxFails(s.reserve.Approve(signer(banker), spender.address(), amount))
	s.assertAllowance(banker.address(), spender.address(), approveAmount)

	// IncreaseAllowance is not allowed while paused.
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(banker), spender.address(), amount))
	s.assertAllowance(banker.address(), spender.address(), approveAmount)

	// DecreaseAllowance is not allowed while paused.
	s.requireTxFails(s.reserve.DecreaseAllowance(signer(banker), spender.address(), amount))
	s.assertAllowance(banker.address(), spender.address(), approveAmount)

	// Unpause.
	s.requireTx(s.reserve.Unpause(s.signer))

	// Transfers are allowed while unpaused.
	s.requireTx(s.reserve.Transfer(signer(banker), recipient.address(), amount))
	s.assertBalance(recipient.address(), amount)

	// Approving is allowed while unpaused.
	s.requireTx(s.reserve.Approve(signer(banker), spender.address(), bigInt(2)))
	s.assertAllowance(banker.address(), spender.address(), bigInt(2))

	// DecreaseAllowance is allowed while unpaused.
	s.requireTx(s.reserve.DecreaseAllowance(signer(banker), spender.address(), approveAmount))
	s.assertAllowance(banker.address(), spender.address(), approveAmount)

	// IncreaseAllowance is allowed while unpaused.
	s.requireTx(s.reserve.IncreaseAllowance(signer(banker), spender.address(), approveAmount))
	s.assertAllowance(banker.address(), spender.address(), bigInt(2))

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: banker.address(), To: spender.address(), Value: approveAmount},
		{From: banker.address(), To: spender.address(), Value: bigInt(2)},
		{From: banker.address(), To: spender.address(), Value: bigInt(1)},
		{From: banker.address(), To: spender.address(), Value: bigInt(2)},
	})
	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(banker.address(), amount),
		{From: banker.address(), To: recipient.address(), Value: amount},
	})
}

func (s *ReserveDollarSuite) TestFreezeTransferOut() {
	target := s.account[1]
	recipient := s.account[2]

	// Give target funds.
	amount := bigInt(1)
	s.requireTx(s.reserve.Mint(s.signer, target.address(), amount))

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Frozen account shouldn't be able to transfer.
	s.requireTxFails(s.reserve.Transfer(signer(target), recipient.address(), amount))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Unfrozen account should be able to transfer again.
	s.requireTx(s.reserve.Transfer(signer(target), recipient.address(), amount))
	s.assertBalance(recipient.address(), amount)

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(target.address(), amount),
		{From: target.address(), To: recipient.address(), Value: amount},
	})

	s.assertZeroApprovalEvents()
}

func (s *ReserveDollarSuite) TestFreezeTransferIn() {
	target := s.account[1]
	amount := bigInt(200)

	// Mint initial funds to deployer.
	s.requireTx(s.reserve.Mint(s.signer, s.account[0].address(), amount))

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Frozen account shouldn't be able to receive funds.
	s.requireTxFails(s.reserve.Transfer(s.signer, target.address(), amount))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Frozen account should be able to receive funds again.
	s.requireTx(s.reserve.Transfer(s.signer, target.address(), amount))
	s.assertBalance(target.address(), amount)

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(s.account[0].address(), amount),
		{From: s.account[0].address(), To: target.address(), Value: amount},
	})

	s.assertZeroApprovalEvents()
}

func (s *ReserveDollarSuite) TestFreezeApprovals() {
	target := s.account[1]
	recipient := s.account[2]

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Frozen account shouldn't be able to create approvals.
	s.requireTxFails(s.reserve.Approve(signer(target), recipient.address(), bigInt(1)))
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(target), recipient.address(), bigInt(1)))
	s.assertAllowance(target.address(), recipient.address(), bigInt(0))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Unfrozen account should be able to create approvals again.
	s.requireTx(s.reserve.Approve(signer(target), recipient.address(), bigInt(1)))
	s.requireTx(s.reserve.IncreaseAllowance(signer(target), recipient.address(), bigInt(1)))
	s.assertAllowance(target.address(), recipient.address(), bigInt(2))

	// Freeze recipient.
	s.requireTx(s.reserve.Freeze(s.signer, recipient.address()))

	// Frozen recipient should not be able to receive approvals.
	s.requireTxFails(s.reserve.Approve(signer(target), recipient.address(), bigInt(1)))
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(target), recipient.address(), bigInt(1)))
	s.assertAllowance(target.address(), recipient.address(), bigInt(2))

	// Unfreeze recipient.
	s.requireTx(s.reserve.Unfreeze(s.signer, recipient.address()))

	// Unfrozen account should be able to receive approvals again.
	s.requireTx(s.reserve.Approve(signer(target), recipient.address(), bigInt(11)))
	s.requireTx(s.reserve.IncreaseAllowance(signer(target), recipient.address(), bigInt(7)))
	s.assertAllowance(target.address(), recipient.address(), bigInt(18))

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: target.address(), To: recipient.address(), Value: bigInt(1)},
		{From: target.address(), To: recipient.address(), Value: bigInt(2)},
		{From: target.address(), To: recipient.address(), Value: bigInt(11)},
		{From: target.address(), To: recipient.address(), Value: bigInt(18)},
	})
}

func (s *ReserveDollarSuite) TestFreezeTransferFrom() {
	target := s.account[1]
	recipient := s.account[2]
	middleman := s.account[3]

	// Approve target and middleman to transfer funds.
	initialAmount := bigInt(12)
	s.requireTx(s.reserve.Mint(s.signer, s.account[0].address(), initialAmount))
	s.requireTx(s.reserve.Approve(s.signer, target.address(), initialAmount))
	s.requireTx(s.reserve.Approve(s.signer, middleman.address(), initialAmount))
	s.assertAllowance(s.account[0].address(), target.address(), initialAmount)
	s.assertAllowance(s.account[0].address(), middleman.address(), initialAmount)

	////////////////////////////////////
	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Frozen account shouldn't be able to call transferFrom.
	s.requireTxFails(s.reserve.TransferFrom(signer(target), s.account[0].address(), recipient.address(), initialAmount))
	s.assertBalance(recipient.address(), bigInt(0))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Unfrozen account should now be able to call transferFrom.
	s.requireTx(s.reserve.TransferFrom(signer(target), s.account[0].address(), recipient.address(), bigInt(2)))
	s.assertBalance(recipient.address(), bigInt(2))

	////////////////////////////////////
	// Freeze middleman
	s.requireTx(s.reserve.Freeze(s.signer, middleman.address()))

	// Frozen account shouldn't be able to call transferFrom.
	s.requireTxFails(s.reserve.TransferFrom(signer(middleman), s.account[0].address(), recipient.address(), bigInt(5)))
	s.assertBalance(recipient.address(), bigInt(2))

	// Unfreeze middleman.
	s.requireTx(s.reserve.Unfreeze(s.signer, middleman.address()))

	// Unfrozen account should now be able to call transferFrom.
	s.requireTx(s.reserve.TransferFrom(signer(middleman), s.account[0].address(), recipient.address(), bigInt(5)))
	s.assertBalance(recipient.address(), bigInt(7))

	////////////////////////////////////
	// Freeze recipient.
	s.requireTx(s.reserve.Freeze(s.signer, recipient.address()))

	// Shouldn't be able to call transferFrom to a frozen account.
	s.requireTxFails(s.reserve.TransferFrom(signer(target), s.account[0].address(), recipient.address(), initialAmount))
	s.assertBalance(recipient.address(), bigInt(7))

	// Unfreeze recipient.
	s.requireTx(s.reserve.Unfreeze(s.signer, recipient.address()))

	// Unfrozen account should now be able to call transferFrom.
	s.requireTx(s.reserve.TransferFrom(signer(target), s.account[0].address(), recipient.address(), bigInt(3)))
	s.assertBalance(recipient.address(), bigInt(10))

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(s.account[0].address(), initialAmount),
		{From: s.account[0].address(), To: recipient.address(), Value: bigInt(2)},
		{From: s.account[0].address(), To: recipient.address(), Value: bigInt(5)},
		{From: s.account[0].address(), To: recipient.address(), Value: bigInt(3)},
	})

	// Check for Approval events.

	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: s.account[0].address(), To: target.address(), Value: initialAmount},
		{From: s.account[0].address(), To: middleman.address(), Value: initialAmount},
		{From: s.account[0].address(), To: target.address(), Value: bigInt(12 - 2)},
		{From: s.account[0].address(), To: middleman.address(), Value: bigInt(12 - 5)},
		{From: s.account[0].address(), To: target.address(), Value: bigInt(12 - 2 - 3)},
	})
}

func (s *ReserveDollarSuite) TestFreezeApprove() {
	target := s.account[1]

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Should not be able to approve frozen target.
	s.requireTxFails(s.reserve.Approve(s.signer, target.address(), bigInt(1)))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Should be able to approve unfrozen target.
	s.requireTx(s.reserve.Approve(s.signer, target.address(), bigInt(1)))
}

func (s *ReserveDollarSuite) TestFreezeIncreaseAllowance() {
	target := s.account[1]
	owner := s.account[2]

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Should not be able to increase allowance frozen target.
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(owner), target.address(), bigInt(1)))
	s.assertAllowance(owner.address(), target.address(), bigInt(0))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Should be able to increase allowance unfrozen target.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), target.address(), bigInt(1)))
	s.assertAllowance(owner.address(), target.address(), bigInt(1))

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: target.address(), Value: bigInt(1)},
	})
}

func (s *ReserveDollarSuite) TestFreezeDecreaseAllowance() {
	spender := s.account[1]
	owner := s.account[2]

	// Increase allowance to set up for decrease.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), bigInt(6)))

	// Freeze spender.
	s.requireTx(s.reserve.Freeze(s.signer, spender.address()))

	// Should not be able to decrease allowance frozen spender.
	s.requireTxFails(s.reserve.DecreaseAllowance(signer(owner), spender.address(), bigInt(4)))
	s.assertAllowance(owner.address(), spender.address(), bigInt(6))

	// Unfreeze spender.
	s.requireTx(s.reserve.Unfreeze(s.signer, spender.address()))

	// Should be able to decrease allowance unfrozen spender.
	s.requireTx(s.reserve.DecreaseAllowance(signer(owner), spender.address(), bigInt(1)))
	s.assertAllowance(owner.address(), spender.address(), bigInt(5))

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: owner.address(), To: spender.address(), Value: bigInt(6)},
		{From: owner.address(), To: spender.address(), Value: bigInt(5)},
	})
}

func (s *ReserveDollarSuite) TestWiping() {
	target := s.account[1]

	// Give target funds.
	amount := bigInt(100)
	s.requireTx(s.reserve.Mint(s.signer, target.address(), amount))

	// Should not be able to wipe zero address
	s.requireTx(s.reserve.Freeze(s.signer, common.BigToAddress(bigInt(0))))
	s.requireTxFails(s.reserve.Wipe(s.signer, target.address()))

	// Should not be able to wipe target before freezing them.
	s.requireTxFails(s.reserve.Wipe(s.signer, target.address()))

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Target should still have funds.
	s.assertBalance(target.address(), amount)

	// Should not be able to immediately wipe target.
	s.requireTxFails(s.reserve.Wipe(s.signer, target.address()))

	if simulatedBackend, ok := s.node.(backend); ok {
		// Simulate advancing time.
		simulatedBackend.AdjustTime(24 * time.Hour * 40)

		// Should be able to wipe target now.
		s.requireTx(s.reserve.Wipe(s.signer, target.address()))

		// Target should have zero funds.
		s.assertBalance(target.address(), bigInt(0))
	} else {
		fmt.Fprintln(os.Stderr, "\nCan't simulate advancing time in coverage mode -- not testing wiping after a delay.")
		fmt.Fprintln(os.Stderr)
	}
}

func (s *ReserveDollarSuite) TestMintingBurningChain() {
	// Mint to recipient.
	recipient := s.account[1]
	amount := bigInt(100)

	s.requireTx(s.reserve.Mint(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), amount)
	s.assertTotalSupply(amount)

	// Approve signer for burning.
	s.requireTx(s.reserve.Approve(signer(recipient), s.account[0].address(), amount))

	// Burn from recipient.
	s.requireTx(s.reserve.BurnFrom(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), bigInt(0))
	s.assertTotalSupply(bigInt(0))
}

func (s *ReserveDollarSuite) TestMintingTransferBurningChain() {
	recipient := s.account[1]
	amount := bigInt(100)

	// Mint to recipient.
	s.requireTx(s.reserve.Mint(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), amount)
	s.assertTotalSupply(amount)

	// Transfer to target.
	target := s.account[2]
	s.requireTx(s.reserve.Transfer(signer(recipient), target.address(), amount))

	s.assertBalance(target.address(), amount)
	s.assertBalance(recipient.address(), bigInt(0))

	// Approve signer for burning.
	s.requireTx(s.reserve.Approve(signer(target), s.account[0].address(), amount))

	// Burn from target.
	s.requireTx(s.reserve.BurnFrom(s.signer, target.address(), amount))

	s.assertBalance(target.address(), bigInt(0))
	s.assertBalance(recipient.address(), bigInt(0))
	s.assertTotalSupply(bigInt(0))

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(recipient.address(), amount),
		{From: recipient.address(), To: target.address(), Value: amount},
		{From: target.address(), To: common.BigToAddress(bigInt(0)), Value: amount},
	})

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: target.address(), To: s.account[0].address(), Value: amount},
		{From: target.address(), To: s.account[0].address(), Value: bigInt(0)},
	})
}

func (s *ReserveDollarSuite) TestBurnFromWouldUnderflow() {
	// Mint to recipient.
	recipient := s.account[1]
	amount := bigInt(100)
	causesUnderflowAmount := bigInt(101)

	s.assertTotalSupply(bigInt(0))
	s.requireTx(s.reserve.Mint(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), amount)
	s.assertTotalSupply(amount)

	// Approve signer for burning.
	s.requireTx(s.reserve.Approve(signer(recipient), s.account[0].address(), amount))

	// Burn from recipient.
	s.requireTxFails(s.reserve.BurnFrom(s.signer, recipient.address(), causesUnderflowAmount))

	s.assertBalance(recipient.address(), amount)
	s.assertTotalSupply(amount)
}

func (s *ReserveDollarSuite) TestTransferFrom() {
	sender := s.account[1]
	middleman := s.account[2]
	recipient := s.account[3]

	amount := bigInt(1)
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), amount))
	s.assertBalance(sender.address(), amount)
	s.assertBalance(middleman.address(), bigInt(0))
	s.assertBalance(recipient.address(), bigInt(0))
	s.assertTotalSupply(amount)

	// Approve middleman to transfer funds from the sender.
	s.requireTx(s.reserve.Approve(signer(sender), middleman.address(), amount))
	s.assertAllowance(sender.address(), middleman.address(), amount)

	// transferFrom allows the msg.sender to send an existing approval to an arbitrary destination.
	s.requireTx(s.reserve.TransferFrom(signer(middleman), sender.address(), recipient.address(), amount))
	s.assertBalance(sender.address(), bigInt(0))
	s.assertBalance(middleman.address(), bigInt(0))
	s.assertBalance(recipient.address(), amount)

	// Check for Transfer event.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(sender.address(), amount),
		{From: sender.address(), To: recipient.address(), Value: amount},
	})

	// Allowance should have been decreased by the transfer
	s.assertAllowance(sender.address(), middleman.address(), bigInt(0))
	// transfers should not change totalSupply.
	s.assertTotalSupply(amount)
}

func (s *ReserveDollarSuite) TestTransferFromWouldUnderflow() {
	sender := s.account[1]
	middleman := s.account[2]
	recipient := s.account[3]

	approveAmount := bigInt(2)
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), approveAmount))
	s.assertBalance(sender.address(), approveAmount)
	s.assertBalance(middleman.address(), bigInt(0))
	s.assertBalance(recipient.address(), bigInt(0))
	s.assertTotalSupply(approveAmount)

	// Approve middleman to transfer funds from the sender.
	s.requireTx(s.reserve.Approve(signer(sender), middleman.address(), approveAmount))
	s.assertAllowance(sender.address(), middleman.address(), approveAmount)

	// now reduce the approveAmount in the sender's account to less than the approval for the middleman
	s.requireTx(s.reserve.Transfer(signer(sender), recipient.address(), bigInt(1)))

	// Attempt to transfer more funds than the sender's current balance, but
	// passing the approval checks. Should fail when subtracting value from
	// sender's current balance.
	s.requireTxFails(s.reserve.TransferFrom(signer(middleman), sender.address(), recipient.address(), approveAmount))

	s.assertBalance(sender.address(), bigInt(1))
	s.assertBalance(middleman.address(), bigInt(0))
	s.assertBalance(recipient.address(), bigInt(1))

	// Allowance should not have been changed
	s.assertAllowance(sender.address(), middleman.address(), approveAmount)
	// should not change totalSupply.
	s.assertTotalSupply(approveAmount)

	// Check for Transfer events.
	s.assertTransferEvents([]ReserveDollarTransfer{
		mintingTransfer(sender.address(), approveAmount),
		{From: sender.address(), To: recipient.address(), Value: bigInt(1)},
	})

	// Check for Approval events.
	s.assertApprovalEvents([]ReserveDollarApproval{
		{From: sender.address(), To: middleman.address(), Value: approveAmount},
	})
}

///////////////////////

func (s *ReserveDollarSuite) TestPauseFailsForNonPauser() {
	s.requireTxFails(s.reserve.Pause(signer(s.account[2])))
}

func (s *ReserveDollarSuite) TestUnpauseFailsForNonPauser() {
	s.requireTx(s.reserve.Pause(s.signer))
	s.requireTxFails(s.reserve.Unpause(signer(s.account[1])))
}

func (s *ReserveDollarSuite) TestChangePauserFailsForNonPauser() {
	s.requireTxFails(s.reserve.ChangePauser(signer(s.account[2]), s.account[1].address()))
}

///////////////////////

func (s *ReserveDollarSuite) TestFreezeFailsForNonFreezer() {
	criminal := common.BigToAddress(bigInt(1))
	s.requireTxFails(s.reserve.Freeze(signer(s.account[2]), criminal))
}

func (s *ReserveDollarSuite) TestUnfreezeFailsForNonFreezer() {
	criminal := common.BigToAddress(bigInt(1))
	s.requireTx(s.reserve.Freeze(s.signer, criminal))
	s.requireTxFails(s.reserve.Unfreeze(signer(s.account[1]), criminal))
}

func (s *ReserveDollarSuite) TestChangeFreezerFailsForNonFreezer() {
	s.requireTxFails(s.reserve.ChangeFreezer(signer(s.account[2]), s.account[1].address()))
}

func (s *ReserveDollarSuite) TestWipeFailsForNonFreezer() {
	criminal := common.BigToAddress(bigInt(1))
	s.requireTxFails(s.reserve.Wipe(signer(s.account[2]), criminal))
}

///////////////////////

func (s *ReserveDollarSuite) TestMintFailsForNonMinter() {
	recipient := common.BigToAddress(bigInt(1))
	s.requireTxFails(s.reserve.Mint(signer(s.account[2]), recipient, bigInt(7)))
}

func (s *ReserveDollarSuite) TestChangeMinterFailsForNonMinter() {
	s.requireTxFails(s.reserve.ChangeMinter(signer(s.account[2]), s.account[1].address()))
}

///////////////////////

func (s *ReserveDollarSuite) TestUpgrade() {
	recipient := s.account[1]
	amount := big.NewInt(100)

	// Mint to recipient.
	s.requireTx(s.reserve.Mint(s.signer, recipient.address(), amount))

	// Deploy new contract.
	newKey := s.account[2]
	newTokenAddress, tx, newToken, err := abi.DeployReserveDollarV2(signer(newKey), s.node)
	s.requireTx(tx, err)

	// Make the switch.
	s.requireTx(s.reserve.NominateNewOwner(s.signer, newTokenAddress))
	s.requireTx(newToken.CompleteHandoff(signer(newKey), s.reserveAddress))

	// Old token should not be functional.
	s.requireTxFails(s.reserve.Mint(s.signer, recipient.address(), big.NewInt(1500)))
	s.requireTxFails(s.reserve.Transfer(signer(recipient), s.account[3].address(), big.NewInt(10)))
	s.requireTxFails(s.reserve.Pause(s.signer))
	s.requireTxFails(s.reserve.Unpause(s.signer))

	// assertion function for new token
	assertBalance := func(address common.Address, amount *big.Int) {
		balance, err := newToken.BalanceOf(nil, address)
		s.NoError(err)
		s.Equal(amount.String(), balance.String()) // assert.Equal can mis-compare big.Ints, so compare strings instead
	}

	// New token should be functional.
	assertBalance(recipient.address(), amount)
	s.requireTx(newToken.ChangeMinter(signer(newKey), newKey.address()))
	s.requireTx(newToken.ChangePauser(signer(newKey), newKey.address()))
	s.requireTx(newToken.Mint(signer(newKey), recipient.address(), big.NewInt(1500)))
	s.requireTx(newToken.Transfer(signer(recipient), s.account[3].address(), big.NewInt(10)))
	s.requireTx(newToken.Pause(signer(newKey)))
	s.requireTx(newToken.Unpause(signer(newKey)))
	assertBalance(recipient.address(), big.NewInt(100+1500-10))
	assertBalance(s.account[3].address(), big.NewInt(10))
}

//////////////// Utility

func maxUint256() *big.Int {
	z := bigInt(1)
	z = z.Lsh(z, 256)
	z = z.Sub(z, bigInt(1))
	return z
}

func maxUint160() *big.Int {
	z := bigInt(1)
	z = z.Lsh(z, 160)
	z = z.Sub(z, bigInt(1))
	return z
}

func minInt160AsUint160() *big.Int {
	z := bigInt(1)
	z = z.Lsh(z, 159)
	return z
}

func bigInt(n uint32) *big.Int {
	return big.NewInt(int64(n))
}

func mintingTransfer(to common.Address, amount *big.Int) ReserveDollarTransfer {
	return ReserveDollarTransfer{From: common.BigToAddress(bigInt(0)), To: to, Value: amount}
}
