package collactor

import (
	"context"
	"errors"
	"fmt"
	retry "github.com/avast/retry-go"
	tmclient "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/light"
	lightp "github.com/tendermint/tendermint/light/provider"
	lighthttp "github.com/tendermint/tendermint/light/provider/http"
	dbs "github.com/tendermint/tendermint/light/store/db"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"time"
)

func lightError(err error) error { return fmt.Errorf("light client: %w", err) }


var (
	logger = light.Logger(log.NewTMLogger(log.NewSyncWriter(ioutil.Discard)))

	ErrDatabase = errors.New("database failure")
)

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
		// NOTE: This requires adding them to the chain configs
		[]lightp.Provider{prov},
		dbs.New(db, ""),
		logger)
}

// ValidateLightInitialized returns an error if the light client isn't initialized or there is a problem
// interacting with the light client.
func (c *Chain) ValidateLightInitialized() error {
	height, err := c.GetLatestLightHeight()
	if err != nil {
		return fmt.Errorf("encountered issue with off chain light client for chain (%s): %v", c.ChainID, err)
	}

	// height will return -1 when the client has not been initialized
	if height == -1 {
		return fmt.Errorf("please initialize an off chain light client for chain (%s)", c.ChainID)
	}

	return nil
}

// GetLatestLightHeight returns the latest height of the light client.
func (c *Chain) GetLatestLightHeight() (int64, error) {
	db, df, err := c.NewLightDB()
	if err != nil {
		return -1, err
	}
	defer df()

	client, err := c.LightClient(db)
	if err != nil {
		return -1, err
	}

	return client.LastTrustedHeight()
}

// LightClient initializes the light client for a given chain from the trusted store in the database
// this should be call for all other light client usage
func (c *Chain) LightClient(db dbm.DB) (*light.Client, error) {
	prov := c.LightHTTP()
	return light.NewClientFromTrustedStore(
		c.ChainID,
		c.GetTrustingPeriod(),
		prov,
		// TODO: provide actual witnesses!
		// NOTE: This requires adding them to the chain configs
		[]lightp.Provider{prov},
		dbs.New(db, ""),
		logger,
		light.PruningSize(0),
	)
}



// UpdateLightClient updates the tendermint light client by verifying the current
// header against a trusted header.
func (c *Chain) UpdateLightClient() (*tmtypes.LightBlock, error) {
	// create database connection
	db, df, err := c.NewLightDB()
	if err != nil {
		return nil, lightError(err)
	}
	defer df()

	client, err := c.LightClient(db)
	if err != nil {
		return nil, lightError(err)
	}

	lightBlock, err := client.Update(context.Background(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to update off-chain light client for chain %s: %w", c.ChainID, err)
	}

	// new clients, cannot be updated without trusted starting state
	if lightBlock == nil {
		lightBlock, err = client.TrustedLightBlock(0)
		if err != nil {
			return nil, lightError(err)
		}
	}

	return lightBlock, nil
}


// UpdateLightClients updates the off-chain tendermint light clients concurrently.
func UpdateLightClients(src, dst *Chain) (srcLB, dstLB *tmtypes.LightBlock, err error) {
	var eg = new(errgroup.Group)
	eg.Go(func() error {
		srcLB, err = src.UpdateLightClient()
		return err
	})
	eg.Go(func() error {
		dstLB, err = dst.UpdateLightClient()
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	return srcLB, dstLB, nil
}

// GetLightSignedHeaderAtHeight returns a signed header at a particular height (0 - the latest).
func (c *Chain) GetLightSignedHeaderAtHeight(height int64) (*tmclient.Header, error) {
	// create database connection
	db, df, err := c.NewLightDB()
	if err != nil {
		return nil, err
	}
	defer df()

	client, err := c.LightClient(db)
	if err != nil {
		return nil, err
	}

	if height == 0 {
		height, err = client.LastTrustedHeight()
		if err != nil {
			return nil, err
		}
	}

	// VerifyLightBlock will return the header at provided height if it already exists in store,
	// otherwise it retrieves from primary and verifies against trusted store before returning.
	sh, err := client.VerifyLightBlockAtHeight(context.Background(), height, time.Now())
	if err != nil {
		return nil, err
	}

	protoVal := tmtypes.NewValidatorSet(sh.ValidatorSet.Validators)
	if err != nil {
		return nil, err
	}

	return &tmclient.Header{SignedHeader: sh.SignedHeader, ValidatorSet: protoVal}, nil
}

// GetLatestLightHeights returns both the src and dst latest height in the local client
func GetLatestLightHeights(src, dst *Chain) (srch int64, dsth int64, err error) {
	var eg = new(errgroup.Group)
	eg.Go(func() error {
		srch, err = src.GetLatestLightHeight()
		return err
	})
	eg.Go(func() error {
		dsth, err = dst.GetLatestLightHeight()
		return err
	})
	if err = eg.Wait(); err != nil {
		return
	}
	return
}


// MustGetLatestLightHeight returns the latest height of the light client. If
// an error occurs due to a database failure, we keep trying with a delayed
// re-attempt. Otherwise, we panic.
func (c *Chain) MustGetLatestLightHeight() uint64 {
	height, err := c.GetLatestLightHeight()
	if err != nil {
		if errors.Is(err, ErrDatabase) {
			// XXX: Sleep and try again if the database is unavailable. This can easily
			// happen if two distinct resources try to access the database at the same
			// time. To avoid causing a corrupted or lost packet, we keep trying as to
			// not halt the relayer.
			//
			// ref: https://github.com/cosmos/relayer/issues/444
			c.logger.Error("failed to get latest height due to a database failure; trying again...", "err", err)
			time.Sleep(time.Second)
			c.MustGetLatestLightHeight()
		} else {
			panic(err)
		}
	}

	return uint64(height)
}