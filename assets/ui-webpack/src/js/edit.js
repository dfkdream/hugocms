/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */

require('../css/edit.css');
require('../css/style.css');

require('spectre.css/dist/spectre.min.css');

require('../css/all.min.css');

import Editor from "@toast-ui/editor";
import "codemirror/lib/codemirror.css";
import "@toast-ui/editor/dist/toastui-editor.css";

import colorSyntax from "@toast-ui/editor-plugin-color-syntax";
import "tui-color-picker/dist/tui-color-picker.css"

import codeSyntaxHighlight from "@toast-ui/editor-plugin-code-syntax-highlight";
import hljs from "highlight.js";
import 'highlight.js/styles/github.css';
import popups from "./popup";
import filesPopup from "./files-popup";

import publish from "./publish";

import {i18n} from "../i18n";

const filepath = require('./filepath');

const fpath = require('path');

i18n().then(t => {
    const popup = new popups(t);

    const path = location.pathname.replace(/^(\/admin\/edit)/, "");
    const endpoint = filepath.join("/admin/api/post", path);

    function files(editor) {
        const addFileMarkdown = () => {
            filesPopup(filepath.clean(filepath.join(path, "..")),t)
                .then(res => {
                    if (res !== false) {
                        const fn = fpath.relative(/^_?index(\..+)?\.(md|html|htm)$/i.test(fpath.basename(path)) ? fpath.dirname(path) : path, res);
                        switch (filepath.ext(fn).toLowerCase()) {
                            case "jpg":
                            case "jpeg":
                            case "png":
                            case "bmp":
                            case "gif":
                            case "tiff":
                            case "svg":
                            case "webp":
                                editor.insertText(`\n![${fn.split("/").pop()}](${fn})\n`);
                                break;
                            default:
                                editor.insertText(`\n[${fn.split("/").pop()}](${fn})\n`);
                                break;
                        }
                    }
                })
        };

        if (!editor.isViewer()&&editor.getUI().name==="default"){
            const toolbar = editor.getUI().getToolbar();
            toolbar.addItem("divider");
            const buttonEl = document.createElement("button");
            buttonEl.className="tui-hugocms-file";

            toolbar.addItem({
                type: "button",
                options: {
                    command: 'fileClicked',
                    tooltip: t("addFile"),
                    el: buttonEl
                }
            });

            editor.addCommand('markdown',{
                name: 'fileClicked',
                exec(){
                    addFileMarkdown();
                }
            });

            editor.addCommand('wysiwyg',{
                name: 'fileClicked',
                exec(){
                    popup.alert(document.body, t("error"), t("errFileLinkMarkdownOnly"))
                }
            });
        }
    }

    const editor = new Editor({
        el: document.getElementById("editor"),
        initialEditType: 'markdown',
        previewStyle: 'vertical',
        height: '80vh',
        //exts: ['colorSyntax', 'scrollSync', 'files'],
        plugins: [colorSyntax, [codeSyntaxHighlight, {hljs}],files],
        usageStatistics: false
    });

    class optionsList {
        constructor(target) {
            this.target = target;
            this.list = [];
        }

        append(value) {
            this.list.push(value);
            let op = document.createElement("option");
            op.text = value;
            this.target.add(op);
        }

        delete(idx) {
            this.list.splice(idx, 1);
            this.target.remove(idx);
        }

        fromList(list) {
            list.forEach(i => this.append(i));
        }

        selectedIndex() {
            return this.target.selectedIndex;
        }
    }

    document.getElementById("location").innerText = path;

    const title = document.getElementById("title");
    const subtitle = document.getElementById("subtitle");
    const date = document.getElementById("date");
    const author = document.getElementById("author");
    const attachments = new optionsList(document.getElementById("attachments"));
    const showReadingTime = document.getElementById("show-reading-time");
    const showLanguages = document.getElementById("show-languages");
    const showAuthor = document.getElementById("show-author");
    const showDate = document.getElementById("show-date");

    document.getElementById("date-now").onclick = () => {
        date.valueAsDate = new Date();
    };

    document.getElementById("attachment-add").onclick = () => {
        filesPopup(filepath.clean(filepath.join(path, "..")), t)
            .then(res => {
                if (res !== false) {
                    attachments.append(
                        fpath.relative(
                            /^_?index(\..+)?\.(md|html|htm)$/i.test(fpath.basename(path)) ? fpath.dirname(path) : path,
                            res));
                }
            })
    };

    document.getElementById("attachment-delete").onclick = () => {
        if (attachments.list.length > 0) {
            attachments.delete(attachments.selectedIndex());
        }
    };

    document.getElementById("raw").setAttribute("href", filepath.join("/admin/api/blob", path));

    let publicPath = fpath.dirname(path);
    let filename = fpath.basename(path);
    if (!/^_?index(\..+)?\.(md|html|htm)$/i.test(filename)){
        publicPath = fpath.join(publicPath,filename.split(".")[0]);
    }
    publicPath=fpath.join("/",publicPath);

    document.getElementById("page").setAttribute("href",publicPath);

    document.getElementById("save").onclick = () => {
        fetch(endpoint, {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
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
            })
        })
            .then((resp) => {
                if (!resp.ok) return Promise.reject(resp.json());
                popup.alert(document.body, t("success"), t("saved"));
            })
            .catch(err => {
                if (err instanceof Promise) {
                    err.then(json => {
                        popup.alert(document.body, t("error"), `${json.code} ${json.message}`);
                    })
                        .catch(() => {
                            popup.alert(document.body, t("error"), t("errUnknown"));
                        })
                } else {
                    console.log(err);
                    popup.alert(document.body, t("error"), t("errUnknown"));
                }
            })
    };

    document.getElementById("publish").onclick = () => {
        publish(t);
    };

    fetch(endpoint)
        .then(resp => {
            if (!resp.ok) return Promise.reject(resp.json());
            return resp.json()
        })
        .then(data => {
            title.value = data.frontMatter.title;
            subtitle.value = data.frontMatter.subtitle;
            date.valueAsDate = new Date(data.frontMatter.date);
            author.value = data.frontMatter.author;
            attachments.fromList(data.frontMatter.attachments);
            showReadingTime.checked = data.frontMatter.showReadingTime;
            showLanguages.checked = data.frontMatter.showLanguages;
            showAuthor.checked = data.frontMatter.showAuthor;
            showDate.checked = data.frontMatter.showDate;
            editor.setMarkdown(data.body);
        })
        .catch(err => {
            if (err instanceof Promise) {
                err.then(json => {
                    if (json.code === 404) {
                        fetch("/admin/api/whoami")
                            .then(res => res.json())
                            .then(json => {
                                author.value = json.username;
                            });
                        return;
                    }
                    popup.alert(document.body, t("error"), `${json.code} ${json.message}`);
                })
                    .catch(() => {
                        popup.alert(document.body, t("error"), t("errUnknown"));
                    });
            } else {
                console.log(err);
                popup.alert(document.body, t("error"), t("errUnknown"));
            }
        });

});