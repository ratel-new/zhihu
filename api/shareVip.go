package api

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
	"zhihu/tool"
)

func ShareVipCore(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.RequestURI()
	query, err := url.Parse(string(body))
	if err != nil {
		ctx.Error("Request Body Failed to parse", 400)
		return
	}
	dealUrl, err := url.Parse(query.Query().Get("url"))
	if err != nil {
		ctx.Error("url parameter abnormality", 400)
		return
	}
	IncURL := dealUrl.Scheme + "://" + dealUrl.Host + dealUrl.Path
	value, b := tool.GetHtmlCache(`vip` + IncURL)
	if b {
		fillStruct := tool.Fill{
			Title: value.Title,
			Core:  value.Core,
			Style: value.Style,
		}
		err := tempProcessors.ExecuteTemplate(ctx, "index.html", fillStruct)
		if err != nil {
			ctx.Error(`template exception`, 500)
			return
		}
		ctx.SetContentType("text/html")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
	if strings.HasPrefix(IncURL, "https://www.zhihu.com") && strings.Contains(IncURL, "question") && strings.Contains(IncURL, "answer") {
		browserExecutionVip(IncURL, `.QuestionAnswer-content`, `.QuestionHeader-title`, ctx)
	} else if strings.HasPrefix(IncURL, "https://zhuanlan.zhihu.com") && strings.Contains(IncURL, "/p/") {
		browserExecutionVip(IncURL, `.Post-RichTextContainer`, `.Post-Title`, ctx)
	}
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func browserExecutionVip(IncURL, contentDomSel, titleDomSel string, ctx *fasthttp.RequestCtx) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	logCtx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithDebugf(log.Printf),
	)

	defer cancel()

	chromeCtx, cancel := context.WithTimeout(logCtx, 10*time.Second)
	defer cancel()

	js, _ := os.ReadFile("executionJs/share.js")
	var content, zhihuDOMContentLoaded, title string

	err := chromedp.Run(chromeCtx,
		tool.SetCookie(),
		chromedp.Navigate(IncURL),
		chromedp.Evaluate(string(js), nil),
		chromedp.WaitVisible(`#zhihuDOMContentLoaded`),
		chromedp.TextContent(`#zhihuDOMContentLoaded`, &zhihuDOMContentLoaded),
		chromedp.Text(titleDomSel, &title),
		chromedp.OuterHTML(contentDomSel, &content),
	)
	if err != nil {
		ctx.Error(`run err`, 500)
		return
	}

	fillStruct := tool.Fill{
		Title: title,
		Core:  template.HTML(content),
		Style: template.CSS(zhihuDOMContentLoaded),
	}

	tool.SetCache("vip"+IncURL, &fillStruct)

	err = tempProcessors.ExecuteTemplate(ctx, "index.html", fillStruct)

	if err != nil {
		ctx.Error(`template exception`, 500)
		return
	}

	ctx.SetContentType("text/html")
	ctx.SetStatusCode(fasthttp.StatusOK)

}
