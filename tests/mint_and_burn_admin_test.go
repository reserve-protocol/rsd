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
	s.Require().NoError(s.node.(backend).AdjustTime(14 * time.Hour))

	// Trying to confirm it should now succeed.
	s.requireTx(s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true))

	// The mint should have happened.
	s.assertBalance(recipient, amount)

	// Trying to confirm a second time should fail.
	_, err = s.adminContract.Confirm(s.adminSigner, common.Big0, recipient, amount, true)
	s.Error(err)
}
