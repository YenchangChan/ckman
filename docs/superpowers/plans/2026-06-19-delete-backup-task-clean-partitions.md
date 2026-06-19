# 删除备份任务时同步删除已备份分区 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在「备份管理 → 任务列表」删除任务时,弹窗中列出关联表及各表已备份分区个数,并提供「同步删除已备份分区」勾选框,勾选后真正删除 S3/本地盘上的备份数据。

**Architecture:** 纯前端实现,后端零改动。复用现有三个 API:`listRunsByTable`(聚合已成功备份分区 + 检测进行中 run)、`deletePartitionRecords`(删已备份数据 + 清台账)、`deletePolicy`(软删策略)。将 `task-list.vue` 现有的 `$confirm` 纯文本删除确认,替换为内联 `el-dialog`,承载表格 + 合计 + 勾选框 + 进行中阻断。

**Tech Stack:** Vue 2 + Element UI(`el-dialog` / `el-table` / `el-checkbox` / `el-alert`),i18n via `vue-i18n`(`src/services/i18n.ts`)。

## Global Constraints

- 所有前端改动在工作副本 `/data/root/go/src/github.com/housepower/ckman-fe/`,提交到 `main` 分支;**绝不在 `ckman/frontend` submodule 内编辑或提交**。
- 本前端仓库**无单元测试框架**(`package.json` 无 `test`/`jest`/`vitest`)。每个任务的验证 = 在 `ckman-fe` 执行 `make lint` + `make build`,外加结构化手工验证清单。
- 后端零改动,不新增任何 Go 接口。
- 勾选框「同步删除已备份分区」**默认不勾选**。
- 任务下存在进行中(`status` 为 `queued` 或 `running`)的 run 时,**禁用确认按钮,阻止删除并提示**。
- 字段名(verbatim,来自现有代码):policy 对象 `cluster_name` / `database` / `table` / `policy_id`;run 对象 `operation`(`'backup'`/`'restore'`)/ `status` / `partitions`;partition 对象 `partition` / `status`(`'success'` 等);API 成功判定 `res.data.retCode === '0000'`,返回体 `res.data.entity`,清理告警 `res.data.entity.warnings`。

---

## File Structure

| 文件 | 责任 |
|---|---|
| `src/services/i18n.ts`(修改) | 新增 `history.*` 命名空间下删除任务弹窗所需的中英文文案 key |
| `src/views/data-manage/component/task-list.vue`(修改) | 将 `deleteTask` 从 `$confirm` 改为内联删除确认 `el-dialog`;新增聚合逻辑、进行中检测、勾选同步删除的执行逻辑 |

---

## Task 1: 新增 i18n 文案 key(中英文)

**Files:**
- Modify: `src/services/i18n.ts`(en `history` 区块 `Task Delete Result Partial` 之后,约 739 行;zh `history` 区块对应位置,约 1722 行)

**Interfaces:**
- Produces(供 Task 2 的模板/方法引用,key 名 verbatim):
  - `history.Delete Task Title`
  - `history.Delete Task Tip`(占位符 `{name}`、`{count}`,含 HTML)
  - `history.Delete Task Col Table`
  - `history.Delete Task Col Backed Partitions`
  - `history.Delete Task Total`(占位符 `{tables}`、`{partitions}`)
  - `history.Delete Task Clean Checkbox`
  - `history.Delete Task Inflight Warning`
  - `history.Delete Task Fetch Error`
  - `history.Task Delete Clean Warnings`(占位符 `{count}`)
- 复用既有 key(无需新增):`common.Cancel`、`history.Confirm Delete Btn`、`history.Task Delete Result OK`、`history.Task Delete Result Partial`。

- [ ] **Step 1: 在英文 `history` 区块插入新 key**

定位英文区块中这一行(约 739 行):

```js
      'Task Delete Result Partial': 'Deleted {success} policy/policies, {failed} failed',
```

在其**后面**插入:

```js
      'Delete Task Title': 'Delete Task',
      'Delete Task Tip': 'Deleting task <b>{name}</b> ({count} table(s)). This removes the backup configuration and run history permanently.',
      'Delete Task Col Table': 'Table',
      'Delete Task Col Backed Partitions': 'Backed-up Partitions',
      'Delete Task Total': 'Total: {tables} table(s) / {partitions} backed-up partition(s)',
      'Delete Task Clean Checkbox': 'Also delete backed-up partitions (permanently removes backup files on S3 / local disk, cannot be recovered)',
      'Delete Task Inflight Warning': 'This task has running or queued backup/restore jobs. Please wait until they finish before deleting.',
      'Delete Task Fetch Error': 'Failed to load backup records for some tables; deleting backed-up partitions is disabled.',
      'Task Delete Clean Warnings': 'Task deleted, but cleanup of remote data reported {count} warning(s)',
```

- [ ] **Step 2: 在中文 `history` 区块插入新 key**

定位中文区块中这一行(约 1722 行):

```js
      'Task Delete Result Partial': '已删除 {success} 张表的备份配置，{failed} 张失败',
```

在其**后面**插入:

```js
      'Delete Task Title': '删除任务',
      'Delete Task Tip': '即将删除任务 <b>{name}</b>（{count} 张表）。这将永久删除备份配置与运行历史记录。',
      'Delete Task Col Table': '表',
      'Delete Task Col Backed Partitions': '已备份分区数',
      'Delete Task Total': '合计：{tables} 张表 / {partitions} 个已备份分区',
      'Delete Task Clean Checkbox': '同步删除已备份分区（将永久删除 S3/本地盘上的备份文件，不可恢复）',
      'Delete Task Inflight Warning': '任务下有进行中（排队/运行中）的备份或恢复，请等待结束后再删除。',
      'Delete Task Fetch Error': '部分表的备份记录加载失败，已禁用「同步删除已备份分区」。',
      'Task Delete Clean Warnings': '任务已删除，但远端数据清理有 {count} 条告警',
```

- [ ] **Step 3: lint 校验**

Run: `cd /data/root/go/src/github.com/housepower/ckman-fe && make lint`
Expected: `DONE  No lint errors found!`

- [ ] **Step 4: 确认 key 已成对存在(中英文各 9 个新 key)**

Run: `cd /data/root/go/src/github.com/housepower/ckman-fe && grep -c "Delete Task Title\|Delete Task Tip\|Delete Task Col Table\|Delete Task Col Backed Partitions\|Delete Task Total\|Delete Task Clean Checkbox\|Delete Task Inflight Warning\|Delete Task Fetch Error\|Task Delete Clean Warnings" src/services/i18n.ts`
Expected: `18`(en 9 + zh 9)

- [ ] **Step 5: 提交**

```bash
cd /data/root/go/src/github.com/housepower/ckman-fe
git add src/services/i18n.ts
git commit -m "i18n(backup): 新增删除任务弹窗(同步删除已备份分区)文案"
```

---

## Task 2: task-list.vue 删除任务弹窗 + 同步删除逻辑

**Files:**
- Modify: `src/views/data-manage/component/task-list.vue`(template 根 `div` 内新增 `el-dialog`;`data()` 新增字段;新增 computed;替换 `deleteTask` 并新增方法)

**Interfaces:**
- Consumes:
  - `DataManageApi.listRunsByTable(clusterName, database, table)` → `res.data.entity` 为 run 数组,run 含 `operation` / `status` / `partitions`(每个 partition 含 `partition` / `status`)。
  - `DataManageApi.deletePartitionRecords(clusterName, database, table, { partitions: string[], clean_remote: boolean })` → `res.data.retCode`、`res.data.entity.warnings`。
  - `DataManageApi.deletePolicy(policyId)` → `res.data.retCode`。
  - 既有方法 `displayName(t)`(已存在)、Task 1 新增的 i18n key。
- Produces:组件内部状态,不对外暴露新接口;删除完成后沿用既有 `this.$emit('refresh')`。

- [ ] **Step 1: 在 `data()` 返回对象中新增删除弹窗状态字段**

定位 `task-list.vue` 的 `data()`(约 128-138 行),将其 `return { ... }` 内容替换为(在原有字段后追加新字段):

```js
  data() {
    return {
      filterEnabled: 'all',
      searchKey: '',
      currentPage: 1,
      pageSize: 20,
      latestRunMap: {},
      sortField: '',
      sortOrder: '',
      // 删除任务弹窗
      deleteDialogVisible: false,
      deleteTaskTarget: null,
      deleteTableSummaries: [],
      deleteCleanRemote: false,
      deleteLoading: false,
      deleting: false,
      deleteHasInflight: false,
      deleteFetchError: false,
    };
  },
```

- [ ] **Step 2: 新增 `deleteTotalPartitions` computed**

在 `computed: { ... }` 中(在 `pagedTasks()` 之后、闭合 `}` 之前)新增:

```js
    deleteTotalPartitions() {
      return this.deleteTableSummaries.reduce((sum, s) => sum + (s.error ? 0 : s.count), 0);
    },
```

- [ ] **Step 3: 替换 `deleteTask` 方法,并新增 `loadDeleteSummary` / `confirmDeleteTask` / `onDeleteDialogClosed`**

将现有 `deleteTask(t)` 方法(约 287-311 行,整段 `async deleteTask(t) { ... }`)替换为以下四个方法:

```js
    deleteTask(t) {
      this.deleteTaskTarget = t;
      this.deleteCleanRemote = false;
      this.deleteTableSummaries = [];
      this.deleteHasInflight = false;
      this.deleteFetchError = false;
      this.deleteDialogVisible = true;
      this.loadDeleteSummary(t);
    },
    async loadDeleteSummary(t) {
      this.deleteLoading = true;
      this.deleteHasInflight = false;
      this.deleteFetchError = false;
      try {
        const summaries = await Promise.all(t.policies.map(async (p) => {
          const entry = {
            policy_id: p.policy_id,
            cluster_name: p.cluster_name,
            database: p.database,
            table: p.table,
            partitions: [],
            count: 0,
            error: false,
          };
          try {
            const res = await DataManageApi.listRunsByTable(p.cluster_name, p.database, p.table);
            if (res.data.retCode === '0000') {
              const runs = res.data.entity || [];
              const succ = new Set();
              for (const run of runs) {
                if (run.status === 'queued' || run.status === 'running') this.deleteHasInflight = true;
                if (run.operation === 'backup') {
                  for (const part of (run.partitions || [])) {
                    if (part.status === 'success') succ.add(part.partition);
                  }
                }
              }
              entry.partitions = [...succ];
              entry.count = succ.size;
            } else {
              entry.error = true;
              this.deleteFetchError = true;
            }
          } catch {
            entry.error = true;
            this.deleteFetchError = true;
          }
          return entry;
        }));
        // 防止用户在加载期间切换到另一任务时把结果错配
        if (this.deleteTaskTarget === t) {
          this.deleteTableSummaries = summaries;
        }
      } finally {
        this.deleteLoading = false;
      }
    },
    async confirmDeleteTask() {
      const t = this.deleteTaskTarget;
      if (!t) return;
      this.deleting = true;
      try {
        const warnings = [];
        if (this.deleteCleanRemote) {
          for (const s of this.deleteTableSummaries) {
            if (s.error || s.count === 0) continue;
            try {
              const res = await DataManageApi.deletePartitionRecords(
                s.cluster_name, s.database, s.table,
                { partitions: s.partitions, clean_remote: true }
              );
              if (res.data.retCode !== '0000') {
                warnings.push(`${s.database}.${s.table}: ${res.data.retMsg || 'failed'}`);
              } else {
                const w = (res.data.entity && res.data.entity.warnings) || [];
                warnings.push(...w);
              }
            } catch (e) {
              warnings.push(`${s.database}.${s.table}: ${e.message || 'error'}`);
            }
          }
        }
        const results = await Promise.allSettled(
          t.policies.map(p => DataManageApi.deletePolicy(p.policy_id))
        );
        const success = results.filter(r => r.status === 'fulfilled' && r.value.data.retCode === '0000').length;
        const failed = results.length - success;
        if (failed > 0) {
          this.$message.warning(this.$t('history.Task Delete Result Partial', { success, failed }));
        } else if (warnings.length > 0) {
          this.$message.warning(this.$t('history.Task Delete Clean Warnings', { count: warnings.length }));
        } else {
          this.$message.success(this.$t('history.Task Delete Result OK', { success }));
        }
        this.deleteDialogVisible = false;
        this.$emit('refresh');
      } finally {
        this.deleting = false;
      }
    },
    onDeleteDialogClosed() {
      this.deleteTaskTarget = null;
      this.deleteTableSummaries = [];
      this.deleteCleanRemote = false;
    },
```

- [ ] **Step 4: 在 template 中新增删除确认 `el-dialog`**

在 `task-list.vue` 模板里,`el-pagination` 之后、根 `</div>` 之前(约 113-114 行之间)插入:

```html
    <el-dialog
      :visible.sync="deleteDialogVisible"
      :title="$t('history.Delete Task Title')"
      width="560px"
      @closed="onDeleteDialogClosed"
    >
      <div v-loading="deleteLoading">
        <p
          style="margin:0 0 10px"
          v-html="$t('history.Delete Task Tip', {
            name: deleteTaskTarget ? displayName(deleteTaskTarget) : '',
            count: deleteTaskTarget ? deleteTaskTarget.policies.length : 0,
          })"
        ></p>

        <el-alert
          v-if="deleteHasInflight"
          :title="$t('history.Delete Task Inflight Warning')"
          type="error"
          :closable="false"
          show-icon
          style="margin-bottom:10px"
        />
        <el-alert
          v-else-if="deleteFetchError"
          :title="$t('history.Delete Task Fetch Error')"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom:10px"
        />

        <el-table :data="deleteTableSummaries" size="small" border max-height="260" style="width:100%">
          <el-table-column :label="$t('history.Delete Task Col Table')" min-width="220" show-overflow-tooltip>
            <template #default="{ row }">{{ row.database }}.{{ row.table }}</template>
          </el-table-column>
          <el-table-column :label="$t('history.Delete Task Col Backed Partitions')" width="140" align="right">
            <template #default="{ row }">
              <span v-if="row.error" class="muted">-</span>
              <span v-else>{{ row.count }}</span>
            </template>
          </el-table-column>
        </el-table>

        <div class="muted" style="margin:10px 0">
          {{ $t('history.Delete Task Total', { tables: deleteTableSummaries.length, partitions: deleteTotalPartitions }) }}
        </div>

        <el-checkbox
          v-model="deleteCleanRemote"
          :disabled="deleteHasInflight || deleteFetchError"
        >{{ $t('history.Delete Task Clean Checkbox') }}</el-checkbox>
      </div>

      <span slot="footer">
        <el-button @click="deleteDialogVisible = false">{{ $t('common.Cancel') }}</el-button>
        <el-button
          type="danger"
          :loading="deleting"
          :disabled="deleteLoading || deleteHasInflight"
          @click="confirmDeleteTask"
        >{{ $t('history.Confirm Delete Btn') }}</el-button>
      </span>
    </el-dialog>
```

- [ ] **Step 5: lint 校验**

Run: `cd /data/root/go/src/github.com/housepower/ckman-fe && make lint`
Expected: `DONE  No lint errors found!`

- [ ] **Step 6: 生产构建校验**

Run: `cd /data/root/go/src/github.com/housepower/ckman-fe && make build`
Expected: `DONE  Build complete. The dist directory is ready to be deployed.`

- [ ] **Step 7: 手工验证(开发服务器 `make dev`,逐项确认)**

| # | 操作 | 期望 |
|---|---|---|
| 1 | 任务含多张表、各有已备份分区,点「删除」 | 弹窗逐行列出 `db.table`,每行显示该表已成功备份分区个数;底部合计「N 张表 / M 个分区」正确 |
| 2 | **不勾选**直接确认 | 仅软删 policy;S3/本地盘备份文件仍在;列表刷新后任务消失;提示「已删除 N 张表的备份配置」 |
| 3 | **勾选**「同步删除已备份分区」后确认 | 各表已备份分区的远端数据被清除、台账被清理;policy 软删;列表刷新 |
| 4 | 任务下有进行中(running/queued)run | 顶部红色提示;勾选框禁用;「删除」确认按钮禁用,无法提交 |
| 5 | 某表 0 已备份分区 | 该行显示 `0`;勾选确认时跳过该表(不报错);其余表正常清理 |
| 6 | 模拟某表 runs 拉取失败(可临时断网/改错表名场景) | 该行个数显示 `-`;黄色提示;勾选框禁用;仍可不勾选执行软删 |
| 7 | 关闭弹窗再打开另一任务 | 勾选框重置为未勾选;表格为新任务数据,无残留 |

- [ ] **Step 8: 提交**

```bash
cd /data/root/go/src/github.com/housepower/ckman-fe
git add src/views/data-manage/component/task-list.vue
git commit -m "feat(backup): 删除任务弹窗列出关联表/分区数并支持同步删除已备份分区"
```

---

## Task 3: 推送前端并更新父仓 submodule 指针

**Files:**
- Modify: `ckman/frontend`(submodule 指针)

- [ ] **Step 1: 推送 ckman-fe**

Run:
```bash
cd /data/root/go/src/github.com/housepower/ckman-fe && git push
```
Expected: 推送成功,本地 `main` 与远端一致。

- [ ] **Step 2: 更新父仓 submodule 指针**

Run(将 `<sha>` 替换为 Task 2 提交后 `ckman-fe` 的 `git rev-parse HEAD`):
```bash
cd /data/root/go/src/github.com/housepower/ckman/frontend && git fetch origin && git checkout <sha>
```
Expected: `HEAD is now at <sha> feat(backup): 删除任务弹窗...`

- [ ] **Step 3: 父仓提交 submodule 指针**

Run:
```bash
cd /data/root/go/src/github.com/housepower/ckman
git add frontend
git commit -m "chore(frontend): 更新 submodule 指针 → <sha>(删除任务支持同步删除已备份分区)"
```
Expected: 提交成功(仅 `frontend` 一处变更)。

---

## Self-Review

**1. Spec coverage:**
- 列出关联表 + 每表已备份分区个数 → Task 2 Step 4 表格 + Task 2 Step 3 `loadDeleteSummary` 聚合 ✓
- 「同步删除已备份分区」勾选框,默认不勾选 → Task 2 Step 1(`deleteCleanRemote: false`)+ Step 4 复选框 ✓
- 勾选后真正删除远端数据 → Task 2 Step 3 `confirmDeleteTask` 调 `deletePartitionRecords(clean_remote: true)` ✓
- 进行中 run 阻止删除并提示 → Task 2 Step 3(`deleteHasInflight`)+ Step 4(`el-alert` + 按钮 `:disabled`)✓
- 不勾选维持现状(仅软删)→ Task 2 Step 3(`deleteCleanRemote` 为 false 时跳过清理,直接 `deletePolicy`)✓
- 边界:0 分区跳过、拉取失败禁用勾选、legacy policy → Task 2 Step 3 + Step 7 手工项 5/6 ✓
- 后端零改动 → 全计划无 Go 文件 ✓

**2. Placeholder scan:** 无 TBD/TODO;Task 3 的 `<sha>` 是运行期产物,已在步骤中说明如何取得,非占位代码。✓

**3. Type consistency:** `deleteTableSummaries` 元素结构(`{policy_id, cluster_name, database, table, partitions, count, error}`)在 `loadDeleteSummary`(生产)、`deleteTotalPartitions`(消费)、模板表格(消费)、`confirmDeleteTask`(消费)中字段名一致;`deleteCleanRemote` / `deleteHasInflight` / `deleteFetchError` / `deleting` / `deleteLoading` 命名在 data、模板、方法间一致。✓
