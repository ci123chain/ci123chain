../build/cid boot-gen --chain-pre=node --node-num=3 --output-dir=.

mkdir gateway
docker-compose -f part1.yaml up -d
sleep 15
docker-compose -f part22.yaml up -d