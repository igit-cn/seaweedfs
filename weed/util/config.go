package util

import (
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/chrislusf/seaweedfs/weed/glog"
)

type Configuration interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetStringSlice(key string) []string
	SetDefault(key string, value interface{})
}

func LoadConfiguration(configFileName string, required bool) (loaded bool) {

	// find a filer store
	viper.SetConfigName(configFileName)              // name of config file (without extension)
	viper.AddConfigPath(".")                         // optionally look for config in the working directory
	viper.AddConfigPath("$HOME/.seaweedfs")          // call multiple times to add many search paths
	viper.AddConfigPath("/usr/local/etc/seaweedfs/") // search path for bsd-style config directory in
	viper.AddConfigPath("/etc/seaweedfs/")           // path to look for the config file in

	glog.V(1).Infof("Reading %s.toml from %s", configFileName, viper.ConfigFileUsed())

	if err := viper.MergeInConfig(); err != nil { // Handle errors reading the config file
		if strings.Contains(err.Error(), "Not Found") {
			glog.V(1).Infof("Reading %s: %v", viper.ConfigFileUsed(), err)
		} else {
			glog.Fatalf("Reading %s: %v", viper.ConfigFileUsed(), err)
		}
		if required {
			glog.Fatalf("Failed to load %s.toml file from current directory, or $HOME/.seaweedfs/, or /etc/seaweedfs/"+
				"\n\nPlease use this command to generate the default %s.toml file\n"+
				"    weed scaffold -config=%s -output=.\n\n\n",
				configFileName, configFileName, configFileName)
		} else {
			return false
		}
	}

	return true
}

type ViperProxy struct {
	*viper.Viper
	sync.Mutex
}

var (
	vp = &ViperProxy{}
)

func (vp *ViperProxy) SetDefault(key string, value interface{}) {
	vp.Lock()
	defer vp.Unlock()
	vp.Viper.SetDefault(key, value)
}

func (vp *ViperProxy) GetString(key string) string {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetString(key)
}

func (vp *ViperProxy) GetBool(key string) bool {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetBool(key)
}

func (vp *ViperProxy) GetInt(key string) int {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetInt(key)
}

func (vp *ViperProxy) GetStringSlice(key string) []string {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetStringSlice(key)
}

func GetViper() *ViperProxy {
	vp.Lock()
	defer vp.Unlock()

	if vp.Viper == nil {
		vp.Viper = viper.GetViper()
		vp.AutomaticEnv()
		vp.SetEnvPrefix("weed")
		vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	return vp
}
