require("../css/publish.css");
require('spectre.css/dist/spectre-exp.min.css');

module.exports = () => {
    return new Promise((resolve, reject) => {
        let fragment = document.createDocumentFragment();
        let popup = document.createElement("div");
        popup.setAttribute("class", "modal active");
        popup.innerHTML = require("../html/publish.html");
        fragment.append(popup);

        const title = fragment.getElementById("title");
        const result = fragment.getElementById("result");
        const close = fragment.getElementById("close");

        fetch("/admin/api/build", {method: "POST"})
            .then(resp => {
                if (!resp.ok) return Promise.reject(`${resp.status} ${resp.statusText}`);
                return resp.json();
            })
            .then(json => {
                title.innerText = json.code === 0 ? "Published" : "Publish Error";
                result.innerHTML = "";
                let res = document.createElement("pre");
                res.setAttribute("class", "console");
                res.innerText = json.result;
                result.appendChild(res);
                close.removeAttribute("disabled");
            })
            .catch(err=>{
                reject(err);
            });

        close.onclick = ()=>{
            document.body.removeChild(popup);
            resolve();
        };

        document.body.appendChild(fragment);
    });
};