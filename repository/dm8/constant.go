package dm8

const (
	DM8PersistentName          string = "dm8"
	DM8_PORT_DEFAULT           int    = 5236
	DM8_USER_DEFAULT           string = "ckman"
	DM8_PASSWD_DEFAULT         string = "123456789"
	DM8_MAX_IDLE_CONNS_DEFAULT int    = 10
	DM8_MAX_OPEN_CONNS_DEFAULT int    = 100
	DM8_MAX_LIFETIME_DEFAULT   int    = 3600
	DM8_MAX_IDLE_TIME_DEFAULT  int    = 10

	// DM8_TBL_USER 必须用复数 tbl_users:DM8 其它表都靠 gorm 默认命名被复数化
	// (tbl_clusters/tbl_backup_policies…),而 vendored gorm-dm8 的 HasIndex 在
	// 判断索引是否已存在时,会把表名再过一遍 NamingStrategy.TableName()——它会做
	// 复数化。若这里固定成单数 tbl_user,DM 里真实表是 TBL_USER,HasIndex 却按
	// 复数 TBL_USERS 去查 USER_INDEXES,查不到 → 误判索引不存在 → 每次重启重复
	// CREATE UNIQUE INDEX → 报 -2140「索引已存在」导致 init persistent 失败。
	// 用复数名后再复数化是幂等的,索引检测才能命中。切勿改回单数。
	DM8_TBL_USER string = "tbl_users"
)
