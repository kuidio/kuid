/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/henderiw/apiserver-store/pkg/db/badgerdb"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/pkg/registry/options"
)

var (
	configDir      = "/config"
	configFilePath = "/etc/kuid/config.json"
)

type StorageType string

const (
	StorageType_Memory   StorageType = "memory"
	StorageType_File     StorageType = "file"
	StorageType_Git      StorageType = "git"
	StorageType_Badgerdb StorageType = "badgerdb"
	StorageType_Etcd     StorageType = "etcd"
	StorageType_PostGres StorageType = "postgres"
)

type KuidGroupConfig struct {
	Group   string `json:"group"`
	Enabled bool   `json:"enabled"`
	Sync    bool   `json:"sync"` // Sync or Async -> only possible with
}

type KuidConfig struct {
	Storage StorageType        `json:"storage"`
	Groups  []*KuidGroupConfig `json:"groups"`
}

func GetKuidConfig() (*KuidConfig, error) {
	if !isFile(configFilePath) {
		return getDefaultConfig(), nil
	}
	b, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	cfg := &KuidConfig{}
	if err := json.Unmarshal([]byte(b), cfg); err != nil {
		return nil, err
	}
	// TODO need to add some validation
	// sync with etcd is not possible
	return cfg, nil
}

func getDefaultConfig() *KuidConfig {
	return &KuidConfig{
		Storage: StorageType_Badgerdb,
		Groups: []*KuidGroupConfig{
			{Group: "infra.kuid.dev", Enabled: true, Sync: true},
			{Group: "as.be.kuid.dev", Enabled: true, Sync: true},
			{Group: "vlan.be.kuid.dev", Enabled: true, Sync: true},
			{Group: "ipam.be.kuid.dev", Enabled: true, Sync: true},
			{Group: "genid.be.kuid.dev", Enabled: true, Sync: true},
			{Group: "extcomm.be.kuid.dev", Enabled: true, Sync: true},
		},
	}
}

func GetRegistryOptions(ctx context.Context, typ StorageType) (*options.Options, error) {
	log := log.FromContext(ctx)
	switch typ {
	case StorageType_Badgerdb:
		db, err := badgerdb.OpenDB(ctx, configDir)
		if err != nil {
			log.Error("cannot open db", "err", err.Error())
			return nil, err
		}

		return &options.Options{
			Prefix: configDir,
			Type:   options.StorageType_KV,
			DB:     db,
		}, nil
	case StorageType_Etcd:
		return nil, nil
	default:
		return &options.Options{
			Prefix: configDir,
			Type:   options.StorageType_Memory,
		}, nil
	}
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false // File does not exist
		}
		return false
	}
	return info.Mode().IsRegular() // Returns true if it is a regular file
}
