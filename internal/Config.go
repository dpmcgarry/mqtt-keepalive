/*
Copyright © 2024 Don P. McGarry

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package internal

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type MQTTDestination struct {
	Host   string
	Topics []string
}

type GlobalConfig struct {
	Interval          int
	PublishTimeout    int
	DisconnectTimeout int
}

func LoadServerConfig() ([]MQTTDestination, error) {
	var destinations []MQTTDestination
	if !viper.IsSet("servers") {
		log.Error().Msg("No Servers Configured")
		return nil, errors.New("no servers set in viper config")
	}
	genericConf := viper.Get("servers")
	// Type assertion to convert 'any' to 'map[string]interface{}'
	serverConf, ok := genericConf.(map[string]interface{})
	if !ok {
		log.Error().Msg("Conversion failed: Server config is not formatted correctly")
		return nil, errors.New("server configuration formatting invalid")
	}

	for k, _ := range serverConf {
		dest := MQTTDestination{}
		log.Debug().Msgf("Server: %v", k)
		dest.Host = k
		if !viper.IsSet("servers." + k + ".topics") {
			log.Error().Msgf("No Topics Configured for host %v", k)
			return nil, fmt.Errorf("no topics set for host %v", k)
		}
		topics := viper.GetStringSlice("servers." + k + ".topics")
		for _, topic := range topics {
			log.Debug().Msgf("Topic %v", topic)
			dest.Topics = append(dest.Topics, topic)
		}
		destinations = append(destinations, dest)
	}
	return destinations, nil
}

func LoadGlobalConfig() (GlobalConfig, error) {
	globalConf := GlobalConfig{}
	if !viper.IsSet("interval") {
		log.Error().Msg("Interval not configured")
		return GlobalConfig{}, errors.New("interval not set")
	}
	globalConf.Interval = viper.GetInt("interval")
	if !(globalConf.Interval > 0) {
		log.Error().Msgf("Interval set to invalid value: %v", globalConf.Interval)
		return GlobalConfig{}, fmt.Errorf("interval set to invalid value %v", globalConf.Interval)
	}
	log.Debug().Msgf("Interval Set to: %v", globalConf.Interval)
	if !viper.IsSet("publishtimeout") {
		log.Error().Msg("Publish Timeout not configured")
		return GlobalConfig{}, errors.New("publishtimeout not set")
	}
	globalConf.PublishTimeout = viper.GetInt("publishtimeout")
	if !(globalConf.PublishTimeout > 0) {
		log.Error().Msgf("Publish Timeout set to invalid value: %v", globalConf.PublishTimeout)
		return GlobalConfig{}, fmt.Errorf("publishtimeout set to invalid value %v", globalConf.PublishTimeout)
	}
	log.Debug().Msgf("Publish Timeout Set to: %v", globalConf.PublishTimeout)
	if !viper.IsSet("disconnecttimeout") {
		log.Error().Msg("Disconnect Timeout not configured")
		return GlobalConfig{}, errors.New("disconnecttimeout not set")
	}
	globalConf.DisconnectTimeout = viper.GetInt("disconnecttimeout")
	if !(globalConf.DisconnectTimeout > 0) {
		log.Error().Msgf("Disconnect Timeout set to invalid value: %v", globalConf.DisconnectTimeout)
		return GlobalConfig{}, fmt.Errorf("disconnecttimeout set to invalid value %v", globalConf.DisconnectTimeout)
	}
	log.Debug().Msgf("Disconnect Timeout Set to: %v", globalConf.DisconnectTimeout)
	return globalConf, nil
}
