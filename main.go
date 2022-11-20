package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	"github.com/tnek/notes-site/config"
	"github.com/tnek/notes-site/model"
	"github.com/tnek/notes-site/serv"

	// OMIT
	"github.com/tnek/notes-site/admin"
	// OMITEND
)

func main() {
	DOMAIN := flag.String("domain", "", "External domain of the site")
	ASSET_PATH := flag.String("assets", "assets/templates/*", "Path to templates of the site")
	FIREFOX_FLAG := flag.String("firefox_flag", "flag{too_bad_chrome_cant_hang", "flag")
	PORT := flag.Int("p", 8080, "Port to listen off of")

	flag.Parse()

	glob, err := filepath.Glob(*ASSET_PATH)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.AppConfig{
		Host:      "0.0.0.0",
		Port:      *PORT,
		Templates: glob,
		Domain:    *DOMAIN,
	}

	if err := model.Conn(cfg.DBPath); err != nil {
		log.Fatal(err)
	}

	app, err := serv.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// OMIT
	admin, err := admin.Init(ctx, cfg, *FIREFOX_FLAG)
	if err != nil {
		log.Fatal(err)
	}

	app.Admin = admin

	if ok, _ := model.UserExists("test"); !ok {
		if err := app.DoRegister("test", "test"); err != nil {
			log.Fatal(err)
		}
	}
	// OMITEND

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
