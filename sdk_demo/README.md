#SDK使用与API请求介绍

##Transfer-普通转账交易
```
1.在线交易
a.调用sdk API接口实现在线转账交易
sdk.HttpTransferTx(from, to, gas, nonce, amount, priv, proxy)

2.离线交易
线下签名，将签名后的交易通过sdk的同步或异步API发送出去。调用SDK的签名方法签名交易
a.转账给同一条链的其他账户：
sdk.signTransferMsg(from,to,amount,gas,nonce,priv,isFabric)
sdk.SendTransaction(tx, async, isIBC)
```

##Shard-分片交易
```
1.在线交易
a.添加分片交易
sdk.HttpAddShardTx(from, gas, nonce, t, name, height, priv, proxy)

2.离线交易
线下签名，将签名后的交易通过sdk的同步或异步API发送出去。调用SDK的签名方法签名交易
a.添加分片交易：
sdk.signTransferMsg(from,gas,nonce,type,name,height,priv)
sdk.SendTransaction(tx, async, isIBC)
```

##IBC-跨链交易
```
目前IBC跨链交易(假设A链-> B链)只支持离线交易，线下签名。通过同步或异步交易API发出。
跨链交易有多个步骤，每个步骤的交易类别不同：
1.发起IBC跨链交易（from,privateKey必须匹配，这个交易由A链上用户发起）
sdk.SignIBCTransferMsg(from, to, amount, gas, nonce, privateKey)
sdk.SendTransaction(tx, async, isIBC)

2.申请处理IBC跨链交易（from,privateKey必须匹配，这个交易由observer发起）
sdk.SignApplyIBCMsg(from, uniqueID, observerID, gas, nonce, privateKey)
sdk.SendTransaction(tx, async, isIBC)

3.发起扣款交易（A链上用户将交易款项暂时转账到临时生成的账户，from,privateKey必须匹配，
这个交易由observer发起）
sdk.SignIBCBankSendMsg(from, raw, gas, nonce, privateKey)
sdk.SendTransaction(tx, async, isIBC)

4.接收回执交易，最终转账（B链转账成功，A链从临时账户将交易款项转账到银行账户，
from,privateKey必须匹配，这个交易由observer发起）
sdk.SignIBCReceiptMsg(from, raw, gas, nonce, privateKey)
sdk.SendTransaction(tx, async, isIBC)
```