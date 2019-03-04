pragma solidity ^0.5.4;

import "EchidnaTestBase.sol";

// Test that 3 users that can transfer between themselves only cannot increase
// or decrease the total supply
contract EchidnaTotalSupplyTests is WrappedReserveDollar {
    constructor() public {
        require(r.owner() == address(this), "document expected contract ownership during setup");
        r.changeMinter(address(this));
        // add some balance for users at 0x1, 0x2, and 0x3 addresses
        r.mint(address(0x1), 123);
        r.mint(address(0x2), 1234);
        r.mint(address(0x3), 12345);

        // don't allow this contract to mint because it will be passed on as msg.sender
        r.testFixtureChangeOwner(address(0x900));
        r.changeMinter(address(0x900));
    }

    // Disable the `burnFrom` function to prevent the total supply from
    // changing from extracted ways during the transfer testing
    function burnFrom(address, uint256) external { return; }
    // Disable the `transfer` and `transferFrom` function to force echidna to
    // only transfer between our three accounts that correspond to the 3
    // possible senders of messages during the test
    function transfer(address, uint256) external returns (bool) { return false; }
    function transferFrom(address, address, uint256) external returns (bool) { return false; }

    // Function stubs to let echidna only tranfer into account 1
    function transferToOne(uint256 value) public returns (bool) {
        return r.transfer(address(0x1), value);
    }

    function transferToTwo(uint256 value) public returns (bool) {
        return r.transfer(address(0x2), value);
    }

    function transferToThree(uint256 value) public returns (bool) {
        return r.transfer(address(0x3), value);
    }

    function echidna_constant_supply() public view returns (bool) {
        return r.totalSupply() == 13702;
    }

    function echidna_zero_sum_transfers() public view returns (bool) {
        return (r.balanceOf(address(0x1)) + r.balanceOf(address(0x2)) + r.balanceOf(address(0x3))) == 13702;
    }

    // Check that the owner has not been changed by transfers just in case.
    function echidna_owner_unchanged_by_transfers() public view returns (bool) {
        return r.owner() == address(0x900);
    }
}
