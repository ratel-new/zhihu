try {

    let DOMContentLoadedCode = document.createElement('div')
    DOMContentLoadedCode.id = 'zhihuDOMContentLoaded'
    document.body.appendChild(DOMContentLoadedCode)

    const mainDom = document.querySelector('.Post-Main.Post-NormalMain')
    const titleDom = document.querySelector('.QuestionHeader-title')

    Post-Title
    if (mainDom) {

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