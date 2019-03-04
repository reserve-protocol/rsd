
fuzz: fuzz-total-supply fuzz-ownership

fuzz-total-supply: setup
	docker run --rm -t -v `pwd`:/src -w /src reserveprotocol/echidna echidna-test EchidnaTransferTests.sol EchidnaTotalSupplyTests --config multisender.yaml

fuzz-ownership: setup
	docker run --rm -t -v `pwd`:/src -w /src reserveprotocol/echidna echidna-test EchidnaOwnershipTests.sol EchidnaConstantOwnerSimple --config multisender.yaml



setup: ReserveDollar.sol \
       ReserveDollarEternalStorage.sol \
       zeppelin/SafeMath.sol

ReserveDollar.sol: ../contracts/ReserveDollar.sol
	ln $< $@

ReserveDollarEternalStorage.sol: ../contracts/ReserveDollarEternalStorage.sol
	echo "\033[31;1m[WARNING] Removing onlyOwner modifier from eternal storage.\nechidna\033[0m"
	sed -e 's|require(msg.sender == owner, "onlyOwner");||' $< > $@

zeppelin/SafeMath.sol: ../contracts/zeppelin/SafeMath.sol
	mkdir -p zeppelin
	ln $< $@

clean:
	rm ReserveDollar.sol \
       ReserveDollarEternalStorage.sol \
       zeppelin/SafeMath.sol

.PHONY: fuzz fuzz-total-supply fuzz-ownership
