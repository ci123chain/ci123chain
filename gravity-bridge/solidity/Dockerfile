FROM node:15.11-alpine3.13

RUN apk update
RUN apk add --no-cache python3 make g++ curl

COPY . .
RUN npm ci
RUN chmod -R +x scripts

RUN npm run typechain

CMD npx ts-node \
    contract-deployer.ts \
    --cosmos-node="http://host.docker.internal:26657" \
    --eth-node="http://host.docker.internal:8545" \
    --eth-privkey="0xa8a54b2d8197bc0b19bb8a084031be71835580a01e70a45a13babd16c9bc1563" \
    --contract=artifacts/contracts/Gravity.sol/Gravity.json \
    --test-mode=true