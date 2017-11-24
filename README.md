## GoSymbols

Windows PDB self-service symbol server by Golang and Vuejs.

Index Page
![Index](/web/static/1.png?raw=true "Index")

Symbols Page
![Download](/web/static/2.png?raw=true "Download PDB")


## Config

The web portal has integrated with Windows Azure AD OAuth Authorization.

Replace the correct AppID/AppKey and RedirectURI in the `config.ini`.

``` ini
[base]
SYMSTORE_EXE    = "C:\Program Files (x86)\Windows Kits\8.1\Debuggers\x86\symstore.exe"
BUILD_SOURCE    = "Z:\BuildServer"
DESTINATION     = "D:\\SymbolServer"
LATEST_BUILD    = latestbuild.txt
EXCLUDE_LIST    = vc120.pdb,zlib10.pdb
DEBUG_ZIP       = debug.zip
LOG_PATH        = 

[app]
CLIENT_ID       = <Your AppId>	  # Windows Azure AD Application ID
CLIENT_KEY      = <Your AppKey>	  # Application Key
REDIRECT_URI    = http://localhost:8010/api/auth/authorize  # Windows AD OAuth redirect URL
GRAPH_SCOPE     = https://graph.microsoft.com/User.Read		
```


## Build

Golang
``` bash
go build
```

Vuejs
``` bash
# web src folder
cd web

# install dependencies
npm install

# build for production with minification
npm run build
```


## Run

``` bash
GoSymbols serve
```
