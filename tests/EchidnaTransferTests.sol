pragma solidity ^0.5.4;

import "ReserveDollar.sol";

contract DisabledTransferReserve is ReserveDollar {
    // disable the `burn` function to prevent the total supply from changing
    // from extected ways during the transfer testing
    function burn(uint256) external {}
    
    // disable the `transfer` function to force echidna to only transfer
    // between our three accounts that correspond to the 3 possible senders of
    // messages during the test
    function transfer(address, uint256) public returns (bool) {
        return false;
    }

    // disable the `transferFrom` function to force echidna to only transfer
    // between our three accounts that correspond to the 3 possible senders of
    // messages during the test
    function transferFrom(address, address, uint256) public returns (bool) {
        return false;
    }
}

// Test that 3 users that can transfer between themselves only cannot increase
// or decrease the total supply
contract EchidnaTotalSupplyTests is DisabledTransferReserve {
    constructor() public {
        // add some balance for users at 0x1, 0x2, and 0x3 addresses
        _totalSupply = _totalSupply.add(123);
        super.data.addBalance(address(0x1), 123);
        emit Transfer(address(0), address(0x1), 123);
    }

    // Function stubs to let echidna only tranfer into account 1
    function transferToOne(uint256 value) public returns (bool) {
        return super.transfer(address(0x1), value);
    }
    
    function transferToTwo(uint256 value) public returns (bool) {
        return super.transfer(address(0x2), value);
    }

    function transferToThree(uint256 value) public returns (bool) {
        return super.transfer(address(0x3), value);
    }
    
    function echidna_constant_supply() public view returns (bool) {
        return totalSupply() == 13702;
    }
    
    function echidna_zero_sum_transfers() public view returns (bool) {
        return (balanceOf(address(0x1)) + balanceOf(address(0x2)) + balanceOf(address(0x3))) == 13702;
    }
    
    // Check that the owner has not been changed
    function echidna_owner_unchanged_by_transfers() public view returns (bool) {
        return this.owner == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}

contract EchidnaAllowance is DisabledTransferReserve {

}

// Use the multisender.yaml to send messages from 0x1, 0x2, 0x3, but not the
// deployer address. Non owners should be unable to change the owner
contract EchidnaConstantOwnerSimple is ReserveDollar {
    // Check that the owner has not been changed
    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}

// Variation where attacking users can transfer any amounts they want. This
// lets us get random non zero values into the contract balances
contract EchidnaConstantOwnerFreeMoney is ReserveDollar {
    function transfer(address to, uint256 value) public returns (bool) {
        super._mint(msg.sender, value);
        return transfer(to, value);
    }

    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        super._mint(msg.sender, value);
        return transferFrom(from, to, value);
    }

    function echidna_owner_unchanged() public view returns (bool) {
        return this.owner == address(0x00a329c0648769A73afAc7F9381E08FB43dBEA72);
    }
}
