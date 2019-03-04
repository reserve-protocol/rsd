
fuzz: fuzz-total-supply fuzz-ownership

fuzz-total-supply: setup
	docker run --rm -t -v `pwd`/echidna:/src -w /src reserveprotocol/echidna echidna-test EchidnaTransferTests.sol EchidnaTotalSupplyTests --config multisender.yaml

fuzz-ownership: setup
	docker run --rm -t -v `pwd`/echidna:/src -w /src reserveprotocol/echidna echidna-test EchidnaOwnershipTests.sol EchidnaConstantOwnerSimple --config multisender.yaml



setup: echidna/ReserveDollar.sol \
       echidna/ReserveDollarEternalStorage.sol \
       echidna/zeppelin/SafeMath.sol

echidna/ReserveDollar.sol: ../contracts/ReserveDollar.sol
	cp ../contracts/ReserveDollar.sol echidna/

echidna/ReserveDollarEternalStorage.sol: ../contracts/ReserveDollarEternalStorage.sol
	echo "\033[31;1m[WARNING] Removing onlyOwner modifier from eternal storage.\nechidna\033[0m"
	sed -e 's|require(msg.sender == owner, "onlyOwner");||' $< > $@

echidna/zeppelin/SafeMath.sol: ../contracts/zeppelin/SafeMath.sol
	mkdir -p echidna/zeppelin
	cp ../contracts/zeppelin/SafeMath.sol echidna/zeppelin/

clean:
	rm echidna/ReserveDollar.sol \
       echidna/ReserveDollarEternalStorage.sol \
       echidna/zeppelin/SafeMath.sol

.PHONY: fuzz fuzz-total-supply fuzz-ownership
