# 删除备份任务时支持「同步删除已备份分区」设计文档

- 日期:2026-06-19
- 模块:备份管理 → 任务列表(前端为主,后端零改动)
- 状态:已与用户确认,待评审

## 1. 背景与现状

### 1.1 概念关系
前端「任务(Task)」是一个**虚拟概念**,没有独立的数据库表或 struct,由共享同一 `task_id` 的多个 `BackupPolicy` 聚合而成。

```
前端虚拟 Task(无 DB 表,按 BackupPolicy.task_id 聚合)
  └─ BackupPolicy(tbl_backup_policy)        一个 Task 含 1~N 个 Policy,每张表一个
       └─ BackupRun(tbl_backup_run)         一个 Policy 每次触发生成一条 Run
            └─ BackupRunPartition(嵌套 JSON) 一条 Run 含 0~N 个分区执行记录
```

### 1.2 当前删除任务的行为
`task-list.vue` 的 `deleteTask(t)`:
1. 弹出纯文本 `this.$confirm(...)`(无勾选框、无表格);
2. 确认后对任务下每个 policy 调 `DataManageApi.deletePolicy(p.policy_id)`;
3. 后端 `DeletePolicy` **仅软删 Policy**(`deleted=true, enabled=false`),**不删 Run 记录,不动 S3/本地盘的备份文件**。

### 1.3 已具备、可复用的能力
- 查询某表全历史 Run:`GET .../backup/table/:cluster/:db/:table/runs?days=0`(前端 `listRunsByTable`)。Run 内含分区级状态,可聚合「已成功备份分区个数」。
- 真正删除已备份分区数据:`POST .../backup/table/:cluster/:db/:table/partitions/delete`,body `{ partitions: []string, clean_remote: bool }`(前端 `deletePartitionRecords`)。内部:
  - 清理台账(从 Run 的 Partitions 数组剔除目标分区,Run 清空则物理删 Run);
  - `clean_remote=true` 时通过 `Storage.CleanPartition` 真正删 S3 对象 / 本地盘文件(best-effort);
  - **存在进行中(queued/running)Run 时会直接拒绝**,防止竞态。
- 软删策略:`deletePolicy`。
- 参考样板:`partition-list-dialog.vue` 已有「同步清理远端数据」勾选框与分区聚合逻辑。

## 2. 目标

在「任务列表」删除任务时:
1. 弹窗中**列出该任务关联的表**,以及**每张表已成功备份的分区个数**;
2. 增加一个**「同步删除已备份分区」勾选框**,默认不勾选;
3. 勾选后,删除任务时**同步删除这些表已备份的分区数据**(S3/本地盘文件 + 台账);
4. 不勾选时维持现状(仅软删 Policy,保留备份文件)。

## 3. 已确认的关键决策

| 决策 | 选择 | 理由 |
|---|---|---|
| 勾选框默认状态 | **默认不勾选** | 删除 S3/本地备份文件不可恢复,默认保守避免误删 |
| 存在进行中 Run 时 | **阻止删除并提示** | 后端 `DeletePartitionRecords` 本就拒绝进行中 run;前端提前阻断,避免竞态与部分失败 |
| 实现位置 | **纯前端复用现有接口,后端零改动** | 现有 API 已覆盖全部需求,YAGNI |

## 4. 架构

纯前端实现,复用三个现成接口:`listRunsByTable`(聚合 + 进行中检测)、`deletePartitionRecords`(删已备份数据)、`deletePolicy`(软删策略)。后端不改动。

## 5. 交互流程

1. 点击「删除任务」→ 打开自定义 `el-dialog`(`$confirm` 无法容纳表格+勾选框,参考 `partition-list-dialog.vue`)。
2. 打开时对任务下每个 policy 调 `listRunsByTable`,前端聚合:
   - 每张表「已成功备份分区个数」=去重后所有 `operation=backup && status=success` 的分区数;
   - 是否存在进行中 Run(任一 run `status ∈ {queued, running}`)。
3. dialog 展示:
   - 表格:列 `数据库.表名`、`已备份分区个数`;
   - 底部合计:「共 N 张表 / M 个已备份分区」;
   - 勾选框「同步删除已备份分区(将永久删除 S3/本地盘上的备份文件,不可恢复)」,默认不勾选;
   - 确认 / 取消按钮。
4. 若检测到进行中 Run:**禁用勾选框 + 禁用确认按钮**,顶部红色提示「任务下有进行中的备份/恢复,请等待结束后再删除」。
5. 点确认删除:
   - **勾选同步删除**:对每张「已备份分区个数 > 0」的表,先调
     `deletePartitionRecords(cluster, db, table, { partitions: <该表已成功分区列表>, clean_remote: true })`;
     全部处理完后再逐个软删 policy(`deletePolicy`)。
   - **未勾选**:维持现状,仅逐个软删 policy。
   - 用 `Promise.allSettled` 汇总,部分失败给出 warning 提示,最后 `this.$emit('refresh')` 刷新列表。

## 6. 边界与错误处理

- 某表 0 个已备份分区 → 显示 `0`;勾选同步删除时跳过该表(无分区可删)。
- 拉取 runs 失败 → 个数显示 `-`;此时**禁止勾选**同步删除(拿不到分区列表),但允许仅软删。
- legacy policy(`task_id` 为空)→ 任务聚合逻辑不变,删除路径相同。
- 删除顺序:先 `deletePartitionRecords`(需要 run 数据定位 storagePrefix)再软删 policy;软删 policy 不影响 run 表查询。
- 语义说明:`deletePartitionRecords` 按 `(cluster, db, table)` 删该表**全历史**这些分区的台账与远端数据(与现有「按表删分区」对话框一致)。通常一表一策略,无副作用;若同表被多策略备份,会一并清理该分区——属预期且与现有行为一致。

## 7. 改动文件(均在 `ckman-fe`)

| 文件 | 改动 |
|---|---|
| `src/views/data-manage/component/task-list.vue` | `deleteTask` 改为打开新 dialog;新增聚合(每表已成功分区 + 进行中检测)与删除逻辑(勾选则先删分区数据再软删) |
| 新增删除确认 dialog | 可内联于 `task-list.vue`,或抽独立组件;含表格 + 合计 + 勾选框 + 进行中提示 |
| `src/services/i18n.ts` | `history.*` 命名空间新增中英文 key(标题、表头、合计、勾选框文案、进行中提示、结果汇总) |
| `src/apis/dataManage.ts` | 确认已有 `listRunsByTable` / `deletePartitionRecords` / `deletePolicy`(预期已有,无需新增) |

## 8. 测试与验证

- 前端 `make lint` + `make build` 通过(在 `ckman-fe` 执行)。
- 手工验证场景:
  1. 任务含多表、各有已备份分区 → 弹窗正确展示每表分区数与合计;
  2. 不勾选删除 → 仅软删,S3/本地文件仍在;
  3. 勾选删除 → 远端文件被清除、台账被清理、policy 软删;
  4. 存在进行中 run → 勾选框与确认按钮禁用,提示正确;
  5. 某表 0 分区 / runs 拉取失败 → 展示与禁用逻辑符合预期。

## 9. 不在本次范围(YAGNI)

- 不新增后端聚合/删除接口。
- 不改动「删除任务不勾选时」的既有软删语义。
- 不处理 Run 孤儿记录的额外清理(维持现状)。
