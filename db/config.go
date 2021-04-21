/*
Copyright © 2020 Marvin

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
package db

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Cfg struct {
	OracleDB OracleDB `toml:"oracle" json:"oracle"`
}

type OracleDB struct {
	Username      string   `toml:"username" json:"username"`
	Password      string   `toml:"password" json:"password"`
	ConnectString string   `toml:"connect-string",json:"connect-string"`
	SessionParams []string `toml:"session-params" json:"session-params"`
	Timezone      string   `toml:"timezone" json:"timezone"`
}

func ReadConfigFile(file string) (*Cfg, error) {
	cfg := &Cfg{}
	if err := cfg.configFromFile(file); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c *Cfg) configFromFile(file string) error {
	if _, err := toml.DecodeFile(file, c); err != nil {
		return fmt.Errorf("failed decode toml config file %s: %v", file, err)
	}
	return nil
}

func (c *Cfg) String() string {
	cfg, err := json.Marshal(c)
	if err != nil {
		return "<nil>"
	}
	return string(cfg)
}
