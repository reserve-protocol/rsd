pragma solidity ^0.5.4;

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
    function transferOwnership(address newOwner) public onlyOwner {
        emit OwnershipTransferred(owner, newOwner);
        owner = newOwner;
    }

    constructor() public {
        owner = msg.sender;
    }

    /**
     * DATA
     */

    // _balance

    mapping (address => uint256) private _balance;

    function balance(address key) public view onlyOwner returns (uint256) {
        return _balance[key];
    }

    function setBalance(address key, uint256 value) public onlyOwner {
        _balance[key] = value;
    }

    // _allowed

    mapping (address => uint256) private _allowed;

    function allowed(address key) public view onlyOwner returns (uint256) {
        return _allowed[key];
    }

    function setAllowed(address key, uint256 value) public onlyOwner {
        _allowed[key] = value;
    }
}
