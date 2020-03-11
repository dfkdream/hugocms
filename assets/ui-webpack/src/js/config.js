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
    .then((res)=>{
        if (!res.ok) return Promise.reject(res.json());
        return res.text()
    })
    .then(text=>editor.setValue(text, -1))
    .catch(err=>{
        if (err instanceof Promise) {
            err.then(json => {
                if (json.code === 403) popup.alert(document.body, "Session Expired", "Please sign in again").then(() => location.reload());
                else popup.alert(document.body, "Error", `${json.code} ${json.message}`);
            })
                .catch(err => {
                    console.log(err);
                    popup.alert(document.body, "Error", "Unknown error occurred. Please reload.");
                });
        }else{
            console.log(err);
            popup.alert(document.body, "Error", "Unknown error occurred. Please reload.");
        }
    });

document.getElementById("save").onclick=()=>{
    fetch("/admin/api/config",{
        method: "POST",
        body: editor.getValue()
    })
        .then(res=>{
            if (!res.ok) return Promise.reject(res.json());
            popup.alert(document.body,"Config Upload","Upload Succeed");
        })
        .catch((err)=>{
            err.then(json=>{
                popup.alert(document.body,"Config Upload",`Upload Failed: ${json.code} ${json.message}`);
            })
                .catch(err=>{
                    console.log(err);
                    popup.alert(document.body,"Config Upload","Upload Failed")
                });
        });
};