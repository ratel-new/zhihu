try {
    let DOMContentLoadedCode = document.createElement('div')
    DOMContentLoadedCode.id = 'zhihuDOMContentLoaded'
    if (location.href.startsWith(`https://www.zhihu.com`)&& location.href.includes(`question`)&& location.href.includes(`answer`)){
        const mainDom = document.querySelector('div.QuestionAnswer-content')
        const titleDom = document.querySelector('.QuestionHeader-title')
        if (mainDom && titleDom) {
            const newTitleDom = titleDom.cloneNode()
            newTitleDom.className = 'QuestionHeader-title'
            newTitleDom.textContent = titleDom.textContent
            const mainDomDiv = mainDom.querySelector('div')
            mainDom.insertBefore(newTitleDom, mainDomDiv)
            document.querySelector('.ContentItem-meta').remove()
            document.querySelector('.ContentItem-actions').remove()
        }
    }else {
        let DOMContentLoadedCode = document.createElement('div')
        DOMContentLoadedCode.id = 'zhihuDOMContentLoaded'
        const mainDom = document.querySelector('.Post-Main.Post-NormalMain .Post-RichTextContainer')
        const titleDom = document.querySelector('.Post-Title')
        if (mainDom) {
            const newTitleDom = titleDom.cloneNode()
            newTitleDom.textContent = titleDom.textContent
            const mainDomDiv = mainDom.querySelector('div')
            mainDom.insertBefore(newTitleDom, mainDomDiv)
        }
    }
    let yanxuanStyle = ''
    document.querySelectorAll('style').forEach(s => s.outerText.indexOf(`DynamicFonts`) !== -1 && (yanxuanStyle += s.outerText + '\n'))
    DOMContentLoadedCode.innerText = yanxuanStyle;
    document.body.appendChild(DOMContentLoadedCode)
} catch (error) {
    console.error('zhihu error :' + error)
}