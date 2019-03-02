pragma solidity ^0.5.4;

import "./ReserveDollar.sol";

/**
 * @title Admin contract for the Reserve Dollar
 *
 * @dev Time-delayed admin contract.
 * Allows performing actions with a 12-hour delay.
 */
contract MintAndBurnAdmin {
    ReserveDollar public reserve;
    uint256 public delay = 12 hours;
    address public admin;

    constructor(address reserveDollar) public {
        reserve = ReserveDollar(reserveDollar);
        admin = msg.sender;
    }

    struct Proposal {
        address addr;
        uint256 value;
        uint256 time;
        bool isMint;
    }
    uint256 public nextProposal;
    mapping (uint256 => Proposal) public proposals;
    mapping (uint256 => bool) public completed;

    event ProposalCreated(uint256 index, address addr, uint256 value, bool isMint, uint256 delayUntil);
    event ProposalConfirmed(uint256 index, address addr, uint256 value, bool isMint);

    /**
     * @dev Propose a new mint or burn, which can be confirmed after 12 hours.
     */
    function propose(address addr, uint256 value, bool isMint) public {
        require(msg.sender == admin, "must be admin");

        // Delay by at least 12 hours.
        // We are relying on block.timestamp for this, and aware of the possibility of its
        // manipulation by miners. But given the in-protocol bounds on the change in block.timestamp
        // and the way we are using it, we are satisfied with this choice.
        uint256 delayUntil = now + delay; // solium-disable-line security/no-block-members

        proposals[nextProposal] = Proposal({
            addr: addr,
            value: value,
            isMint: isMint,
            time: delayUntil
        });

        emit ProposalCreated(nextProposal, addr, value, isMint, delayUntil);

        nextProposal++;
    }

    /**
     * @dev Confirm a proposed mint or burn.
     *
     * If enough time has passed since the proposal, the owner
     * of the admin contract can confirm it.
     */
    function confirm(uint256 index, address addr, uint256 value, bool isMint) public {
        // Ensure proposal is authorized.
        require(msg.sender == admin, "must be admin");
        require(index < nextProposal, "no such proposal");
        // See commentary above about using `now`.
        require(proposals[index].time < now, "too early"); // solium-disable-line security/no-block-members
        require(!completed[index], "already completed");

        // Sanity-check inputs.
        require(proposals[index].addr == addr, "addr mismatched");
        require(proposals[index].value == value, "value mismatched");
        require(proposals[index].isMint == isMint, "isMint mismatched");

        // Proceed with action.
        if (proposals[index].isMint) {
            reserve.mint(addr, value);
        } else {
            reserve.burnFrom(addr, value);
        }

        // Record completion.
        completed[index] = true;

        emit ProposalConfirmed(index, addr, value, isMint);
    }
}
