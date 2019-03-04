pragma solidity ^0.5.4;

import "ReserveDollar.sol";

// Use the multisender.yaml to send messages from 0x1, 0x2, 0x3, but not the
// deployer address. Non owners should be unable to change the owner
contract EchidnaConstantOwnerSimple is ReserveDollar {
    // Check that the owner has not been changed
    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}
