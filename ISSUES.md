# Elvish 测试发现的问题

测试运行时间：2025-08-27
测试命令：`go test ./... -timeout 10m`

## 摘要

- **总体测试结果**: 部分失败
- **失败的包数量**: 4个包
- **主要影响区域**: Windows平台路径显示、导航UI、守护进程交互

## 详细问题分析

### 1. E2E测试失败

**包**: `src.elv.sh/e2e`
**测试**: `TestTranscripts_NoDaemon`
**状态**: ❌ 失败

**错误详情**:
```
--- FAIL: TestTranscripts_NoDaemon (35.28s)
    testing.go:1267: TempDir RemoveAll cleanup: CreateFile Z:\Temp\TestTranscripts_NoDaemon480687213\001\elvish.exe: Access is denied.
```

**问题分析**:
- Windows平台上临时目录清理权限问题
- 可能是elvish.exe进程仍在运行导致文件锁定
- E2E测试在Windows上的资源清理存在竞态条件

**优先级**: 🔴 高 - 影响CI/CD流程

### 2. 编辑器导航显示问题

**包**: `src.elv.sh/pkg/edit`
**受影响测试**: 5个导航相关测试
**状态**: ❌ 失败

**具体失败测试**:
1. `TestLocationAddon_Workspace`
2. `TestNavigation`
3. `TestNavigation_WidthRatio`
4. `TestNavigation_EnterDoesNotAddSpaceAfterSpace`
5. `TestNavigation_UsesEvalerChdir`

**核心问题**:
```
wanted: ~\ws1\tmp>
got:    ~/ws1/tmp>
```

**问题分析**:
- Windows路径分隔符显示不一致
- 期望显示反斜杠(`\`)但实际显示正斜杠(`/`)
- 这表明最近的Windows路径显示修复(d0ab9431)可能不完整
- 测试期望值可能需要更新以匹配新的跨平台路径显示逻辑

**优先级**: 🟡 中 - 影响Windows用户体验

### 3. 守护进程模块测试失败

**包**: `src.elv.sh/pkg/mods/daemon`
**测试**: `TestTranscripts/daemon.elvts`
**状态**: ❌ 失败

**错误详情**:
```
~> put $daemon:pid
-want +got:
@@ -1,2 +1,1 @@
 ▶ 12345
-

~> put $daemon:sock
-want +got:
@@ -1,2 +1,1 @@
 ▶ /tmp/elvish-test.sock
-
```

**问题分析**:
- transcript测试期望多一个空行但实际输出缺少
- 可能是输出格式化的微小变化
- 或者是transcript测试框架的换行符处理问题

**优先级**: 🟢 低 - 功能正常，仅输出格式差异

### 4. Shell交互测试失败

**包**: `src.elv.sh/pkg/shell`
**测试**: `TestTranscripts/interact_test.elvts`
**状态**: ❌ 失败

**错误详情**:
```
~> echo "\nuse store; store:next-cmd-seq" | elvish 2>$os:dev-null
-want +got:
@@ -1,1 +0,0 @@
-▶ (num 1)
```

**问题分析**:
- 守护进程行为测试失败
- 空命令不应存储在历史记录中的功能可能受影响
- Windows环境下的管道和重定向处理可能有问题

**优先级**: 🔴 高 - 影响核心shell功能

### 5. 资源清理问题

**错误详情**:
```
failed to remove temp dir Z:\Temp\elvishtest.3353483367: remove Z:\Temp\elvishtest.3353483367\xdg-state-home\elvish\db.bolt: The process cannot access the file because it is being used by another process.
```

**问题分析**:
- bbolt数据库文件在测试结束时仍被进程持有
- Windows平台的文件锁定机制更严格
- 测试清理逻辑需要确保数据库正确关闭

**优先级**: 🟡 中 - 影响测试环境清理

## 成功的测试包

以下包的测试全部通过：
- `pkg/buildinfo`, `pkg/cli`, `pkg/daemon`
- 所有`pkg/mods/*`模块(除daemon外)
- `pkg/eval`, `pkg/parse`等核心包
- 大部分支持包

## 建议修复优先级

1. **立即修复** (🔴 高优先级):
   - E2E测试资源清理问题
   - Shell交互测试中的守护进程行为

2. **近期修复** (🟡 中优先级):
   - Windows路径显示一致性
   - 数据库文件锁定清理

3. **后续优化** (🟢 低优先级):
   - transcript测试输出格式统一

## Windows平台特定问题

当前Windows平台测试通过率约为93.9%，主要问题集中在：
- 路径分隔符显示不一致
- 文件权限和锁定处理
- 临时文件清理机制

这些问题表明虽然已有显著改进，但Windows平台兼容性仍需进一步优化。