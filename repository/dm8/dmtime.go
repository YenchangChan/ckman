package dm8

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// dmZeroSentinel 是 Go 零值时间在达梦列里的替身。
//
// 达梦不接受 Go 的零值时间 0001-01-01T00:00:00Z,直接写入会报
// `Error -6118: 非法的时间日期类型数据`。用 Unix 纪元 1970-01-01T00:00:00Z
// 作为哨兵替身(达梦合法日期),而不是写 NULL,有两个好处:
//  1. DESC 排序时自然排在最后,无需依赖 `NULLS LAST` 这类方言语法;
//  2. `started_at < ?` 这类范围过滤仍能命中,分页/过滤行为与 MySQL/Postgres
//     等后端(存 0001-01-01,同样排最后且能被 < 命中)保持一致。
//
// 真实业务中的备份时间都是近期值,几乎不可能恰好等于 1970-01-01,故哨兵不会与
// 真实数据冲突;即便冲突,读回时也只是被当作零值时间还原(见 Scan),而 started_at
// 列只作索引/排序/过滤用,权威值始终存于 JSON Run blob,无副作用。
var dmZeroSentinel = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// DmTime 是达梦时间列的统一包装类型,把 Go 零值时间在 DB 边界上 ⇄ 哨兵 对称转换,
// 使所有时间字段(现有的和将来新增的)都免疫 -6118。
//
// 用法:把 dm8 model 里所有 time.Time 列改成 DmTime;在 Updates(map) 这类绕过
// 结构体字段类型的写入处,也用 DmTime(t) 包一层,使 driver.Valuer 生效。
// GormDataType 返回 "time",保证 AutoMigrate 生成的列类型(TIMESTAMP WITH TIME
// ZONE)与未包装前的 time.Time 列逐字节一致,不会触发列变更。
type DmTime time.Time

// GormDataType 让 GORM 把该字段识别为时间类型,迁移时映射到达梦的
// TIMESTAMP WITH TIME ZONE,与未包装前的 time.Time 列完全等价。
func (DmTime) GormDataType() string { return "time" }

// Value 实现 driver.Valuer:零值时间写哨兵纪元,其余按原值写入。
func (t DmTime) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return dmZeroSentinel, nil
	}
	return tt, nil
}

// Scan 实现 sql.Scanner:哨兵纪元(及历史遗留的 NULL)还原为 Go 零值时间。
//
// 达梦驱动对 TIMESTAMP 列返回的是 time.Time(改造前普通 time.Time 字段能正常
// 读出已证明这一点),故只需处理 time.Time 与 NULL 两种来源。
func (t *DmTime) Scan(v any) error {
	switch x := v.(type) {
	case nil:
		*t = DmTime(time.Time{})
	case time.Time:
		if x.Equal(dmZeroSentinel) {
			*t = DmTime(time.Time{})
		} else {
			*t = DmTime(x)
		}
	default:
		return fmt.Errorf("dm8: cannot scan %T into DmTime", v)
	}
	return nil
}

// Time 返回底层 time.Time,便于读出后赋值给 model 层的 time.Time 字段。
func (t DmTime) Time() time.Time { return time.Time(t) }
