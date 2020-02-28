/*!
 * HugoCMS
 * https://github.com/dfkdream/HugoCMS
 *
 * Copyright 2020 dfkdream
 * Released under MIT License
 */
module.exports = {
    confirm: (target, title, msg)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../html/confirm.html");
            fragment.append(popup);

            fragment.getElementById("title").innerText = title;
            fragment.getElementById("content").innerText = msg;

            let yes = fragment.getElementById("yes");
            yes.onclick = ()=>{
                target.removeChild(popup);
                resolve(true);
            };

            fragment.getElementById("no").onclick = ()=>{
                target.removeChild(popup);
                resolve(false);
            };

            fragment.getElementById("close-overlay").onclick = ()=>{
                target.removeChild(popup);
                resolve(false);
            };

            target.appendChild(fragment);

            yes.focus();
        });
    },

    prompt: (target,title,msg)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../html/prompt.html");
            fragment.append(popup);

            let input = fragment.getElementById("input");

            fragment.getElementById("title").innerText = title;
            fragment.getElementById("content").innerText = msg;

            let ok = fragment.getElementById("ok");
            ok.onclick = ((input)=>{
                return ()=>{
                    let v = input.value;
                    target.removeChild(popup);
                    resolve(v);
                };
            })(input);

            fragment.getElementById("cancel").onclick = ()=>{
                target.removeChild(popup);
                resolve(null);
            };

            fragment.getElementById("close-overlay").onclick = ()=>{
                target.removeChild(popup);
                resolve(null);
            };

            fragment.getElementById("form").onsubmit = ()=>{
                ok.onclick(null);
            };

            target.appendChild(fragment);

            input.focus();
        });
    },

    alert:(target,title,msg)=>{
        return new Promise((resolve)=>{
            let fragment = document.createDocumentFragment();
            let popup = document.createElement("div");
            popup.setAttribute("class","modal active");
            popup.innerHTML = require("../html/alert.html");
            fragment.append(popup);

            fragment.getElementById("title").innerText = title;
            fragment.getElementById("content").innerText = msg;
            let ok = fragment.getElementById("ok");
            ok.onclick = ()=>{
                target.removeChild(popup);
                resolve();
            };

            fragment.getElementById("close-overlay").onclick = ()=>{
                target.removeChild(popup);
                resolve();
            };

            target.appendChild(fragment);

            ok.focus();
        });
    },
};