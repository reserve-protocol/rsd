# Reserve Dollar

>The Reserve Dollar (RSVD) is a fiat-backed stablecoin from [Reserve](https://reserve.org).

The Reserve Dollar is designed to maintain a stable price on the open market, by offering issuance and redemption of Reserve Dollars for fiat currency: 1 RSVD per 1 USD.

The Reserve Dollar is *not* the Reserve or Reserve Share token described in the Reserve [whitepaper](https://reserve.org/whitepaper). The whitepaper design is more fully decentralized, and more complicated. Rather, the Reserve Dollar is an initial stablecoin, suitable for building into products to test, being the first in a series of steps towards total decentralization, and to be an early-stage asset in the vault of the whitepaper design.

## What does it do?
The Reserve Dollar offers normal [ERC-20](http://eips.ethereum.org/EIPS/eip-20) behavior.

Reserve administers this token and offers issuance and redemption to verified users. Once a user has registered with our web portal (not in this repository) and entered information necessary for KYC/AML checks, the user will be eligible to buy and sell RSVD for fiat currency, directly from Reserve. To this end, Reserve controls external admin accounts. The admin accounts can mint and burn tokens to represent the on-chain side of these transactions. Various admin accounts can also freeze accounts, wipe long-frozen accounts, pause the entire contract, upgrade the contract itself, and update admin addresses.

## How does it fit together?
There are three main smart contracts in `contracts/`: `ReserveDollar`, `ReserveDollarEternalStorage`, and `MintAndBurnAdmin`.

* `ReserveDollar` is the main token implementation. It provides ERC-20 and administration interfaces.
* `ReserveDollarEternalStorage` is a static implementation of the [Eternal Storage][] pattern for `ReserveDollar` Any non-constant-sized data for `ReserveDollar` is stored there. However, because the token's storage is relatively simple, and the dynamic form of the EternalStorage pattern introduces difficulties in analysis, `ReserveDollarEternalStorage` provides accessors for _specific_ maps -- `balance`, `allowed`, and `frozenTime`.
* `MintAndBurnAdmin` is intended to hold the `minter` role for `ReserveDollar`. We use this to give ourselves time to respond and recover in case our operational keys are stolen.

[EternalStorage]: https://fravoll.github.io/solidity-patterns/eternal_storage.html

# Development environment

- Install Go >= 1.12
    - Recommended approach: [download the latest binary distribution](https://golang.org/dl/) and then [follow the installation instructions](https://golang.org/doc/install#install)
    - Check with `go version`

- Install node and npm
    - Recommended approach: [download from nodejs.org](https://nodejs.org/en/)

- Run `npm install`

## Dockerized environment

There is also a dockerized version of the development environment. You can open it with `make run-dev-container`. It's not intended to handle all development workflows, but you should be able to successfully run `make test` in it, and use it to troubleshoot your host environment if necessary.

# Running tests

Test with `make test`.

To get a coverage report, run `make coverage`. Note that it has a few obvious false negatives, like the ReserveDollar constructor, interface definitions, and the `_;` line in modifiers.
