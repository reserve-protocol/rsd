pragma solidity ^0.5.4;

import "ReserveDollar.sol";

contract ReserveDollarInner is ReserveDollar {
    constructor(address msgSender) public {
        owner = msgSender;
    }

    function testFixtureChangeOwner(address newOwner) public {
        owner = newOwner;
    }
}

contract WrappedReserveDollar {
    ReserveDollarInner r;

    constructor() public {
        r = new ReserveDollarInner(address(this));
    }

    function getEternalStorageAddress() public view returns(address) {
        return r.getEternalStorageAddress();
    }

    function name() public view returns (string memory) {
        return r.name();
    }

    function symbol() public view returns (string memory) {
        return r.symbol();
    }

    function decimals() public view returns (uint8) {
        return r.decimals();
    }

    function paused() public view returns (bool) {
        return r.paused();
    }

    function owner() public view returns (address) {
        return r.owner();
    }

    function minter() public view returns (address) {
        return r.minter();
    }

    function pauser() public view returns (address) {
        return r.pauser();
    }

    function freezer() public view returns (address) {
        return r.freezer();
    }

    function nominatedOwner() public view returns (address) {
        return r.nominatedOwner();
    }

    function changeMinter(address newMinter) external {
        return r.changeMinter(newMinter);
    }

    function changePauser(address newPauser) external {
        return r.changePauser(newPauser);
    }

    function changeFreezer(address newFreezer) external {
        return r.changeFreezer(newFreezer);
    }

    function nominateNewOwner(address nominee) external {
        return r.nominateNewOwner(nominee);
    }

    function acceptOwnership() external {
        return r.acceptOwnership();
    }

    function renounceOwnership() external {
        return r.renounceOwnership();
    }

    function transferEternalStorage(address newOwner) external {
        return r.transferEternalStorage(newOwner);
    }

    function changeName(string memory newName, string memory newSymbol) public {
        return r.changeName(newName, newSymbol);
    }

    function pause() external {
        return r.pause();
    }

    function unpause() public {
        return r.unpause();
    }

    function freeze(address who) external {
        return r.freeze(who);
    }

    function unfreeze(address who) external {
        return r.unfreeze(who);
    }

    function wipe(address who) external {
        return r.wipe(who);
    }

    function totalSupply() external view returns (uint256) {
        return r.totalSupply();
    }

    function balanceOf(address _owner) external view returns (uint256) {
        return r.balanceOf(_owner);
    }

    function allowance(address _owner, address spender) external view returns (uint256) {
        return r.allowance(_owner, spender);
    }

    function transfer(address to, uint256 value) external returns (bool) {
        return r.transfer(to, value);
    }

    function approve(address spender, uint256 value) external returns (bool) {
        return r.approve(spender, value);
    }

    function transferFrom(address from, address to, uint256 value) external returns (bool) {
        return r.transferFrom(from, to, value);
    }

    function increaseAllowance(address spender, uint256 addedValue) external returns (bool) {
        return r.increaseAllowance(spender, addedValue);
    }

    function decreaseAllowance(address spender, uint256 subtractedValue) external returns (bool) {
        return r.decreaseAllowance(spender, subtractedValue);
    }

    function mint(address account, uint256 value) external {
        return r.mint(account, value);
    }

    function burnFrom(address account, uint256 value) external {
        return r.burnFrom(account, value);
    }
}
