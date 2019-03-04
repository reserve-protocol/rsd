pragma solidity ^0.5.4;

import "ReserveDollar.sol";

// Check that there is no combination of messages an owner can send to change the owner other than
// `renounceOwnership`
contract EchidnaConstantOwnerSimple is ReserveDollar {
    // prevent ownership from being renounced during this test
    function renounceOwnership() external {
        return;
    }

    // Check that the owner has not been changed
    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}

// Try to find overflows by allowing all possible transfer amounts. Attacking
// users can transfer any amounts they want. This lets us get random non zero
// values into any contract balances.
contract EchidnaConstantOwnerFreeMoney is ReserveDollar {
    constructor() public {
        minter = address(this);
    }
    
    function testFixtureIncreaseAllowance(address from, address spender, uint256 addedValue) public returns (bool) {
        _approve(from, spender, addedValue);
        return true;
    }

    function testFixtureMint(address account, uint256 value) internal notPaused notFrozen(account) only(minter) {
        require(account != address(0), "can't mint to address zero");

        totalSupply = totalSupply.add(value);
        data.addBalance(account, value);
        emit Transfer(address(0), account, value);
    }


    function transfer(address to, uint256 value) public returns (bool) {
        testFixtureMint(msg.sender, value);
        return transfer(to, value);
    }

    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        testFixtureMint(from, value);
        return transferFrom(from, to, value);
    }

    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}
