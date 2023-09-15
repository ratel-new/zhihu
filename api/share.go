package api

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"net/url"
	"path"
	"strings"
	"time"
	jsonparser "zhihu/json"
	"zhihu/tool"
)

var (
	options []chromedp.ExecAllocatorOption
)

func init() {
	options = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag(`headless`, true),
		chromedp.DisableGPU,
		chromedp.Flag(`enable-automation`, false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`),
	)
}

func ShareCore(ctx *fasthttp.RequestCtx) {
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

	if strings.HasPrefix(IncURL, "https://www.zhihu.com") && strings.Contains(IncURL, "question") && strings.Contains(IncURL, "answer") {
		browserExecution(
			IncURL,
			ctx,
			[]string{"initialState", "entities", "answers", path.Base(IncURL), "content"},
			[]string{"initialState", "entities", "answers", path.Base(IncURL), "question", "title"},
			[]chromedp.Action{
				tool.SetCookie(),
				chromedp.Navigate(IncURL),
			}...,
		)
		return
	} else if strings.HasPrefix(IncURL, "https://zhuanlan.zhihu.com") && strings.Contains(IncURL, "/p/") {
		browserExecution(
			IncURL,
			ctx,
			[]string{"initialState", "entities", "articles", path.Base(IncURL), "content"},
			[]string{"initialState", "entities", "articles", path.Base(IncURL), "title"},
			[]chromedp.Action{
				tool.SetCookie(),
				chromedp.Navigate(IncURL),
			}...,
		)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func browserExecution(IncURL string, ctx *fasthttp.RequestCtx, framework []string, titleFramework []string, action ...chromedp.Action) {
	value, b := tool.GetHtmlCache(IncURL)

	if b {
		fillStruct := tool.Fill{
			Title: value.Title,
			Core:  value.Core,
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

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	logCtx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithDebugf(log.Printf),
	)

	defer cancel()

	chromeCtx, cancel := context.WithTimeout(logCtx, 10*time.Second)
	defer cancel()

	var js_initialData string
	action = append(action, []chromedp.Action{
		chromedp.TextContent(`#js-initialData`, &js_initialData),
	}...)

	err := chromedp.Run(chromeCtx, action...)
	if err != nil {
		ctx.Error(`run err`, 500)
		return
	}

	//Get body
	body, err := jsonparser.GetString([]byte(js_initialData), framework...)
	if err != nil {
		ctx.Error(`Get body err`, 500)
		return
	}
	//Get title
	title, err := jsonparser.GetString([]byte(js_initialData), titleFramework...)
	if err != nil {
		ctx.Error(`Get body title`, 500)
		return
	}

	fillStruct := tool.Fill{
		Title: title,
		Core:  template.HTML(body),
	}

	tool.SetCache(IncURL, &fillStruct)

	err = tempProcessors.ExecuteTemplate(ctx, "index.html", fillStruct)

	if err != nil {
		ctx.Error(`template exception`, 500)
		return
	}

	ctx.SetContentType("text/html")
	ctx.SetStatusCode(fasthttp.StatusOK)
	return
}
