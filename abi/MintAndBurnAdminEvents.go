// This file is auto-generated. Do not edit.

package abi

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (c *MintAndBurnAdminFilterer) ParseLog(log *types.Log) (fmt.Stringer, error) {
	var event fmt.Stringer
	var eventName string
	switch log.Topics[0].Hex() {
	case "0x3732302b0efc3e1e883bb80d83c641031dc1e32223cb406c3e4d5de68208c4e9": // AllProposalsCancelled
		event = new(MintAndBurnAdminAllProposalsCancelled)
		eventName = "AllProposalsCancelled"
	case "0xc1ea9ad7fe3cfb48a741fc229353411aabb3b135d9446697bf6db7c197a9ac0f": // ProposalCancelled
		event = new(MintAndBurnAdminProposalCancelled)
		eventName = "ProposalCancelled"
	case "0xc398e86b1dfd2596a48f97df67ac573ef31251ea5b65d73e4096be478df18f57": // ProposalConfirmed
		event = new(MintAndBurnAdminProposalConfirmed)
		eventName = "ProposalConfirmed"
	case "0xd1d2eb762bbbecbc03b8a9dd22368874018771d0c93d855cd08c5a8fa6086b96": // ProposalCreated
		event = new(MintAndBurnAdminProposalCreated)
		eventName = "ProposalCreated"
	default:
		return nil, fmt.Errorf("no such event hash for MintAndBurnAdmin: %v", log.Topics[0])
	}

	err := c.contract.UnpackLog(event, eventName, *log)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e MintAndBurnAdminAllProposalsCancelled) String() string {
	return fmt.Sprintf("MintAndBurnAdmin.AllProposalsCancelled()")
}

func (e MintAndBurnAdminProposalCancelled) String() string {
	return fmt.Sprintf("MintAndBurnAdmin.ProposalCancelled(%v, %v, %v, %v)", e.Index, e.Addr.Hex(), e.Value, e.IsMint)
}

func (e MintAndBurnAdminProposalConfirmed) String() string {
	return fmt.Sprintf("MintAndBurnAdmin.ProposalConfirmed(%v, %v, %v, %v)", e.Index, e.Addr.Hex(), e.Value, e.IsMint)
}

func (e MintAndBurnAdminProposalCreated) String() string {
	return fmt.Sprintf("MintAndBurnAdmin.ProposalCreated(%v, %v, %v, %v, %v)", e.Index, e.Addr.Hex(), e.Value, e.IsMint, e.DelayUntil)
}
