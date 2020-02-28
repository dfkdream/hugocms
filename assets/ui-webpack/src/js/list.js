/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
const filepath = require('./filepath');

const fileList = require('./filelist');

const popup = require('./popup');

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

async function uploadFiles(path,files,callback){
    for (let i=0;i<files.length;i++) {
        let file = files[i];
        callback(i+1,files.length,file.name);
        let resp = await fetch(filepath.join(filepath.join("/admin/api/blob", path), filepath.clean("/" + file.name)), {
            method: "POST",
            body: file
        });

        if (!resp.ok) {
            let data = await resp.json();
            alert(`Upload error: ${file.name} ${data.code} ${data.message}`);
        }
    }
}

const path = location.pathname.replace(/^(\/admin\/list)/,"");

const rebuildModal = document.getElementById("rebuild-modal");
const rebuildModalTitle = document.getElementById("rebuild-modal-title");
const rebuildModalResult = document.getElementById("rebuild-modal-result");
const rebuildModalClose = document.getElementById("rebuild-modal-close");

const uploadModal = document.getElementById("upload-modal");
const uploadModalTitle = document.getElementById("upload-modal-title");
const uploadModalContent = document.getElementById("upload-modal-content");
const uploadModalUpload = document.getElementById("upload-modal-upload");
const uploadModalClose = document.getElementById("upload-modal-close");

const locationHeader = document.getElementById("location");
locationHeader.innerText=path;

const f = new fileList({
    path: path,
    target: document.getElementById("list-tbody"),
    onclickCallback: (file)=> {
        if (file.isDir) {
            f.navigate(filepath.join(f.path, file.name));
            locationHeader.innerText = f.path;
            history.pushState(f.path, "HugoCMS - " + f.path, filepath.join("/admin/list", f.path));
            document.title = "HugoCMS - " + f.path;
        } else {
            location.href = filepath.join(filepath.join(fileToEndpoint(file), f.path), file.name);
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
        }
    ]
});

document.getElementById("new-directory").onclick=()=>{
    popup.prompt(document.body,"Create Directory","Enter directory name")
        .then(inp=>{
            if (!inp) return;
            fetch(filepath.join(filepath.join("/admin/api/list",f.path),inp),{method:"POST"})
                .then(()=>f.reload());
        });
};

document.getElementById("new-post").onclick=()=>{
    popup.prompt(document.body,"New Post","Enter post filename")
        .then(inp=>{
            if (!inp) return;
            location.href = filepath.join(filepath.join("/admin/edit",f.path),inp);
        });
};

document.getElementById("rebuild").onclick=()=>{
    rebuildModalClose.setAttribute("disabled","");
    rebuildModal.setAttribute("class","modal active");
    rebuildModalTitle.innerText = "Rebuilding...";
    rebuildModalResult.innerHTML = '<progress class="progress" max="100"/>';

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

document.getElementById("upload-file").onclick=()=>{
    uploadModalUpload.setAttribute("disabled","");
    uploadModalClose.removeAttribute("disabled");
    uploadModalTitle.innerText = "Upload File";
    uploadModalContent.innerHTML="";
    let inputFile = document.createElement("input");
    inputFile.setAttribute("type","file");
    inputFile.setAttribute("class","form-input");
    inputFile.setAttribute("multiple","");
    uploadModalContent.appendChild(inputFile);
    uploadModal.setAttribute("class","modal active");

    inputFile.onchange = ()=>{
        uploadModalUpload.removeAttribute("disabled");
    };

    uploadModalUpload.onclick = ()=>{
        uploadModalTitle.innerText = "Uploading...";
        uploadModalContent.innerHTML="";
        uploadModalClose.setAttribute("disabled","");
        uploadModalUpload.setAttribute("disabled","");
        let statusProgress = document.createElement("progress");
        statusProgress.setAttribute("class","progress");
        statusProgress.setAttribute("max","100");
        uploadModalContent.appendChild(statusProgress);

        let statusP = document.createElement("p");
        uploadModalContent.appendChild(statusP);

        uploadFiles(f.path,inputFile.files,(idx,len,name)=>{
            statusProgress.value = idx/len*100;
            statusP.innerText = `(${idx}/${len}) ${name}`;
            statusP.style.margin=".8rem 0 0";
        })
            .then(()=>{
                uploadModal.setAttribute("class","modal");
                f.reload();
            });
    };
};

uploadModalClose.onclick = ()=>{
    uploadModal.setAttribute("class","modal");
};

history.replaceState(f.path,"HugoCMS - "+f.path,filepath.join("/admin/list",f.path));
document.title = "HugoCMS - "+f.path;

window.onpopstate = (e)=>{
    f.navigate(e.state);
    locationHeader.innerText = f.path;
};