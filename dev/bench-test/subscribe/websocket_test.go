package subscribe

import (
	"context"
	"fmt"
	"testing"
	"time"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

func TestSubscribe(t *testing.T)  {

	client := rpchttp.NewHTTP("tcp://localhost:26657", "/websocket")
	err := client.Start()
	if err != nil {
		panic(err)
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	//query := "tm.event = 'Tx' AND tx.height = 3"
	query := "write_db.msg = 'address1'"

	txs, err := client.Subscribe(ctx, "test-clients", query)
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range txs {
			for _, eve := range e.Data.(types.EventDataTx).Result.Events {
				fmt.Println("got event ", eve)
			}
		}
	}()
	select {
	}
}