DOCKER?=docker
DOCKER_RUN_ECHIDNA=$(DOCKER) run --rm -t -v `pwd`/echidna:/src -w /src reserveprotocol/echidna echidna-test

fuzz: fuzz-total-supply \
	fuzz-ownership \
	fuzz-constant-owner-under-free-money \
	fuzz-freezing


fuzz-total-supply: setup
	$(DOCKER_RUN_ECHIDNA) EchidnaTransferTests.sol EchidnaTotalSupplyTests --config multisender.yaml

fuzz-ownership: setup
	$(DOCKER_RUN_ECHIDNA) EchidnaOwnershipTests.sol EchidnaConstantOwnerSimple --config multisender_with_owner.yaml


fuzz-constant-owner-under-free-money: setup 
	$(DOCKER_RUN_ECHIDNA) EchidnaOwnershipTests.sol EchidnaConstantOwnerFreeMoney --config multisender.yaml

fuzz-freezing: setup
	$(DOCKER_RUN_ECHIDNA) EchidnaFreezeTests.sol EchidnaFrozenAccount --config multisender.yaml

setup: echidna/ReserveDollar.sol \
       echidna/ReserveDollarEternalStorage.sol \
       echidna/zeppelin/SafeMath.sol

echidna/ReserveDollar.sol: ../contracts/ReserveDollar.sol
	rm $@ || true
	ln $< $@

echidna/ReserveDollarEternalStorage.sol: ../contracts/ReserveDollarEternalStorage.sol
	echo "\033[31;1m[WARNING] Removing onlyOwner modifier from eternal storage.\nechidna\033[0m"
	sed -e 's|require(msg.sender == owner, "onlyOwner");||' $< > $@

echidna/zeppelin/SafeMath.sol: ../contracts/zeppelin/SafeMath.sol
	mkdir -p echidna/zeppelin
	rm $@ || true
	ln $< $@

clean:
	rm echidna/ReserveDollar.sol \
       echidna/ReserveDollarEternalStorage.sol \
       echidna/zeppelin/SafeMath.sol

.PHONY: fuzz fuzz-total-supply fuzz-ownership fuzz-constant-owner-under-free-money \
	fuzz-freezing
