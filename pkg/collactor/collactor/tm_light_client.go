package collactor

import (
	"context"
	"fmt"
	retry "github.com/avast/retry-go"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/light"
	lightp "github.com/tendermint/tendermint/light/provider"
	lighthttp "github.com/tendermint/tendermint/light/provider/http"
	dbs "github.com/tendermint/tendermint/light/store/db"
	dbm "github.com/tendermint/tm-db"
	"io/ioutil"
)
var logger = light.Logger(log.NewTMLogger(log.NewSyncWriter(ioutil.Discard)))

// NewLightDB returns a new instance of the lightclient database connection
// CONTRACT: must close the database connection when done with it (defer df())
func (c *Chain) NewLightDB() (db *dbm.GoLevelDB, df func(), err error) {
	if err := retry.Do(func() error {
		db, err = dbm.NewGoLevelDB(c.ChainID, lightDir(c.HomePath))
		if err != nil {
			return fmt.Errorf("can't open light client database: %w", err)
		}
		return nil
	}, rtyAtt, rtyDel, rtyErr); err != nil {
		return nil, nil, err
	}

	df = func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}

	return
}

// LightHTTP returns the http client for light clients
func (c *Chain) LightHTTP() lightp.Provider {
	cl, err := lighthttp.New(c.ChainID, c.RPCAddr)
	if err != nil {
		panic(err)
	}
	return cl
}


// LightClientWithoutTrust querys the latest header from the chain and initializes a new light client
// database using that header. This should only be called when first initializing the light client
func (c *Chain) LightClientWithoutTrust(db dbm.DB) (*light.Client, error) {
	var (
		height int64
		err    error
	)
	prov := c.LightHTTP()

	if err := retry.Do(func() error {
		height, err = c.QueryLatestHeight()
		switch {
		case err != nil:
			return err
		case height == 0:
			return fmt.Errorf("shouldn't be here")
		default:
			return nil
		}
	}, rtyAtt, rtyDel, rtyErr); err != nil {
		return nil, err
	}

	lb, err := prov.LightBlock(context.Background(), height)
	if err != nil {
		return nil, err
	}
	return light.NewClient(
		context.Background(),
		c.ChainID,
		light.TrustOptions{
			Period: c.GetTrustingPeriod(),
			Height: height,
			Hash:   lb.SignedHeader.Hash(),
		},
		prov,
		// TODO: provide actual witnesses!
		// NOTE: This requires adding them to the chain config
		[]lightp.Provider{prov},
		dbs.New(db, ""),
		logger)
}
