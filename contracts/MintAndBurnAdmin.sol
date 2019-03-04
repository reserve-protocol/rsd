pragma solidity ^0.5.4;

import "./ReserveDollar.sol";

/// @title Time-delayed admin contract for the Reserve Dollar
/// Can only execute proposals after a 12-hour confirmation delay.
contract MintAndBurnAdmin {

    // TYPES

    struct Proposal {
        address addr;
        uint256 value;
        bool isMint;
        uint256 time;
        bool closed;
    }

    // DATA

    ReserveDollar public reserve;
    uint256 public constant delay = 12 hours;
    address public admin;

    Proposal[] public proposals;

    // EVENTS

    event ProposalCreated(uint256 index, address indexed addr, uint256 value, bool isMint, uint256 delayUntil);
    event ProposalConfirmed(uint256 index, address indexed addr, uint256 value, bool isMint);
    event ProposalCancelled(uint256 index, address indexed addr, uint256 value, bool isMint);
    event AllProposalsCancelled();

    // FUNCTIONALITY

    constructor(address reserveDollar) public {
        reserve = ReserveDollar(reserveDollar);
        admin = msg.sender;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "must be admin");
        _;
    }

    /// Propose a new mint or burn, which can be confirmed after 12 hours.
    function propose(address addr, uint256 value, bool isMint) external onlyAdmin {
        // Delay by at least 12 hours.
        // We are relying on block.timestamp for this, and aware of the possibility of its
        // manipulation by miners. But given the in-protocol bounds on the change in
        // block.timestamp and the way we are using it, we are satisfied with this choice.
        // solium-disable-next-line security/no-block-members
        uint256 delayUntil = now + delay;

        proposals.push(Proposal({
            addr: addr,
            value: value,
            isMint: isMint,
            time: delayUntil,
            closed: false
        }));

        emit ProposalCreated(proposals.length - 1, addr, value, isMint, delayUntil);
    }

    /// Cancel a proposed mint or burn.
    function cancel(uint256 index, address addr, uint256 value, bool isMint) external onlyAdmin {
        // Check authorization.
        requireMatchingOpenProposal(index, addr, value, isMint);

        // Cancel proposal.
        proposals[index].closed = true;
        emit ProposalCancelled(index, addr, value, isMint);
    }

    /// Cancel all proposals.
    function cancelAll() external onlyAdmin {
        proposals.length = 0;
        emit AllProposalsCancelled();
    }

    /// Confirm and execute a proposed mint or burn, if enough time has passed since the proposal.
    function confirm(uint256 index, address addr, uint256 value, bool isMint) external onlyAdmin {
        // Check authorization.
        requireMatchingOpenProposal(index, addr, value, isMint);

        // See commentary above about using `now`.
        // solium-disable-next-line security/no-block-members
        require(proposals[index].time < now, "too early");

        // Record execution of proposal.
        proposals[index].closed = true;
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
        require(!proposals[index].closed, "proposal already closed");

        // Slither reports "dangerous strict equality" for each of these, but it's OK.
        // These equalities are to confirm that the proposal entered is equal to the
        // matching previous proposal. We're vetting data entry; strict equality is appropriate.
        require(proposals[index].addr == addr, "addr mismatched");
        require(proposals[index].value == value, "value mismatched");
        require(proposals[index].isMint == isMint, "isMint mismatched");
    }
}
