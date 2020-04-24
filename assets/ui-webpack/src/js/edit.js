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

const fpath = require('path');

const filesPopup = require('./files-popup');

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

document.getElementById("date-now").onclick=()=>{
    date.valueAsDate = new Date();
};

document.getElementById("attachment-add").onclick = ()=>{
    filesPopup(filepath.clean(filepath.join(path,"..")))
        .then(res=>{
            if (res!==false) {
                attachments.append(
                    fpath.relative(
                        /^_?index(\..+)?\.(md|html|htm)$/i.test(fpath.basename(path))?fpath.dirname(path):path,
                        res));
            }
        })
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
            if (!resp.ok) return Promise.reject(resp.json());
            popup.alert(document.body,"Save","Saved.");
        })
        .catch(err=>{
            if (err instanceof Promise){
                err.then(json=>{
                    popup.alert(document.body,"Save Error",`${json.code} ${json.message}`);
                })
                    .catch(()=>{
                        popup.alert(document.body,"Save Error","Unknown error occurred.");
                    })
            }else{
                console.log(err);
                popup.alert(document.body,"Save Error","Unknown error occurred.");
            }
        })
};

fetch(endpoint)
    .then(resp=>{
        if (!resp.ok) return Promise.reject(resp.json());
        return resp.json()
    })
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
    })
    .catch(err=>{
        if (err instanceof Promise){
            err.then(json=>{
                if (json.code===404) {
                    fetch("/admin/api/whoami")
                        .then(res=>res.json())
                        .then(json=>{
                            author.value = json.username;
                        });
                    return;
                }
                popup.alert(document.body,"Error",`${json.code} ${json.message}`);
            })
                .catch(()=>{
                    popup.alert(document.body,"Error",`Unknown error occurred. Please reload.`);
                });
        }else{
            console.log(err);
            popup.alert(document.body,"Error",`Unknown error occurred. Please reload.`);
        }
    });

