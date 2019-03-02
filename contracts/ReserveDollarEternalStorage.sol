pragma solidity ^0.5.4;

import "./zeppelin/SafeMath.sol";

/**
 * @title Eternal Storage for the Reserve Dollar
 *
 * @dev Eternal Storage facilitates future upgrades.
 *
 * If Reserve chooses to release an upgraded contract for the
 * Reserve Dollar in the future, Reserve will have the option
 * of reusing the deployed version of this data contract to
 * simplify migration.
 *
 * The use of this contract does not imply that Reserve will choose
 * to do a future upgrade, nor that any future upgrades will
 * necessarily re-use this storage. It merely provides option value.
 */
contract ReserveDollarEternalStorage {

    /**
     * MATH
     */
    using SafeMath for uint256;

    /**
     * OWNERSHIP
     */

    address public owner; // TODO: is it cheaper if this is private with an explicit public accessor?

    event OwnershipTransferred(address oldOwner, address newOwner);

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "onlyOwner");
        _;
    }

    /**
     * @dev Allows the current owner to transfer control of the contract to a newOwner.
     * @param newOwner The address to transfer ownership to.
     */
    function transferOwnership(address newOwner) external onlyOwner {
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    constructor() public {
        owner = msg.sender;
    }

    /**
     * DATA
     */

    // balance

    mapping (address => uint256) public balance;

    function addBalance(address key, uint256 value) external onlyOwner {
        balance[key] = balance[key].add(value);
    }

    function subBalance(address key, uint256 value) external onlyOwner {
        balance[key] = balance[key].sub(value);
    }

    function setBalance(address key, uint256 value) external onlyOwner {
        balance[key] = value;
    }

    // allowed

    mapping (address => mapping (address => uint256)) public allowed;

    function setAllowed(address from, address to, uint256 value) external onlyOwner {
        allowed[from][to] = value;
    }

    // frozenTime

    mapping (address => uint256) public frozenTime;

    function setFrozenTime(address who, uint256 time) external onlyOwner {
        frozenTime[who] = time;
    }
}
