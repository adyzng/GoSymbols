package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "gopkg.in/clog.v1"
	ini "gopkg.in/ini.v1"
)

var (
	Debug     bool
	AppPath   string
	AppName   string
	IsWindows bool

	WebRoot string // website assets folder
	Address string // website listen address
	Port    uint   // website listen port

	ClientID    string //
	ClientKey   string //
	ADAuthURI   string
	RedirectURI string
	GraphScope  string

	LogPath         string
	SymStoreExe     string
	Destination     string // pdb server destination
	BuildSource     string // pdb source folder
	PDBZipFile      string // pdb zip file, default `debug.zip`
	LatestBuildFile string // latest build trigger file `latestbuild.txt`
	ScheduleTime    string // default trigger time in 24H, eg: 5:00 => 5:00AM
	SymExcludeList  []string
)

func init() {
	if err := LoadConfig(); err != nil {
		LoadConfig("..\\config.ini")
	}
}

func exePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	return filepath.Abs(file)
}

// LoadConfig ...
//
func LoadConfig(files ...interface{}) error {
	var file = "config.ini"
	if len(files) > 0 {
		file = files[0].(string)
		files = files[1:]
	}
	if _, err := os.Stat(file); err != nil {
		return err
	}

	AppName = "GoSymbols"
	AppPath, _ = os.Getwd() //exePath()
	IsWindows = runtime.GOOS == "windows"
	log.Info("[Config] App path %s.", AppPath)

	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowBooleanKeys: true,
		Insensitive:      true,
		Loose:            true,
	},
		file,
		files...,
	)

	if err != nil {
		log.Fatal(2, "[Config] load config failed : %v.", err)
		return err
	}

	base := cfg.Section("base")
	Debug, _ = base.Key("Debug").Bool()

	SymStoreExe = base.Key("SYMSTORE_EXE").String()
	if SymStoreExe == "" {
		fmt.Println("[Config] SYMSTORE_EXE is missing.")
		log.Fatal(2, "[Config] SYMSTORE_EXE is missing.")
	}

	Destination = base.Key("DESTINATION").String()
	if Destination == "" {
		fmt.Println("[Config] DESTINATION is missing.")
		log.Fatal(2, "[Config] DESTINATION is missing.")
	}

	BuildSource = base.Key("BUILD_SOURCE").String()
	if BuildSource == "" {
		fmt.Println("[Config] BUILD_SOURCE is missing.")
		log.Fatal(2, "[Config] BUILD_SOURCE is missing.")
	}

	PDBZipFile = base.Key("DEBUG_ZIP").String()
	if PDBZipFile == "" {
		fmt.Println("[Config] BUILD_SOURCE is missing.")
		log.Fatal(2, "[Config] BUILD_SOURCE is missing.")
	}
	LatestBuildFile = base.Key("LATEST_BUILD").String()
	if LatestBuildFile == "" {
		fmt.Println("[Config] BUILD_SOURCE is missing.")
		log.Fatal(2, "[Config] BUILD_SOURCE is missing.")
	}
	LogPath = base.Key("LOG_PATH").String()
	if LogPath == "" {
		LogPath = "logs"
	}
	SymExcludeList = base.Key("EXCLUDE_LIST").Strings(",")
	for index, v := range SymExcludeList {
		SymExcludeList[index] = strings.ToLower(v)
	}

	appSec := cfg.Section("app")
	ClientID = appSec.Key("CLIENT_ID").String()
	ClientKey = appSec.Key("CLIENT_KEY").String()
	GraphScope = appSec.Key("GRAPH_SCOPE").String()
	RedirectURI = appSec.Key("REDIRECT_URI").String()

	web := cfg.Section("web")
	Address = web.Key("ADDRESS").String()
	WebRoot = web.Key("WEB_ROOT").String()
	if WebRoot == "" {
		WebRoot = "."
	}
	Port, _ = web.Key("PORT").Uint()
	if Port == 0 {
		Port = 8080
	}

	return nil
}

// GetTriggerTime ...
//
func GetTriggerTime() (hour, min int) {
	fmt.Sscanf(ScheduleTime, "%d:%d", &hour, &min)
	return
}
