pragma solidity ^0.5.4;

import "TestFixtureReserveDollar.sol";

// Use the multisender.yaml to send messages from 0x1, 0x2, 0x3, but not the
// deployer address. Non owners should be unable to change the owner
contract EchidnaFrozenAccount is TestFixtureReserveDollar {
    bool transferredFromFrozen;

    constructor() public {
        minter = msg.sender;
        freezer = msg.sender;
        
        testFixtureMint(address(0x1), 123);
        testFixtureMint(address(0x2), 1234);
        testFixtureMint(address(0x3), 12345);

        // Create an allowance for transfering into the frozen address.
        testFixtureIncreaseAllowance(address(0x2), address(0x1), 17);
        // Create an allowance for transfering out of the frozen address.
        testFixtureIncreaseAllowance(address(0x1), address(0x2), 19);

        // freeze the 0x1 address;
        testFixtureFreeze(address(0x1));
        
        transferredFromFrozen = false;
    }
    
    // Only allow allowance transfers into 0x3 to check that frozen account
    // cannot send funds to a third account from an account that gives frozen
    // account an allowance.
    function transfer(address to, uint256 value) external returns (bool)
    {
        require(to != address(0x3));
        return this.transfer(to, value);
    }

    ///////////////////////////////// 
    

    function echidna_frozen_account_constant_balance() public view returns (bool) {
        return this.balanceOf(address(0x1)) == 123;
    }
    
    function echidna_prevent_frozen_msg_sender_transferFrom() public view returns (bool) {
        return (this.balanceOf(address(0x3)) == 12345);
    }
    
    function echidna_prevent_transfers_into_frozen_account() public view returns (bool) {
        return ((this.balanceOf(address(0x2)) == 1234)
            && (this.allowance(address(0x2), address(0x1)) == 17));
    }

    function echidna_prevent_transfers_from_frozen_account() public view returns (bool) {
        return ((this.balanceOf(address(0x2)) == 1234)
            && (this.allowance(address(0x2), address(0x1)) == 17));
    }

    function echidna_freezer_role_unchanged() public view returns (bool) {
        return this.freezer() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}
