# HugoCMS
CMS for Hugo
## About
HugoCMS is Simple CMS for Self-Hosted Hugo Website.
## Features
* Web-based file browser
* WYSIWYG / Markdown Editor (based on [TOAST UI Editor](https://github.com/nhn/tui.editor))
* [Plugin support](https://github.com/dfkdream/HugoCMS/tree/master/plugin)
## Build
1. `git clone https://github.com/dfkdream/HugoCMS.git`
2. `docker build --tag hugo-cms .`
## Usage
```shell script
docker run -d --restart always \
    -v /your/website/directory:/website \
    -v /your/data/directory:/data \
    -e "DIR=/website" \
    -e "BOLT=/data/bolt.db" \
    -p 8080:80 \
    hugo-cms
```
Your website will be available on http://localhost:8080.

If you're running HugoCMS for first time, access http://localhost:8080/admin to setup root account.

On first run, your website may return `404 Not Found` Error. Click Rebuild button on admin page to build website.