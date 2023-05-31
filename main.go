package main

import (
	"os"
	"path"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/o98k-ok/dian/file"
	"github.com/o98k-ok/lazy/v2/alfred"
	"github.com/o98k-ok/lazy/v2/lark/doc"
)

func query(root string, args []string) *alfred.Items {
	if len(args) == 0 {
		args = append(args, "")
	}

	items := alfred.NewItems()
	files, err := file.Walk(root)
	if err != nil {
		return items
	}

	for _, item := range files {
		if !strings.Contains(item, args[0]) {
			continue
		}
		_, name := path.Split(item)
		items.Append(&alfred.Item{
			Arg:   item,
			Title: name,
		})
	}
	return items
}

type Config struct {
	Parent   string
	AppId    string
	Secret   string
	LocalDir string
}

func NewConfig(envs alfred.Envs) *Config {
	parent := envs["FOLDER_TOKEN"]
	appId := envs["APP_ID"]
	sec := envs["APP_KEY"]
	localDir := envs["LOCAL"]
	if len(parent) == 0 || len(appId) == 0 || len(sec) == 0 || len(localDir) == 0 {
		return nil
	}
	return &Config{
		Parent:   parent,
		AppId:    appId,
		Secret:   sec,
		LocalDir: localDir,
	}
}

func NewFromEnv() *Config {
	parent := os.Getenv("FOLDER_TOKEN")
	appId := os.Getenv("APP_ID")
	sec := os.Getenv("APP_KEY")
	localDir := os.Getenv("LOCAL")
	if len(parent) == 0 || len(appId) == 0 || len(sec) == 0 || len(localDir) == 0 {
		return nil
	}
	return &Config{
		Parent:   parent,
		AppId:    appId,
		Secret:   sec,
		LocalDir: localDir,
	}
}

func main() {
	cli := alfred.NewApp("lark toolkit for importing tasks")

	envs, err := alfred.GetFlowEnv()
	if err != nil {
		alfred.InputErrItems(err.Error()).Show()
		return
	}

	config := NewConfig(envs)
	if config == nil {
		alfred.InputErrItems("missing required envs, please check").Show()
		return
	}

	cli.Bind("query", func(s []string) {
		query(config.LocalDir, s).Show()
	})

	cli.Bind("import", func(s []string) {
		lark := doc.NewLarkDocer(config.AppId, config.Secret, config.Parent)

		dat, err := os.ReadFile(s[0])
		if err != nil {
			alfred.ErrItems("read source file failed", err).Show()
			return
		}
		_, name := path.Split(s[0])
		token, err := lark.Upload(name, dat)
		if err != nil {
			alfred.ErrItems("upload file to lark failed", err).Show()
			return
		}

		newName := strings.Join([]string{strings.Trim(name, "md"), "v", gofakeit.AppVersion()}, "_")
		if err := lark.FormatEdit(token, newName); err != nil {
			alfred.ErrItems("change file format failed", err).Show()
			return
		}
	})
	cli.Run(os.Args)
}
