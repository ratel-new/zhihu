package api

import (
	"github.com/chromedp/chromedp"
	"html/template"
)

var (
	tempProcessors = template.Must(template.ParseGlob("templates/*.html"))
	options        []chromedp.ExecAllocatorOption
)

func init() {
	options = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag(`headless`, true),
		chromedp.DisableGPU,
		chromedp.Flag(`enable-automation`, false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36`),
	)

}
