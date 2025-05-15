package config

import (
	"fmt"
	"os"
	"log"

	"gopkg.in/yaml.v2"
)

type Tenant struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type DatabaseConfig struct {
	Host    string            `yaml:"host"`
	Port    int               `yaml:"port"`
	Name    string            `yaml:"name"`
	Tenants map[string]Tenant `yaml:"tenants"`
}

type CognitoConfig struct {
	UserPoolId string `yaml:"userPoolId"`
	ClientId   string `yaml:"clientId"`
	TokenType  string `yaml:"tokenType"`
}

type AppConfig struct {
	Server    struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	AwsCognito CognitoConfig `yaml:"awsCognito"`
}


var EndpointToTableMap = map[string]string{
	"/dataLPI":     "public_dev.datalpi_filtered_view",
	"/dataHeader": "public_dev.dataheader_filtered_view",
	"/dataHeight": "public_dev.dataheight_filtered_view",
	"/dataGap": "public_dev.datagap_filtered_view",
	"/dataHorizontalFlux": "public_dev.datahorizontalflux_filtered_view",
	"/dataPlotCharacterization": "public_dev.dataplotcharacterization_filtered_view",
	"/dataSoilHorizons": "public_dev.datasoilhorizons_filtered_view",
	"/dataSoilStability": "public_dev.datasoilstability_filtered_view",
	"/dataSpeciesInventory": "public_dev.dataspeciesinventory_filtered_view",
	"/geoIndicators": "public_dev.geoindicators_filtered_view",
	"/geoSpecies": "public_dev.geospecies_filtered_view",
	"/tblRHEM": "public_dev.tblrhem_filtered_view",
	"/tblProject": "public_dev.tblproject_filtered_view",
	"/dataAeroSummary": "aero_data.aero_summary",
}

var Config AppConfig


func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	log.Println("Configuration loaded successfully")
	return nil
}
