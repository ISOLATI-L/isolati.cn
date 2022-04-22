let EditormdView = editormd.markdownToHTML("editormd_view", {
    htmlDecode: true,
    emoji: false,
    taskList: true,
    tocm: true,
    tex: true,  // 默认不解析
    flowChart: true,  // 默认不解析
    sequenceDiagram: true,  // 默认不解析
});
