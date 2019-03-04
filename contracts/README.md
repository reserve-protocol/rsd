# Smart Contracts for the Reserve Dollar

The Reserve Dollar (RSVD) is an [ERC-20][] token. RSVD is designed to maintain a stable price on the open market, through issuance and redemption of Reserve Dollars for fiat currency: 1 RSVD per 1 USD.

There are three main smart contracts in `contracts/`: `ReserveDollar`, `ReserveDollarEternalStorage`, and `MintAndBurnAdmin`.

* `ReserveDollar` is the main token implementation. It provides ERC-20 and administration interfaces.
* `ReserveDollarEternalStorage` is a static implementation of the [Eternal Storage][] pattern for `ReserveDollar`. Any non-constant-sized data for `ReserveDollar` is stored there, like mapping from addresses to balances.
* `MintAndBurnAdmin` is intended to hold the `minter` role for `ReserveDollar`. We use this to give ourselves time to respond and recover in case our operational keys are stolen.

[Eternal Storage]: https://fravoll.github.io/solidity-patterns/eternal_storage.html
[ERC-20]: http://eips.ethereum.org/EIPS/eip-20
## Philosophy
> In secure code, boring is better than magical.

We have aimed to prioritize clarity, robustness, static checks, and security over concision, abstraction, and programmer convenience. The latter items are _good_, but in this code base, the former items are _better_.

We have a set of needed features, and our primary optimization target is our ability to reach justified confidence in the correctness of our implementation of those features. This has led to a somewhat different style in implementing these smart contracts than seems fashionable.

Some consequences of this philosophy:

- We have avoided the `DelegateCall` upgrading pattern. After a great deal of consideration, we have also avoided other more-static proxying designs.
- Our [Eternal Storage][] implementation is a set of concrete maps, rather than a single, dynamic key-value store as in the paradigm implementation.
- Our system for roles is similarly concrete and non-dynamic. Role names are represented in static code, not as data.
- We deliberately use far less contract inheritance than is fashionable.
- We're avoiding JavaScript in our testing and deployment frameworks, even though many Ethereum tools are centered in JavaScript. We prefer Go for its quite-explicit static typing. (I _yearn_ for stronger types than that.)

We've aimed to keep these choices from yielding *tedious* code -- tedium makes it hard to focus, so tedium leads to mistakes. Nonetheless, in secure code, boring is better than magical.

## Features of the RSVD Token

### User-Facing Features

RSVD offers the [ERC-20][] interface for compatibility with standard wallet software. So that clients can, in many cases, more efficiently work around the ERC-20 [approval API bug][API Bug], RSVD also implements `increaseAllowance` and `decreaseAllowance`.

[API Bug]: https://docs.google.com/document/d/1YLPtQxZu1UAvO9cZ1O2RPXBbT0mooh4DYKjA_jp-RLM/edit

These are the only non-view API features of the RSVD token that an Ethereum account should be able to call without special authorization.

### Administration Roles and Features
The RSVD contacts have a number of features that should be usable only with special authorization. An account has a special authorization if a contract has that account's address assigned to some _role_.

Roles in the `ReserveDollar` contract:

- `owner` is the master key for the entire RSVD token. In our current plans, `owner` is securely stored off-site, and only accessed in deployment, contract upgrade, and slow system recovery operations. `owner` is authorized to:
    - Change the name and ticker symbol of the Reserve Dollar. (`changeName`)
    - Change the addresses of any other roles. (`changeMinter`, `changePauser`, `changeFreezer`)
    - Nominate a new owner. (`nominateNewOwner`).
    - Make a different (presumably new) Reserve Dollar contract the owner of the Eternal Storage contract. (`transferEternalStorage`)
- `nominatedOwner` is usually zero or `owner`. If `nominatedOwner` is a valid address, that nominee can become `owner` by calling `acceptOwnership()`.
- `minter` is expected to be the address of a deployed `MintAndBurnAdmin` contract. `minter` is authorized to:
    - Change the address of `minter`. (`changeMinter`)
    - Mint new tokens to some address. (`mint`)
    - Burn a user's tokens after that account has `approve`d the user to do so.
- `pauser` is an on-site but infrequently accessed role. `pauser` is authorized to:
    - Change the address of `pauser`
    - Pause and unpause the token. (`pause`, `unpause`)
- `freezer` is an on-site but infrequently accessed role. `freezer` is authorized to:
    - Freeze and unfreeze an account's tokens (`freeze`, `unfreeze`)
    - Burn an user's tokens after the user has been frozen for at least four weeks. (`wipe`)

Roles in the `ReserveDollarEternalStorage` contract:

- Without authorization, any user can get values.
- `owner` is expected to be the current version of the `ReserveDollar` contract. `owner` is authorized to:
    - Set stored values (`addBalance`, `subBalance`, `setBalance`, `setAllowed`, `setFrozenTime`)
    - Change the `owner` address. (`transferOwnership`)
- `escapeHatch` is expected to be a Reserve-held external account, to be used only if setting `owner` is mishandled during an upgrade. `escapeHatch` is authorized to:
    - Change the `owner` address. (`transferOwnership`)
    - Change the `escapeHatch` address. (`transferEscapeHatch`)

The only role in the `MintAndBurnAdmin` contract is `admin`. Only the `admin` account can call any functions of `MintAndBurnAdmin`. The `admin` is authorized to:

- Propose new minting and burning actions. (`propose`)
- Cancel proposals before they're executed. (`cancel`)
- Confirm proposals once they've been proposed for 12 hours. (`confirm`)

#### Minting and Burning
The `minter` role in the `ReserveDollar` contract is expected to be held by a deployed `MintAndBurnAdmin` contract. This is a somewhat complicated choice, so it's worth noting the rationale:

Minting and burning are frequent and necessary administration actions in the life of our contract. We'll need to do this perhaps several times a day during normal operations. That means that the private key with minting-and-burning authorization will need to be on-site and relatively accessible. *That* means that it's a prime target for a potential attacker. If that key was able to arbitrarily and quickly mint tokens, this would leave us with substantial weakness to exfiltration or insider attack.

We should expect that this key will get stolen _sometimes_. We have to have a useful response to such theft in some productive way.

The plan, then, is to hold the private key named by `MintAndBurnAdmin.admin`, and for the `ReserveDollar.minter` to be the `MintAndBurnAdmin` contract address. `MintAndBurnAdmin` can take arbitrary minting and burning actions by proposing the action, waiting at least 12 hours without cancelling the action, and then confirming the action.

So, if an attacker gets to be `MintAndBurnAdmin.admin`, they can easily DoS our minting and burning actions, but we can easily DoS their minting and burning actions. The `MintAndBurnAdmin.admin` does _not_ have the ability to transfer the `admin` address, so we can keep this up until we recover the `ReserveDollar.owner` key from its off-site location. We can deploy a new `MintAndBurnAdmin` contract with a new `admin` address, and once we have the `owner` key, we make the new contract the new `minter`.

#### Transferring Ownership
`ReserveDollar.owner` is the master key for the overall token. In the case where we need to update `ReserveDollar.owner`, we would undergo substantial expense if we accidentally, say, change the `owner` to an address we do not actually control. Since this will be a highly infrequent update, we should expect bitrot and human error.

To mitigate this risk, we require an accepting transaction from the newly-nominated `owner` address before actually setting `owner`. We encode this handshake with the functions `nominateNewOwner` and `acceptOwnership`:
- Only owner can call `nominateNewOwner(nominee)`, which just sets `nominatedOwner = nominee`.
- Only `owner` or `nominatedOwner` can call `acceptOwnership`, which clears `nominatedOwner`, sets `owner = msg.sender`, and emits an `OwnerChanged` event if the owner changed.

#### Transferring Eternal Storage / Upgrading
