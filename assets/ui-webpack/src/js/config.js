/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
window.onload=()=>{
    const instance = new tui.Editor({
        el: document.getElementById("editorSection"),
        initialEditType: 'markdown',
        previewStyle: 'vertical',
        height: '80vh',
        usageStatistics: false
    });
};