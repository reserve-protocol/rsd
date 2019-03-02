pragma solidity ^0.5.4;

import "./ReserveDollar.sol";
import "./ReserveDollarEternalStorage.sol";

/**
 * @dev A version of Reserve Dollar for testing upgrades.
 */
contract ReserveDollarV2 is ReserveDollar {
    constructor() public {
        paused = true;
    }

    function completeHandoff(address previousImplementation) public only(owner) {
        ReserveDollar previous = ReserveDollar(previousImplementation);
        data = ReserveDollarEternalStorage(previous.getEternalStorageAddress());
        previous.acceptOwnership();

        // Take control of Eternal Storage.
        previous.transferEternalStorage(address(this));
        previous.changePauser(address(this));

        // Old contract off, new contract on.
        previous.pause();
        unpause();

        // Burn the bridge behind us.
        previous.changeMinter(address(0));
        previous.changePauser(address(0));
        previous.changeFreezer(address(0));
        previous.renounceOwnership();
    }
}
