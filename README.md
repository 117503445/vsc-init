# vsc-init

Code Server 安装扩展、修改设置

## 使用方法

安装

```sh
go install github.com/117503445/vsc-init@latest
```

运行

```sh
vsc-init

EXTS=golang.go,njzy.stats-bar vsc-init # 使用 EXTS 环境变量安装额外的扩展
```

`vsc-init` 会根据 `pkg/assets/assets.go`，对于本地已安装的 Code Server，进行

- 下载安装拓展
- 写入 Settings 配置
- 写入 Keybindings 配置
