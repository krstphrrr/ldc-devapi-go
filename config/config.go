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
	"/dataLPI":     "datalpi_filtered_view",
	"/dataHeader": "dataheader_filtered_view",
	"/dataHeight": "dataheight_filtered_view",
	"/dataGap": "datagap_filtered_view",
	"/dataHorizontalFlux": "datahorizontalflux_filtered_view",
	"/dataPlotCharacterization": "dataplotcharacterization_filtered_view",
	"/dataSoilHorizons": "datasoilhorizons_filtered_view",
	"/dataSoilStability": "datasoilstability_filtered_view",
	"/dataSpeciesInventory": "dataspeciesinventory_filtered_view",
	"/geoIndicators": "geoindicators_filtered_view",
	"/geoSpecies": "geospecies_filtered_view",
	"/tblRHEM": "tblrhem_filtered_view",
	"/tblProject": "tblproject_filtered_view",
	"/aeroSummary": "aero_summary",
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
