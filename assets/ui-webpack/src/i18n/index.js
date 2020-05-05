import i18next from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import en from "./translations/en.json";
import ko from "./translations/ko.json";

export function i18n() {
    return new Promise((resolve, reject) => {
        i18next.use(LanguageDetector).init({
            fallbackLng: 'en',
            resources: {
                en: {translation: en},
                ko: {translation: ko},
            }
        })
            .then(t => {
                i18nDOM(t, document.body);
                resolve(t);
            })
            .catch(r => {
                reject(r);
            })
    })
}

export function i18nDOM(t, fragment) {
    fragment.querySelectorAll('[data-i18n]').forEach(v => {
        v.innerText = t(v.dataset.i18n);
    })
}

export default {i18n, i18nDOM};