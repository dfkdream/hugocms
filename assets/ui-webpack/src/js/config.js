require("spectre.css/dist/spectre.css");

require("../css/all.min.css");

require("../css/style.css");
require("../css/config.css");

import popups from "./popup";
import {i18n} from "../i18n";

const ace = require("ace-builds/src-noconflict/ace");

require("ace-builds/src-noconflict/mode-yaml");

i18n().then(t => {
    const popup = new popups(t);
    const configCode = document.getElementById("config-code");

    let editor = ace.edit(configCode, {
        mode: "ace/mode/yaml"
    });

    fetch("/admin/api/config")
        .then((res) => {
            if (!res.ok) return Promise.reject(res.json());
            return res.text()
        })
        .then(text => editor.setValue(text, -1))
        .catch(err => {
            if (err instanceof Promise) {
                err.then(json => {
                    if (json.code === 403) popup.alert(document.body, t("error"), t("errSessionExpired")).then(() => location.reload());
                    else popup.alert(document.body, t("error"), `${json.code} ${json.message}`);
                })
                    .catch(err => {
                        console.log(err);
                        popup.alert(document.body, t("error"), t("errUnknown"));
                    });
            } else {
                console.log(err);
                popup.alert(document.body, t("error"), t("errUnknown"));
            }
        });

    document.getElementById("save").onclick = () => {
        fetch("/admin/api/config", {
            method: "POST",
            body: editor.getValue()
        })
            .then(res => {
                if (!res.ok) return Promise.reject(res.json());
                popup.alert(document.body, t("success"), t("uploadSuccess"));
            })
            .catch((err) => {
                err.then(json => {
                    popup.alert(document.body, t("error"), `Upload Failed: ${json.code} ${json.message}`);
                })
                    .catch(err => {
                        console.log(err);
                        popup.alert(document.body, t("error"), t("errUploadFailed"))
                    });
            });
    };
});