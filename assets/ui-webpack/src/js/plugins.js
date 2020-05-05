/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */

require('../css/style.css');
require('../css/plugins.css');

require('spectre.css/dist/spectre.min.css');

require("../css/all.min.css");

import popups from "./popup";

import {i18n} from "../i18n";

i18n().then(t => {
    const popup = new popups(t);
    fetch("/admin/api/plugins")
        .then((res) => {
            if (!res.ok) return Promise.reject(res.json());
            return res.json();
        })
        .then(json => {
            let fragment = document.createDocumentFragment();
            json.forEach(p => {
                let f = document.createDocumentFragment();
                let d = document.createElement("div");
                d.setAttribute("class", "column col-6 col-md-12 card-container");
                d.innerHTML = require('../html/plugin-card.html');
                f.appendChild(d);
                f.getElementById("plugin-title").innerText = p.info.name + "@" + p.info.version;
                f.getElementById("plugin-author").innerText = p.info.author;
                f.getElementById("plugin-description").innerText = p.info.description;
                let chip = f.getElementById("plugin-live");
                chip.innerText = p.isLive ? "Live" : 'Down';
                chip.setAttribute("class", p.isLive ? "chip bg-success" : "chip bg-error");
                fragment.appendChild(f);
            });
            document.getElementById("plugin-cards").appendChild(fragment);
        })
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

}).catch(e => console.log(e));