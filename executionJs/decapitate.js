try {
    let DOMContentLoadedCode = document.createElement('div')
    DOMContentLoadedCode.id = 'js-initialData-processed'
    let js_initialData = JSON.parse(document.querySelector("#js-initialData").textContent);
    let js_initialData_processed = [];

    for (const key in js_initialData['initialState']['entities']['answers']) js_initialData_processed.push({
        url: `https://www.zhihu.com/question/${js_initialData['initialState']['entities']['answers'][key]['question']['id']}/answer/${js_initialData['initialState']['entities']['answers'][key]['id']}`,
        content: js_initialData['initialState']['entities']['answers'][key]['content'],
        title: js_initialData['initialState']['entities']['answers'][key]['question']['title']
    });
    //专栏,暂时取消,软广太多
    // for (const key in js_initialData['initialState']['entities']['articles']) js_initialData_processed.push({
    //     url: `ttps://zhuanlan.zhihu.com/p/${js_initialData['initialState']['entities']['articles'][key]['id']}`,
    //     content: js_initialData['initialState']['entities']['articles'][key]['content'],
    //     title: js_initialData['initialState']['entities']['articles'][key]['title']
    // });
    for (const jsInitialDataProcessed of js_initialData_processed)  DOMContentLoadedCode.innerText  +=  `<div class="root-card card container heti">
        <div class="body">
            <div class="card-header-title title-centered">
                <a  href="${jsInitialDataProcessed.url}">${jsInitialDataProcessed.title}</a>
            </div>
            <content class="card-content">
                <div class="content">${jsInitialDataProcessed.content}</div>
            </content>
        </div>
    </div>\n`
    document.body.appendChild(DOMContentLoadedCode)
} catch (error) {
    console.error('zhihu-photo-sharing error :' + error)
}