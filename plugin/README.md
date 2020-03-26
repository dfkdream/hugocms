# HugoCMS Plugin
Plugin library for HugoCMS
## Adding Plugins to HugoCMS
Use environment variable `PLUGINS` to add plugins to HugoCMS.
Plugin addresses are separated by `,` (comma)
* Docker
```shell script
-e "PLUGINS=plugin1,plugin2:port,https://plugin3:port"
```
## Installation
```shell script
go get github.com/dfkdream/hugocms/plugin
```
