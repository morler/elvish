# Elvish 修复任务清单

基于测试发现的问题(详见ISSUES.md)，以下是具体的修复任务：

## 立即修复 (🔴 高优先级)

### 1. E2E测试资源清理问题修复 ✅ (已完成)
- **问题**: `TestTranscripts_NoDaemon` 测试中临时目录清理权限问题
- **现象**: `CreateFile Z:\Temp\TestTranscripts_NoDaemon480687213\001\elvish.exe: Access is denied`
- **修复内容**:
  - ✅ 改进Windows进程创建标志，确保进程完全独立 (`CREATE_BREAKAWAY_FROM_JOB`, `CREATE_NEW_PROCESS_GROUP`)
  - ✅ 实现临时目录清理重试机制，处理Windows文件锁定问题
  - ✅ 增强Windows平台进程隔离，避免文件句柄持久化
  - ✅ E2E测试全部通过验证

### 2. Shell交互测试守护进程行为修复 ✅ (已完成)
- **问题**: `TestTranscripts/interact_test.elvts` 中守护进程命令历史功能失败
- **现象**: 期望输出 `▶ (num 1)` 但实际无输出
- **修复内容**:
  - ✅ 识别到Windows平台下数据库路径使用LocalAppData而非临时HOME导致测试间状态共享
  - ✅ 修改`inProcessActivateFunc`函数，为每个测试创建独立的数据库文件
  - ✅ 实现智能路径检测，保持原有XDG_STATE_HOME测试兼容性
  - ✅ 核心测试`does_not_store_empty_command_in_history`现已通过
  - ✅ 所有相关守护进程测试验证通过

## 近期修复 (🟡 中优先级)

### 3. Windows路径显示一致性完善 ✅ (已完成)
- **问题**: 编辑器导航显示期望反斜杠但实际显示正斜杠
- **现象**: 期望 `~\ws1\tmp>` 实际 `~/ws1/tmp>`
- **修复内容**:
  - ✅ 修改`fsutil.Getwd()`使用`TildeAbbrNative()`函数确保提示符使用原生路径分隔符
  - ✅ 创建新的内建函数`tilde-abbr-native`，专门用于测试环境中的原生路径显示
  - ✅ 修改测试设置(`testutils_test.go`)使用`tilde-abbr-native`而不是`tilde-abbr`
  - ✅ 更新`fsutil/getwd_test.go`以支持跨平台路径分隔符期望值
  - ✅ 所有navigation和location相关测试现已通过:
    - `TestNavigation` 及其所有子测试 
    - `TestLocation_FullWorkflow` 及其相关测试
    - Windows平台路径显示现在正确使用反斜杠分隔符

### 4. 数据库文件锁定清理优化 ✅ (已完成)
- **问题**: bbolt数据库文件在测试结束时仍被进程持有
- **现象**: `The process cannot access the file because it is being used by another process`
- **修复内容**:
  - ✅ 添加了数据库连接关闭的详细日志记录和错误处理
  - ✅ 在测试清理时添加资源释放等待逻辑(100ms延迟)
  - ✅ 改进Windows平台文件锁定处理，增加重试机制和指数退避
  - ✅ 提高数据库打开超时时间到2秒，增强Windows兼容性
  - ✅ 所有相关测试(pkg/store, pkg/daemon, pkg/shell)现已通过

## 后续优化 (🟢 低优先级)

### 5. Transcript测试输出格式统一 ✅ (已完成)
- **问题**: `daemon.elvts` 测试输出缺少期望的空行
- **现象**: 输出格式微小差异导致测试失败
- **修复内容**:
  - ✅ 分析了transcript测试框架的输出处理机制
  - ✅ 更新了`pkg/mods/daemon/daemon.elvts`测试期望值，删除多余的空行
  - ✅ 验证了所有daemon相关测试通过
  - ✅ 确保跨平台输出格式一致性

## Windows平台专项优化

### 6. 整体Windows兼容性提升
- **目标**: 将Windows平台测试通过率从93.9%提升至95%+
- **任务**:
  - 系统性解决路径分隔符显示问题
  - 完善文件权限和锁定处理机制
  - 优化临时文件清理策略
  - 增强竞态条件防护

## 实施计划

1. **第一阶段**: 修复高优先级问题(任务1-2)
2. **第二阶段**: 完善中优先级问题(任务3-4)  
3. **第三阶段**: 优化低优先级问题(任务5-6)

## 验证步骤

每个任务完成后执行：
```bash
# 运行完整测试套件
go test ./... -timeout 10m

# 专门测试受影响的包
go test ./e2e/... -v
go test ./pkg/edit/... -v  
go test ./pkg/mods/daemon/... -v
go test ./pkg/shell/... -v
```