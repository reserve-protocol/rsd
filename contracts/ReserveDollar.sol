pragma solidity ^0.5.4;

import "./zeppelin/SafeMath.sol";

import "./ReserveDollarEternalStorage.sol";

/**
 * @title ERC20 interface
 * @dev see https://eips.ethereum.org/EIPS/eip-20
 *
 * This interface serves as a compile-time check that ReserveDollar
 * implements the ERC-20 interface.
 */
interface IERC20 {
    function transfer(address to, uint256 value) external returns (bool);
    function approve(address spender, uint256 value) external returns (bool);
    function transferFrom(address from, address to, uint256 value) external returns (bool);
    function totalSupply() external view returns (uint256);
    function balanceOf(address who) external view returns (uint256);
    function allowance(address addr, address spender) external view returns (uint256);
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed _owner, address indexed spender, uint256 value);
}

/**
 * @title ReserveDollar
 *
 * @dev An ERC-20 token with minting, burning, pausing, and blacklisting.
 *
 * Based on OpenZeppelin's implementation:
 *
 * https://github.com/OpenZeppelin/openzeppelin-solidity/blob/41aa39afbc13f0585634061701c883fe512a5469/contracts/token/ERC20/ERC20.sol
 *
 * Some data is held in an Eternal Storage contract to facilitate potential future upgrades.
 */
contract ReserveDollar is IERC20 {
    using SafeMath for uint256;

    ReserveDollarEternalStorage private data;

    uint256 private _totalSupply;

    string public name = "Reserve Dollar";
    string public symbol = "RSVD";
    uint8 public constant decimals = 18;

    bool public paused;

    mapping (address => bool) frozen;

    address private owner;
    address minter;
    address pauser;
    address freezer;

    constructor() public {
        data = new ReserveDollarEternalStorage();
        owner = msg.sender;
        minter = msg.sender;
        pauser = msg.sender;
        freezer = msg.sender;
    }

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

    event NameChanged(string newName, string newSymbol);

    function changeName(string memory newName, string memory newSymbol) public onlyOwner {
        name = newName;
        symbol = newSymbol;
        emit NameChanged(newName, newSymbol);
    }

    event Paused(address account);
    event Unpaused(address account);

    function pause() public onlyRole(pauser) {
        paused = true;
        emit Paused(pauser);
    }

    function unpause() public onlyRole(pauser) {
        paused = false;
        emit Unpaused(pauser);
    }

    modifier whenNotPaused() {
        require(!paused, "paused");
        _;
    }

    event Frozen(address indexed freezer, address indexed account);
    event Unfrozen(address indexed freezer, address indexed account);
    event Wiped(address indexed freezer, address indexed wiped);

    function freeze(address who) public onlyRole(freezer) {
        require(!frozen[who], "account already frozen");
        frozen[who] = true;
        emit Frozen(freezer, who);
    }

    function unfreeze(address who) public onlyRole(freezer) {
        require(frozen[who], "account not frozen");
        frozen[who] = false;
        emit Unfrozen(freezer, who);
    }

    modifier notFrozen(address account) {
        require(!frozen[account], "account frozen");
        _;
    }

    function wipe(address who) public onlyRole(freezer) {
        require(frozen[who], "cannot wipe unfrozen account");
        _burn(who, data.balance(who));
        emit Wiped(freezer, who);
    }

    /**
     * @dev Total number of tokens in existence
     */
    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    /**
     * @dev Gets the balance of the specified address.
     * @param _owner The address to query the balance of.
     * @return An uint256 representing the amount owned by the passed address.
     */
    function balanceOf(address _owner) public view returns (uint256) {
        return data.balance(_owner);
    }

    /**
     * @dev Function to check the amount of tokens that an _owner allowed to a spender.
     * @param _owner address The address which owns the funds.
     * @param spender address The address which will spend the funds.
     * @return A uint256 specifying the amount of tokens still available for the spender.
     */
    function allowance(address _owner, address spender) public view returns (uint256) {
        return data.allowed(_owner, spender);
    }

    /**
     * @dev Transfer token to a specified address
     * @param to The address to transfer to.
     * @param value The amount to be transferred.
     */
    function transfer(address to, uint256 value)
        public
        whenNotPaused
        notFrozen(msg.sender)
        notFrozen(to)
        returns (bool)
    {
        _transfer(msg.sender, to, value);
        return true;
    }

    /**
     * @dev Approve the passed address to spend the specified amount of tokens on behalf of msg.sender.
     * Beware that changing an allowance with this method brings the risk that someone may use both the old
     * and the new allowance by unfortunate transaction ordering. One possible solution to mitigate this
     * race condition is to first reduce the spender's allowance to 0 and set the desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     * Another workaround is to use increaseAllowance/decreaseAllowance.
     * @param spender The address which will spend the funds.
     * @param value The amount of tokens to be spent.
     */
    function approve(address spender, uint256 value)
        public
        whenNotPaused
        notFrozen(msg.sender)
        notFrozen(spender)
        returns (bool)
    {
        _approve(msg.sender, spender, value);
        return true;
    }

    /**
     * @dev Transfer tokens from one address to another.
     * Note that while this function emits an Approval event, this is not required as per the specification,
     * and other compliant implementations may not emit the event.
     * @param from address The address which you want to send tokens from
     * @param to address The address which you want to transfer to
     * @param value uint256 the amount of tokens to be transferred
     */
    function transferFrom(address from, address to, uint256 value)
        public
        whenNotPaused
        notFrozen(msg.sender)
        notFrozen(from)
        notFrozen(to)
        returns (bool)
    {
        _transfer(from, to, value);
        _approve(from, msg.sender, data.allowed(from, msg.sender).sub(value));
        return true;
    }

    /**
     * @dev Increase the amount of tokens that an _owner allowed to a spender.
     * approve should be called when data.allowed(msg.sender, spender) == 0. To increment
     * allowed value is better to use this function to avoid 2 calls (and wait until
     * the first transaction is mined)
     * From MonolithDAO Token.sol
     * Emits an Approval event.
     * @param spender The address which will spend the funds.
     * @param addedValue The amount of tokens to increase the allowance by.
     */
    function increaseAllowance(address spender, uint256 addedValue)
        public
        whenNotPaused
        notFrozen(msg.sender)
        notFrozen(spender)
        returns (bool)
    {
        _approve(msg.sender, spender, data.allowed(msg.sender, spender).add(addedValue));
        return true;
    }

    /**
     * @dev Decrease the amount of tokens that an _owner allowed to a spender.
     * approve should be called when data.allowed(msg.sender, spender) == 0. To decrement
     * allowed value is better to use this function to avoid 2 calls (and wait until
     * the first transaction is mined)
     * From MonolithDAO Token.sol
     * Emits an Approval event.
     * @param spender The address which will spend the funds.
     * @param subtractedValue The amount of tokens to decrease the allowance by.
     */
    function decreaseAllowance(address spender, uint256 subtractedValue)
        public
        whenNotPaused
        notFrozen(msg.sender)
        notFrozen(spender)
        returns (bool)
    {
        _approve(msg.sender, spender, data.allowed(msg.sender, spender).sub(subtractedValue));
        return true;
    }

    /**
     * @dev Transfer token for a specified addresses
     * @param from The address to transfer from.
     * @param to The address to transfer to.
     * @param value The amount to be transferred.
     */
    function _transfer(address from, address to, uint256 value) internal {
        require(to != address(0), "can't transfer to address zero");

        data.subBalance(from, value);
        data.addBalance(to, value);
        emit Transfer(from, to, value);
    }

    /**
     * @dev Function that mints an amount of the token and assigns it to
     * an account. This encapsulates the modification of balances such that the
     * proper events are emitted.
     * @param account The account that will receive the created tokens.
     * @param value The amount that will be created.
     */
    function mint(address account, uint256 value) public whenNotPaused onlyRole(minter) {
        require(account != address(0), "can't mint to address zero");

        _totalSupply = _totalSupply.add(value);
        data.addBalance(account, value);
        emit Transfer(address(0), account, value);
    }

    /**
     * @dev Internal function that burns an amount of the token of a given
     * account.
     * @param account The account whose tokens will be burnt.
     * @param value The amount that will be burnt.
     */
    function _burn(address account, uint256 value) internal {
        require(account != address(0), "can't burn from address zero");

        _totalSupply = _totalSupply.sub(value);
        data.subBalance(account, value);
        emit Transfer(account, address(0), value);
    }

    /**
     * @dev Approve an address to spend another addresses' tokens.
     * @param _owner The address that owns the tokens.
     * @param spender The address that will spend the tokens.
     * @param value The number of tokens that can be spent.
     */
    function _approve(address _owner, address spender, uint256 value) internal {
        require(spender != address(0), "spender cannot be address zero");
        require(_owner != address(0), "_owner cannot be address zero");

        data.setAllowed(_owner, spender, value);
        emit Approval(_owner, spender, value);
    }
}
