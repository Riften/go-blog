package common

import (
	"log"
	"os"
	"path/filepath"
)

const DEFAULT_CFG_DIR = ".RiftenGoBlog"
const ENV_CFG_DIR = "GOBLOG_CFG"
const ENV_RESOURCE_DIR = "GO_BLOG_RES"

// PathCfgDir return the path of repo directory.
// It would be $HOME/.RiftenGoBlog by default.
// It also can be set through os environment GO_BLOG_CFG
// TODO: Change Cfg to Repo
func PathCfgDir() string {
	dir := os.Getenv(ENV_CFG_DIR)
	if dir != "" {
		return dir
	}
	homeDir, err := Home()
	if err != nil {
		log.Fatal("can not fetch home directory")
	}
	return filepath.Join(homeDir, DEFAULT_CFG_DIR)
}

func PathCfgFile() string {
	dir:= PathCfgDir()
	return filepath.Join(dir, "config.json")
}

func PathResDir() string {
	dir := os.Getenv(ENV_RESOURCE_DIR)
	if dir != "" {
		return dir
	}
	cfgDir := PathCfgDir()
	return filepath.Join(cfgDir, "res")
}