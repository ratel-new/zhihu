package main

import (
	"bufio"
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strings"
)

type cookie struct {
	Domain   string
	HttpOnly bool
	Path     string
	Secure   bool
	Expiry   string
	Name     string
	Value    string
}

var (
	cookies = make([]cookie, 0)
)

func init() {
	txt, err := os.Open("zhihu.com_cookies.txt")
	if err != nil {
		return
	}
	defer func(txt *os.File) {
		err := txt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(txt)
	scanner := bufio.NewScanner(txt)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		sts := strings.Split(scanner.Text(), "\t")
		cookies = append(cookies, cookie{Domain: sts[0], HttpOnly: sts[1] == "true", Path: sts[2], Secure: sts[3] == "true", Name: sts[5], Value: sts[6]})
	}
}

func SetCookie() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		for _, v := range cookies {
			_ = network.SetCookie(v.Name, v.Value).
				WithDomain(v.Domain).
				WithPath(v.Path).
				WithHTTPOnly(v.HttpOnly).
				WithSecure(v.Secure).
				Do(ctx)
		}
		return nil
	})
}
