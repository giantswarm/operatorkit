package harness

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultConfigFile = "config.yaml"

	DefaultKubeConfig = "/workdir/.shipyard/config"
)

type Harness struct {
	logger micrologger.Logger
	fs     afero.Fs
	cfg    Config
}

type Config struct {
	ExistingCluster bool `yaml:"existingCluster"`
	RemoteCluster   bool `yaml:"remoteCluster"`
}

func New(logger micrologger.Logger, fs afero.Fs, cfg Config) *Harness {
	return &Harness{
		logger: logger,
		fs:     fs,
		cfg:    cfg,
	}
}

// Init initializes the harness.
func (h *Harness) Init(ctx context.Context) error {
	h.logger.Log("info", "starting harness initialization")
	baseDir, err := BaseDir()
	if err != nil {
		return microerror.Mask(err)
	}
	workDir := filepath.Join(baseDir, "workdir")
	err = h.fs.MkdirAll(workDir, 0777)
	if err != nil {
		return microerror.Mask(err)
	}

	shipyardDir := filepath.Join(workDir, ".shipyard")
	err = h.fs.MkdirAll(shipyardDir, 0777)
	if err != nil {
		return microerror.Mask(err)
	}

	// circumvent umask settings, by assigning the right
	// permissions to workdir and its parent
	for _, d := range []string{baseDir, workDir, shipyardDir} {
		err = h.fs.Chmod(d, 0777)
		if err != nil {
			return microerror.Mask(err)
		}
	}
	h.logger.Log("info", "finished harness initialization")
	return nil
}

// WriteConfig is a Task that persists the current config to a file.
func (h *Harness) WriteConfig(ctx context.Context) error {
	dir, err := BaseDir()
	if err != nil {
		return microerror.Mask(err)
	}

	content, err := yaml.Marshal(&h.cfg)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, defaultConfigFile), []byte(content), 0644)

	return microerror.Mask(err)
}

// ReadConfig populates a Config struct data read
// from a default file location.
func (h *Harness) ReadConfig() (Config, error) {
	dir, err := BaseDir()
	if err != nil {
		return Config{}, microerror.Mask(err)
	}

	afs := &afero.Afero{Fs: h.fs}
	content, err := afs.ReadFile(filepath.Join(dir, defaultConfigFile))
	if err != nil {
		return Config{}, microerror.Mask(err)
	}

	c := &Config{}

	if err := yaml.Unmarshal(content, c); err != nil {
		return Config{}, microerror.Mask(err)
	}

	return *c, nil
}

func BaseDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", microerror.Mask(err)
	}
	return filepath.Join(dir, ".e2e-harness"), nil
}

func GetProjectName() string {
	if os.Getenv("CIRCLE_PROJECT_REPONAME") != "" {
		return os.Getenv("CIRCLE_PROJECT_REPONAME")
	}
	dir, err := os.Getwd()
	if err != nil {
		return "e2e-harness"
	}
	return filepath.Base(dir)
}

func GetProjectTag() string {
	if os.Getenv("CIRCLE_SHA1") != "" {
		return os.Getenv("CIRCLE_SHA1")
	}
	return "latest"
}
