pragma solidity ^0.5.4;

import "TestFixtureReserveDollar.sol";

// Use the multisender.yaml to send messages from 0x1, 0x2, 0x3, but not the
// deployer address. Non owners should be unable to change the owner
contract EchidnaPausedContract is TestFixtureReserveDollar {
    constructor() public {
        minter = msg.sender;
        freezer = msg.sender;
        pauser = msg.sender;
        
        testFixtureMint(address(0x1), 123);
        testFixtureMint(address(0x2), 1234);
        testFixtureMint(address(0x3), 12345);

        testFixtureIncreaseAllowance(address(0x1), address(0x2), 19);
        testFixtureIncreaseAllowance(address(0x1), address(0x3), 31);
        testFixtureIncreaseAllowance(address(0x2), address(0x3), 53);
        
        // pause the contract
        testFixturePause();
        pauser = address(0x900);
    }
    // prevent the owner from pausing. We want the the owner to be able to do
    // everything but unpause.
    function changePauser(address) external {
        return;
    }
    
    function echidna_while_paused_balances_are_fixed() public view returns (bool) {
        return (this.balanceOf(address(0x1)) == 123)
            &&( this.balanceOf(address(0x2)) == 1234)
            &&( this.balanceOf(address(0x3)) == 12345);
    }

    function echidna_while_paused_allowance_are_fixed() public view returns (bool) {
        return (this.allowance(address(0x1), address(0x2)) == 19)
            && (this.allowance(address(0x1), address(0x3)) == 31)
            && (this.allowance(address(0x2), address(0x3)) == 53);
    }
    
    function echidna_constant_supply() public view returns (bool) {
        return this.totalSupply() == 13702;
    }
    
    function echidna_cannot_be_unpaused_without_pauser_role() public view returns (bool) {
        return this.pauser() == address(0x900);
    }
    
    function echidna_cannot_be_unpaused() public view returns (bool) {
        return this.paused();
    }
}

