package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/o98k-ok/dian/builder"
	"github.com/o98k-ok/dian/file"
	"github.com/o98k-ok/lazy/v2/alfred"
)

func query(root string, args []string) *alfred.Items {
	if len(args) == 0 {
		return alfred.InputErrItems("bad arg length")
	}

	items := alfred.NewItems()
	for _, item := range file.Grep(root, args[0]) {
		filename := strings.TrimLeft(item.Filename, root)
		items.Append(&alfred.Item{
			Arg:      builder.Openfile(filename, item.Line),
			Title:    filename,
			SubTitle: item.Content,
		})
	}
	return items
}

func main() {
	cli := alfred.NewApp("obsidian commandline tools")

	envs, err := alfred.GetFlowEnv()
	if err != nil {
		alfred.InputErrItems(err.Error()).Show()
		// return
		envs = make(alfred.Envs)
	}

	cli.Bind("query", func(s []string) {
		fmt.Println(strings.Replace(query(envs["root"], s).Encode(), `\u0026`, "&", -1))
	})
	cli.Run(os.Args)
}
