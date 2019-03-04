pragma solidity ^0.5.4;

import "ReserveDollar.sol";
import "EchidnaTestBase.sol";

// Use the multisender.yaml to send messages from 0x1, 0x2, 0x3, but not the
// deployer address. Non owners should be unable to change the owner
contract EchidnaFrozenAccount is WrappedReserveDollar {
    constructor() public {
        require(r.owner() == address(this), "documents the expected owner");
        r.changeMinter(address(this));
        r.changeFreezer(address(this));

        // add some balance for users at 0x1, 0x2, and 0x3 addresses
        r.mint(address(0x1), 123);
        r.testFixtureIncreaseAllowance(address(0x1), address(0x3), 19);
        r.mint(address(0x2), 1234);
        r.mint(address(0x3), 12345);
        // freeze the 0x1 address;
        r.freeze(address(0x1));

        // don't allow this contract to mint because this contract address will
        // be passed on as msg.sender
        r.changeMinter(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
        r.changeFreezer(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
        r.testFixtureChangeOwner(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
    }

    function echidna_frozen_account_constant_balance() public view returns (bool) {
        return r.balanceOf(address(0x1)) == 123;

    }

    function echidna_frozen_account_unchanged_allowance() public view returns (bool) {
        return (r.allowance(address(0x1), address(0x2)) == 0) &&
            (r.allowance(address(0x1), address(0x3)) == 19);
    }

    function echidna_freezer_role_unchanged() public view returns (bool) {
        return this.freezer() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}

