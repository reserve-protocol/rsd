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

    /**
     * DATA
     */

    // ETERNAL STORAGE DATA
    ReserveDollarEternalStorage private data;

    // BASIC DATA
    string public name = "Reserve Dollar";
    string public symbol = "RSVD";
    uint8 public constant decimals = 18;

    address private owner;

    constructor() public {
        owner = msg.sender;
        data = new ReserveDollarEternalStorage();
    }

    /**
     * EVENTS
     */

    event NameChanged(string newName, string newSymbol);

    /**
     * FUNCTIONALITY
     */

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "onlyOwner");
        _;
    }

    // Name change functionality
    function changeName(string calldata newName, string calldata newSymbol) external onlyOwner {
        name = newName;
        symbol = newSymbol;
        emit NameChanged(newName, newSymbol);
    }

    // BASIC FUNCTIONALITY

    function balanceOf(address who) public view returns (uint256) {
        return data.balance(who);
    }
}
