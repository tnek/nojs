package admin

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/tebeka/selenium"
	"golang.org/x/net/publicsuffix"
)

func (v *Admin) adminCookie(ctx context.Context, adminName string) ([]*http.Cookie, error) {
	host := v.Host
	if v.Domain != "" {
		host = v.Domain
	}

	jar, err := cookiejar.New(
		&cookiejar.Options{PublicSuffixList: publicsuffix.List},
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	site := fmt.Sprintf("http://%v:%v", host, v.Port)
	path := fmt.Sprintf("%v/login", site)

	u, err := url.Parse(site)
	if err != nil {
		return nil, err
	}

	resp, err := client.PostForm(path, url.Values{
		"username": {adminName},
		"password": {kAdminPassword},
	})
	if err != nil {
		return nil, err
	}

	r := jar.Cookies(u)
	for _, cookie := range r {
		log.Printf("Found cookie: %v %v", cookie.Domain, cookie)
	}
	defer resp.Body.Close()
	return r, nil
}

func (v *Admin) toSeleniumCookie(c *http.Cookie) *selenium.Cookie {
	return &selenium.Cookie{
		Name:   c.Name,
		Value:  c.Value,
		Path:   "/",
		Domain: v.Domain,
		Secure: c.Secure,
		Expiry: math.MaxUint32,
	}
}

func (v *Admin) Auth(ctx context.Context, adminName string) ([]*selenium.Cookie, error) {
	cookies, err := v.adminCookie(ctx, adminName)
	if err != nil {
		return nil, err
	}
	conv := []*selenium.Cookie{}
	for _, cookie := range cookies {
		conv = append(conv, v.toSeleniumCookie(cookie))
	}
	return conv, nil
}
