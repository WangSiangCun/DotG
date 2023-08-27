2023/8/26 20:36 by wangxiangkun

### 构建:

查看环境:

`go env `

```
set GO111MODULE=on
set GOARCH=amd64
set GOBIN=
set GOCACHE=C:\Users\wxk17\AppData\Local\go-build
set GOENV=C:\Users\wxk17\AppData\Roaming\go\env
set GOEXE=.exe
set GOEXPERIMENT=
set GOFLAGS=
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GOINSECURE=
set GOMODCACHE=E:\Soft\JDK\GoJDK\pkg\mod
set GONOPROXY=
set GONOSUMDB=
set GOOS=windows
set GOPATH=E:\Soft\JDK\GoJDK
...
```

`GOOS` 应当为对应平台

切换到对应平台

`go env -w GOOS=windows`
`go env -w GOOS=linux`

构建server:

cd到server文件夹后：
`go build server.go `

或指定文件名：` go build -o dotg_server.exe server.go `