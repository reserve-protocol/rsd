pragma solidity ^0.5.4;

import "./ReserveDollarEternalStorage.sol";

/**
 * @title ReserveDollar
 *
 * @dev An ERC-20 token with minting, burning, pausing, and blacklisting.
 *
 * Some data is held in an Eternal Storage contract to facilitate potential future upgrades.
 */
contract ReserveDollar {
    ReserveDollarEternalStorage private data;

    constructor() public {
        data = new ReserveDollarEternalStorage();
    }

    function balanceOf(address who) public view returns (uint256) {
        return data.balance(who);
    }
}
