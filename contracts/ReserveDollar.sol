pragma solidity ^0.5.4;

import "./zeppelin/SafeMath.sol";

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
     * MATH
     */
    using SafeMath for uint256;

    /**
     * DATA
     */

    // ETERNAL STORAGE DATA
    ReserveDollarEternalStorage private data;

    // BASIC DATA
    string public name = "Reserve Dollar";
    string public symbol = "RSVD";
    uint8 public constant decimals = 18;
    uint256 public totalSupply = 0;

    address private owner;

    // ROLES
    address minter;

    constructor() public {
        data = new ReserveDollarEternalStorage();
        owner = msg.sender;
        minter = msg.sender;
    }

    /**
     * EVENTS
     */

    event NameChanged(string newName, string newSymbol);
    event TokensMinted(address indexed to, uint256 quantity);
    event Transfer(address indexed from, address indexed to, uint256 quantity);

    /**
     * FUNCTIONALITY
     */

    // AUTHORIZATION FUNCTIONALITY

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "onlyOwner");
        _;
    }

    /**
     * @dev Throws if called by any account other than role.
     */
    modifier onlyRole(address role) {
        require(msg.sender == role, "onlyRole");
        _;
    }

    // Name change functionality
    function changeName(string memory newName, string memory newSymbol) public onlyOwner {
        name = newName;
        symbol = newSymbol;
        emit NameChanged(newName, newSymbol);
    }

    // BASIC FUNCTIONALITY

    function balanceOf(address who) public view returns (uint256) {
        return data.balance(who);
    }

    // MINTING FUNCTIONALITY

    function mint(address to, uint256 value) public onlyRole(minter) {
        totalSupply = totalSupply.add(value);
        data.setBalance(to, data.balance(to).add(value));
        emit TokensMinted(to, value);
        emit Transfer(address(0), to, value);
    }
}
