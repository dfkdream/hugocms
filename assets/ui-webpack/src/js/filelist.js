/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
const filepath = require('./filepath');

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

module.exports = class fileList{
    constructor(config){
        this.path = config.path?config.path:"/";
        this.onclickCallback = config.onclickCallback?
            config.onclickCallback:file=>this.navigate(filepath.clean(filepath.join(this.path,file.name)));
        this.endpoint = config.endpoint?config.endpoint:"/admin/api/list";
        this.target = config.target;
        this.actions = config.actions;

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
            .then(resp=>resp.json())
            .then(data=>{
                if (data.code) {
                    alert(`${data.code} ${data.message}`);
                    return
                }
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
            });
    }
};