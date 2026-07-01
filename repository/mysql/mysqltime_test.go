package mysql

import (
	"testing"
	"time"
)

func TestMysqlTimeValueZeroToSentinel(t *testing.T) {
	v, err := MysqlTime(time.Time{}).Value()
	if err != nil {
		t.Fatal(err)
	}
	if !v.(time.Time).Equal(mysqlZeroSentinel) {
		t.Fatalf("零值应映射为哨兵 1970-01-01,got %v", v)
	}
	// 非零值原样返回
	now := time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)
	v, _ = MysqlTime(now).Value()
	if !v.(time.Time).Equal(now) {
		t.Fatalf("非零值应原样,got %v", v)
	}
}

func TestMysqlTimeScanSentinelToZero(t *testing.T) {
	var mt MysqlTime
	if err := mt.Scan(mysqlZeroSentinel); err != nil {
		t.Fatal(err)
	}
	if !mt.Time().IsZero() {
		t.Fatalf("哨兵应还原为零值,got %v", mt.Time())
	}
	if err := mt.Scan(nil); err != nil || !mt.Time().IsZero() {
		t.Fatalf("NULL 应还原为零值,got %v err=%v", mt.Time(), err)
	}
}
