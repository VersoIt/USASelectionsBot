package cfg

import (
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"sync"
	"time"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	ParseUrl                 string        `yaml:"parse-url"`
	FirstCandidateContainer  string        `yaml:"first-candidate-container"`
	SecondCandidateContainer string        `yaml:"second-candidate-container"`
	UpdateDelay              time.Duration `yaml:"update-delay"`
	BotToken                 string        `yaml:"bot-token"`
	ChannelId                string        `yaml:"channel-id"`
	FirstCandidateName       string        `yaml:"first-candidate-name"`
	SecondCandidateName      string        `yaml:"second-candidate-name"`
}

func Get() Config {
	once.Do(func() {
		f, err := os.Open("config.yml")
		if err != nil {
			panic(err)
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		data, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}

		if err = yaml.Unmarshal(data, &cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}
