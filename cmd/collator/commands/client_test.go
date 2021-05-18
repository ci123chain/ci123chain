package commands

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"github.com/ci123chain/ci123chain/pkg/collactor/helper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	file, _ := ioutil.ReadFile("./testdata/config.yaml")
	config = &Config{}
	// unmarshall them into the struct
	if err := yaml.Unmarshal(file, config); err != nil {
		panic("test config invalid")
	}
	validateConfig(config)

	m.Run()
}

func setupClients() (*collactor.Chain, *collactor.Chain) {
	c, src, dst, err := config.ChainsFromPath("demopath")

	_, err = c[src].GetAddress()
	if err != nil {
		panic(err)
	}

	_, err = helper.InitLight(c[dst], true)
	if err != nil {
		panic(err)
	}
	_, err = helper.InitLight(c[src], true)
	if err != nil {
		panic(err)
	}

	return c[src], c[dst]
}

func TestCreateClientState(t *testing.T)  {
	src, dst := setupClients()
	//c[src].PathEnd.ClientID = "07-tendermint-5"
	//res, err := c[src].QueryClientState(0)
	//fmt.Println(res)

	_, err := src.CreateClients(dst)
	require.Nil(t, err)
}

func TestUpdateClientState(t *testing.T)  {
	src, dst := setupClients()
	_, err := src.CreateClients(dst)
	require.Nil(t, err)
	err = src.UpdateClients(dst)
	require.Nil(t, err)
}

func TestCreateOpenConnection(t *testing.T)  {
	src, dst := setupClients()
	_, err := src.CreateClients(dst)
	require.Nil(t, err)
	_, err = src.CreateOpenConnections(dst, 5, 5 * time.Second)
	require.Nil(t, err)

	con , err := src.QueryConnection(0)
	require.Equal(t, con.Connection.ClientId, src.PathEnd.ClientID)
}

func TestCreateOpenChannel(t *testing.T)  {
	src, dst := setupClients()
	_, err := src.CreateClients(dst)
	require.Nil(t, err)

	_, err = src.CreateOpenConnections(dst, 5, 3 * time.Second)
	require.Nil(t, err)

	_, err = src.CreateOpenChannels(dst, 5, 3 * time.Second)
	require.Nil(t, err)
}

func TestStart(t *testing.T)  {
	src, dst := setupClients()
	_, err := src.CreateClients(dst)
	require.Nil(t, err)
	_, err = src.CreateOpenConnections(dst, 5, 3 * time.Second)
	require.Nil(t, err)
	_, err = src.CreateOpenChannels(dst, 5, 3 * time.Second)
	require.Nil(t, err)

	strategy := &collactor.NaiveStrategy{}
	strategy.MaxTxSize = 2 * MB // in MB
	strategy.MaxMsgLength = 5

	_, err = collactor.RunStrategy(src, dst, strategy)
	require.Nil(t, err, "err should be nil")
	select {}
}

func TestSendPacket(t *testing.T)  {
	src, dst := setupClients()

	_, err := src.CreateClients(dst)
	require.Nil(t, err)

	_, err = src.CreateOpenConnections(dst, 5, 3 * time.Second)
	require.Nil(t, err)

	_, err = src.CreateOpenChannels(dst, 5, 3 * time.Second)
	require.Nil(t, err)

	amount := sdk.NewCoin("stack0", sdk.NewInt(10000))

	dstAddr := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	err = src.SendTransferMsg(dst, amount, dstAddr, 10, 0)
	require.Nil(t, err, "err should be nil")

}

func TestQueryCommitPacket(t *testing.T)  {
	src, dst := setupClients()
	_, err := src.CreateClients(dst)
	require.Nil(t, err)
	_, err = src.CreateOpenConnections(dst, 5, 300)
	require.Nil(t, err)
	_, err = src.CreateOpenChannels(dst, 5, 300)
	require.Nil(t, err)

	res, err := src.QueryPacketCommitments(0, 1000, 0)
	require.Nil(t, err)
	require.Equal(t, 1, len(res.Commitments))
}