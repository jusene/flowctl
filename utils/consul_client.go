package utils

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"gitlab.hho-inc.com/devops/flowctl/models"
	"io"
	"os"
	"text/template"
)

type consul struct {
	client *api.Client
}

func NewConsul() *consul {
	address := "consul.hho-inc.com"
	port := "80"
	conf := api.DefaultConfig()
	conf.Address = address + ":" + port

	client, err := api.NewClient(conf)
	cobra.CheckErr(err)

	return &consul{
		client,
	}
}

func (c *consul) GetKV(key string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("未找到%s, 请检测consul", key)
			os.Exit(1)
		}
	}()
	KVPair, _, err := c.client.KV().Get(key, &api.QueryOptions{})
	return KVPair.Value, err
}

func (c *consul) Render2file(key string, tgt io.Writer, attr *models.AppInfo) {
	temp, err := c.GetKV(key)
	cobra.CheckErr(err)

	if t, err := template.New(attr.APP).Parse(string(temp)); err != nil {
		cobra.CheckErr(err)
	} else {
		t.Execute(tgt, attr)
	}
}
