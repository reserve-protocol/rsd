pragma solidity ^0.5.4;

import "ReserveDollar.sol";

// Test that 3 users that can transfer between themselves only cannot increase
// or decrease the total supply
contract EchidnaTotalSupplyTests is ReserveDollar {
    ReserveDollar r;

    constructor() public {
        minter = msg.sender;

        testFixtureMint(address(0x1), 123);
        testFixtureMint(address(0x2), 1234);
        testFixtureMint(address(0x3), 12345);
    }
    function testFixtureMint(address account, uint256 value) internal notPaused notFrozen(account) only(minter) {
        require(account != address(0), "can't mint to address zero");

        totalSupply = totalSupply.add(value);
        data.addBalance(account, value);
        emit Transfer(address(0), account, value);
    }

    // Disable the `burnFrom` function to prevent the total supply from
    // changing during the transfer testing
    function burnFrom(address, uint256) external { return; }
    // Disable the `transferFrom` function to increase the odds for echidna to
    // transfer between our three accounts that correspond to the 3 possible
    // senders of messages during the test
    function transferFrom(address, address, uint256) external returns (bool) { return false; }

    // Increase the chances of echidna creating transfer loops between
    // different users by reducing the search space. By setting the destination
    // to an address that will be active during the tests, we increase the
    // chance of seeing multi account collusion test cases. Otherwise almost
    // all generated test cases would involve transfers to an address that will
    // not be active during the test, thus acting as a sink for tokens.
    function transferToOne(uint256 value) public returns (bool) {
        return this.transfer(address(0x1), value);
    }

    function transferToTwo(uint256 value) public returns (bool) {
        return this.transfer(address(0x2), value);
    }

    function transferToThree(uint256 value) public returns (bool) {
        return this.transfer(address(0x3), value);
    }

    function echidna_constant_supply() public view returns (bool) {
        return this.totalSupply() == 13702;
    }

    function echidna_zero_sum_transfers() public view returns (bool) {
        return (this.balanceOf(address(0x1)) + this.balanceOf(address(0x2)) + this.balanceOf(address(0x3))) == 13702;
    }

    // Check that the owner has not been changed by transfers just in case.
    function echidna_owner_unchanged_by_transfers() public view returns (bool) {
        return this.owner() == 0x00a329c0648769A73afAc7F9381E08FB43dBEA72;
    }
}
