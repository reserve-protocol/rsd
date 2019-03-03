pragma solidity ^0.5.4;

import "./zeppelin/SafeMath.sol";

import "./ReserveDollarEternalStorage.sol";

/**
 * @title The standard ERC20 interface
 * @dev see https://eips.ethereum.org/EIPS/eip-20
 */
interface IERC20 {
    function transfer(address, uint256) external returns (bool);
    function approve(address, uint256) external returns (bool);
    function transferFrom(address, address, uint256) external returns (bool);
    function totalSupply() external view returns (uint256);
    function balanceOf(address) external view returns (uint256);
    function allowance(address, address) external view returns (uint256);
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed from, address indexed to, uint256 value);
}

/**
 * @title The Reserve Dollar
 * @dev An ERC-20 token with minting, burning, pausing, and user freezing.
 * Based on OpenZeppelin's [implementation](https://github.com/OpenZeppelin/openzeppelin-solidity/blob/41aa39afbc13f0585634061701c883fe512a5469/contracts/token/ERC20/ERC20.sol).
 *
 * Non-constant-sized data is held in ReserveDollarEternalStorage, to facilitate potential future upgrades.
 */
contract ReserveDollar is IERC20 {
    using SafeMath for uint256;

    ReserveDollarEternalStorage internal data;

    function getEternalStorageAddress() public view returns(address) {
        return address(data);
    }

    uint256 private _totalSupply;

    string public name = "Reserve Dollar";
    string public symbol = "RSVD";
    uint8 public constant decimals = 18;

    bool public paused;


    // ==== auth ====

    address public owner;
    address public minter;
    address public pauser;
    address public freezer;
    address public nominatedOwner;

    event OwnerChanged(address indexed newOwner);
    event MinterChanged(address indexed newMinter);
    event PauserChanged(address indexed newPauser);
    event FreezerChanged(address indexed newFreezer);

    /// Initialize with absolutely critical fields already set.
    constructor() public {
        data = new ReserveDollarEternalStorage(msg.sender);
        owner = msg.sender;
        pauser = msg.sender;
    }

    /// Modifies a function to only run if sent by `role`.
    modifier only(address role) {
        require(msg.sender == role, "unauthorized");
        _;
    }

    /// Modifies a function to only run if sent by `role` or the contract's `owner`.
    modifier onlyOwnerOr(address role) {
        require(msg.sender == owner || msg.sender == role, "unauthorized and not owner");
        _;
    }

    /// Change who holds the `minter` role
    function changeMinter(address newMinter) external onlyOwnerOr(minter) {
        minter = newMinter;
        emit MinterChanged(newMinter);
    }

    /// Change who holds the `pauser` role
    function changePauser(address newPauser) external onlyOwnerOr(pauser) {
        pauser = newPauser;
        emit PauserChanged(newPauser);
    }

    /// Change the holds the `freezer` role
    function changeFreezer(address newFreezer) external onlyOwnerOr(freezer) {
        freezer = newFreezer;
        emit FreezerChanged(newFreezer);
    }

    /// Nominate a new `owner`.  We want to ensure that `owner` is always valid, so we don't
    /// actually change `owner` to `nominatedOwner` until `nominatedOwner` calls `acceptOwnership`.
    function nominateNewOwner(address nominee) external only(owner) {
        nominatedOwner = nominee;
    }

    /// Accept nomination for ownership.
    /// This completes the `nominateNewOwner` handshake.
    function acceptOwnership() external onlyOwnerOr(nominatedOwner) {
        if (msg.sender != owner) {
            emit OwnerChanged(msg.sender);
        }
        owner = msg.sender;
        nominatedOwner = msg.sender;
    }

    /// Set `owner` to 0.
    /// Only do this if you're abandoning this contract, e.g., for an upgrade.
    function renounceOwnership() external only(owner) {
        owner = address(0);
        emit OwnerChanged(owner);
    }

    /// Make a different address own the EternalStorage contract.
    /// Only do this if you're abandoning this contract, e.g., for an upgrade.
    function transferEternalStorage(address newOwner) external only(owner) {
        data.transferOwnership(newOwner);
    }


    // ==== admin functions ====

    event NameChanged(string newName, string newSymbol);

    /// Change the name and ticker symbol of this token.
    function changeName(string memory newName, string memory newSymbol) public only(owner) {
        name = newName;
        symbol = newSymbol;
        emit NameChanged(newName, newSymbol);
    }

    event Paused(address account);
    event Unpaused(address account);

    /// Pause the contract.
    function pause() external only(pauser) {
        paused = true;
        emit Paused(pauser);
    }

    /// Unpause the contract.
    function unpause() public only(pauser) {
        paused = false;
        emit Unpaused(pauser);
    }

    /// Modifies a function to run only when the contract is not paused.
    modifier notPaused() {
        require(!paused, "contract is paused");
        _;
    }

    event Frozen(address indexed freezer, address indexed account);
    event Unfrozen(address indexed freezer, address indexed account);
    event Wiped(address indexed freezer, address indexed wiped);

    /// Freeze token transactions for a particular address.
    function freeze(address who) external only(freezer) {
        require(data.frozenTime(who) == 0, "account already frozen");
        data.setFrozenTime(who, now); // solium-disable-line security/no-block-members
        emit Frozen(freezer, who);
    }

    /// Unfreeze token transactions for a particular address.
    function unfreeze(address who) external only(freezer) {
        require(data.frozenTime(who) > 0, "account not frozen");
        data.setFrozenTime(who, 0);
        emit Unfrozen(freezer, who);
    }

    /// Modifies a function to run only when the `account` is not frozen.
    modifier notFrozen(address account) {
        require(data.frozenTime(account) == 0, "account frozen");
        _;
    }

    /// Burn the balance of an account that has been frozen for at least 4 weeks.
    function wipe(address who) external only(freezer) {
        require(data.frozenTime(who) > 0, "cannot wipe unfrozen account");
        require(data.frozenTime(who) + 4 weeks < now, "cannot wipe frozen account before 4 weeks");
        _burn(who, data.balance(who));
        emit Wiped(freezer, who);
    }


    // ==== token transfers, allowances, minting, and burning ====

    /// @return how many tokens exist.
    function totalSupply() external view returns (uint256) {
        return _totalSupply;
    }

    /// @return how many tokens are held by `_owner`.
    function balanceOf(address _owner) external view returns (uint256) {
        return data.balance(_owner);
    }

    /// @return how many tokens `_owner` has allowed `spender` to control.
    function allowance(address _owner, address spender) external view returns (uint256) {
        return data.allowed(_owner, spender);
    }

    /// Transfer `value` attotokens from `msg.sender` to `to`.
    function transfer(address to, uint256 value)
        external
        notPaused
        notFrozen(msg.sender)
        notFrozen(to)
        returns (bool)
    {
        _transfer(msg.sender, to, value);
        return true;
    }

    /**
     * Approve `spender` to spend `value` attotkens on behalf of `msg.sender`.
     *
     * Beware that changing a nonzero allowance with this method brings the risk that
     * someone may use both the old and the new allowance by unfortunate transaction ordering. One
     * way to mitigate this risk is to first reduce the spender's allowance
     * to 0, and then set the desired value afterwards, per
     * [this ERC-20 issue](https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729).
     *
     * A simpler workaround is to use `increaseAllowance` or `decreaseAllowance`, below.
     *
     * @param spender address The address which will spend the funds
     * @param value uint256 How many attotokens to allow `spender` to spend.
     */
    function approve(address spender, uint256 value)
        external
        notPaused
        notFrozen(msg.sender)
        notFrozen(spender)
        returns (bool)
    {
        _approve(msg.sender, spender, value);
        return true;
    }

    /**
     * Transfer approved tokens from one address to another.
     *
     * @param from address The address to send tokens from
     * @param to address The address to send tokens to
     * @param value uint256 The amount of tokens to send
     */
    function transferFrom(address from, address to, uint256 value)
        external
        notPaused
        notFrozen(msg.sender)
        notFrozen(from)
        notFrozen(to)
        returns (bool)
    {
        _transfer(from, to, value);
        _approve(from, msg.sender, data.allowed(from, msg.sender).sub(value));
        return true;
    }

    /// Increase `spender`'s allowance of the sender's tokens.
    /// @dev From MonolithDAO Token.sol
    /// @param spender The address which will spend the funds.
    /// @param addedValue How many attotokens tokens to increase the allowance by.
    function increaseAllowance(address spender, uint256 addedValue)
        external
        notPaused
        notFrozen(msg.sender)
        notFrozen(spender)
        returns (bool)
    {
        _approve(msg.sender, spender, data.allowed(msg.sender, spender).add(addedValue));
        return true;
    }

    /// Decrease `spender`'s allowance of the sender's tokens.
    /// @dev From MonolithDAO Token.sol
    /// @param spender The address which will spend the funds.
    /// @param subtractedValue How many attotokens to decrease the allowance by.
    function decreaseAllowance(address spender, uint256 subtractedValue)
        external
        notPaused
        notFrozen(msg.sender)
        notFrozen(spender)
        returns (bool)
    {
        _approve(msg.sender, spender, data.allowed(msg.sender, spender).sub(subtractedValue));
        return true;
    }

    /// @dev Transfer of `value` attotokens from `from` to `to`.
    /// Internal; doesn't cheeck permissions.
    function _transfer(address from, address to, uint256 value) internal {
        require(to != address(0), "can't transfer to address zero");

        data.subBalance(from, value);
        data.addBalance(to, value);
        emit Transfer(from, to, value);
    }

    /// Mint `value` new attotokens to `account`.
    function mint(address account, uint256 value) external notPaused only(minter) {
        require(account != address(0), "can't mint to address zero");

        _totalSupply = _totalSupply.add(value);
        data.addBalance(account, value);
        emit Transfer(address(0), account, value);
    }

    /// Burn `value` attotokens from `account`, if sender has that much allowance from `account`.
    function burnFrom(address account, uint256 value) external notPaused only(minter) {
        _burn(account, value);
        _approve(account, msg.sender, data.allowed(account, msg.sender).sub(value));
    }

    /// @dev Burn `value` attotokens from `account`.
    /// Internal; doesn't check permissions.
    function _burn(address account, uint256 value) internal {
        require(account != address(0), "can't burn from address zero");

        _totalSupply = _totalSupply.sub(value);
        data.subBalance(account, value);
        emit Transfer(account, address(0), value);
    }

    /// @dev Set `spender`'s allowance on `_owner`'s tokens to `value`.
    /// Internal; doesn't check permissions.
    function _approve(address _owner, address spender, uint256 value) internal {
        require(spender != address(0), "spender cannot be address zero");
        require(_owner != address(0), "_owner cannot be address zero");

        data.setAllowed(_owner, spender, value);
        emit Approval(_owner, spender, value);
    }
}
