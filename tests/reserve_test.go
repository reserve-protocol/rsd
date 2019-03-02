package tests

import (
	"fmt"
	"math/big"
	"os/exec"
	"testing"

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

func (s *ReserveDollarSuite) SetupSuite() {
	s.setup()
	if testing.CoverMode() == "" {
		s.createFastNode()
	} else {
		s.createSlowCoverageNode()
	}
}

func (s *ReserveDollarSuite) TearDownSuite() {
	if testing.CoverMode() != "" {
		s.Assert().NoError(s.node.(*soltools.Backend).WriteCoverage())
		s.Assert().NoError(s.node.(*soltools.Backend).Close())

		if out, err := exec.Command("npx", "istanbul", "report", "html").CombinedOutput(); err != nil {
			fmt.Println()
			fmt.Println("I generated coverage information in coverage/coverage.json.")
			fmt.Println("I tried to process it with `istanbul` to turn it into a readable report, but failed.")
			fmt.Println("The error I got when running istanbul was:", err)
			fmt.Println("Istanbul's output was:\n" + string(out))
		}
	}
}

func (s *ReserveDollarSuite) BeforeTest(suiteName, testName string) {
	_, tx, reserve, err := abi.DeployReserveDollar(s.signer, s.node)
	s.requireTx(tx, err)

	s.reserve = reserve

	// Make the deployment account a minter, pauser, and freezer.
	s.requireTx(s.reserve.ChangeMinter(s.signer, s.account[0].address()))
	s.requireTx(s.reserve.ChangePauser(s.signer, s.account[0].address()))
	s.requireTx(s.reserve.ChangeFreezer(s.signer, s.account[0].address()))
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

	// Mint to recipient.
	s.requireTx(s.reserve.Mint(s.signer, recipient, amount))

	// Check that balances are as expected.
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

	// Owner approves spender.
	s.requireTx(s.reserve.Approve(signer(owner), spender.address(), amount))

	// Approval should be reflected in allowance.
	s.assertAllowance(owner.address(), spender.address(), amount)

	// Shouldn't be symmetric.
	s.assertAllowance(spender.address(), owner.address(), common.Big0)

	// Balances shouldn't change.
	s.assertBalance(owner.address(), common.Big0)
	s.assertBalance(spender.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}

func (s *ReserveDollarSuite) TestIncreaseAllowance() {
	owner := s.account[1]
	spender := s.account[2]
	amount := big.NewInt(2000)

	// Owner approves spender through increaseAllowance.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), amount))

	// Approval should be reflected in allowance.
	s.assertAllowance(owner.address(), spender.address(), amount)

	// Shouldn't be symmetric.
	s.assertAllowance(spender.address(), owner.address(), common.Big0)

	// Balances shouldn't change.
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

	// Owner approves spender for initial amount.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// Owner should not be able to increase approval high enough to overflow a uint256.
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(owner), spender.address(), maxUint256()))
}

func (s *ReserveDollarSuite) TestDecreaseAllowance() {
	owner := s.account[1]
	spender := s.account[2]
	initialAmount := big.NewInt(10)
	decrease := big.NewInt(6)
	final := big.NewInt(4)

	// Owner approves spender for initial amount.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), initialAmount))

	// Owner decreases allowance.
	s.requireTx(s.reserve.DecreaseAllowance(signer(owner), spender.address(), decrease))

	// Allowance should be as we expect.
	s.assertAllowance(owner.address(), spender.address(), final)

	// Balances shouldn't change.
	s.assertBalance(owner.address(), common.Big0)
	s.assertBalance(spender.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}

func (s *ReserveDollarSuite) TestPausing() {
	sender := s.account[1]
	amount := big.NewInt(1000)
	approveAmount := common.Big1
	recipient := s.account[2]
	spender := s.account[3]

	// Give sender funds. Minting is allowed while unpaused.
	s.requireTx(s.reserve.Mint(s.signer, sender.address(), amount))
	s.assertBalance(sender.address(), amount)

	// Approve spender to spend senders funds.
	s.requireTx(s.reserve.Approve(signer(sender), spender.address(), approveAmount))
	s.assertAllowance(sender.address(), spender.address(), approveAmount)

	// Pause.
	s.requireTx(s.reserve.Pause(s.signer))

	// Minting is not allowed while paused.
	s.requireTxFails(s.reserve.Mint(s.signer, recipient.address(), amount))

	// Transfers from are not allowed while paused.
	s.requireTxFails(s.reserve.TransferFrom(s.signer, sender.address(), recipient.address(), amount))
	s.assertBalance(recipient.address(), common.Big0)
	s.assertBalance(sender.address(), amount)

	// Transfers are not allowed while paused.
	s.requireTxFails(s.reserve.Transfer(signer(sender), recipient.address(), amount))
	s.assertBalance(recipient.address(), common.Big0)
	s.assertBalance(sender.address(), amount)

	// Approving is not allowed while paused.
	s.requireTxFails(s.reserve.Approve(signer(sender), spender.address(), amount))
	s.assertAllowance(sender.address(), spender.address(), approveAmount)

	// IncreaseAllowance is not allowed while paused.
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(sender), spender.address(), amount))
	s.assertAllowance(sender.address(), spender.address(), approveAmount)

	// DecreaseAllowance is not allowed while paused.
	s.requireTxFails(s.reserve.DecreaseAllowance(signer(sender), spender.address(), amount))
	s.assertAllowance(sender.address(), spender.address(), approveAmount)

	// Unpause.
	s.requireTx(s.reserve.Unpause(s.signer))

	// Transfers are allowed while unpaused.
	s.requireTx(s.reserve.Transfer(signer(sender), recipient.address(), amount))
	s.assertBalance(recipient.address(), amount)

	// Approving is allowed while paused.
	s.requireTx(s.reserve.Approve(signer(sender), spender.address(), common.Big2))
	s.assertAllowance(sender.address(), spender.address(), common.Big2)

	// DecreaseAllowance is allowed while unpaused.
	s.requireTx(s.reserve.DecreaseAllowance(signer(sender), spender.address(), approveAmount))
	s.assertAllowance(sender.address(), spender.address(), approveAmount)

	// IncreaseAllowance is allowed while unpaused.
	s.requireTx(s.reserve.IncreaseAllowance(signer(sender), spender.address(), approveAmount))
	s.assertAllowance(sender.address(), spender.address(), common.Big2)
}

func (s *ReserveDollarSuite) TestFreezeTransferOut() {
	target := s.account[1]
	recipient := s.account[2]

	// Give target funds.
	amount := common.Big1
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
}

func (s *ReserveDollarSuite) TestFreezeTransferIn() {
	target := s.account[1]
	amount := big.NewInt(200)

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
}

func (s *ReserveDollarSuite) TestFreezeApprovals() {
	target := s.account[1]
	recipient := s.account[2]

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Frozen account shouldn't be able to create approvals.
	s.requireTxFails(s.reserve.Approve(signer(target), recipient.address(), common.Big1))
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(target), recipient.address(), common.Big1))
	s.assertAllowance(target.address(), recipient.address(), common.Big0)

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Unfrozen account should be able to create approvals again.
	s.requireTx(s.reserve.Approve(signer(target), recipient.address(), common.Big1))
	s.requireTx(s.reserve.IncreaseAllowance(signer(target), recipient.address(), common.Big1))
	s.assertAllowance(target.address(), recipient.address(), common.Big2)
}

func (s *ReserveDollarSuite) TestFreezeTransferFrom() {
	target := s.account[1]
	recipient := s.account[2]

	// Approve target to transfer funds.
	amount := common.Big1
	s.requireTx(s.reserve.Mint(s.signer, s.account[0].address(), amount))
	s.requireTx(s.reserve.Approve(s.signer, target.address(), amount))
	s.assertAllowance(s.account[0].address(), target.address(), amount)

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Frozen account shouldn't be able to call transferFrom.
	s.requireTxFails(s.reserve.TransferFrom(signer(target), s.account[0].address(), recipient.address(), amount))
	s.assertBalance(recipient.address(), common.Big0)

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Unfrozen account should now be able to call transferFrom.
	s.requireTx(s.reserve.TransferFrom(signer(target), s.account[0].address(), recipient.address(), amount))
	s.assertBalance(recipient.address(), amount)
}

func (s *ReserveDollarSuite) TestFreezeApprove() {
	target := s.account[1]

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Should not be able to approve frozen target.
	s.requireTxFails(s.reserve.Approve(s.signer, target.address(), common.Big1))

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Should be able to approve unfrozen target.
	s.requireTx(s.reserve.Approve(s.signer, target.address(), common.Big1))
}

func (s *ReserveDollarSuite) TestFreezeIncreaseAllowance() {
	target := s.account[1]
	owner := s.account[2]

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Should not be able to increase allowance frozen target.
	s.requireTxFails(s.reserve.IncreaseAllowance(signer(owner), target.address(), common.Big1))
	s.assertAllowance(owner.address(), target.address(), common.Big0)

	// Unfreeze target.
	s.requireTx(s.reserve.Unfreeze(s.signer, target.address()))

	// Should be able to increase allowance unfrozen target.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), target.address(), common.Big1))
	s.assertAllowance(owner.address(), target.address(), common.Big1)
}

func (s *ReserveDollarSuite) TestFreezeDecreaseAllowance() {
	spender := s.account[1]
	owner := s.account[2]

	// Increase allowance to set up for decrease.
	s.requireTx(s.reserve.IncreaseAllowance(signer(owner), spender.address(), common.Big1))

	// Freeze spender.
	s.requireTx(s.reserve.Freeze(s.signer, spender.address()))

	// Should not be able to decrease allowance frozen spender.
	s.requireTxFails(s.reserve.DecreaseAllowance(signer(owner), spender.address(), common.Big1))
	s.assertAllowance(owner.address(), spender.address(), common.Big1)

	// Unfreeze spender.
	s.requireTx(s.reserve.Unfreeze(s.signer, spender.address()))

	// Should be able to decrease allowance unfrozen spender.
	s.requireTx(s.reserve.DecreaseAllowance(signer(owner), spender.address(), common.Big1))
	s.assertAllowance(owner.address(), spender.address(), common.Big0)
}

func (s *ReserveDollarSuite) TestWiping() {
	target := s.account[1]

	// Give target funds.
	amount := big.NewInt(100)
	s.requireTx(s.reserve.Mint(s.signer, target.address(), amount))

	// Should not be able to wipe target before freezing them.
	s.requireTxFails(s.reserve.Wipe(s.signer, target.address()))

	// Freeze target.
	s.requireTx(s.reserve.Freeze(s.signer, target.address()))

	// Target should still have funds.
	s.assertBalance(target.address(), amount)

	// Wipe target.
	s.requireTx(s.reserve.Wipe(s.signer, target.address()))

	// Target should have zero funds.
	s.assertBalance(target.address(), common.Big0)
}

func (s *ReserveDollarSuite) TestMintingBurningChain() {
	// Mint to recipient.
	recipient := s.account[1]
	amount := big.NewInt(100)

	s.requireTx(s.reserve.Mint(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), amount)
	s.assertTotalSupply(amount)

	// Burn from recipient
	s.requireTx(s.reserve.BurnFrom(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}

func (s *ReserveDollarSuite) TestMintingTransferBurningChain() {
	// Mint to recipient.
	recipient := s.account[1]
	amount := big.NewInt(100)

	s.requireTx(s.reserve.Mint(s.signer, recipient.address(), amount))

	s.assertBalance(recipient.address(), amount)
	s.assertTotalSupply(amount)

	// Transfer to target.
	target := s.account[2]
	s.requireTx(s.reserve.Transfer(signer(recipient), target.address(), amount))

	s.assertBalance(target.address(), amount)
	s.assertBalance(recipient.address(), common.Big0)

	// Burn from target.
	s.requireTx(s.reserve.BurnFrom(s.signer, target.address(), amount))

	s.assertBalance(target.address(), common.Big0)
	s.assertBalance(recipient.address(), common.Big0)
	s.assertTotalSupply(common.Big0)
}
