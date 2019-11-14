package context

//import (
//	"fmt"
//	"github.com/spf13/viper"
//	"github.com/tanhuiya/ci123chain/pkg/abci/types"
//	"github.com/tanhuiya/ci123chain/pkg/client/helper"
//)




//
//func (ctx *Context) Sign(msg []byte, addr types.AccAddress) ([]byte, error) {
//	passphrase, err := ctx.GetPassphrase(addr)
//	if err != nil {
//		return nil, err
//	}
//	ks := keystore.NewKeyStore(ctx.HomeDir, keystore.StandardScryptN, keystore.StandardScryptP)
//	acc := accounts.Account{
//		Address: addr.Address,
//	}
//	acct, err := ks.Find(acc)
//	if err != nil {
//		return nil, err
//	}
//	return ks.SignHashWithPassphrase(acct, passphrase, msg)
//}

