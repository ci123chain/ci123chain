#! /bin/bash


rm -rf ~/init

./build/cid init --home ~/init --chain_id ci0 --validator_key=4wttMiieaewLiRYu+y05j0uslBDOX5IA3k4TY9GtQzSdTcXyd5Y982Q3CUdh+h1XcCvtpIUb+5q6rtJ8W4SEFw==

./build/cid add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 100000000000 --home ~/init

./build/cid add-genesis-account 0x505A74675dc9C71eF3CB5DF309256952917E801e 100000000000 --home ~/init
./build/cid add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000000000000000000 --home ~/init
#2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70
./build/cid add-genesis-account 0xb4124cEB3451635DAcedd11767f004d8a28c6eE7 10000000000000000000000000 --home ~/init
#a8a54b2d8197bc0b19bb8a084031be71835580a01e70a45a13babd16c9bc1563

./build/cid add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 100000000000 --home ~/init

./build/cid add-genesis-validator 0xb4124cEB3451635DAcedd11767f004d8a28c6eE7 8000000 nU3F8neWPfNkNwlHYfodV3Ar7aSFG/uauq7SfFuEhBc= 1 50 4 --home ~/init