pragma solidity ^0.5.4;

import "ReserveDollar.sol";

contract TestFixtureReserveDollar is ReserveDollar {
    function testFixtureIncreaseAllowance(address from, address spender, uint256 addedValue) public returns (bool) {
        _approve(from, spender, addedValue);
        return true;
    }

    function testFixtureFreeze(address account) internal only(freezer) {
        require(data.frozenTime(account) == 0, "account already frozen");

        // In `wipe` we use block.timestamp (aka `now`) to check that enough time has passed since
        // this freeze happened. That required time delay -- 4 weeks -- is a long time relative to
        // the maximum drift of block.timestamp, so it is fine to trust the miner here.
        // solium-disable-next-line security/no-block-members
        data.setFrozenTime(account, now);

        emit Frozen(freezer, account);
    }
    

    function testFixtureMint(address account, uint256 value) internal notPaused notFrozen(account) only(minter) {
        require(account != address(0), "can't mint to address zero");

        totalSupply = totalSupply.add(value);
        data.addBalance(account, value);
        emit Transfer(address(0), account, value);
    }

    function testFixturePause() internal only(pauser) {
        paused = true;
        emit Paused(pauser);
    }
}
