require('../css/style.css');
require('../css/profile.css');

require('spectre.css/dist/spectre.min.css');

require("../css/all.min.css");

import popups from "./popup";

import {i18n} from "../i18n";

i18n().then(t=> {
    const popup = new popups(t);

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
                popup.alert(document.body, t("success"), t("usernameUpdated"));
            })
            .catch(err => {
                if (err instanceof Promise) {
                    err.then(json => {
                        popup.alert(document.body, t("error"), `${json.code} ${json.message}`);
                    })
                        .catch(err => {
                            console.log(err);
                            popup.alert(document.body, t("error"), t("errUnknown"))
                        })
                } else {
                    console.log(err);
                    popup.alert(document.body, t("error"), t("errUnknown"))
                }
            });
    };

    passwordForm.onsubmit = (e) => {
        e.preventDefault();
        if (newPassword.value !== confirmPassword.value) {
            popup.alert(document.body, t("error"), t("errPasswordConfirmationFailed"));
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
                popup.alert(document.body, t("success"), t("passwordUpdated"));
                currentPassword.value = "";
                newPassword.value = "";
                confirmPassword.value = "";
            })
            .catch(err => {
                if (err instanceof Promise) {
                    err.then(json => {
                        if (json.code === 403) popup.alert(document.body, t("error"), t("errCurrentPasswordConfirmationFailed"));
                        else popup.alert(document.body, t("error"), `${json.code} ${json.message}`);
                    })
                        .catch(err => {
                            console.log(err);
                            popup.alert(document.body, t("error"), t("errUnknown"));
                        })
                } else {
                    console.log(err);
                    popup.alert(document.body, t("error"), t("errUnknown"));
                }
            });
    };
});