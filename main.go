package main

import (
	"context"
	"github.com/buaazp/fasthttprouter"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"strings"
	"time"
)

var (
	questionAnswerJs string
	zhuanLanJs       string
)

func init() {
	initCookies() //初始化cookies
	questionAnswerJs = initQuestionAnswerJs()
	zhuanLanJs = initZhuanLanJs()
}

func main() {
	router := fasthttprouter.New()

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

		options := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag(`headless`, false),
			chromedp.DisableGPU,
			chromedp.Flag(`disable-extensions`, false),
			chromedp.Flag(`enable-automation`, false),
			chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`),
		)

		if strings.HasPrefix(zhihuUrl, "https://www.zhihu.com") && strings.Contains(zhihuUrl, "question") && strings.Contains(zhihuUrl, "answer") {

			allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
			defer cancel()

			logCtx, cancel := chromedp.NewContext(
				allocCtx,
				chromedp.WithDebugf(log.Printf),
			)

			defer cancel()

			chromeCtx, cancel := context.WithTimeout(logCtx, 10*time.Second)
			defer cancel()

			var example string

			err = chromedp.Run(chromeCtx,
				SetCookie(),
				chromedp.Navigate(zhihuUrl),
				chromedp.Evaluate(questionAnswerJs, nil),
				chromedp.WaitVisible(`#zhihuDOMContentLoaded`),
				chromedp.OuterHTML(`.QuestionAnswer-content`, &example),
			)
			if err != nil {
				log.Println(err)
			}

			if len(example) == 0 {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
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

		if strings.HasPrefix(zhihuUrl, "https://zhuanlan.zhihu.com") && strings.Contains(zhihuUrl, "/p/") {
			allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
			defer cancel()

			logCtx, cancel := chromedp.NewContext(
				allocCtx,
				chromedp.WithDebugf(log.Printf),
			)

			defer cancel()

			chromeCtx, cancel := context.WithTimeout(logCtx, 10*time.Second)
			defer cancel()

			var example string
			var title string

			err = chromedp.Run(chromeCtx,
				SetCookie(),
				chromedp.Navigate(zhihuUrl),
				chromedp.Evaluate(zhuanLanJs, nil),
				chromedp.WaitVisible(`#zhihuDOMContentLoaded`),
				chromedp.OuterHTML(`.Post-RichTextContainer`, &example),
				chromedp.Text(`.Post-Title`, &title),
			)
			if err != nil {
				log.Println(err)
			}

			if len(example) == 0 {
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}

			res := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		 <meta charset="UTF-8">
		 <meta http-equiv="X-UA-Compatible" content="IE=edge">
		 <meta name="viewport" content="width=device-width, initial-scale=1.0">
		 <title>` + title + `</title>
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
