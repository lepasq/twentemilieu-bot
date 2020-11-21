package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

/*
 * Struct for making an API request to get one's address.
 */
type Address struct {
	CompanyCode string
	PostCode    string
	HouseNumber int
	HouseLetter string
}

type ConfigFile struct {
	Whatsapp []map[string]interface{}
	Api      []map[string]interface{}
}

type Config struct {
	Whatsapp *string
	Api      *Address
}

func (c *Config) SetFromBytes() error {
	var config []ConfigFile
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	c.Whatsapp, err = config[0].GetRemoteJid()
	if err != nil {
		return fmt.Errorf("%d", err)
	}

	c.Api, err = config[1].GetAddress()
	if err != nil {
		return fmt.Errorf("%d", err)
	}
	return nil
}

func (config ConfigFile) GetRemoteJid() (*string, error) {
	remoteJid, exists := config.Whatsapp[0]["remoteJid"]
	if !exists {
		return nil, fmt.Errorf("remoteJid does not exist in config.yml")
	}
	str := fmt.Sprintf("%v", remoteJid)
	return &str, nil
}

func (config ConfigFile) GetAddress() (*Address, error) {
	postCode, exists := config.Api[0]["postCode"]
	if !exists {
		return nil, fmt.Errorf("postCode does not exist in config.yml")
	}

	houseNumber, exists := config.Api[1]["houseNumber"]
	if !exists {
		return nil, fmt.Errorf("houseNumber does not exist in config.yml")
	}

	houseLetter, exists := config.Api[2]["houseLetter"]
	if !exists {
		return nil, fmt.Errorf("houseLetter does not exist in config.yml")
	}

	hn, ok := houseNumber.(int)
	if !ok {
		return nil, fmt.Errorf("houseNumber was not an int")
	}

	address := Address{
		CompanyCode: "8d97bb56-5afd-4cbc-a651-b4f7314264b4",
		PostCode:    fmt.Sprintf("%v", postCode),
		HouseNumber: hn,
		HouseLetter: fmt.Sprintf("%v", houseLetter),
	}
	return &address, nil
}
