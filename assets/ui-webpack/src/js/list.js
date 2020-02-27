/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
const filepath = require('./filepath');

require('../css/style.css');
require('../css/list.css');

require('spectre.css/dist/spectre.min.css');
require('spectre.css/dist/spectre-exp.min.css');

require('../css/all.min.css');

function fileToEndpoint(file){
   if (file.isDir) return "/admin/list";
   switch (filepath.ext(file.name)){
       case "md": case "html":
           return "/admin/edit";
       default:
           return "/admin/api/blob";
   }
}

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

window.onload=()=>{
    const path = location.pathname.replace(/^(\/admin\/list)/,"");
    const endpoint = filepath.join("/admin/api/list",path);

    const listBody = document.getElementById("list-tbody");

    const rebuildModal = document.getElementById("rebuild-modal");
    const rebuildModalTitle = document.getElementById("rebuild-modal-title");
    const rebuildModalResult = document.getElementById("rebuild-modal-result");
    const rebuildModalClose = document.getElementById("rebuild-modal-close");

    document.getElementById("location").innerText=path;

    document.getElementById("new-directory").onclick=()=>{
        let inp = prompt("Enter directory name");
        if (!inp) return;
        fetch(filepath.join(filepath.join("/admin/api/list",path),inp),{method:"POST"})
            .then(()=>location.reload());
    };

    document.getElementById("new-post").onclick=()=>{
        let inp = prompt("Enter post filename");
        if (!inp) return;
        location.href = filepath.join(filepath.join("/admin/edit",path),inp);
    };

    document.getElementById("rebuild").onclick=()=>{
        rebuildModalClose.setAttribute("disabled","");
        rebuildModal.setAttribute("class","modal active");
        rebuildModalTitle.innerText = "Rebuilding...";
        rebuildModalResult.innerHTML="";
        let pg = document.createElement("progress");
        pg.setAttribute("class","progress");
        pg.setAttribute("max","100");
        rebuildModalResult.appendChild(pg);

        fetch("/admin/api/build",{method:"POST"})
            .then(resp=>resp.json())
            .then(data=>{
                rebuildModalTitle.innerText = (data.code===0?"Done":"Error");
                rebuildModalResult.innerHTML="";
                let res = document.createElement("pre");
                res.setAttribute("class","console");
                res.innerText = data.result;
                rebuildModalResult.appendChild(res);
                rebuildModalClose.removeAttribute("disabled");
            });
    };

    rebuildModalClose.onclick = ()=>{
        rebuildModal.setAttribute("class","modal");
    };

    fetch(endpoint)
        .then(resp=>resp.json())
        .then(data=>{
            if (path!=="/"){
                data=[{name:"..", isDir:true}].concat(data);
            }
            listBody.innerHTML="";
            let tb = document.createDocumentFragment();
            data.forEach(file=>{
                let row = document.createElement("tr");
                let fileLink = document.createElement("a");
                fileLink.innerText = file.name + (file.isDir?"/":"");
                fileLink.href = fileToEndpoint(file)+filepath.join(path,file.name);
                let fnCell = row.insertCell();
                fnCell.appendChild(fileToIcon(file));
                fnCell.appendChild(fileLink);
                row.insertCell().appendChild(document.createTextNode(file.size?file.size:""));
                row.insertCell().appendChild(document.createTextNode(file.mode?file.mode:""));
                row.insertCell().appendChild(document.createTextNode(file.modTime?file.modTime:""));
                tb.appendChild(row);
            });
            listBody.appendChild(tb);
        })
};