// This file is auto-generated. Do not edit.

package abi

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (c *ReserveDollarV2Filterer) ParseLog(log *types.Log) (fmt.Stringer, error) {
	var event fmt.Stringer
	var eventName string
	switch log.Topics[0].Hex() {
	case "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925": // Approval
		event = new(ReserveDollarV2Approval)
		eventName = "Approval"
	case "0x97feb20f655745f67bbd05465394b86626eaafbbaae4a509a838c19237ea9da3": // FreezerChanged
		event = new(ReserveDollarV2FreezerChanged)
		eventName = "FreezerChanged"
	case "0xf0906ec3b3af5007c736f1174c73ff022e930e45637fbdbc797f05ea613474de": // Frozen
		event = new(ReserveDollarV2Frozen)
		eventName = "Frozen"
	case "0xb6b8f1859c5c352e5ffad07d0f77e384ac725512c015bd3a3ffc885831c8a425": // MinterChanged
		event = new(ReserveDollarV2MinterChanged)
		eventName = "MinterChanged"
	case "0x6c20b91d1723b78732eba64ff11ebd7966a6e4af568a00fa4f6b72c20f58b02a": // NameChanged
		event = new(ReserveDollarV2NameChanged)
		eventName = "NameChanged"
	case "0xa2ea9883a321a3e97b8266c2b078bfeec6d50c711ed71f874a90d500ae2eaf36": // OwnerChanged
		event = new(ReserveDollarV2OwnerChanged)
		eventName = "OwnerChanged"
	case "0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258": // Paused
		event = new(ReserveDollarV2Paused)
		eventName = "Paused"
	case "0xb80482a293ca2e013eda8683c9bd7fc8347cfdaeea5ede58cba46df502c2a604": // PauserChanged
		event = new(ReserveDollarV2PauserChanged)
		eventName = "PauserChanged"
	case "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": // Transfer
		event = new(ReserveDollarV2Transfer)
		eventName = "Transfer"
	case "0x295f47d48ca3de5c5214af57c89859243090803a47bbca8a4bbe6231a77067b4": // Unfrozen
		event = new(ReserveDollarV2Unfrozen)
		eventName = "Unfrozen"
	case "0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa": // Unpaused
		event = new(ReserveDollarV2Unpaused)
		eventName = "Unpaused"
	case "0xaa7ccaa635252f24fc5a4665e7c4a8af5aa588c2e60d326c1196a0e4d2d59f2c": // Wiped
		event = new(ReserveDollarV2Wiped)
		eventName = "Wiped"
	default:
		return nil, fmt.Errorf("no such event hash for ReserveDollarV2: %v", log.Topics[0])
	}

	err := c.contract.UnpackLog(event, eventName, *log)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e ReserveDollarV2Approval) String() string {
	return fmt.Sprintf("ReserveDollarV2.Approval(%v, %v, %v)", e.Holder.Hex(), e.Spender.Hex(), e.Value)
}

func (e ReserveDollarV2FreezerChanged) String() string {
	return fmt.Sprintf("ReserveDollarV2.FreezerChanged(%v)", e.NewFreezer.Hex())
}

func (e ReserveDollarV2Frozen) String() string {
	return fmt.Sprintf("ReserveDollarV2.Frozen(%v, %v)", e.Freezer.Hex(), e.Account.Hex())
}

func (e ReserveDollarV2MinterChanged) String() string {
	return fmt.Sprintf("ReserveDollarV2.MinterChanged(%v)", e.NewMinter.Hex())
}

func (e ReserveDollarV2NameChanged) String() string {
	return fmt.Sprintf("ReserveDollarV2.NameChanged(%q, %q)", e.NewName, e.NewSymbol)
}

func (e ReserveDollarV2OwnerChanged) String() string {
	return fmt.Sprintf("ReserveDollarV2.OwnerChanged(%v)", e.NewOwner.Hex())
}

func (e ReserveDollarV2Paused) String() string {
	return fmt.Sprintf("ReserveDollarV2.Paused(%v)", e.Account.Hex())
}

func (e ReserveDollarV2PauserChanged) String() string {
	return fmt.Sprintf("ReserveDollarV2.PauserChanged(%v)", e.NewPauser.Hex())
}

func (e ReserveDollarV2Transfer) String() string {
	return fmt.Sprintf("ReserveDollarV2.Transfer(%v, %v, %v)", e.From.Hex(), e.To.Hex(), e.Value)
}

func (e ReserveDollarV2Unfrozen) String() string {
	return fmt.Sprintf("ReserveDollarV2.Unfrozen(%v, %v)", e.Freezer.Hex(), e.Account.Hex())
}

func (e ReserveDollarV2Unpaused) String() string {
	return fmt.Sprintf("ReserveDollarV2.Unpaused(%v)", e.Account.Hex())
}

func (e ReserveDollarV2Wiped) String() string {
	return fmt.Sprintf("ReserveDollarV2.Wiped(%v, %v)", e.Freezer.Hex(), e.Wiped.Hex())
}
