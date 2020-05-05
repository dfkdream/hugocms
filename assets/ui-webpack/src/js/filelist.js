/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
const filepath = require('./filepath');

import popup from "./popup";

require('../css/filelist.css');

function fileToIcon(file){
    let i = document.createElement("i");
    if (file.isDir){
        i.setAttribute("class","fa fa-folder");
        return i
    }
    switch (filepath.ext(file.name)){
        case "md":
            i.setAttribute("class","fab fa-markdown");
            return i;
        case "html":
            i.setAttribute("class", "fa fa-file-code");
            return i;
        default:
            i.setAttribute("class","fa fa-file");
            return i;
    }
}

class fileList{
    constructor(config,t){
        this.path = config.path?config.path:"/";
        this.onclickCallback = config.onclickCallback?
            config.onclickCallback:file=>this.navigate(filepath.clean(filepath.join(this.path,file.name)));
        this.endpoint = config.endpoint?config.endpoint:"/admin/api/list";
        this.target = config.target;
        this.actions = config.actions;

        this.t = t;

        this.popup = new popup(t);

        this.build()
    }

    navigate(path){
        this.path=filepath.clean(path);
        this.build();
    }

    reload(){
        this.build();
    }

    build(){
        fetch(filepath.join(this.endpoint,this.path))
            .then(resp=>{
                if (!resp.ok) return Promise.reject(resp.json);
                return resp.json();
            })
            .then(data=>{
                if (this.path!=="/"){
                    data=[{name:"..", isDir:true}].concat(data);
                }

                let fragment = document.createDocumentFragment();
                data.forEach(file=>{
                    let row = document.createElement("tr");
                    let fileLink = document.createElement("a");
                    fileLink.innerText = file.name + (file.isDir?"/":"");
                    fileLink.href = filepath.join("./",file.name);
                    fileLink.onclick = (e)=>{
                        e.preventDefault();
                        this.onclickCallback(file);
                    };
                    let fnCell = row.insertCell();
                    fnCell.appendChild(fileToIcon(file));
                    fnCell.appendChild(fileLink);
                    row.insertCell().appendChild(document.createTextNode(file.size?file.size:""));
                    row.insertCell().appendChild(document.createTextNode(file.mode?file.mode:""));
                    row.insertCell().appendChild(document.createTextNode(file.modTime?file.modTime:""));

                    if (this.actions) {
                        let actionsCell = row.insertCell();
                        if (file.name!=="..") {
                            this.actions.forEach(action => {
                                let icon = document.createElement("i");
                                icon.setAttribute("class", "tooltip tooltip-bottom "+action.icon);
                                icon.dataset.tooltip = action.tooltip;
                                icon.onclick = () => {
                                    action.callback(file)
                                };
                                actionsCell.appendChild(icon);
                            });
                        }
                    }

                    fragment.appendChild(row);
                });

                this.target.innerHTML="";
                this.target.appendChild(fragment);
            })
            .catch(err=>{
                if (err instanceof Promise) {
                    err.then(json => {
                        if (json.code === 404) location.href = "/admin/list/";
                        else if (json.code === 403) location.reload();
                        else this.popup.alert(document.body, this.t("error"), `${json.code} ${json.message}`);
                    })
                        .catch(() => this.popup.alert(document.body, this.t("error"), this.t("errUnknown")));
                }else{
                    console.log(err);
                    this.popup.alert(document.body, this.t("error"), this.t("errUnknown"));
                }
            });
    }
}

export default fileList;