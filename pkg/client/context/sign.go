package context

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/viper"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
)

func (ctx *Context) GetPassphrase(addr types.AccAddress) (string, error) {
	pass := viper.GetString(helper.FlagPassword)
	if pass == "" {
		return ctx.getPassphraseFromStdin(addr)
	}
	return pass, nil
}

// Get passphrase from std input
func (ctx *Context) getPassphraseFromStdin(addr types.AccAddress) (string, error) {
	buf := helper.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", addr.Hex())
	return helper.GetPassword(prompt, buf)
}



func (ctx *Context) Sign(msg []byte, addr types.AccAddress) ([]byte, error) {
	passphrase, err := ctx.GetPassphrase(addr)
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(ctx.HomeDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc := accounts.Account{
		Address: addr.Address,
	}
	acct, err := ks.Find(acc)
	if err != nil {
		return nil, err
	}
	return ks.SignHashWithPassphrase(acct, passphrase, msg)
}

