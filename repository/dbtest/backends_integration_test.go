//go:build dbintegration

// Package dbtest 是持久化后端的真库集成测试。
//
// 这类 bug(达梦 -6118 零值时间 / -2116 重复加列 / -2140 重复建索引,以及 MySQL
// DATETIME 下界、PG 大小写等)只有「对着真库连启多次 + 真写零值/边界数据」才现形,
// mock 单测一律抓不到。故独立成一个带 build tag 的集成测试,平时不跑,给库就跑:
//
//	DBTEST_MYSQL_HOST=... DBTEST_MYSQL_USER=... DBTEST_MYSQL_PASS=... DBTEST_MYSQL_DB=... \
//	DBTEST_DM_HOST=...    DBTEST_DM_USER=...    DBTEST_DM_PASS=...    DBTEST_DM_SCHEMA=... \
//	DBTEST_PG_HOST=...    DBTEST_PG_USER=...    DBTEST_PG_PASS=...    DBTEST_PG_DB=... \
//	go test -tags dbintegration -v ./repository/dbtest/ -run TestBackends
//
// 未设置对应环境变量的后端会被 Skip。
package dbtest

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/housepower/ckman/log"
	"github.com/housepower/ckman/model"
	"github.com/housepower/ckman/repository"
	"github.com/housepower/ckman/repository/dm8"
	"github.com/housepower/ckman/repository/mysql"
	"github.com/housepower/ckman/repository/postgres"
)

func init() { log.InitLoggerConsole() }

func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// backendCase 描述一个待测后端:如何造新实例、连接配置、是否配置齐全。
type backendCase struct {
	name string
	newFn func() repository.PersistentMgr
	cfg  interface{}
	ok   bool
}

func mysqlCase() backendCase {
	host := os.Getenv("DBTEST_MYSQL_HOST")
	return backendCase{
		name:  "mysql",
		newFn: func() repository.PersistentMgr { return mysql.NewMysqlPersistent() },
		cfg: mysql.MysqlConfig{
			Host: host, Port: envInt("DBTEST_MYSQL_PORT", 3306),
			User: os.Getenv("DBTEST_MYSQL_USER"), Password: os.Getenv("DBTEST_MYSQL_PASS"),
			DataBase: os.Getenv("DBTEST_MYSQL_DB"),
		},
		ok: host != "",
	}
}

func dmCase() backendCase {
	host := os.Getenv("DBTEST_DM_HOST")
	return backendCase{
		name:  "dm8",
		newFn: func() repository.PersistentMgr { return dm8.NewDM8Persistent() },
		cfg: dm8.DM8Config{
			Host: host, Port: envInt("DBTEST_DM_PORT", 5236),
			User: os.Getenv("DBTEST_DM_USER"), Password: os.Getenv("DBTEST_DM_PASS"),
			Schema: os.Getenv("DBTEST_DM_SCHEMA"),
		},
		ok: host != "",
	}
}

func pgCase() backendCase {
	host := os.Getenv("DBTEST_PG_HOST")
	return backendCase{
		name:  "postgres",
		newFn: func() repository.PersistentMgr { return postgres.NewPostgresPersistent() },
		cfg: postgres.PostgresConfig{
			Host: host, Port: envInt("DBTEST_PG_PORT", 5432),
			User: os.Getenv("DBTEST_PG_USER"), Password: os.Getenv("DBTEST_PG_PASS"),
			DataBase: os.Getenv("DBTEST_PG_DB"),
		},
		ok: host != "",
	}
}

// TestBackends 对每个配置齐全的后端跑一遍完整套路。
func TestBackends(t *testing.T) {
	for _, bc := range []backendCase{mysqlCase(), dmCase(), pgCase()} {
		bc := bc
		t.Run(bc.name, func(t *testing.T) {
			if !bc.ok {
				t.Skipf("%s 环境变量未配置,跳过", bc.name)
			}
			runBackend(t, bc)
		})
	}
}

func runBackend(t *testing.T, bc backendCase) {
	// ── 雷点 1:连启多次,AutoMigrate 幂等性(达梦 -2116 重复加列 / -2140 重复建索引)──
	// 首次建表 OK 不代表没问题;这类 bug 只在第 2/3 次启动重跑 AutoMigrate 时爆。
	for i := 1; i <= 3; i++ {
		m := bc.newFn()
		if err := m.Init(bc.cfg); err != nil {
			t.Fatalf("boot#%d Init 失败(AutoMigrate 非幂等?): %v", i, err)
		}
	}

	mgr := bc.newFn()
	if err := mgr.Init(bc.cfg); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}
	// 让内部依赖 repository.Ps 的代码路径也指向被测实例
	repository.Ps = mgr

	t.Run("cluster_crud", func(t *testing.T) { testClusterCRUD(t, mgr) })
	t.Run("backup_policy_crud", func(t *testing.T) { testPolicyCRUD(t, mgr) })
	t.Run("backup_run_zero_time", func(t *testing.T) { testBackupRunZeroTime(t, mgr) })
	t.Run("user_crud", func(t *testing.T) { testUserCRUD(t, mgr) })
	t.Run("query_history", func(t *testing.T) { testQueryHistory(t, mgr) })
	t.Run("task_crud", func(t *testing.T) { testTaskCRUD(t, mgr) })
}

const tagPrefix = "dbtest_"

func testClusterCRUD(t *testing.T, mgr repository.PersistentMgr) {
	name := tagPrefix + "cluster"
	_ = mgr.DeleteCluster(name)
	conf := model.CKManClickHouseConfig{
		Cluster: name, Port: 9000, User: "default", Password: "secret",
		Hosts: []string{"1.2.3.4"},
	}
	if err := mgr.CreateCluster(conf); err != nil {
		t.Fatalf("CreateCluster: %v", err)
	}
	defer mgr.DeleteCluster(name)

	got, err := mgr.GetClusterbyName(name)
	if err != nil {
		t.Fatalf("GetClusterbyName: %v", err)
	}
	if got.Cluster != name {
		t.Fatalf("round-trip cluster name = %q", got.Cluster)
	}
	conf.Port = 9440
	if err := mgr.UpdateCluster(conf); err != nil {
		t.Fatalf("UpdateCluster: %v", err)
	}
	if got, _ = mgr.GetClusterbyName(name); got.Port != 9440 {
		t.Fatalf("update not persisted, port = %d", got.Port)
	}
}

func testPolicyCRUD(t *testing.T, mgr repository.PersistentMgr) {
	id := tagPrefix + "policy"
	_ = mgr.HardDeleteBackupPolicy(id)
	p := model.BackupPolicy{
		PolicyID: id, ClusterName: tagPrefix + "cluster", Database: "db", Table: "t",
		ScheduleType: "scheduled", Instance: "10.0.0.1:8808", BackupType: "daily",
		Enabled: true, CreateTime: time.Now(), UpdateTime: time.Now(),
	}
	if err := mgr.CreateBackupPolicy(p); err != nil {
		t.Fatalf("CreateBackupPolicy: %v", err)
	}
	defer mgr.HardDeleteBackupPolicy(id)
	if _, err := mgr.GetBackupPolicy(id); err != nil {
		t.Fatalf("GetBackupPolicy: %v", err)
	}
	p.Enabled = false
	if err := mgr.UpdateBackupPolicy(p); err != nil {
		t.Fatalf("UpdateBackupPolicy: %v", err)
	}
	if _, err := mgr.GetBackupPoliciesByCluster(p.ClusterName); err != nil {
		t.Fatalf("GetBackupPoliciesByCluster: %v", err)
	}
}

// testBackupRunZeroTime 是核心雷点:queued/skipped 状态的 run 的 StartedAt/FinishedAt
// 是 Go 零值时间 0001-01-01。达梦直接报 -6118;MySQL DATETIME 下界是 1000-01-01,
// 0001-01-01 越界也可能报错/被塞成 0000-00-00;PG 时间范围够大通常没事。
func testBackupRunZeroTime(t *testing.T, mgr repository.PersistentMgr) {
	id := tagPrefix + "run"
	_ = mgr.DeleteBackupRun(id)
	run := model.BackupRun{
		RunID: id, PolicyID: tagPrefix + "policy", ClusterName: tagPrefix + "cluster",
		Database: "db", Table: "t", Operation: "backup", TriggerType: "cron",
		Instance: "10.0.0.1:8808", Status: model.BACKUP_STATUS_QUEUED,
		Partitions: []model.BackupRunPartition{{Partition: "20260101", Status: model.BACKUP_PARTITION_STATUS_WAITING}},
		StartedAt:  time.Time{}, // 零值 —— 雷点
		FinishedAt: time.Time{}, // 零值 —— 雷点
		CreateTime: time.Now(),
	}
	if err := mgr.CreateBackupRun(run); err != nil {
		t.Fatalf("CreateBackupRun(零值时间): %v", err)
	}
	defer mgr.DeleteBackupRun(id)

	got, err := mgr.GetBackupRun(id)
	if err != nil {
		t.Fatalf("GetBackupRun: %v", err)
	}
	if got.RunID != id || len(got.Partitions) != 1 {
		t.Fatalf("round-trip run mismatch: %+v", got)
	}
	// before=零值:测试 started_at < ? 之类分页/过滤路径不被零值时间打挂
	if _, err := mgr.GetRunsByPolicy(run.PolicyID, 10, time.Time{}); err != nil {
		t.Fatalf("GetRunsByPolicy(before=zero): %v", err)
	}
	if _, err := mgr.GetRunsByPolicy(run.PolicyID, 10, time.Now()); err != nil {
		t.Fatalf("GetRunsByPolicy(before=now): %v", err)
	}
	// queued → running:把零值 started_at 更新成真实时间
	ok, err := mgr.MarkRunRunningIfQueued(id, "10.0.0.1:8808", time.Now())
	if err != nil {
		t.Fatalf("MarkRunRunningIfQueued: %v", err)
	}
	if !ok {
		t.Fatalf("MarkRunRunningIfQueued 应把 queued 置为 running")
	}
}

func testUserCRUD(t *testing.T, mgr repository.PersistentMgr) {
	name := tagPrefix + "user"
	_ = mgr.DeleteUser(name)
	u := model.CkmanUser{
		Username: name, PasswordHash: "hash", Policy: "ordinary", Enabled: true,
		CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(),
	}
	if err := mgr.CreateUser(u); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	defer mgr.DeleteUser(name)
	if !mgr.UserExists(name) {
		t.Fatalf("UserExists=false 刚建的用户")
	}
	if _, err := mgr.GetUserByName(name); err != nil {
		t.Fatalf("GetUserByName: %v", err)
	}
	u.Policy = "guest"
	if err := mgr.UpdateUser(u); err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
}

func testQueryHistory(t *testing.T, mgr repository.PersistentMgr) {
	cs := tagPrefix + "qh"
	_ = mgr.DeleteQueryHistory(cs)
	qh := model.QueryHistory{Cluster: tagPrefix + "cluster", CheckSum: cs, QuerySql: "SELECT 1", CreateTime: time.Now()}
	if err := mgr.CreateQueryHistory(qh); err != nil {
		t.Fatalf("CreateQueryHistory: %v", err)
	}
	defer mgr.DeleteQueryHistory(cs)
	if _, err := mgr.GetQueryHistoryByCheckSum(cs); err != nil {
		t.Fatalf("GetQueryHistoryByCheckSum: %v", err)
	}
}

func testTaskCRUD(t *testing.T, mgr repository.PersistentMgr) {
	id := tagPrefix + "task"
	_ = mgr.DeleteTask(id)
	task := model.Task{
		TaskId: id, ClusterName: tagPrefix + "cluster", ServerIp: "10.0.0.1",
		TaskType: "deploy", Status: model.TaskStatusWaiting,
		CreateTime: time.Now(), UpdateTime: time.Now(),
	}
	if err := mgr.CreateTask(task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	defer mgr.DeleteTask(id)
	if _, err := mgr.GetTaskbyTaskId(id); err != nil {
		t.Fatalf("GetTaskbyTaskId: %v", err)
	}
}
