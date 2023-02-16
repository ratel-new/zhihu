package main

import (
	"bufio"
	"context"
	"github.com/buaazp/fasthttprouter"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {

	router := fasthttprouter.New()

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag(`headless`, false),
		chromedp.DisableGPU,
		chromedp.Flag(`disable-extensions`, false),
		chromedp.Flag(`enable-automation`, false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	logCtx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithDebugf(log.Printf),
	)

	defer cancel()

	chromeCtx, cancel := context.WithTimeout(logCtx, 1*time.Hour)
	defer cancel()

	router.GET("/", func(ctx *fasthttp.RequestCtx) {

		body := ctx.Request.RequestURI()
		query, err := url.Parse(string(body))
		if err != nil {
			ctx.Error("Request Body Failed to parse", 400)
			return
		}

		if len(query.Query().Get("url")) == 0 {
			ctx.Error("url must be provided", 400)
			return
		}

		zhihuUrl := query.Query().Get("url")

		if strings.HasPrefix(zhihuUrl, "https://www.zhihu.com") && strings.Contains(zhihuUrl, "question") && strings.Contains(zhihuUrl, "answer") {
			var example string

			err = chromedp.Run(chromeCtx,
				SetCookie(),
				chromedp.Navigate(zhihuUrl),
				chromedp.WaitVisible(`.QuestionAnswer-content`),
				chromedp.Evaluate(`
		try {
		 let DOMContentLoadedCode = document.createElement('div')
		 DOMContentLoadedCode.id = 'zhihuDOMContentLoaded'
		 document.body.appendChild(DOMContentLoadedCode)
		 const mainDom = document.querySelector('div.QuestionAnswer-content')
		 const titleDom = document.querySelector('.QuestionHeader-title')
		 if (mainDom && titleDom) {
		     const newTitleDom = titleDom.cloneNode()
		     newTitleDom.textContent = titleDom.textContent
		     const mainDomDiv = mainDom.querySelector('div')
		     mainDom.insertBefore(newTitleDom, mainDomDiv)
		     /**
		      * 清理 <noscript></noscript>
		      */
		     mainDom.querySelectorAll("noscript").forEach((dom) => dom.remove());
		     mainDom.querySelectorAll("img").forEach((dom) => {
		
		         if (dom.getAttribute("data-default-watermark-src")) {
		             dom.setAttribute("src", dom.getAttribute("data-default-watermark-src"));
		         }else if (dom.getAttribute("data-original")) {
		             dom.setAttribute("src", dom.getAttribute("data-original"));
		         }else if (dom.getAttribute("data-actualsrc")) {
		             dom.setAttribute("src", dom.getAttribute("data-actualsrc"));
		         }
		
		         if (dom.getAttribute("data-rawwidth")) {
		             dom.setAttribute("width", dom.getAttribute("data-rawwidth"));
		         }
		         if (dom.getAttribute("data-rawheight")) {
		             dom.setAttribute("height", dom.getAttribute("data-rawheight"));
		         }
		     });
		     document.querySelector('.ContentItem-meta').remove()
		     document.querySelector('.ContentItem-actions').remove()
		 }
		} catch (error) {
		 console.error('zhihu-photo-sharing error :' + error)
		}
		
		`, nil),
				chromedp.OuterHTML(`.QuestionAnswer-content`, &example),
				ShowCookies(),
			)
			if err != nil {
				log.Println(err)
			}
			res := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		 <meta charset="UTF-8">
		 <meta http-equiv="X-UA-Compatible" content="IE=edge">
		 <meta name="viewport" content="width=device-width, initial-scale=1.0">
		 <title>知乎</title>
		</head>
		<body>
		` + example + `
		</body>
		</html>
		`

			ctx.SetContentType("text/html")
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString(res)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusNotFound)

	})

	log.Fatal(fasthttp.ListenAndServe(":12345", router.Handler))

}

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

func ShowCookies() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}
		for i, cookie := range cookies {
			log.Printf("chrome cookie %d: %+v", i, cookie)
		}
		return nil
	})
}
