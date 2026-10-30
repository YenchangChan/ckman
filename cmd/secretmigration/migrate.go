package secretmigration

import (
	"fmt"
	"sort"

	"github.com/housepower/ckman/common"
	"github.com/housepower/ckman/config"
	"github.com/housepower/ckman/repository"
	_ "github.com/housepower/ckman/repository/dm8"
	_ "github.com/housepower/ckman/repository/mysql"
	_ "github.com/housepower/ckman/repository/postgres"
	_ "github.com/housepower/ckman/repository/sqlite"
)

// Run rewrites cluster secrets using SM4 in a single persistent-layer transaction.
// Run this while ckman writers are stopped so a concurrent update cannot be lost.
func Run(configFile string, dryRun bool) (int, error) {
	if err := config.ParseConfigFile(configFile, ""); err != nil {
		return 0, err
	}
	if config.GlobalConfig.EncryptAlgorithm != common.GosyptSM4Algorithm {
		return 0, fmt.Errorf("encrypt_algorithm must be SM4 for migration")
	}
	if err := common.Gsypt.SetAlgorithm(config.GlobalConfig.EncryptAlgorithm); err != nil {
		return 0, err
	}
	if err := common.Gsypt.Unmarshal(&config.GlobalConfig); err != nil {
		return 0, err
	}
	policy := config.GlobalConfig.Server.PersistentPolicy
	ps := repository.GetPersistentByName(policy)
	if ps == nil {
		return 0, fmt.Errorf("unknown persistent policy %q", policy)
	}
	var backendConfig interface{}
	if c, ok := config.GlobalConfig.PersistentConfig[policy]; ok {
		backendConfig = ps.UnmarshalConfig(c)
	}
	if err := ps.Init(backendConfig); err != nil {
		return 0, err
	}
	clusters, err := ps.GetAllClusters()
	if err != nil {
		return 0, err
	}
	if dryRun {
		return len(clusters), nil
	}
	if err := ps.Begin(); err != nil {
		return 0, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = ps.Rollback()
		}
	}()
	names := make([]string, 0, len(clusters))
	for name := range clusters {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := ps.UpdateCluster(clusters[name]); err != nil {
			return 0, fmt.Errorf("update cluster %s: %w", name, err)
		}
	}
	if err := ps.Commit(); err != nil {
		return 0, err
	}
	committed = true
	return len(clusters), nil
}
