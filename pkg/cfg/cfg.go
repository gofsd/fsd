package cfg

import (
	"os"
	"strings"

	"github.com/gofsd/fsd/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	//CliName name of this cli
	CliName = "fsd"
	//ConfigType type of config for `CliName` cli
	ConfigType = "json"
)

// Cfg config
type Cfg struct {
	CmdPath    string
	ConfigPath string
	Viper      *viper.Viper
	RootCmd    *cobra.Command
}

var (
	//Config default config
	Config *Cfg
)

// NewCfg global config
func GetCfg() *Cfg {
	if Config == nil {
		Config = &Cfg{Viper: viper.New()}
	}
	return Config
}

func (c *Cfg) FindCMD(s string) *cobra.Command {
	var (
		cmd *cobra.Command
	)
	var args = strings.Split(s, " ")
	cmd, args, _ = c.RootCmd.Find(args)
	return cmd
}

// Set cmd path for current executing command
func (c *Cfg) SetCommandPath(cmdPath string) {
	c.CmdPath = cmdPath
}

func (c *Cfg) SetRootCMD(cmd *cobra.Command) *Cfg {
	c.RootCmd = cmd
	return c
}

// ReadCfgFromFile read config from file
func (c *Cfg) ReadCfgFromFile(cmd *cobra.Command) *Cfg {
	var (
		configName, configType string
		err                    error
	)
	configName = CliName
	configType = ConfigType
	c.ConfigPath = util.FindRootDir(CliName + "." + configType)
	c.Viper.AutomaticEnv()
	c.Viper.SetConfigName(configName)
	c.Viper.SetConfigType(configType)
	c.Viper.AddConfigPath(c.ConfigPath)
	if err = c.Viper.ReadInConfig(); err == nil {
	} else if err != nil {
		c.Viper.SafeWriteConfig()
		if err = c.Viper.ReadInConfig(); err == nil {
		}
	}

	util.HandleError(err)
	return c
}

func (c *Cfg) SetLog() {
	file, err := os.OpenFile("./log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	util.HandleError(err)
	//log.SetOutput(file)
	c.RootCmd.SetOut(file)
	c.RootCmd.SetErr(file)
}

// WriteCfgToFile write config to file
func (c *Cfg) WriteCfgToFile() *Cfg {
	var err error
	cmd := c.FindCMD(c.CmdPath)
	if cmd == nil {
		err = c.Viper.SafeWriteConfig()
	} else {
		c.CmdPath = strings.Join(strings.Split(cmd.CommandPath(), " "), ".")
		err = c.Viper.WriteConfig()
	}
	util.HandleError(err)

	return c
}

// Set cfg field
func (c *Cfg) Set(k string, val interface{}) {
	c.Viper.Set(k, val)
}

// Get cfg fileed
func (c *Cfg) Get(key string) interface{} {
	if c == nil {
		c = GetCfg()
	}
	return c.Viper.Get(key)
}

// GetString cfg fileed
func (c *Cfg) GetString(key string) string {
	if c == nil {
		c = GetCfg()
	}
	return c.Viper.GetString(key)
}

func (c *Cfg) GetPathWithoutRootCmdName() string {
	if strings.Contains(c.CmdPath, CliName) {
		return strings.Replace(c.CmdPath, CliName, "", 1)
	}
	return c.CmdPath
}

// Bind cfg fileed
func (c *Cfg) Bind(key string) {
	cmd := c.FindCMD(c.GetPathWithoutRootCmdName())
	c.Viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key))
}
