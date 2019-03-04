# Smart Contracts for the Reserve Dollar

There are three main smart contracts in `contracts/`: `ReserveDollar`, `ReserveDollarEternalStorage`, and `MintAndBurnAdmin`.

* `ReserveDollar` is the main token implementation. It provides ERC-20 and administration interfaces.
* `ReserveDollarEternalStorage` is a static implementation of the [Eternal Storage][] pattern for `ReserveDollar`. Any non-constant-sized data for `ReserveDollar` is stored there. However, because the token's storage is relatively simple, and the dynamic form of the EternalStorage pattern introduces difficulties in analysis, `ReserveDollarEternalStorage` provides accessors for _specific_ maps -- `balance`, `allowed`, and `frozenTime`.
* `MintAndBurnAdmin` is intended to hold the `minter` role for `ReserveDollar`. We use this to give ourselves time to respond and recover in case our operational keys are stolen.

[Eternal Storage]: https://fravoll.github.io/solidity-patterns/eternal_storage.html

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



