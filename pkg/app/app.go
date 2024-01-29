package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/zarldev/zarldotdev/pkg/transport/http"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > git.sha"
//go:generate sh -c "date > build.date"
//go:embed git.sha
var gitsha string

//go:embed build.date
var builddate string

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

var logoStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))
var subtitleStyle = lipgloss.NewStyle().
	Padding(0, 22).
	Bold(true).
	Foreground(lipgloss.Color("#2db4ee"))

var shaStyle = lipgloss.NewStyle().
	Padding(0, 12).
	Bold(true).
	Foreground(lipgloss.Color("#3d3d3d"))

const logoStr = `
    _______  ________  ________  _____       _______  ________  ________ 
  _/__     \/        \/        \/     \    _/       \/        \/     /  \
 /         /    /    /    /    /      /   /    /    /       __/     /   /
/       __/         /        _/      /__ /         /       __/\        / 
\________/\___/____/\____/___/\________/ \________/\________/  \______/  				
`
const subtitle = `[service: %s @ v%s]`
const sha = `[git sha: %s]`
const build = `[build date: %s]`

func logo(c Config) {
	fmt.Print(logoStyle.Render(logoStr))
	fmt.Println()
	fmt.Print(subtitleStyle.Render(fmt.Sprintf(subtitle, c.Name, c.Version)))
	fmt.Println()
	fmt.Print(shaStyle.Render(fmt.Sprintf(sha, gitsha)))
	fmt.Println()
	fmt.Print(shaStyle.Render(fmt.Sprintf(build, strings.ReplaceAll(builddate, "\n", ""))))
	fmt.Println()
	fmt.Println()
}
