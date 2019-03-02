Reserve Dollar
---

TODO: overview of the Reserve Dollar

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
