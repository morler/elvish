# Windows使用指南

本文档提供Windows平台使用Elvish的完整指南，包括安装、配置和Windows特定的功能差异。

## 🪟 Windows平台支持状态

Elvish在Windows上提供良好支持，但有一些平台差异需要注意：

- ✅ **完整支持**: 核心shell功能、脚本执行、交互式编辑
- ✅ **跨平台兼容**: 大部分命令和功能与Unix系统一致  
- ⚠️ **部分差异**: 作业控制、权限管理、路径处理有Windows特定行为
- 🔄 **持续改进**: 测试通过率83.3%，持续优化中

## 📦 安装方式

### 方式1: 预编译二进制文件 (推荐)

访问 [elv.sh/get/](https://elv.sh/get/) 下载Windows版本：

1. 下载 `elvish-v0.xx.x-windows.zip`
2. 解压到目标目录 (如 `C:\elvish\`)
3. 将可执行文件目录添加到系统PATH
4. 在命令提示符或PowerShell中运行 `elvish`

### 方式2: 从源码编译

**前置要求**:
- Go >= 1.21.0 (推荐 Go 1.24+，获得最佳兼容性)
- Git (可选，用于克隆源代码)

**在线安装**:
```cmd
# 安装最新发布版本
go install src.elv.sh/cmd/elvish@latest

# 安装最新开发版本
go install src.elv.sh/cmd/elvish@master

# 安装特定版本
go install src.elv.sh/cmd/elvish@v0.21.0
```

**本地编译**:
```cmd
# 克隆源代码
git clone https://github.com/elves/elvish.git
cd elvish

# 编译安装
go install ./cmd/elvish
```

### 安装位置控制

默认安装位置: `%USERPROFILE%\go\bin\elvish.exe`

自定义安装位置:
```cmd
# 设置GOBIN环境变量
set GOBIN=C:\bin
go install src.elv.sh/cmd/elvish@latest
```

## ⚙️ Windows配置

### 设置为默认Shell

**注意**: Windows不同于Unix系统，不建议将Elvish设为系统默认shell。推荐的使用方式：

1. **Windows Terminal中使用**:
   - 打开Windows Terminal设置
   - 添加新的配置文件
   - 设置命令行为Elvish可执行文件路径

2. **PowerShell中启动**:
   ```powershell
   # 在PowerShell配置文件中添加别名
   Set-Alias elv C:\path\to\elvish.exe
   ```

### 环境变量配置

Elvish遵循Windows环境变量约定：

```elvish
# 查看Windows特定环境变量
echo $E:USERPROFILE    # 用户目录
echo $E:COMPUTERNAME   # 计算机名
echo $E:USERNAME       # 用户名
echo $E:TEMP           # 临时目录

# 设置环境变量
E:MY_VAR = 'value'
```

### 路径处理

Elvish在Windows上自动处理路径分隔符：

```elvish
# 这些路径表示法都有效
cd C:\Users\username
cd C:/Users/username
cd \Users\username

# 获取当前工作目录
pwd    # 返回Windows格式路径 (如 C:\Users\username)
```

## 🔧 Windows特定功能

### 文件系统操作

```elvish
# 列出文件和目录
ls C:\          # 列出C盘根目录
ls .            # 当前目录

# 文件操作
copy file.txt backup.txt    # 复制文件
move old.txt new.txt        # 移动/重命名文件
rm file.txt                 # 删除文件
```

### 网络驱动器支持

```elvish
# 访问网络驱动器
ls \\server\share
cd Z:\                      # 映射的网络驱动器
```

### Windows服务和进程

```elvish
# 查看进程 (如果安装了相关工具)
ps | head 10

# Windows特定的作业控制说明见下文
```

## ⚠️ Windows平台限制和差异

### 1. 作业控制

Windows平台的作业控制功能受限：

- ❌ **不支持**: Unix风格的作业挂起/恢复 (`Ctrl+Z`, `fg`, `bg`)
- ✅ **支持**: 进程终止 (`Ctrl+C`)
- ✅ **支持**: 基本的后台进程执行

**解决方案**: 使用Windows Task Manager或PowerShell作业管理功能。

### 2. 权限管理

Windows权限模型与Unix不同：

- 文件权限使用Windows ACL而非Unix模式位
- 管理员权限需要UAC提升
- 某些操作可能需要"以管理员身份运行"

### 3. 路径长度限制

Windows传统路径长度限制为260字符：

- 现代Windows 10/11已支持长路径
- 可能需要启用长路径支持策略
- 使用UNC路径格式避免限制

### 4. 文件名字符限制

Windows文件名不支持某些字符：
```
< > : " | ? * /
```

### 5. 区分大小写

Windows文件系统默认不区分大小写：

```elvish
# 这些指向同一个文件
ls File.txt
ls file.txt
ls FILE.TXT
```

### 6. 换行符差异

Windows使用CRLF (`\r\n`)，而Unix使用LF (`\n`)：

- Elvish会自动处理大部分情况
- 编辑文件时注意编辑器设置

## 🐛 故障排除

### 常见问题

**1. "找不到命令"错误**:
```
解决方案: 确保Elvish可执行文件在PATH环境变量中
```

**2. 路径包含空格**:
```elvish
# 使用引号包围路径
cd "C:\Program Files\MyApp"
ls "Documents and Settings"
```

**3. 权限被拒绝**:
```
解决方案: 以管理员身份运行终端，或检查文件权限
```

**4. Unicode字符显示问题**:
```
解决方案: 
1. 设置终端为UTF-8编码
2. 使用现代终端如Windows Terminal
3. 检查系统区域设置
```

### 性能优化

**1. Windows Defender排除**:
- 将Elvish安装目录添加到Windows Defender排除列表
- 排除Go构建目录 (`%USERPROFILE%\go`)

**2. 选择合适的终端**:
- 推荐: Windows Terminal
- 替代: ConEmu, Cmder
- 避免: 旧版命令提示符 (cmd.exe)

### 调试信息收集

遇到问题时，收集以下信息：

```elvish
# Elvish版本信息
elvish -version

# 系统信息  
echo $E:OS
echo $E:PROCESSOR_ARCHITECTURE

# 环境变量
echo $E:PATH
echo $E:GOPATH
echo $E:GOBIN
```

## 🔄 从其他Shell迁移

### 从PowerShell迁移

|PowerShell概念|Elvish等价|说明|
|---|---|---|
|`Get-ChildItem`|`ls`|列出文件|
|`Set-Location`|`cd`|改变目录|
|`Copy-Item`|`cp`|复制文件|
|`$env:VAR`|`$E:VAR`|环境变量|
|`Write-Host`|`echo`|输出文本|

### 从Batch迁移

|Batch概念|Elvish等价|说明|
|---|---|---|
|`dir`|`ls`|列出文件|
|`cd`|`cd`|改变目录|
|`copy`|`cp`|复制文件|
|`%VAR%`|`$E:VAR`|环境变量|
|`echo`|`echo`|输出文本|

## 📚 进一步学习

- 🌟 [Elvish官方教程](https://elv.sh/learn/)
- 📖 [语言参考](https://elv.sh/ref/)
- 🎯 [Windows Terminal配置](https://docs.microsoft.com/en-us/windows/terminal/)
- 💬 [社区支持](https://elv.sh) (论坛、Discord等)

## 🆘 获取帮助

如果遇到Windows特定问题：

1. 查看[GitHub Issues](https://github.com/elves/elvish/issues)
2. 在[论坛](https://bbs.elv.sh)提问并标注"Windows"
3. 加入[Discord社区](https://discord.gg/jrmuzRBU8D)
4. 提供详细的系统信息和错误日志

---

**最后更新**: 2025-08-25  
**适用版本**: Elvish v0.21+  
**测试状态**: 83.3%测试通过率 (55/66包通过)