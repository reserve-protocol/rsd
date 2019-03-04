# Reserve Dollar

>The Reserve Dollar (RSVD) is a fiat-backed stablecoin from [Reserve](https://reserve.org).

The Reserve Dollar is designed to maintain a stable price on the open market, by offering issuance and redemption of Reserve Dollars for fiat currency: 1 RSVD per 1 USD.

The Reserve Dollar is *not* the Reserve or Reserve Share token described in the Reserve [whitepaper](https://reserve.org/whitepaper). The whitepaper design is more fully decentralized, and more complicated. Rather, the Reserve Dollar is an initial stablecoin, suitable for building into products to test, being the first in a series of steps towards total decentralization, and to be an early-stage asset in the vault of the whitepaper design.

## What does it do?
The Reserve Dollar offers normal [ERC-20](http://eips.ethereum.org/EIPS/eip-20) behavior.

Reserve administers this token and offers issuance and redemption to verified users. Once a user has registered with our web portal (not in this repository) and entered information necessary for KYC/AML checks, the user will be eligible to buy and sell RSVD for fiat currency. When a user buys RSVD for fiat through our system, we mint RSVD to give them; when a user sells RSVD for fiat through our system, we burn those RSVD.

To this end, Reserve controls external admin accounts. The admin accounts can mint and burn tokens to represent the on-chain side of these transactions. Various admin accounts can also freeze accounts, wipe long-frozen accounts, pause the token, upgrade the token, and update admin addresses.

## How does it fit together?
There are three main smart contracts in `contracts/`: `ReserveDollar`, `ReserveDollarEternalStorage`, and `MintAndBurnAdmin`.

* `ReserveDollar` is the main token implementation. It provides ERC-20 and administration interfaces.
* `ReserveDollarEternalStorage` is a static implementation of the [Eternal Storage][] pattern for `ReserveDollar`. Any non-constant-sized data for `ReserveDollar` is stored there. However, because the token's storage is relatively simple, and the dynamic form of the EternalStorage pattern introduces difficulties in analysis, `ReserveDollarEternalStorage` provides accessors for _specific_ maps -- `balance`, `allowed`, and `frozenTime`.
* `MintAndBurnAdmin` is intended to hold the `minter` role for `ReserveDollar`. We use this to give ourselves time to respond and recover in case our operational keys are stolen.

[Eternal Storage]: https://fravoll.github.io/solidity-patterns/eternal_storage.html

# Environment Setup

To build and test everything in our configuration, your development environment will need:

* **Go** -- 1.12 or later, to run many tools: ABI generation, tests, coverage, our contract CLI.
* **Node** and **NPM** -- to run `sol-compiler` and `solium`
* **Docker** -- to run echidna and 0x's a nicely-packaged `geth` node, for coverage and some end-to-end tests.

## Setting up a clean environment

- Install basic utilities for your system: `curl`, `git`, `make`, `gcc`.

- Install Go
    - [Download the latest binary distribution](https://golang.org/dl/)
    - [Follow the installation instructions](https://golang.org/doc/install#install)
    - Check that the version is 1.12 or later with `go version`

- Install Node and npm
    - [Download from nodejs.org](https://nodejs.org/en/)
    - In this repo's working directory, do `npm install` to get local packages

- Install Docker
    - A very standard Docker installation is sufficient:
        - On Windows: Instructions [here](https://docs.docker.com/docker-for-windows/install/)
        - On MacOS: Instructions [here](https://docs.docker.com/docker-for-mac/install/)
        - On Linux: `curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh`

## Dockerized environment

There is also a dockerized version of our development environment on [Docker Hub][]. If you have docker set up, you can open it with `make run-dev-container`. It's not intended to handle all development workflows, but you should be able to successfully run `make test` in it, and use it to troubleshoot your host environment if necessary. (That container is built from the Dockerfile in the root of this repository.)

[Docker Hub]: https://cloud.docker.com/u/reserveprotocol/repository/docker/reserveprotocol/env

# Building and Testing

## Quickstart

With an environment set up as above,

- To build and run unit tests: `make test`
- To build smart contracts and their Go bindings: `make abi/bindings`
- To compute coverage:
    - In a separate long-lived terminal, `make run-devnet`
    - `make coverage` (You will get a lot of warnings saying "We'll just skip that trace", even things are working right.)
- To deploy and interact with the smart contracts
    - In a separate long-lived terminal, `make run-devnet`
    - `make poke`
    - `poke --help`
    - `$(poke deploy)`
    - Interact using poke, e.g.:
    ``` bash
    poke changeMinter @1
    poke mint --from @1 @2 123
    poke balanceOf @2
    ```

## Makefile

The root Makefile provides entry points for building and testing:

- `make fmt`: Use [ethlint][] to lint and format smart contracts.
- `make test`: Run tests for the smart contracts
- `make coverage`: Compute test coverage (where possible) for the smart contracts. Note that it has a few obvious false negatives, like the ReserveDollar constructor, interface definitions, and the `_;` line in modifiers. Needs a local `geth` node listening at `localhost:8545`; see `make run-devnet`.
- `make poke`: Build and install the `poke` CLI, for deploying and exercising with the Reserve Dollar smart contracts. `poke` needs a local `geth` node listening at `localhost:8545`; see `make run-devnet`.
- `make run-devnet`: Run a local ethereum node suitable for testing and coverage.
- `make run-dev-container`: Open dockerized development environment.

## More on `make run-devnet`
Some of our tools -- `poke` and `make coverage` -- expect to interact with a local [geth][go-ethereum] node at `http://localhost:8545`. Notably `make test` does _not_ require this; it uses a faster, in-memory EVM implementation.

The command `make run-devnet` sets up a local `geth` node specialized for testing. `make run-devnet` will run the [`0xorg/devnet` container][devnet] and have it listen on port 8545. This command produces lots of live output to stdout, which is frequently useful. We recommend either running it in its own terminal, or at least piping its output somewhere so that you can `tail -f` it.

[devnet]: https://github.com/0xProject/0x-monorepo/tree/development/packages/devnet
[go-ethereum]: https://github.com/ethereum/go-ethereum/wiki/geth
[ethlint]: https://www.npmjs.com/package/ethlint

## Directory Layout

### Highlights
- `contracts/`: The Reserve Dollar smart contracts
    - `ReserveDollar.sol`: The main token implementation contract.
    -  `ReserveDollarEternalStorage.sol`: Static implementation of the [Eternal Storage][] pattern for `ReserveDollar`
    - `MintAndBurnAdmin.sol`: Contract intended to hold the `minter` role for `ReserveDollar`.
    - `ReserveDollarV2.sol`: A tiny "upgraded" ReserveDollar contract, used here to test migration workflows.
    - `zepplin/SafeMath.sol`: The OpenZepplin SafeMath library.
- `tests/`: Unit tests for the Reserve Dollar smart contracts
- `cmd/poke/` A CLI for interacting with the Reserve Dollar smart contract.
- `soltools/`: Go-to-JavaScript bridge, to wrap 0x's solidity tools.

### The Rest
- `artifacts/`: Build destination for smart contracts
- `abi/`: Build directory Go bindings for the smart contracts get built.
- `README.md`: You're reading it.
- `Makefile`: Entry points for building and testing, as [above](#Makefile).
- `Dockerfile`: Dockerfile for the dockerized dev environment
- ... and a handful of other source-tree configuration files
