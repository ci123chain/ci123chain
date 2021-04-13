package commands

import (
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"github.com/ci123chain/ci123chain/pkg/collactor/helper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
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
	c, src, dst, err := config.ChainsFromPath("demo")

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
	src.UpdateClients(dst)
}

