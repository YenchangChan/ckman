package mysql

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// mysqlZeroSentinel 是 Go 零值时间在 MySQL 列里的替身。
//
// MySQL 严格模式(NO_ZERO_DATE/NO_ZERO_IN_DATE)下,Go 零值时间 0001-01-01 会被
// 驱动渲染成 '0000-00-00',直接写入报 `Error 1292: Incorrect datetime value`;
// 即便非严格模式,DATETIME 下界也只到 1000-01-01,0001-01-01 同样越界。用 Unix
// 纪元 1970-01-01(MySQL 合法 DATETIME)作哨兵,而非 NULL,好处同 dm8 的 DmTime:
//  1. DESC 排序天然排最后,不依赖 NULLS LAST 方言;
//  2. `started_at < ?` 范围过滤仍能命中,分页/过滤行为与其它后端一致。
//
// 真实备份时间都是近期值,几乎不可能恰好等于 1970-01-01;即便冲突,读回时也只是被
// 当作零值还原(见 Scan),而 started_at 列只作索引/排序/过滤,权威值始终存于 JSON
// Run blob,无副作用。
var mysqlZeroSentinel = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// MysqlTime 是 MySQL 时间列的包装类型,把 Go 零值时间在 DB 边界上 ⇄ 哨兵 对称转换,
// 使 queued/skipped 状态 run 的零值 started_at 免疫 Error 1292。
//
// 用法:把可能取零值的时间列(started_at)改成 MysqlTime;在 Updates(map) 这类绕过
// 结构体字段类型的写入处,也用 MysqlTime(t) 包一层使 driver.Valuer 生效。GormDataType
// 返回 "time",迁移生成的列类型(datetime(3))与未包装前的 time.Time 列一致,不触发列变更。
type MysqlTime time.Time

// GormDataType 让 GORM 把该字段识别为时间类型,mysql 方言映射到 datetime(3),
// 与未包装前的 time.Time 列完全等价。
func (MysqlTime) GormDataType() string { return "time" }

// Value 实现 driver.Valuer:零值时间写哨兵纪元,其余按原值写入。
func (t MysqlTime) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() {
		return mysqlZeroSentinel, nil
	}
	return tt, nil
}

// Scan 实现 sql.Scanner:哨兵纪元(及历史遗留 NULL)还原为 Go 零值时间。
func (t *MysqlTime) Scan(v any) error {
	switch x := v.(type) {
	case nil:
		*t = MysqlTime(time.Time{})
	case time.Time:
		if x.Equal(mysqlZeroSentinel) {
			*t = MysqlTime(time.Time{})
		} else {
			*t = MysqlTime(x)
		}
	default:
		return fmt.Errorf("mysql: cannot scan %T into MysqlTime", v)
	}
	return nil
}

// Time 返回底层 time.Time。
func (t MysqlTime) Time() time.Time { return time.Time(t) }
