cp -r node0-template node0
cp -r node1-template node1
cp -r node2-template node2
mkdir gateway
docker-compose -f part1.yaml up -d
sleep 15
docker-compose -f part2.yaml up -d