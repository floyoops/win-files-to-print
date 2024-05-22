package print

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const defaultConfigFile = "c:\\win-files-to-print-config.json"

type ConfigPrinter struct {
	configFile string
	Folder     string `json:"folder"`
	PrintName  string `json:"printName"`
}

func NewConfigPrinter() (*ConfigPrinter, error) {
	return &ConfigPrinter{configFile: defaultConfigFile, Folder: "", PrintName: ""}, nil
}

func (p *ConfigPrinter) SetFolder(folder string) {
	p.Folder = folder
}

func (p *ConfigPrinter) GetFolder() string {
	return p.Folder
}

func (p *ConfigPrinter) SetPrintName(printName string) {
	p.PrintName = printName
}

func (p *ConfigPrinter) GetPrintName() string {
	return p.PrintName
}

func (p *ConfigPrinter) GetConfigFile() string {
	return p.configFile
}

func (p *ConfigPrinter) SaveFile() error {
	configData, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = os.WriteFile(p.configFile, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (p *ConfigPrinter) DeleteFile() error {
	err := os.Remove(p.configFile)
	if err != nil {
		return err
	}
	return nil
}

func (p *ConfigPrinter) LoadConfig() error {
	configData, err := os.ReadFile(p.configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %s %v", p.configFile, err)
	}
	err = json.Unmarshal(configData, &p)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
	return nil
}

func (p *ConfigPrinter) Render() string {
	return fmt.Sprintf("configFile %s\n"+
		"Folder %s\n"+
		"Print Name: %s\n",
		p.configFile, p.Folder, p.PrintName,
	)
}
