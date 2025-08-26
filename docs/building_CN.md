# 从源码构建Elvish

要从源码构建Elvish，您需要：

-   支持的操作系统：Linux、{Free,Net,Open}BSD、macOS或Windows 10+。

    **📋 Windows用户**：请参阅[Windows用户指南](windows_CN.md)了解平台特定安装说明和配置建议。

-   Go >= 1.21.0。

要从源码构建Elvish，运行以下命令之一：

```sh
go install src.elv.sh/cmd/elvish@master # 安装最新提交
go install src.elv.sh/cmd/elvish@latest # 安装最新发布版本  
go install src.elv.sh/cmd/elvish@v0.18.0 # 安装特定版本
```

## 控制安装位置

[`go install`](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies)
命令将Elvish安装到`$GOBIN`；二进制文件名为`elvish`。您可以通过覆盖`$GOBIN`来控制安装位置，例如在`go install`命令前添加`env GOBIN=...`。

如果未设置`$GOBIN`，安装位置默认为`$GOPATH/bin`，如果`$GOPATH`也未设置，则默认为`~/go/bin`。

安装目录可能不在您操作系统的默认`$PATH`中。您应该将其添加到`$PATH`中，或手动将Elvish二进制文件复制到已在`$PATH`中的目录。

## 构建替代入口点

除了`src.elv.sh/cmd/elvish`（对应仓库中的[`cmd/elvish`](./cmd/elvish)目录）之外，还有一些替代入口点，都名为`cmd/*/elvish`，具有略微不同的功能集。（从Go的角度来看，这些只是不同的`main`包。）

例如，安装`cmd/withpprof/elvish`入口点以获得[性能分析支持](https://pkg.go.dev/runtime/pprof)（更改`@`后的部分以获得不同版本）：

```sh
go install src.elv.sh/cmd/withpprof/elvish@master
```

## 从本地源码树构建

如果您正在修改Elvish的源代码，您将希望克隆Elvish的Git仓库并从本地源码树构建Elvish。要做到这一点，请从源码树的根目录运行：

```sh
go install ./cmd/elvish
```

无需像`@master`那样指定版本；在源码树内时，`go install`将始终使用存在的任何源代码。

有关贡献者的更多说明，请参阅[contributing_CN.md](contributing_CN.md)。

## 使用实验性插件支持构建

Elvish对构建和导入插件（用Go编写的模块）具有实验性支持。它依赖于Go的[插件支持](https://pkg.go.dev/plugin)，该功能仅在少数平台上可用。

插件支持需要使用[cgo](https://pkg.go.dev/cmd/cgo)构建Elvish。官方[预构建二进制文件](https://elv.sh/get)为了兼容性和可重现性而不使用cgo构建，但默认情况下Go工具链启用cgo构建。

如果您在支持插件的平台上从源码构建了Elvish，您的Elvish构建可能已经支持插件。要在构建Elvish时强制使用cgo，您可以执行以下操作：

```sh
env CGO_ENABLED=1 go install ./cmd/elvish
```

要构建插件，请参阅此[示例](https://github.com/elves/sample-plugin)。

## 语言切换 / Language

- [English Version](building.md) - 英文版本
- [中文版本](building_CN.md) - 当前页面