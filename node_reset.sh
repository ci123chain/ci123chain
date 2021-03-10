#! /bin/bash


rm -rf ~/.ci123

./build/cid init --home ~/.ci123 --chain_id ci0 --validator_key=oQLmM5pM5wL78a6LJntQY8tPGQPpp050udIA5YZMkCc=

./build/cid add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 100000000000 --home ~/.ci123

./build/cid add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000000000000000000000 --home ~/.ci123
#2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70

./build/cid add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 100000000000 --home ~/.ci123

./build/cid add-genesis-validator 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 8000000 AuZzwlxbJVV+Cc5MJFdr4M8330EA4MZ5fPApBiSO5vfe 1 40 5 --home ~/.ci123
