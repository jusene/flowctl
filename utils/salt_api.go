package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

type SaltApi struct {
	saltUser string
	saltPass string
	saltUrl  string
	client   *http.Client
}

func NewSaltApi() *SaltApi {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
	}

	return &SaltApi{
		saltUrl: "https://salt-api.hho-inc.com",
		saltUser: "saltapi",
		saltPass: "saltapi",
		client: client,
	}
}

func (s *SaltApi) getToken() (string, error) {
	jsonData := map[string]string{
		    "eauth": "pam",
			"username": s.saltUser,
			"password": s.saltPass,
	}

	jsonBody, err:= json.Marshal(jsonData)
	cobra.CheckErr(err)

	request, err:= http.NewRequest("POST", strings.Join([]string{s.saltUrl, "login"}, "/"), bytes.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Accept", "application/json")
	cobra.CheckErr(err)

	resp, err := s.client.Do(request)
	cobra.CheckErr(err)

	data, err:= ioutil.ReadAll(resp.Body)
	cobra.CheckErr(err)

	v := viper.New()
	v.SetConfigType("json")
	v.ReadConfig(bytes.NewReader(data))
	token := v.GetString("return.0.token")
	return token, nil
}

func (s *SaltApi) RemoteExec(tgt, arg string) (map[string]interface{}, error) {
	jsonData := map[string]string{
		"client": "local",
		"tgt": tgt,
		"fun": "cmd.run",
		"arg": arg,
	}

	jsonBody, err:= json.Marshal(jsonData)
	cobra.CheckErr(err)
	token, _ := s.getToken()
	request, err := http.NewRequest("POST", s.saltUrl, bytes.NewReader(jsonBody))
	cobra.CheckErr(err)
	request.Header.Set("X-Auth-Token", token)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(request)
	cobra.CheckErr(err)

	data, err:= ioutil.ReadAll(resp.Body)
	cobra.CheckErr(err)

	v := viper.New()
	v.SetConfigType("json")
	v.ReadConfig(bytes.NewReader(data))

	return v.AllSettings(), nil
}