res
---

A CLI for interacting with the Reserve Dollar smart contract.

This is designed for testing purposes. The goal is to make it easier to run small experiments
on the Reserve Dollar from the command line, without needing to write any code.

When we deploy the Reserve Dollar for real, we will use similar code, but it will go through
a hardware wallet.

The CLI includes access to all of the public functions on the Reserve Dollar.

The CLI is written assuming that it is being run against a local Ethereum node, available
on http://localhost:8545, with the same-prefunded accounts as the 0xorg/devnet docker image.
To run the 0xorg/devnet docker image, use the command:

    docker run -it --rm -p 8545:8501 0xorg/devnet

To deploy a new copy of the Reserve Dollar locally, run:

    $(res deploy)

Running this command inside `$(...)` will cause your shell to execute the output of the
command, which will set the appropriate environment variable for your next calls to `res`
to run against the contract you just deployed.

To see the owner of the contract you just deployed, run:

    res owner

This should show `0x5409ED021D9299bf6814279A6A1411A7e866A631`, the 0th pre-funded account
in the 0xorg/devnet image. You can check that with `res address`:

    res address @0

Anywhere you need to supply an address or a private key to the res tool, you can use
the special strings `@0` - `@9` to get the corresponding address or key from the ten
pre-funded accounts in the 0xorg/devnet image.

For paid mutator calls, `res` will default to using account `@0`. To override this default
per-command, you can use the `-F` (aka `--from`) flag, like so:

	res --from @1 transfer @2 200.5

You can also switch the default for the remainder of the current terminal session with
`res account`:

	$(res account @3)
