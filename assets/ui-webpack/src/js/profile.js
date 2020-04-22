require('../css/style.css');
require('../css/profile.css');

require('spectre.css/dist/spectre.min.css');

require("../css/all.min.css");

const popup = require("./popup");

const profileForm = document.getElementById("form-profile");
const username = document.getElementById("username");

const passwordForm = document.getElementById("form-password");
const currentPassword = document.getElementById("current-password");
const newPassword = document.getElementById("new-password");
const confirmPassword = document.getElementById("confirm-password");

profileForm.onsubmit = (e) => {
    e.preventDefault();

    fetch("/admin/api/whoami", {
        method: "POST",
        body: JSON.stringify({
            username: username.value
        })
    })
        .then(res => {
            if (!res.ok) return Promise.reject(res.json());
            else return res.json();
        })
        .then(() => {
            popup.alert(document.body, "Success", "Username updated.");
        })
        .catch(err => {
            if (err instanceof Promise) {
                err.then(json => {
                    popup.alert(document.body, "Error", `${json.code} ${json.message}`);
                })
                    .catch(err => {
                        console.log(err);
                        popup.alert(document.body, "Error", "Unknown error occurred. Please reload.")
                    })
            } else {
                console.log(err);
                popup.alert(document.body, "Error", "Unknown error occurred. Please reload.")
            }
        });
};

passwordForm.onsubmit = (e) => {
    e.preventDefault();
    if (newPassword.value !== confirmPassword) {
        popup.alert(document.body, "Confirmation Failed", "Password confirmation failed.");
        return;
    }

    fetch("/admin/api/whoami", {
        method: "POST",
        body: JSON.stringify({
            currentPassword: currentPassword.value,
            newPassword: newPassword.value
        })
    })
        .then(res => {
            if (!res.ok) return Promise.reject(res.json());
            else return res.json();
        })
        .then(() => {
            popup.alert(document.body, "Success", "Password updated.");
            currentPassword.value = "";
            newPassword.value = "";
            confirmPassword.value = "";
        })
        .catch(err => {
            if (err instanceof Promise) {
                err.then(json => {
                    if (json.code === 403) popup.alert(document.body, "Error", "Failed to confirm current password");
                    else popup.alert(document.body, "Error", `${json.code} ${json.message}`);
                })
                    .catch(err => {
                        console.log(err);
                        popup.alert(document.body, "Error", "Unknown error occurred. Please reload.")
                    })
            } else {
                console.log(err);
                popup.alert(document.body, "Error", "Unknown error occurred. Please reload.")
            }
        });
};