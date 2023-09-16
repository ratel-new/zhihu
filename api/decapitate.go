package api

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"os"
	"time"
	"zhihu/tool"
)

func DecapitateCore(ctx *fasthttp.RequestCtx) {

	var allocCtx context.Context
	var cancel context.CancelFunc

	allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), options...)

	logCtx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithDebugf(log.Printf),
	)

	defer cancel()

	chromeCtx, cancel := context.WithTimeout(logCtx, 20*time.Second)
	defer cancel()
	js, _ := os.ReadFile("executionJs/decapitate.js")
	var js_initialData_processed string
	err := chromedp.Run(chromeCtx,
		tool.SetCookie(),
		chromedp.Navigate(`https://www.zhihu.com/`),
		chromedp.Evaluate(string(js), nil),
		chromedp.WaitVisible("#js-initialData-processed"),
		chromedp.TextContent(`#js-initialData-processed`, &js_initialData_processed))
	if err != nil {
		ctx.Error(fmt.Sprint(err), 500)
		return
	}

	if js_initialData_processed == "" {
		ctx.Error(`can't find the data`, 500)
		return
	}

	fillStruct := tool.Fill{
		Title: `title`,
		Core:  template.HTML(js_initialData_processed),
	}

	err = tempProcessors.ExecuteTemplate(ctx, "decapitate.html", fillStruct)
	if err != nil {
		ctx.Error(`template exception`, 500)
		return
	}

	ctx.SetContentType("text/html")
	ctx.SetStatusCode(fasthttp.StatusOK)
}
