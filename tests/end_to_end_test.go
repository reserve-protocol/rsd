package tests

import (
	_ "encoding/hex"
	_ "math"
	_ "math/big"
	"testing"

	_ "github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	_ "github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/core"
	_ "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	_ "github.com/reserve-protocol/reserve-dollar/abi"
)

func TestEndToEnd(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

type EndToEndSuite struct {
	TestSuite
}

var (
	// Compile-time check that EndToEndSuite implements the interfaces we think it does.
	// If it does not implement these interfaces, then the corresponding setup and teardown
	// functions will not actually run.
	_ suite.BeforeTest       = &EndToEndSuite{}
	_ suite.SetupAllSuite    = &EndToEndSuite{}
	_ suite.TearDownAllSuite = &EndToEndSuite{}
)

func (s *EndToEndSuite) SetupSuite() {

}

func (s *EndToEndSuite) TearDownSuite() {

}

func (s *EndToEndSuite) BeforeTest(suiteName, testName string) {

}
