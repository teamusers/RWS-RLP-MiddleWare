package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	system "lbe/log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type CmdConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	Database    DatabaseConfig `yaml:"database"`
	Redis       RedisConfig    `yaml:"redis"`
	Chain       []ChainConfig  `yaml:"chain"`
	Log         LogConfig      `yaml:"log"`
	AllStart    int            `yaml:"allStart"`
	Cmd         CmdConfig      `yaml:"cmd"`
	Http        HttpConfig     `yaml:"http"`
	ProxyEnable bool           `yaml:"proxyEnable"`
	Smtp        SmtpConfig     `yaml:"smtp"`
	// other fields you already have...
	Api struct {
		Memberservice struct {
			Host   string `yaml:"host"`
			AppID  string `yaml:"appid"`
			Secret string `yaml:"secret"`
		} `yaml:"memberservice"`
		Rlp struct {
			Host   string `yaml:"host"`
			ApiKey string `yaml:"apikey"`
		} `yaml:"rlp"`
		Eeid struct {
			Host                   string `yaml:"host"`
			AuthHost               string `yaml:"authhost"`
			TenantID               string `yaml:"tenantid"`
			ClientID               string `yaml:"clientid"`
			ClientSecret           string `yaml:"clientsecret"`
			UserIdLinkExtensionKey string `yaml:"userIdLinkExtensionKey"`
		} `yaml:"eeid"`
		Cms struct {
			Host     string `yaml:"host"`
			ApiKey   string `yaml:"apikey"`
			SystemID string `yaml:"SystemID"`
		} `yaml:"cms"`
	} `yaml:"api"`
	Application struct {
		RLPNumberingFormat struct {
			MaxAttempts  int    `yaml:"maxAttempts"`
			RLPNODefault string `yaml:"rlpNoDefault"`
		} `yaml:"rlpNumberingFormat"`
	} `yaml:"application"`
}

// DatabaseConfig holds the database connection parameters.
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Instance string `yaml:"instance"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"TimeZone"`
}

type SmtpConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

// RedisConfig holds the Redis connection parameters.
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

// ChainConfig holds the Solana chain RPC endpoints.
type ChainConfig struct {
	Name         string   `yaml:"name"`
	WsRpc        string   `yaml:"wsRpc"`
	QueryRpc     []string `yaml:"queryRpc"`
	SlotParallel int      `yaml:"slotParallel"`
	TxDetal      int      `yaml:"txDetal"`
	RangeRound   int      `yaml:"rangeRound"`
	Rpcs         []RpcMapper
	RpcMap       map[string]int
}

// LogConfig holds the logging directory and file name.
type LogConfig struct {
	Path string `yaml:"path"`
}

type RpcMapper struct {
	Rpc   string
	Quote int
}

type AveConf struct {
	Avekey  string `yaml:"avekey"`
	Aveauth string `yaml:"aveauth"`
	Tpl     string `yaml:"tpl"`
}

func initRpcs(ts []ChainConfig) {
	for i := range ts {
		ts[i].initRpc()
	}
}

func (t *ChainConfig) initRpc() {
	if len(t.QueryRpc) == 0 {
		log.Fatal("error rpc config")
	}

	t.RpcMap = make(map[string]int)
	for _, r := range t.QueryRpc {
		v := strings.Split(r, "||")
		num := 0
		if len(v) == 2 {
			numStr := v[1]
			var err error
			num, err = strconv.Atoi(numStr)
			if err != nil {
				log.Fatal("error rpc inner format with Quote:", err)
			}
		}

		t.Rpcs = append(t.Rpcs, RpcMapper{
			Rpc:   v[0],
			Quote: num,
		})
		t.RpcMap[v[0]] = num
	}
}

func GetRpcConfig(code string) *ChainConfig {
	for _, v := range systemConfig.Chain {
		if v.Name == code {
			return &v
		}
	}
	return nil
}

func (t ChainConfig) GetRpc() []string {
	r := make([]string, 0)
	for _, v := range t.Rpcs {
		r = append(r, v.Rpc)
	}
	return r
}

func (t ChainConfig) GetRpcMapper() []RpcMapper {
	return t.Rpcs
}

func (t *ChainConfig) GetSlotParallel() int {
	if t.SlotParallel > 0 {
		return t.SlotParallel
	}
	return 1
}

func (t *ChainConfig) GetTxDelay() int {
	if t.TxDetal > 0 {
		return t.TxDetal
	}
	return 0
}

var systemConfig = &Config{}

func GetConfig() Config {
	return *systemConfig
}

func findProjectRoot(currentDir, rootIndicator string) (string, error) {
	if _, err := os.Stat(filepath.Join(currentDir, rootIndicator)); err == nil {
		return currentDir, nil
	}
	parentDir := filepath.Dir(currentDir)
	if currentDir == parentDir {
		return "", os.ErrNotExist
	}
	return findProjectRoot(parentDir, rootIndicator)
}

func init() {
	var confFilePath string

	if configFilePathFromEnv := os.Getenv("DALINK_GO_CONFIG_PATH"); configFilePathFromEnv != "" {
		confFilePath = configFilePathFromEnv
		log.Println("DALINK_GO_CONFIG_PATH:", confFilePath)
	} else {
		_, filename, _, _ := runtime.Caller(0)
		testDir := filepath.Dir(filename)
		log.Println("filename:", filename)
		log.Println("Looking for marker in:", testDir)
		confFilePath = testDir
		//confFilePath, _ = findProjectRoot(testDir, "__mark__")
		if len(confFilePath) > 0 {
			confFilePath += "/dev.yml"
		}
	}

	log.Println("Found project root as:", confFilePath)
	if len(confFilePath) == 0 {
		log.Fatal("System root directory setting error.123")
	}
	log.Println("current config file ", confFilePath)

	configFileToCheck := confFilePath
	if _, err := os.Stat(configFileToCheck); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist at path: %s", configFileToCheck)
	}
	log.Println("Current config file exists at:", configFileToCheck)

	viper.SetConfigFile(confFilePath)

	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration file: %s", err)
	}

	err = viper.Unmarshal(&systemConfig)
	if err != nil {
		log.Fatalf("Unable to parse configuration: %s", err)
	}
	initRpcs(systemConfig.Chain)

	system.InitLogger(systemConfig.Log.Path)

	_ = godotenv.Load()
}
