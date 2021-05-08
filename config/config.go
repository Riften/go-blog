// Package config provides interface and basic implementation for config tool.
package config

import (
	"encoding/json"
	"go-blog/common"
	"io/ioutil"
	"os"
	"sync"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("config")

// TODO: Implement this through database.
type Config interface {
	Host() string // Return the Served Host of server.
	SetHost(string)
	Port() uint16 // Return the Served PortNumber of server.
	// 0-65535
	SetPort(uint16)
	Resource() string // Return the path of resource directory.
	SetResource(string)
	RunningConfig() RunningConfig // Derive an RunningConfig from Config.
	Reset(string, uint16)
	WriteBack() error // Write config back to file or database
}

// Config implementation through config file.
type fileConfig struct {
	HostName   string `json:"host"`
	PortNumber uint16 `json:"port"`
	ResDir     string `json:"resource"`

	src    string // file path
	rwLock sync.Mutex
}

// DO NOT USE rwLock HERE
func OpenFileConfig() (Config, error) {
	filePath := common.PathCfgFile()

	res := &fileConfig{
		src: filePath,
	}

	// Read config back if file exist.
	if common.FileExist(filePath) {
		err := res.readFromFile(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, &common.ErrCfgNotExists{Path: filePath}
	}
	return res, nil
}

func NewFileConfig() Config {
	return &fileConfig{
		src: common.PathCfgFile(),
	}
}

func (c *fileConfig) Reset(hostName string, port uint16) {
	c.HostName = hostName
	c.PortNumber = port
	c.ResDir = common.PathResDir()
}

func (c *fileConfig) readFromFile(filePath string) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0444) // read only, -r--r--r--
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	return err
}

func (c *fileConfig) WriteBack() error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	var err error
	defer func() {
		if err != nil {
			log.Error("Error when write fileConfig back: ", err)
		}
	}()
	f, err := os.OpenFile(c.src, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664) // create + truncate + write only, -rw-rw-r--
	// create if not exists
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (c *fileConfig) RunningConfig() RunningConfig {
	return &rConfig{
		host:     c.HostName,
		port:     c.PortNumber,
		hostOnly: c.PortNumber != 80,
	}
}

func (c *fileConfig) Host() string {
	return c.HostName
}

func (c *fileConfig) Port() uint16 {
	return c.PortNumber
}

func (c *fileConfig) Resource() string {
	return c.ResDir
}

func (c *fileConfig) SetHost(h string) {
	c.HostName = h
}

func (c *fileConfig) SetPort(p uint16) {
	c.PortNumber = p
}

func (c *fileConfig) SetResource(r string) {
	c.ResDir = r
}
