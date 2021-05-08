package cmd

import (
	"go-blog/common"
	"go-blog/config"
	"os"
)

func cmdInit(initHost string, initPort uint16, overWrite bool) error {
	// Init repo dir
	var err error
	dir := common.PathCfgDir()
	if !common.DirectoryExist(dir) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		isEmpty, err := common.IsDirEmpty(dir)
		if err != nil {
			return err
		}
		if !isEmpty {
			if overWrite {
				err = os.RemoveAll(dir)
				if err != nil {
					return err
				}
				err = os.Mkdir(dir, os.ModePerm)
				if err != nil {
					return err
				}
			} else {
				return &common.ErrDirectoryNotEmpty{Path: dir}
			}
		}
	}

	// Init resource dir
	err = common.CopyDir("resources", common.PathResDir())
	if err != nil {
		return err
	}

	cfg := config.NewFileConfig()
	cfg.Reset(initHost, initPort)
	err = cfg.WriteBack()
	if err != nil {
		return err
	}
	return nil
}
