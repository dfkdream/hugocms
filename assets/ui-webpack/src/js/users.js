require('../css/style.css');
require('../css/users.css');

require('spectre.css/dist/spectre.min.css');

require('../css/all.min.css');

require('../css/filelist.css');

const popups = require('./popup');

const showAddUserPopup = () => {
    return new Promise((resolve, reject) => {
        let fragment = document.createDocumentFragment();
        let popup = document.createElement("div");

        popup.setAttribute("class", "modal active");
        popup.innerHTML = require('../html/user-popup.html');
        fragment.append(popup);

        const id = fragment.getElementById("id");
        const username = fragment.getElementById("username");
        const password = fragment.getElementById("password");
        const permissions = fragment.getElementById("permissions");

        const closePopup = () => {
            document.body.removeChild(popup);
        };

        fragment.getElementById("cancel").onclick = () => {
            closePopup();
            resolve(false);
        };

        fragment.getElementById("modal-form").onsubmit = (e) => {
            e.preventDefault();

            fetch("/admin/api/users", {
                method: "POST",
                body: JSON.stringify({
                    id: id.value,
                    username: username.value,
                    password: password.value,
                    permissions: permissions.value.split("\n").map(v => v.trim()).filter(v => v !== "")
                })
            })
                .then(resp => {
                    if (!resp.ok) {
                        if (resp.status === 409) popups.alert(document.body, "Error", "ID Conflict");
                        else popups.alert(document.body, "Error", `Error ${resp.status}`);
                    } else {
                        closePopup();
                        resolve(true);
                    }
                })
                .catch((err) => {
                    closePopup();
                    reject(err);
                })
        };

        document.body.appendChild(fragment);
    })
};

const showEditUserPopup = (uid) => {
    return new Promise((resolve, reject) => {
        let fragment = document.createDocumentFragment();
        let popup = document.createElement("div");

        popup.setAttribute("class", "modal active");
        popup.innerHTML = require('../html/user-popup.html');
        fragment.append(popup);

        fragment.getElementById("title").innerText = "Edit User";

        const id = fragment.getElementById("id");
        const username = fragment.getElementById("username");
        const password = fragment.getElementById("password");
        const permissions = fragment.getElementById("permissions");

        id.setAttribute("disabled", "");
        password.removeAttribute("required");

        const closePopup = () => {
            document.body.removeChild(popup);
        };

        fragment.getElementById("cancel").onclick = () => {
            closePopup();
            resolve(false);
        };

        fetch("/admin/api/user/" + uid)
            .then(resp => {
                if (!resp.ok) return Promise.reject(`${resp.status} ${resp.statusText}`);
                return resp.json();
            })
            .then(json => {
                id.value = uid;
                username.value = json.username;
                permissions.value = json.permissions.join("\n");
            })
            .catch((err) => {
                console.log(err);
                closePopup();
                reject(err);
            });

        fragment.getElementById("modal-form").onsubmit = (e) => {
            e.preventDefault();

            fetch("/admin/api/user/"+id,{
                method:"POST",
                body:JSON.stringify({
                    username: username.value,
                    password: password.value,
                    permissions: permissions.value.split("\n").map(v => v.trim()).filter(v => v !== "")
                })
            })
                .then(resp=>{
                    if (!resp.ok) popups.alert(document.body, "Error", `Error ${resp.status}`);
                    else {
                        closePopup();
                        resolve(true);
                    }
                })
        };

        document.body.appendChild(fragment);
    })
};

const usersTable = document.getElementById("users-table");

const loadUsers = ()=>{
    usersTable.innerHTML="";
    fetch("/admin/api/users")
        .then(resp=>{
            if (!resp.ok) return Promise.reject(`${resp.status} ${resp.statusText}`);
            return resp.json();
        })
        .then(users=>{
            let fragment = document.createDocumentFragment();
            users.forEach(user=>{
                let row = document.createElement("tr");
                row.insertCell().innerText = user.id;
                row.insertCell().innerText = user.username;

                let actionsCell = row.insertCell();

                let editBtn = document.createElement("i");
                editBtn.setAttribute("class","tooltip tooltip-bottom fas fa-edit");
                editBtn.dataset.tooltip = "Edit";

                editBtn.onclick = ()=>{
                    showEditUserPopup(user.id)
                        .then(ok=>{if (ok) loadUsers()});
                };

                let deleteBtn = document.createElement("i");
                deleteBtn.setAttribute("class","tooltip tooltip-bottom fas fa-trash-alt");
                deleteBtn.dataset.tooltip = "Delete";

                deleteBtn.onclick = ()=>{
                    popups.confirm(document.body,"Confirm Delete",`Delete User ${user.id}?`)
                        .then(ok=>{
                            if (ok){
                                fetch("/admin/api/user/"+user.id,{method:"DELETE"})
                                    .then(resp=>{
                                        if (resp.ok) loadUsers();
                                        else popups.alert(document.body,"Error",`${resp.status} ${resp.statusText}`)
                                    })
                            }
                        })
                };

                actionsCell.appendChild(deleteBtn);
                actionsCell.appendChild(editBtn);

                fragment.appendChild(row);
            });

            usersTable.appendChild(fragment);
        }
    )
};

loadUsers();

document.getElementById("add-user").onclick = ()=>{
    showAddUserPopup()
        .then(ok=>{
            if (ok) loadUsers();
        })
};