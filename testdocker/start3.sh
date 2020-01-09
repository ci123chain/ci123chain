cp -r node0-template node0
docker-compose -f part1.yaml up -d
sleep 10
docker-compose -f part2.yaml up -d
sleep 8
docker-compose -f part3.yaml up -d