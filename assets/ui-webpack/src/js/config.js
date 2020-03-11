require("spectre.css/dist/spectre.css");

require("../css/all.min.css");

require("../css/style.css");
require("../css/config.css");

const popup = require("./popup");

const ace = require("ace-builds/src-noconflict/ace");

require("ace-builds/src-noconflict/mode-yaml");

const configCode = document.getElementById("config-code");

let editor = ace.edit(configCode,{
    mode: "ace/mode/yaml"
});

fetch("/admin/api/config")
    .then(res=>res.text())
    .then(text=>editor.setValue(text, -1))
    .catch(err=>console.log(err));

document.getElementById("save").onclick=()=>{
    fetch("/admin/api/config",{
        method: "POST",
        body: editor.getValue()
    })
        .then(res=>{
            if (!res.ok) {
                res.json()
                    .then(data => popup.alert(document.body,"Config Upload",`Upload Failed: ${data.code} ${data.message}`))
                    .catch(() => popup.alert(document.body,"Config Upload","Upload Failed"));
                return;
            }
            popup.alert(document.body,"Config Upload","Upload Succeed");
        })
        .catch(()=>popup.alert(document.body,"Config Upload","Upload Failed"));
};