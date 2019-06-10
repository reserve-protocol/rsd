// This file is auto-generated. Do not edit.

package abi

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

func (c *ReserveDollarEternalStorageFilterer) ParseLog(log *types.Log) (fmt.Stringer, error) {
	var event fmt.Stringer
	var eventName string
	switch log.Topics[0].Hex() {
	case "0x089af7288b55770a7c1dfd40b9d9e464c64031c45326c0916854814b6c16da28": // EscapeHatchTransferred
		event = new(ReserveDollarEternalStorageEscapeHatchTransferred)
		eventName = "EscapeHatchTransferred"
	case "0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0": // OwnershipTransferred
		event = new(ReserveDollarEternalStorageOwnershipTransferred)
		eventName = "OwnershipTransferred"
	default:
		return nil, fmt.Errorf("no such event hash for ReserveDollarEternalStorage: %v", log.Topics[0])
	}

	err := c.contract.UnpackLog(event, eventName, *log)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e ReserveDollarEternalStorageEscapeHatchTransferred) String() string {
	return fmt.Sprintf("ReserveDollarEternalStorage.EscapeHatchTransferred(%v, %v)", e.OldEscapeHatch.Hex(), e.NewEscapeHatch.Hex())
}

func (e ReserveDollarEternalStorageOwnershipTransferred) String() string {
	return fmt.Sprintf("ReserveDollarEternalStorage.OwnershipTransferred(%v, %v)", e.OldOwner.Hex(), e.NewOwner.Hex())
}
