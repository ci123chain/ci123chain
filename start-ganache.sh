#!/bin/bash

mnemonic='explain tackle mirror kit van hammer degree position ginger unfair soup bonus'

rm -rf gravity-ganache
mkdir gravity-ganache
set -e;
npx ganache-cli -m "${mnemonic}" -i 15 -l 100000000 --db gravity-ganache --e 10000 -h 0.0.0.0

