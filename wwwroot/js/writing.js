let titleTextarea = document.querySelector("#title");
let contentTextarea = document.querySelector("#content");
let dialogMask = document.querySelector(".dialogMask");
let uploadDialog = document.querySelector(".dialog");
let titleInput = document.querySelector("#titleInput");
let information = document.querySelector("#information");

let markdown_editor = editormd("markdown_editor", {
    width: "100%",
    height: "100%",
    path: "/editormd/lib/",
    theme: "dark",
    previewTheme: "dark",
    editorTheme: "pastel-on-dark",
    markdown: "",
    placeholder: "typing...",
    codeFold: true,
    //syncScrolling : false,
    saveHTMLToTextarea: false,    // 保存 HTML 到 Textarea
    searchReplace: true,
    //watch : false,                // 关闭实时预览
    //htmlDecode: "style,script,iframe|on*",            // 开启 HTML 标签解析，为了安全性，默认不开启
    htmlDecode: true,
    //toolbar  : false,             //关闭工具栏
    //previewCodeHighlight : false, // 关闭预览 HTML 的代码块高亮，默认开启
    emoji: false,
    taskList: true,
    tocm: true,         // Using [TOCM]
    tex: true,                   // 开启科学公式TeX语言支持，默认关闭
    flowChart: true,             // 开启流程图支持，默认关闭
    sequenceDiagram: true,       // 开启时序/序列图支持，默认关闭,
    //dialogLockScreen : false,   // 设置弹出层对话框不锁屏，全局通用，默认为true
    //dialogShowMask : false,     // 设置弹出层对话框显示透明遮罩层，全局通用，默认为true
    //dialogDraggable : false,    // 设置弹出层对话框不可拖动，全局通用，默认为true
    //dialogMaskOpacity : 0.4,    // 设置透明遮罩层的透明度，全局通用，默认值为0.1
    //dialogMaskBgColor : "#000", // 设置透明遮罩层的背景颜色，全局通用，默认为#fff
    imageUpload: false,
    //imageFormats: ["jpg", "jpeg", "gif", "png", "bmp", "webp"],
    //imageUploadURL: "./php/upload.php",
    toolbarIcons: [
        "undo",
        "redo",
        "|",
        "bold",
        "del",
        "italic",
        "quote",
        "ucwords",
        "uppercase",
        "lowercase",
        "|",
        "h1",
        "h2",
        "h3",
        "h4",
        "h5",
        "h6",
        "|",
        "list-ul",
        "list-ol",
        "hr",
        "|",
        "link",
        "reference-link",
        "image",
        "code",
        "preformatted-text",
        "code-block",
        "table",
        "datetime",
        // "emoji",
        "html-entities",
        "pagebreak",
        "|",
        "goto-line", // fa-crosshairs
        "watch",
        "preview",
        "fullscreen",
        "clear",
        "search",
        "|",
        "upload",
        // "|",
        // "help",
        // "info",
    ],
    toolbarIconTexts: {
        upload: "上传"
    },
    toolbarHandlers: {
        upload: function (cm, icon, cursor, selection) {
            let title = titleTextarea.value;
            titleInput.value = title;
            information.innerHTML = "";
            dialogMask.style.display = "block";
            uploadDialog.style.display = "block";
        }
    },
    onload: function () {
        console.log('onload', this);
        //this.fullscreen();
        //this.unwatch();
        //this.watch().fullscreen();

        //this.setMarkdown("#PHP");
        //this.width("100%");
        //this.height(480);
        //this.resize("100%", 640);
    }
});

function cancelUpload() {
    dialogMask.style.display = "none";
    uploadDialog.style.display = "none";
}

function upload() {
    let title = titleInput.value;
    if (title === "") {
        information.innerHTML = "请输入标题";
        return
    }
    let content = contentTextarea.value;
    let paragraphs = {
        "title": title,
        "content": content,
    };
    information.innerHTML = "上传中";
    post("/admin/api/paragraph" + window.location.search, JSON.stringify(paragraphs), "application/json").then(
        function (res) {
            console.log(res);
            if (res.status === 200) {
                information.innerHTML = "上传成功";
            } else {
                information.innerHTML = "上传失败";
            }
        },
        function (res) {
            console.log(res);
            information.innerHTML = "上传失败";
        });
}
