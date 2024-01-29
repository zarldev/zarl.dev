package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/zarldev/zarldotdev/pkg/transport/http"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > VERSION.txt"
//go:embed VERSION.txt
var Commit string

type App struct {
	HTTPServer *http.Server
	Config     Config
}

type Config struct {
	HTTP    http.Config
	Name    string
	Version string
}

func New(config Config) (*App, error) {
	httpsvc, err := http.New(config.HTTP)
	if err != nil {
		return nil, err
	}
	return &App{
		Config:     config,
		HTTPServer: httpsvc,
	}, nil
}

func (a *App) Start() error {
	logo(a.Config)
	return a.HTTPServer.Start()
}

func (a *App) Stop() error {
	return a.HTTPServer.Stop()
}

func LoadConfig(filepath string) Config {
	config, err := loadConfigFromFile(filepath)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}
	return config
}

var ErrLoadingConfig = fmt.Errorf("error loading config from file")

func loadConfigFromFile(filepath string) (Config, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrLoadingConfig, err)
	}
	defer jsonFile.Close()
	var config Config
	jsonParser := json.NewDecoder(jsonFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %v", ErrLoadingConfig, err)
	}
	fmt.Println("Loaded config from file")
	return config, nil
}

func (c Config) WriteToFile(filepath string) error {
	return writeConfigToFile(c)
}

var ErrWritingConfig = fmt.Errorf("error writing config to file")

func writeConfigToFile(config Config) error {
	json, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWritingConfig, err)
	}
	err = os.WriteFile("./config.json", json, 0644)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWritingConfig, err)
	}
	return nil
}

const logoStr = `
   ________   ________  ________  _____       _______  ________  ________ 
  /__       \/        \/        \/     \    _/       \/        \/    /   \
 /         /    /    /    /    /      /   /    /    /       __/    /    /
/       __/         /        _/      /__ /         /       __/\        / 
\________/\___/____/\____/___/\________/ \________/\________/  \______/  
                =[%s]= =[version: %s]=
				
`

func logo(c Config) {
	fmt.Printf(logoStr, c.Name, c.Version)
}
