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

    /**
     * @dev Gets the balance of the specified address.
     * @param addr The address to query the balance of.
     * @return A uint256 representing the amount owned by the passed address.
     */
    function balanceOf(address addr) public view returns (uint256) {
        return data.balance(addr);
    }

    /**
     * @dev Transfer token to a specified address
     * @param to The address to transfer to.
     * @param value The amount to be transferred.
     */
    function transfer(address to, uint256 value) public returns (bool) {
        require(to != address(0), "cannot transfer to address zero");

        data.subBalance(msg.sender, value);
        data.addBalance(to, value);
        emit Transfer(msg.sender, to, value);
        return true;
    }

    // MINTING FUNCTIONALITY

    function mint(address to, uint256 value) public onlyRole(minter) {
        totalSupply = totalSupply.add(value);
        data.addBalance(to, value);
        emit TokensMinted(to, value);
        emit Transfer(address(0), to, value);
    }
}
