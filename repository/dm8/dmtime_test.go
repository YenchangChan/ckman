package dm8

import (
	"testing"
	"time"
)

// 零值时间必须落成哨兵纪元(而非 0001-01-01),否则达梦报 -6118 非法的时间日期类型数据。
func TestDmTimeValue_ZeroToSentinel(t *testing.T) {
	v, err := DmTime(time.Time{}).Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := v.(time.Time)
	if !ok || !got.Equal(dmZeroSentinel) {
		t.Fatalf("zero time should map to sentinel %v, got %v", dmZeroSentinel, v)
	}
}

// 非零时间原样写入。
func TestDmTimeValue_NonZero(t *testing.T) {
	now := time.Date(2026, 6, 29, 11, 55, 0, 0, time.UTC)
	v, err := DmTime(now).Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := v.(time.Time)
	if !ok || !got.Equal(now) {
		t.Fatalf("want %v, got %v", now, v)
	}
}

// 哨兵读回还原为 Go 零值时间,与 model 层 time.Time 语义一致。
func TestDmTimeScan_SentinelToZero(t *testing.T) {
	var dt DmTime
	if err := dt.Scan(dmZeroSentinel); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dt.Time().IsZero() {
		t.Fatalf("sentinel should scan back to zero time, got %v", dt.Time())
	}
}

// 历史遗留的 NULL 也兼容,还原为零值时间。
func TestDmTimeScan_NullToZero(t *testing.T) {
	var dt DmTime
	if err := dt.Scan(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dt.Time().IsZero() {
		t.Fatalf("NULL should scan to zero time, got %v", dt.Time())
	}
}

func TestDmTimeScan_Time(t *testing.T) {
	now := time.Date(2026, 6, 29, 11, 55, 0, 0, time.UTC)
	var dt DmTime
	if err := dt.Scan(now); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dt.Time().Equal(now) {
		t.Fatalf("want %v, got %v", now, dt.Time())
	}
}

// 零值时间经 Value→Scan 往返后仍是零值,保证写入/读回闭环。
func TestDmTime_ZeroRoundTrip(t *testing.T) {
	v, err := DmTime(time.Time{}).Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var dt DmTime
	if err := dt.Scan(v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dt.Time().IsZero() {
		t.Fatalf("zero time should round-trip to zero, got %v", dt.Time())
	}
}
