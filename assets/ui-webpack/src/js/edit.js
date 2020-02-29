/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
require('tui-editor/dist/tui-editor-extScrollSync');
require('tui-editor/dist/tui-editor-extColorSyntax');

require('codemirror/lib/codemirror.css');
require('tui-editor/dist/tui-editor.min.css');
require('tui-editor/dist/tui-editor-contents.min.css');
require('tui-color-picker/dist/tui-color-picker.min.css');
require('highlight.js/styles/github.css');

require('../css/edit.css');
require('../css/style.css');

require('spectre.css/dist/spectre.min.css');

require('../css/all.min.css');

const Editor = require('tui-editor');

const filepath = require('./filepath');

const popup = require('./popup');

const fileList = require('./filelist');

const fpath = require('path');

const editor = new Editor({
    el: document.getElementById("editor"),
    initialEditType: 'markdown',
    previewStyle: 'vertical',
    height: '80vh',
    exts: ['colorSyntax','scrollSync'],
    usageStatistics: false
});

class optionsList{
    constructor(target){
        this.target = target;
        this.list = [];
    }

    append(value){
        this.list.push(value);
        let op = document.createElement("option");
        op.text = value;
        this.target.add(op);
    }

    delete(idx){
        this.list.splice(idx,1);
        this.target.remove(idx);
    }

    fromList(list){
        list.forEach(i=>this.append(i));
    }

    selectedIndex(){
        return this.target.selectedIndex;
    }
}

const path = location.pathname.replace(/^(\/admin\/edit)/,"");
const endpoint = filepath.join("/admin/api/post", path);

document.getElementById("location").innerText=path;

const title = document.getElementById("title");
const subtitle = document.getElementById("subtitle");
const date = document.getElementById("date");
const author = document.getElementById("author");
const attachments = new optionsList(document.getElementById("attachments"));
const showReadingTime = document.getElementById("show-reading-time");
const showLanguages = document.getElementById("show-languages");
const showAuthor = document.getElementById("show-author");
const showDate = document.getElementById("show-date");

const attachmentModal = document.getElementById("attachment-modal");
const attachmentModalPath = document.getElementById("attachment-modal-path");

document.getElementById("date-now").onclick=()=>{
    date.valueAsDate = new Date();
};

document.getElementById("attachment-add").onclick = ()=>{
    attachmentModal.setAttribute("class","modal active");
};

document.getElementById("attachment-modal-close").onclick = ()=>{
    attachmentModal.setAttribute("class","modal");
};

document.getElementById("attachment-modal-overlay").onclick = ()=>{
    attachmentModal.setAttribute("class","modal");
};

document.getElementById("attachment-delete").onclick = ()=>{
    if (attachments.list.length>0){
        attachments.delete(attachments.selectedIndex());
    }
};

document.getElementById("raw").setAttribute("href",filepath.join("/admin/api/blob", path));

document.getElementById("save").onclick=()=>{
    fetch(endpoint,{
        method: "POST",
        headers:{
            'Content-Type': 'application/json'
        },
        body:JSON.stringify({
            frontMatter: {
                title: title.value,
                subtitle: subtitle.value,
                date: date.valueAsDate,
                author: author.value,
                attachments: attachments.list,
                showReadingTime: showReadingTime.checked,
                showLanguages: showLanguages.checked,
                showAuthor: showAuthor.checked,
                showDate: showDate.checked
            },
            body: editor.getMarkdown()
        })})
        .then((resp)=>{
            if (resp.ok){
                popup.alert(document.body,"Save","Saved.");
            }else{
                resp.json()
                    .then(err=>{
                        popup.alert(document.body,"Save Error",`${err.code} ${err.message}`);
                    })
            }
        });
};

fetch(endpoint)
    .then(resp=>resp.json())
    .then(data=>{
        title.value = data.frontMatter.title;
        subtitle.value = data.frontMatter.subtitle;
        date.valueAsDate = new Date(data.frontMatter.date);
        author.value = data.frontMatter.author;
        attachments.fromList(data.frontMatter.attachments);
        showReadingTime.checked=data.frontMatter.showReadingTime;
        showLanguages.checked=data.frontMatter.showLanguages;
        showAuthor.checked=data.frontMatter.showAuthor;
        showDate.checked=data.frontMatter.showDate;
        editor.setMarkdown(data.body);
    });


const f = new fileList({
    path: filepath.clean(filepath.join(path,"..")),
    target: document.getElementById("attachment-modal-list"),
    onclickCallback: (file)=>{
        if (file.isDir){
            f.navigate(filepath.join(f.path,file.name));
            attachmentModalPath.innerText = f.path;
        }else{
            switch(filepath.ext(file.name)){
                case "md": case "html":
                    popup.alert(document.body,"Error","Markdown or HTML files cannot be attached");
                    break;
                default:
                    attachments.append(
                        fpath.relative(
                            /^_?index(\..+)?\.(md|html|htm)$/i.test(fpath.basename(path))?fpath.dirname(path):path,
                            fpath.join(f.path,file.name)));
                    attachmentModal.setAttribute("class","modal");
            }
        }
    },
    actions: [
        {
            icon: "fas fa-trash-alt",
            tooltip: "Delete",
            callback:file=>{
                popup.confirm(document.body,"Confirm Delete", `Delete ${file.name}?`)
                    .then(confirm=>{
                        if (confirm){
                            if (file.isDir){
                                fetch(filepath.join("/admin/api/list/",filepath.join(f.path,file.name)),{
                                    method: "DELETE",
                                })
                                    .then(res=>{
                                        if (!res.ok){
                                            alert("Error delete directory");
                                        }
                                        f.reload();
                                    })
                            }else{
                                fetch(filepath.join("/admin/api/blob/",filepath.join(f.path,file.name)),{
                                    method: "DELETE",
                                })
                                    .then(res=>{
                                        if (!res.ok){
                                            alert("Error rename file");
                                        }
                                        f.reload();
                                    })
                            }
                        }
                    });
            }
        }, {
            icon: "fas fa-edit",
            tooltip: "Rename",
            callback:file=>{
                popup.prompt(document.body,"Rename",`Rename ${file.name} to`)
                    .then(fn=>{
                        if (fn){
                            if (file.isDir){
                                fetch(filepath.join("/admin/api/list/",filepath.join(f.path,file.name)),{
                                    method: "PUT",
                                    body: JSON.stringify(filepath.join(f.path,fn))
                                })
                                    .then(res=>{
                                        if (!res.ok){
                                            alert("Error rename directory");
                                        }
                                        f.reload();
                                    })
                            }else{
                                fetch(filepath.join("/admin/api/blob/",filepath.join(f.path,file.name)),{
                                    method: "PUT",
                                    body: JSON.stringify(filepath.join(f.path,fn))
                                })
                                    .then(res=>{
                                        if (!res.ok){
                                            alert("Error rename file");
                                        }
                                        f.reload();
                                    })
                            }
                        }

                    });
            }
        }]
});

document.getElementById("upload-file").onclick = ()=>{
    popup.upload(document.body,filepath.join("/admin/api/blob",f.path))
        .then(()=>f.reload());
};

attachmentModalPath.innerText = f.path;