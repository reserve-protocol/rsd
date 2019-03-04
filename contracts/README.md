# Smart Contracts for the Reserve Dollar

The Reserve Dollar (RSVD) is an [ERC-20][] token. RSVD is designed to maintain a stable price on the open market, through issuance and redemption of Reserve Dollars for fiat currency: 1 RSVD per 1 USD.

There are three main smart contracts in `contracts/`: `ReserveDollar`, `ReserveDollarEternalStorage`, and `MintAndBurnAdmin`.

* `ReserveDollar` is the main token implementation. It provides ERC-20 and administration interfaces.
* `ReserveDollarEternalStorage` is a static implementation of the [Eternal Storage][] pattern for `ReserveDollar`. Any non-constant-sized data for `ReserveDollar` is stored there, like mapping from addresses to balances.
* `MintAndBurnAdmin` is intended to hold the `minter` role for `ReserveDollar`. We use this to give ourselves time to respond and recover in case our operational keys are stolen.

[Eternal Storage]: https://fravoll.github.io/solidity-patterns/eternal_storage.html
[ERC-20]: http://eips.ethereum.org/EIPS/eip-20
# Philosophy
> In secure code, boring is better than magical.

We have aimed to prioritize clarity, robustness, static checks, and security over concision, abstraction, and programmer convenience. The latter items are _good_, but in this code base, the former items are _better_.

We have a set of needed features, and our primary optimization target is our ability to reach justified confidence in the correctness of our implementation of those features. This has led to a somewhat different style in implementing these smart contracts than seems fashionable.

Some consequences of this philosophy:

- We have avoided the `DelegateCall` upgrading pattern. After a great deal of consideration, we have also avoided other more-static proxying designs.
- Our [Eternal Storage][] implementation is a set of concrete maps, rather than a single, dynamic key-value store as in the paradigm implementation.
- Our system for roles is similarly concrete and non-dynamic. Role names are represented in static code, not as data.
- We deliberately use far less contract inheritance than is fashionable.
- We're avoiding JavaScript in our testing and deployment frameworks, even though many Ethereum tools are centered in JavaScript. We prefer Go, not least for its explicitness and static typing.

We've aimed to keep these choices from yielding *tedious* code -- tedium makes it hard to focus, so tedium leads to mistakes. Nonetheless, in secure code, boring is better than magical.

# Features of the RSVD Token

## User-Facing Features

RSVD offers the [ERC-20][] interface for compatibility with standard wallet software. So that clients can more efficiently work around the ERC-20 [approval API bug][API Bug], RSVD also implements `increaseAllowance` and `decreaseAllowance`.

[API Bug]: https://docs.google.com/document/d/1YLPtQxZu1UAvO9cZ1O2RPXBbT0mooh4DYKjA_jp-RLM/edit

These are the only non-view API features of the RSVD token that an Ethereum account should be able to call without special authorization.

## Administration Roles and Features
The RSVD contacts have a number of features that should be usable only with special authorization. An account has a special authorization if a contract has that account's address assigned to some _role_.

### Roles in `ReserveDollar`

- `owner` is the master key for the entire RSVD token. The private key for `owner` is securely stored off-site, and only accessed in deployment, contract upgrade, and slow system recovery operations. `owner` is authorized to:
    - Change the name and ticker symbol of the Reserve Dollar. (`changeName`)
    - Change the addresses of any other roles. (`changeMinter`, `changePauser`, `changeFreezer`)
    - Nominate a new owner. (`nominateNewOwner`)
    - Make another contract the owner of the Eternal Storage contract. (`transferEternalStorage`)
- `nominatedOwner` is usually zero or `owner`. The `nominatedOwner` can become `owner` by calling `acceptOwnership()`.
- `minter` is expected to be the address of a deployed `MintAndBurnAdmin` contract. `minter` is authorized to:
    - Change the address of `minter`. (`changeMinter`)
    - Mint new tokens to some address. (`mint`)
    - Burn a user's tokens after the user set `allowance` for `minter` to do so. (`burn`)
- `pauser` is an infrequently accessed admin role, with private key stored on-site. `pauser` is authorized to:
    - Change the address of `pauser`. (`changePauser`)
    - Pause and unpause the token (`pause`, `unpause`), in case of [emergency](https://consensys.github.io/smart-contract-best-practices/software_engineering/#circuit-breakers-pause-contract-functionality).
- `freezer` is an infrequently accessed admin role with private key stored on-site. In practice, it might be equal to `pauser`. `freezer` is authorized to:
    - Change the address of `freezer`. (`changeFreezer`)
    - Freeze and unfreeze an account's tokens (`freeze`, `unfreeze`)
    - Burn an user's tokens after the user has been frozen for at least four weeks. (`wipe`)

### Roles in `ReserveDollarEternalStorage`

- Without authorization, any user can get values.
- `owner` is expected to be the current version of the `ReserveDollar` contract. `owner` is authorized to:
    - Set stored values. (`addBalance`, `subBalance`, `setBalance`, `setAllowed`, `setFrozenTime`)
    - Change the `owner` address. (`transferOwnership`)
- `escapeHatch` is expected to be a Reserve-held external account, to be used only if setting `owner` is mishandled during an upgrade. `escapeHatch` is authorized to:
    - Change the `owner` address. (`transferOwnership`)
    - Change the `escapeHatch` address. (`transferEscapeHatch`)

### Roles in `MintAndBurnAdmin`
The only role in the `MintAndBurnAdmin` contract is `admin`, and only the `admin` account can call any functions of `MintAndBurnAdmin`. The `admin` is authorized to:

- Propose new minting and burning actions. (`propose`)
- Cancel proposals before they're executed. (`cancel`)
- Cancel all current proposals. (`cancelAll`)
- Confirm proposals once they've been proposed for 12 hours. (`confirm`)

# Some Technical Workflows
These explanations are intended both to document our less-obvious workflow plans, and to aid in understanding those elements of our smart contracts set up to facilitate those flows.

## Minting and Burning
The `minter` role in the `ReserveDollar` contract is expected to be held by a deployed `MintAndBurnAdmin` contract. This is a somewhat complicated choice, so it's worth noting the rationale:

Minting and burning are frequent and necessary administration actions in the life of our contract. We'll need to do this perhaps several times a day during normal operations. That means that the private key with minting-and-burning authorization will need to be on-site and relatively accessible. *That* means that it's a prime target for a potential attacker. If that key was able to arbitrarily and quickly mint tokens, this would leave us with substantial weakness to exfiltration or insider attack.

We should _expect_ that this key will get stolen. We need a useful response to that theft. Here's the plan:

`ReserveDollar.minter` will be the address of a  `MintAndBurnAdmin` contract. We will hold the private key for the `MintAndBurnAdmin.admin` account.  `MintAndBurnAdmin` performs arbitrary minting and burning when its `admin` proposes an action, waits 12 hours without cancelling the proposal, and then confirms the proposal. If the `admin` cancels a proposal, it cannot later be confirmed.

So, suppose an attacker gets the private key to `MintAndBurnAdmin.admin`, both we and they can act as that `admin`. What can they do?
- If they try to mint or burn tokens, we can easily cancel their proposals.
- If they issue many minting or burning proposals, we can cheaply cancel their proposals once every few hours with `cancelAll`.
- They can cancel all of our proposals, causing a denial of service on new mints and burns. This will remain true until we deploy a new `MintAndBurnAdmin` contract with a new `admin` address, and make that new contract's address the new `ReserveDollar.minter`. We can do once we recover our `owner` key from off-site. This may take longer than 12 hours, but will not be longer than 7 days.
- They can _not_ transfer the `admin` address to some other account, because `MintAndBurnAdmin.admin` does _not_ have the ability to change the `admin` address. Thus, the attacker can't thereby deny our ability to carry out the above attack responses.

## Transferring Ownership
`ReserveDollar.owner` is the master key for the overall token. If we need to update `owner`, we would take a serious loss if we accidentally changed the `owner` to an address that we do not control. Since changing `owner` should be highly infrequent, we should expect bit-rot and human error in these changes.

To mitigate this risk, `ReserveDollar` require an transaction from the newly-nominated `owner` before it actually changes `owner`. We encode this two-step handshake in functions `nominateNewOwner` and `acceptOwnership`:
- Only `owner` can call `nominateNewOwner(nominee)`, which just sets `nominatedOwner = nominee`.
- Only `owner` or `nominatedOwner` can call `acceptOwnership()`, which clears `nominatedOwner`, sets `owner = msg.sender`, and emits an `OwnerChanged` event if the owner changed.

## Upgrading
We do not have specific intentions to upgrade or migrate the `ReserveDollar` contract, but we intend to be prepared to do so in case of unforeseen failure modes. We've decided not to use the `delegatecall` pattern for this purpose; it's a central example of code being magic instead of boring. (And it requires consistency properties that cross-cut standard programmer models of the implementation language! Worst kind of magic! Panic! Panic! Considered Harmful; Run Away!)

The three major smart contracts in the system now serve three *contract roles*:

- `ReserveDollar` fills the Logic role.
- `ReserveDollarEternalStorage` fills the Storage role.
- `MintAndBurnAdmin` fills the MinterAdmin role.

There should be bidirectional references between the Logic role and each of the other roles. In the current implementation, those references are:

- `ReserveDollar.data` <--> `ReserveDollarEternalStorage.owner`
- `ReserveDollar.minter` <--> `MintAndBurnAdmin.reserve`

Also, the Logic contract should be referred to by [ENS](https://ens.domains/), so that clients are able to double-check that they have an up-to-date address for that contract. (Certainly, this isn't yet standardized among clients, but it is The Right Way™.)

### Upgrading MinterAdmin
Suppose we need to upgrade the contract in the MinterAdmin role to `NewAdmin`. The upgrade path is pretty simple:

1. As `admin`, deploy `NewAdmin`. Initialize its `reserve` pointer to the `ReserveDollar` contract address. Call the address of the `NewAdmin` deployment `newAdminAddr`.
2. As `ReserveDollar.owner`, call `changeMinter(newAdminAddr)`.

### Upgrading Logic
Suppose, instead, we need to upgrade the contract in the Logic role to `NewLogic`. This is rather more complicated, but still doable in constant time:

1. As `owner`, deploy `NewLogic`. It should be paused on deployment. Let its new address be `newLogicAddr`.
2. As `owner`, initialize the administration roles and Storage address for `NewLogic` as they're set in `ReserveDollar`.
3. As `pauser`, call `ReserveDollar.pause()`.
4. As `owner`, call `ReserveDollar.transferEternalStorage(newLogicAddr)`.
5. On ENS, point our entry for the Reserve Dollar contract to `newLogicAddr`.
6. Wait for confirmation of (4) and (5), and double-check that the on-chain deployment is as intended.
7. As `pauser`, call `NewLogic.unpause()`.

### Upgrading Storage
Don't.

We aim to never need to upgrade the Storage contract. To this end, `ReserveDollarEternalStorage` is carefully built and extremely simple.

If a Logic upgrade requires us to expand the vocabulary of the Storage contract to new arrays or maps, then we will deploy a second Storage contract to store just that new vocabulary, and have the new Logic contract operate over both.

If we *really must* upgrade the Storage contract, then we'll need to carry out a full on-chain data migration -- which is possible, but slow, expensive, and beyond the scope of this document.

# Some Contract Properties

In `ReserveDollar`:

- `ReserveDollar` emits a `Transfer` event any time token balances change.
- `ReserveDollar` emits an `Approval` event any time allowances change.
- Balances only change inside the functions `mint`, `_transfer`, and `_burn`.
- Total supply only changes inside `mint` and `_burn`.
- Allowances only change inside `_approve`.
- `balanceOf[addr]` does not increase when `frozenTime[addr]` is nonzero.
- `balanceOf[addr]` does not decrease when `frozenTime[addr]` is nonzero, _except_ through `wipe`.
- Balances and allowances never change when `paused` is true.
