package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Peterpig/mini_godis/lib/logger"
	"github.com/Peterpig/mini_godis/lib/utils"
)

type ServerProperties struct {
	RunID      string `cfg:"runid"`
	Bind       string `cfg:"bind"`
	Port       int    `cfg:"port"`
	Dir        string `cfg:"dir"`
	AppendOnly bool   `cfg:"appendonly"`

	MaxClients int `cfg:"maxclients"`

	// config file path
	CfPath string `cfg:"cf,omitempty"`
}

type ServerInfo struct {
	StartUpTime time.Time
}

var Properties *ServerProperties
var EachTimeServerInfo *ServerInfo

func init() {
	EachTimeServerInfo = &ServerInfo{
		StartUpTime: time.Now(),
	}
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6399,
		AppendOnly: false,
		RunID:      utils.RandString(40),
	}
}

func parse(configFile *os.File) *ServerProperties {
	config := &ServerProperties{}
	rawMap := make(map[string]string)

	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过注释行
		if len(line) > 0 && strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}

		pivot := strings.IndexAny(line, " ")
		// 空格分隔符
		if (pivot > 0) && (pivot < len(line)-1) {
			key := strings.ToLower(line[0:pivot])
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[key] = value
		}
	}

	logger.Info("%v", rawMap)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// 反射配置文件
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	n := t.Elem().NumField()

	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldValue := v.Elem().Field(i)

		key, ok := field.Tag.Lookup("cfg")
		if !ok || strings.TrimLeft(key, " ") == "" {
			key = field.Name
		}

		value, ok := rawMap[strings.ToLower(key)]
		if !ok {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				fieldValue.SetInt(intValue)
			}
		case reflect.Bool:
			boolValue := value == "yes"
			fieldValue.SetBool(boolValue)
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				slice := strings.Split(value, ",")
				fieldValue.Set(reflect.ValueOf(slice))
			}
		}

	}

	return config
}

func SetupConfig(configFilename string) {
	file, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	Properties = parse(file)
	logger.Info("%v", Properties)

	Properties.RunID = utils.RandString(40)
	configFilePath, err := filepath.Abs(configFilename)
	if err != nil {
		return
	}

	Properties.CfPath = configFilePath
	if Properties.Dir == "" {
		Properties.Dir = "."
	}
}
