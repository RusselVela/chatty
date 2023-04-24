package app

import (
	"embed"
	"os"
	"path"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	kfs "github.com/knadh/koanf/providers/fs"
)

const (
	configFileName         = "application.yaml"
	defaultLocalConfigPath = "/config"
)

type configOptions struct {
	embeddedFileSystem struct {
		efs  *embed.FS
		path string
	}
	localConfigPath string
}

// NewDefaultKoanf creates a new koanf configuration with default options
func NewDefaultKoanf(efs *embed.FS, configPath string) (*koanf.Koanf, error) {
	options := configOptions{
		localConfigPath: defaultLocalConfigPath,
	}
	if configPath != "" {
		options.localConfigPath = configPath
	}
	options.embeddedFileSystem.efs = efs
	options.embeddedFileSystem.path = "embedded"
	return NewKoanf(options)
}

func NewKoanf(options configOptions) (*koanf.Koanf, error) {
	k := koanf.New(".")

	if options.embeddedFileSystem.efs != nil {
		if err := k.Load(kfs.Provider(options.embeddedFileSystem.efs, path.Join(options.embeddedFileSystem.path, configFileName)), yaml.Parser()); err != nil {
			return nil, err
		}
	}

	localConfigPath := path.Join(options.localConfigPath, configFileName)
	if _, err := os.Stat(localConfigPath); err == nil {
		if err := k.Load(file.Provider(localConfigPath), yaml.Parser()); err != nil {
			return nil, err
		}
	}

	if err := k.Load(env.Provider("", "_", strings.ToLower), nil); err != nil {
		return nil, err
	}

	return k, nil
}
