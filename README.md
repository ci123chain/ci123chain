
# 开始使用

### 1.编译
`make build`

### 2.初始化数据
`./build/cid init --home ~/.ci123 --chain_id ci0 --validator_key=oQLmM5pM5wL78a6LJntQY8tPGQPpp050udIA5YZMkCc=`


### 3. 添加初始账户

创建本地账户 并输入密码

`./build/cli new`

将输出的账户地址保存， 如 `0x5D0506ac0411b7E32b2526e2014794Ac418518AB`

将该账户作为创世账户
`./build/cid add-genesis-account 0x5D0506ac0411b7E32b2526e2014794Ac418518AB 100000000000 --home ~/.ci123`

### 4. 启动状态数据库 couchdb
`docker-compose -f testdocker/couchdb-single.yaml up -d`

### 5. 启动节点

couchdb 地址填上面启动的地址

`CI_STATEDB=couchdb://admin:password@#couchdb地址#/ci1 ./build/cid start`

### 6. 转账
新建账户B => `0x6207826Ee35e69e7aAFD3C4049f6863Cd91dEd1b `

`./build/cli transfer --address=0x5D0506ac0411b7E32b2526e2014794Ac418518AB --to=0x6207826Ee35e69e7aAFD3C4049f6863Cd91dEd1b --amount=5000000 --gas=10000`

### 7. 查询余额
`./build/cli balance --address=0x5D0506ac0411b7E32b2526e2014794Ac418518AB`

### 删除数据
`docker-compose -f testdocker/couchdb-single.yaml down`
`rm -rf ~/.ci123`



