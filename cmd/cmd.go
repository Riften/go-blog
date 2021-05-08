package cmd

import (
	logging "github.com/ipfs/go-log"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type cmdMap map[string]func() error

var log = logging.Logger("cmd")

func Run() error {
	appCmd := kingpin.New("go-blog", "go-blog is a blog web framework implemented in pure go.")
	cmds := make(cmdMap)

	initCmd := appCmd.Command("init", "Initialize blog.")
	initHost := initCmd.Flag("host", "The host of router. It would be 127.0.0.1 by default.").Default("127.0.0.1").String()
	initPort := initCmd.Flag("port", "The port running web.").Default("8080").Uint16()
	initOverwrite := initCmd.Flag("overwrite", "").Bool()
	cmds[initCmd.FullCommand()] = func() error {
		return cmdInit(*initHost, *initPort, *initOverwrite)
	}

	startCmd := appCmd.Command("start", "Start the server.")
	cmds[startCmd.FullCommand()] = func () error {
		return cmdStart()
	}

	mdCmd := appCmd.Command("markdown", "Markdown related command. Mainly for debug.")
	mdRenderCmd := mdCmd.Command("render", "Render markdown to html.")
	mdRenderInput := mdRenderCmd.Arg("input", "The input file path.").Required().String()
	mdRenderOutput := mdRenderCmd.Arg("output", "The output file path. It would be input.html by default.").String()
	cmds[mdRenderCmd.FullCommand()] = func() error {
		return cmdMdRender(*mdRenderInput, *mdRenderOutput)
	}

	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
	}
	return nil
}
