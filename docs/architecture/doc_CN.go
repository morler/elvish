/*
# 概述

此文件在高层次上记录了Elvish代码库的结构。您可以在代码编辑器中阅读它，或在godoc查看器中阅读它，例如
https://pkg.go.dev/src.elv.sh@master/docs/architecture。

Elvish是一个Go项目。如果您不熟悉Go代码的组织方式，请从[如何编写Go代码]开始。

Elvish仓库中的Go代码位于两个目录下：

  - cmd目录包含Elvish的入口点，但包含的代码很少。
  - pkg目录包含Elvish的大部分Go代码。它有很多子目录，所以仅通过浏览文件树可能有点难以找到方向。

我们将首先介绍cmd目录，然后专注于pkg下最重要的子目录。

Elvish仓库还包含其他目录。它们在技术上不是Go程序的一部分，所以我们不会在这里介绍它们。阅读它们各自的README文件以了解更多信息。

# 模块、包和符号名称

Elvish的模块名称是[src.elv.sh]。您可以将其视为代码实际托管位置的别名（当前为[github.com/elves/elvish]）。

所有包的导入路径都以模块名称[src.elv.sh]开头。例如，pkg/parse中包的导入路径是[src.elv.sh/pkg/parse]。

当引用包中的符号时，我们将仅使用包导入路径的最后一个组件。例如，[src.elv.sh/pkg/eval]包中的Evaler类型简单地称为[eval.Evaler]。（这与Go的语法一致。）

# 入口点（cmd/elvish和pkg/prog）

Elvish的默认入口点是[src.elv.sh/cmd/elvish]。它有一个main函数执行以下操作：

  - 从多个子程序组装一个"复合程序"，最重要的是[shell.Program]。
  - 调用[prog.Run]。

您可以在[src.elv.sh/pkg/prog]的godoc中了解此方法的优势。

还有其他main包，如[src.elv.sh/cmd/withpprof/elvish]。它们遵循相同的结构，仅在包含的子程序方面有所不同。

# shell子程序（pkg/shell）

shell子程序有两种略有不同的"模式"，交互式和非交互式，取决于命令行参数。[shell.Program]的文档包含更多详细信息。

在这两种模式中，shell子程序都使用在[src.elv.sh/pkg/eval]中实现的解释器来评估代码。

在交互式模式中，shell还使用在[src.elv.sh/pkg/edit]中实现的行编辑器来交互式读取命令。编辑器的某些功能依赖于持久存储；shell子程序也负责使用[src.elv.sh/pkg/daemon]初始化它。

# 解释器（pkg/eval）

[src.elv.sh/pkg/eval]包可能是Elvish中最重要的包，因为它实现了Elvish语言和内置模块。

解释器由[eval.Evaler]表示，通过[eval.NewEvaler]创建。方法[eval.Evaler.Eval]（是的，有3个"eval"）评估Elvish代码，并分几个步骤执行：

 1. 调用解析器获取AST。
 2. 将AST编译为"操作树"。
 3. 运行操作树。

选择这种方法主要是因为其简单性。它可能性能不是很好。

每个AST节点编译为其相应操作节点的过程，以及每个操作节点的运行方式，在几个compile_*.go文件中定义。这些文件是实现大部分语言语义的地方。

此包的另一大块是各种builtin_fn_*.go文件，它们实现内置模块的函数。这些可能在将来会移动到不同的包中。

对解释器重要的其他一些包有：

  - [src.elv.sh/pkg/eval/vals]为Elvish值实现标准操作集。
  - [src.elv.sh/pkg/persistent]实现Elvish的列表和映射，模仿[Clojure的向量和映射]。
  - [src.elv.sh/pkg/mods]的子目录实现各种内置模块。

# 解析器（pkg/parse）

[src.elv.sh/pkg/parse]包实现Elvish代码的解析，以[parse.Parse]作为入口点。

解析算法是手写的[递归下降]算法，具有没有单独标记化阶段的略微不寻常的特性。有关更多详细信息，请阅读包的godoc。

# 编辑器（pkg/edit）

[src.elv.sh/pkg/edit]包包含Elvish的交互式行编辑器，由[edit.Editor]表示。传统术语"行编辑器"有点用词不当；现代行编辑器（包括Elvish的）类似于像Vim这样的完整TUI应用程序，只是它们通常将自己限制在终端的最后N行而不是整个屏幕。

编辑器构建在更低级的[src.elv.sh/pkg/cli]包之上（这也有点用词不当），特别是[cli.App]类型。

整个TUI堆栈很快就需要重写。

编辑器依赖持久存储来实现目录历史和命令历史等功能。如上所述，存储的初始化在pkg/shell中使用pkg/daemon完成。

# 存储守护进程（pkg/daemon）

持久存储支持目前由存储守护进程提供。[src.elv.sh/pkg/daemon]包实现两个东西：

  - 实现存储守护进程的子程序（[daemon.Program]）。
  - 与守护进程通信的客户端（由[daemon.Activate]返回）。

守护进程按需启动和终止：

  - 第一个交互式Elvish shell启动守护进程。
  - 后续的交互式shell与同一个守护进程通信。
  - 当最后一个交互式Elvish shell退出时，守护进程也会退出。

在内部，守护进程使用[bbolt]作为数据库引擎。

在将来（需要评估），Elvish可能会获得自定义数据库，守护进程可能会消失。

# 结语

这应该让您对Elvish实现的最重要部分有一个粗略的了解。实现优先考虑可读性，大多数导出的符号都有文档，所以请随意深入源代码！

如果您有问题，请随时在用户组中询问或私信xiaq。

[如何编写Go代码]: https://go.dev/doc/code
[github.com/elves/elvish]: https://github.com/elves/elvish
[src.elv.sh]: https://src.elv.sh
[Clojure的向量和映射]: https://clojure.org/reference/data_structures
[递归下降]: https://en.wikipedia.org/wiki/Recursive_descent_parser
[bbolt]: https://github.com/etcd-io/bbolt
*/
package architecture

import (
	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/daemon"
	"src.elv.sh/pkg/edit"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/parse"
	"src.elv.sh/pkg/prog"
	"src.elv.sh/pkg/shell"
)

var (
	_ = new(shell.Program)
	_ = prog.Run

	_ = eval.NewEvaler
	_ = (*eval.Evaler).Eval

	_ = parse.Parse

	_ = new(edit.Editor)
	_ = new(cli.App)

	_ = new(daemon.Program)
	_ = daemon.Activate
)
