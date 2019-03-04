pragma solidity ^0.5.4;

// import "ReserveDollar.sol";
import "EchidnaTestBase.sol";

// contract DisabledTransferReserve {
//     ReserveDollar r;
// 
//     constructor() public {
//         r = new ReserveDollar();
//         require(r.owner() == 0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
//         // add some balance for users at 0x1, 0x2, and 0x3 addresses
//         // r.mint(address(0x1), 123);
//     }
// }

// Test that 3 users that can transfer between themselves only cannot increase
// or decrease the total supply
contract EchidnaTotalSupplyTests is WrappedReserveDollar {
    constructor() public {
        require(r.owner() == address(this));
        r.changeMinter(address(this));
        // add some balance for users at 0x1, 0x2, and 0x3 addresses
        r.mint(address(0x1), 123);
        r.mint(address(0x2), 1234);
        r.mint(address(0x3), 12345);
        
        // don't allow this contract to mint because it will be passed on as msg sender
        r.testFixtureChangeOwner(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
        r.changeMinter(address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72));
    }
    
    // Disable the `burnFrom` function to prevent the total supply from
    // changing from extracted ways during the transfer testing
    function burnFrom(address, uint256) external {}
    // Disable the `transfer` and `transferFrom` function to force echidna to
    // only transfer between our three accounts that correspond to the 3
    // possible senders of messages during the test
    function transfer(address, uint256) external returns (bool) { return false; }
    function transferFrom(address, address, uint256) external returns (bool) { return false; }
    // Disable `renounceOwnership`
    

    // Function stubs to let echidna only tranfer into account 1
    function transferToOne(uint256 value) public returns (bool) {
        return r.transfer(address(0x1), value);
    }
    
    function transferToTwo(uint256 value) public returns (bool) {
        return r.transfer(address(0x2), value);
    }

    function transferToThree(uint256 value) public returns (bool) {
        return r.transfer(address(0x3), value);
    }
    
    function echidna_constant_supply() public view returns (bool) {
        return r.totalSupply() == 13702;
    }
    
    function echidna_zero_sum_transfers() public view returns (bool) {
        return (r.balanceOf(address(0x1)) + r.balanceOf(address(0x2)) + r.balanceOf(address(0x3))) == 13702;
    }
    
    // Check that the owner has not been changed by transfers just in case.
    function echidna_owner_unchanged_by_transfers() public view returns (bool) {
        return r.owner() == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}

// Variation where attacking users can transfer any amounts they want. This
// lets us get random non zero values into the contract balances
contract EchidnaConstantOwnerFreeMoney is WrappedReserveDollar {
    constructor() public {
        require(r.owner() == address(this));
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
