package cmd

import (
	"go-blog/common"
	"os"
)

func cmdMdRender(input string, output string) error {
	info, err := os.Stat(input)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return common.MdRenderRecursively(input, output, true, true)
	} else {
		return common.MdRenderFile(input, output)
	}
}
