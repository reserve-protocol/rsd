pragma solidity ^0.5.4;

import "ReserveDollar.sol";
import "EchidnaTestBase.sol";

// Check that there is no combination of messages an owner can send to change the owner other than
// `renounceOwnership`
contract EchidnaConstantOwnerSimple is ReserveDollar {
    // prevent ownership from being renounced during this test
    function renounceOwnership() external {
        return;
    }

    // Check that the owner has not been changed
    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}

// Try to find overflows by allowing all possible transfer amounts. Attacking
// users can transfer any amounts they want. This lets us get random non zero
// values into any contract balances.
contract EchidnaConstantOwnerFreeMoney is WrappedReserveDollar {
    constructor() public {
        require(r.owner() == address(this), "document expected contract ownership during setup");
        r.changeMinter(address(this));
        // add some balance for users at 0x1, 0x2, and 0x3 addresses
        r.mint(address(0x1), 123);
        r.mint(address(0x2), 1234);
        r.mint(address(0x3), 12345);

        // don't allow this contract to mint because this contract address will
        // be passed on as msg.sender
        r.testFixtureChangeOwner(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
        r.changeMinter(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
    }

    function transfer(address to, uint256 value) public returns (bool) {
        r.testFixtureChangeOwner(address(this));
        r.changeMinter(address(this));
        r.mint(msg.sender, value);
        bool result = transfer(to, value);
        r.changeMinter(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
        r.testFixtureChangeOwner(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));

        return result;
    }

    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        r.testFixtureChangeOwner(address(this));
        r.changeMinter(address(this));
        r.mint(from, value);
        bool result = transferFrom(from, to, value);
        r.changeMinter(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
        r.testFixtureChangeOwner(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));

        return result;
    }

    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}
