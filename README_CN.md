# Elvish（中文版）

[![CI状态](https://github.com/elves/elvish/workflows/CI/badge.svg)](https://github.com/elves/elvish/actions?query=workflow%3ACI)
[![FreeBSD & gccgo测试状态](https://img.shields.io/cirrus/github/elves/elvish?logo=Cirrus%20CI&label=CI2)](https://cirrus-ci.com/github/elves/elvish/master)
[![测试覆盖率](https://img.shields.io/codecov/c/github/elves/elvish/master.svg?logo=Codecov&label=coverage)](https://app.codecov.io/gh/elves/elvish/tree/master)
[![Go参考](https://pkg.go.dev/badge/src.elv.sh@master.svg)](https://pkg.go.dev/src.elv.sh@master)
[![打包状态](https://repology.org/badge/tiny-repos/elvish.svg)](https://repology.org/project/elvish/versions)

[![论坛](https://img.shields.io/badge/forum-bbs.elv.sh-5b5.svg?logo=discourse)](https://bbs.elv.sh)
[![Twitter](https://img.shields.io/badge/twitter-@ElvishShell-blue.svg?logo=x)](https://twitter.com/ElvishShell)

[![Telegram群组](https://img.shields.io/badge/telegram-Elvish-blue.svg?logo=telegram&logoColor=white)](https://t.me/+Pv5ZYgTXD-YaKwcP)
[![Discord服务器](https://img.shields.io/badge/discord-Elvish-blue.svg?logo=discord&logoColor=white)](https://discord.gg/jrmuzRBU8D)
[![#users:elv.sh](https://img.shields.io/badge/matrix-%23users:elv.sh-blue.svg?logo=matrix)](https://matrix.to/#/#users:elv.sh)
[![#elvish on libera.chat](https://img.shields.io/badge/libera.chat-%23elvish-blue.svg?logo=liberadotchat&logoColor=white)](https://web.libera.chat/#elvish)
[![Gitter](https://img.shields.io/badge/gitter-elves%2Felvish-blue.svg?logo=gitter)](https://gitter.im/elves/elvish)

（聊天室全部通过[Matrix](https://matrix.org)连接。）

Elvish是：

-   一个强大的脚本语言。

-   一个带有有用交互特性的shell。

-   一个适用于Linux、BSD、macOS或Windows的静态链接二进制文件。

Elvish目前是1.0版本之前的状态。这意味着仍然会不时发生破坏性变更，但它已经足够稳定，可用于脚本编写和交互式使用。

## 文档

[![用户文档](https://img.shields.io/badge/User_Docs-37a779?style=for-the-badge)](https://elv.sh)

用户文档托管在Elvish的网站[elv.sh](https://elv.sh)上。包括[如何安装Elvish](https://elv.sh/get/)，[教程](https://elv.sh/learn/)，[参考页面](https://elv.sh/ref/)，以及[新闻](https://elv.sh/blog/)。

[![开发文档](https://img.shields.io/badge/Development_Docs-blue?style=for-the-badge)](./docs)

开发文档位于[./docs](./docs)。

[![Awesome Elvish](https://img.shields.io/badge/Awesome_Elvish-orange?style=for-the-badge)](https://github.com/elves/awesome-elvish)

支持Elvish的精彩包和工具。

## 语言切换

- [English Version](README.md) - 英文版本
- [中文版本](README_CN.md) - 当前页面

## 许可证

除以下文件外，所有源文件均使用BSD 2条款许可证（请参阅[LICENSE](LICENSE)）：

-   [pkg/diff](pkg/diff)和[pkg/rpc](pkg/rpc)中的文件使用BSD 3条款许可证发布，因为它们派生自[Go源代码](https://github.com/golang/go)。请参阅[pkg/diff/LICENSE](pkg/diff/LICENSE)和[pkg/rpc/LICENSE](pkg/rpc/LICENSE)。

-   [pkg/persistent](pkg/persistent)及其子目录中的文件使用EPL 1.0发布，因为它们部分派生自[Clojure源代码](https://github.com/clojure/clojure)。请参阅[pkg/persistent/LICENSE](pkg/persistent/LICENSE)。

-   [pkg/md/spec](pkg/md/spec)中的文件使用Creative Commons CC-BY-SA 4.0许可证发布，因为它们派生自[CommonMark规范](https://github.com/commonmark/commonmark-spec)。请参阅[pkg/md/spec/LICENSE](pkg/md/spec/LICENSE)。