package admin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/tnek/ctf-browser-visitor/ctfvisitor"
	"github.com/tnek/notes-site/config"
	"github.com/tnek/notes-site/model"
)

const (
	kSeleniumPath     = "/usr/local/bin/selenium-server.jar"
	kGeckodriverPath  = "/usr/local/bin/geckodriver"
	kChromedriverPath = "/usr/local/bin/chromedriver"
	kQueueSize        = 10000
	kNumWorkers       = 4
)

var (
	kChromeConfig = &ctfvisitor.Config{
		QueueSize:    kQueueSize,
		SeleniumPath: kSeleniumPath,
		BrowserPath:  kChromedriverPath,
		Browser:      ctfvisitor.CHROME,
		MinPort:      10000,
		MaxPort:      20000,
	}
	kFirefoxConfig = &ctfvisitor.Config{
		QueueSize:    kQueueSize,
		SeleniumPath: kSeleniumPath,
		BrowserPath:  kGeckodriverPath,
		Browser:      ctfvisitor.FIREFOX,
		MinPort:      20001,
		MaxPort:      30000,
	}
)

type Admin struct {
	Host        string
	Domain      string
	Port        int
	FirefoxFlag string

	firefox *ctfvisitor.Dispatch
}

func Init(ctx context.Context, ac *config.AppConfig, firefoxFlag string) (*Admin, error) {
	firefox, err := ctfvisitor.Init(kFirefoxConfig)
	if err != nil {
		return nil, err
	}

	go firefox.LoopWithRestart(ctx, kNumWorkers)

	return &Admin{
		FirefoxFlag: firefoxFlag,

		Host:   ac.Host,
		Domain: ac.Domain,
		Port:   ac.Port,

		firefox: firefox,
	}, nil
}

func (a *Admin) visit(q *ctfvisitor.Dispatch, adminName string) error {
	host := a.Host
	if a.Domain != "" {
		host = a.Domain
	}

	s := fmt.Sprintf("http://%v:%v/", strings.Clone(host), a.Port)
	cookies, err := a.Auth(context.Background(), adminName)
	if err != nil {
		return err
	}

	log.Printf("visiting %v %v", s, cookies)
	site := &ctfvisitor.Site{Path: s, Cookies: cookies}

	if err := q.Queue(site); err != nil {
		return fmt.Errorf("visit as '%v' failed: %w", adminName, err)
	}
	return nil
}

func (a *Admin) Visit(u *model.User, recipient string) error {
	log.Printf("calling visitor...")
	if firefoxAdminName(u) == recipient {
		return a.visit(a.firefox, firefoxAdminName(u))
	} else {
		log.Printf("attempting to visit an invalidly named admin?: %v", recipient)
		// Avoid exposing information to the competitor
		return errors.New("something went wrong!")
	}
	return nil
}
