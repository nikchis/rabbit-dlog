// Copyright (c) 2020 Nikita Chisnikov <chisnikov@gmail.com>
// Distributed under the MIT/X11 software license

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

const (
	appFilesDir = `.rabbit-dlog`
	cfgFilename = `rabbit-dlog.conf`
)

type config struct {
	Url      string `json:"url"`
	Exchange string `json:"exchange"`
	QueInfo  string `json:"que_info"`
	QueWarn  string `json:"que_warn"`
	QueError string `json:"que_error"`
	KeyInfo  string `json:"key_info"`
	KeyWarn  string `json:"key_warn"`
	KeyError string `json:"key_error"`
	ConInfo  string `json:"con_info"`
	ConWarn  string `json:"con_warn"`
	ConError string `json:"con_error"`
	Id       string `json:"id"`
	Token    string `json:"token"`
	Filepath string `json:"-"`
}

func newConfig(appFilesDirPath string) (cfg *config, err error) {
	var usr *user.User
	cfg = &config{}
	if appFilesDirPath == "" {
		if usr, err = user.Current(); err != nil {
			return
		}
		appFilesDirPath = fmt.Sprintf("%s/%s", usr.HomeDir, appFilesDir)
	}
	if err = cfg.read(appFilesDirPath); err != nil {
		cfg.init()
		err = cfg.create(appFilesDirPath)
	}
	return
}

func (cfg *config) create(appFilesDirPath string) (err error) {
	var buf []byte
	_, err = os.Stat(appFilesDirPath)
	if os.IsNotExist(err) {
		if err = os.Mkdir(appFilesDirPath, 0750); err != nil {
			return
		}
	} else if err != nil {
		return
	}
	if buf, err = json.MarshalIndent(cfg, "", "	"); err != nil {
		return
	}
	cfg.Filepath = fmt.Sprintf("%s/%s", appFilesDirPath, cfgFilename)
	if err = ioutil.WriteFile(cfg.Filepath, buf, 0600); err != nil {
		return
	}
	return
}

func (cfg *config) init() {
	cfg.Url = "amqp://user:password@host:5672/vhost"
	cfg.Exchange = "bv.ex.logs"
	cfg.QueInfo = "q.logs.discord.info"
	cfg.QueWarn = "q.logs.discord.warn"
	cfg.QueError = "q.logs.discord.error"
	cfg.KeyInfo = "k.logs.discord.info"
	cfg.KeyWarn = "k.logs.discord.warn"
	cfg.KeyError = "k.logs.discord.error"
	cfg.ConInfo = "c.logs.discord.info"
	cfg.ConWarn = "c.logs.discord.warn"
	cfg.ConError = "c.logs.discord.error"
	cfg.Id = "discord-webhook-id"
	cfg.Token = "discord-webhook-token"
}

func (cfg *config) read(appFilesDirPath string) (err error) {
	var buf []byte
	fpath := fmt.Sprintf("%s/%s", appFilesDirPath, cfgFilename)
	if buf, err = ioutil.ReadFile(fpath); err != nil {
		return err
	}
	if err = json.Unmarshal(buf, &cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *config) String() (result string) {
	if b, err := json.MarshalIndent(cfg, "", "	"); err == nil {
		result = string(b)
	}
	return result
}
