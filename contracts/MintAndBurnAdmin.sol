pragma solidity ^0.5.4;

import "./ReserveDollar.sol";

/// @title Time-delayed admin contract for the Reserve Dollar
/// Can only execute proposals after a 12-hour confirmation delay.
contract MintAndBurnAdmin {

    // TYPES

    struct Proposal {
        address addr;
        uint256 value;
        uint256 time;
        bool isMint;
    }

    // DATA

    ReserveDollar public reserve;
    uint256 public constant delay = 12 hours;
    address public admin;

    uint256 public nextProposal;
    mapping(uint256 => Proposal) public proposals;
    mapping(uint256 => bool) public closed;

    // EVENTS

    event ProposalCreated(uint256 index, address indexed addr, uint256 value, bool isMint, uint256 delayUntil);
    event ProposalConfirmed(uint256 index, address indexed addr, uint256 value, bool isMint);
    event ProposalCancelled(uint256 index, address indexed addr, uint256 value, bool isMint);

    // FUNCTIONALITY

    constructor(address reserveDollar) public {
        reserve = ReserveDollar(reserveDollar);
        admin = msg.sender;
    }

    /// Propose a new mint or burn, which can be confirmed after 12 hours.
    function propose(address addr, uint256 value, bool isMint) external {
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

    /// Cancel a proposed mint or burn.
    function cancel(uint256 index, address addr, uint256 value, bool isMint) external {
        // Check authorization.
        require(msg.sender == admin, "must be admin");
        requireMatchingOpenProposal(index, addr, value, isMint);

        // Cancel proposal.
        closed[index] = true;
        emit ProposalCancelled(index, addr, value, isMint);
    }

    /// Confirm and execute a proposed mint or burn, if enough time has passed since the proposal.
    function confirm(uint256 index, address addr, uint256 value, bool isMint) external {
        // Check authorization.
        require(msg.sender == admin, "must be admin");
        requireMatchingOpenProposal(index, addr, value, isMint);

        // See commentary above about using `now`.
        // solium-disable-next-line security/no-block-members
        require(proposals[index].time < now, "too early");

        // Record execution of proposal.
        closed[index] = true;
        emit ProposalConfirmed(index, addr, value, isMint);

        // Proceed with execution of proposal.
        if (proposals[index].isMint) {
            reserve.mint(addr, value);
        } else {
            reserve.burnFrom(addr, value);
        }
    }

    /// Throw unless the given proposal exists and matches `addr`, `value`, and `isMint`.
    function requireMatchingOpenProposal(uint256 index, address addr, uint256 value, bool isMint) private view {
        require(index < nextProposal, "no such proposal");
        require(!closed[index], "proposal already closed");

        // Slither reports "dangerous strict equality" for each of these, but it's OK.
        // These equalities are to confirm that the proposal entered is equal to the
        // matching previous proposal. We're vetting data entry; strict equality is appropriate.
        require(proposals[index].addr == addr, "addr mismatched");
        require(proposals[index].value == value, "value mismatched");
        require(proposals[index].isMint == isMint, "isMint mismatched");
    }
}
