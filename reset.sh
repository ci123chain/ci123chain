#! /bin/bash

rm -rf ~/.shard1*
rm -rf ~/.shard2*

./build/cid init --home ~/.shard1 --chain-id Shard1

./build/cid add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 100000 --home ~/.shard1

./build/cid add-genesis-account 0x505A74675dc9C71eF3CB5DF309256952917E801e 100000 --home ~/.shard1
./build/cid add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 100000 --home ~/.shard1
#2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70

./build/cid add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 100000 --home ~/.shard1

./build/cid init --home ~/.shard2 --chain-id asdjqj