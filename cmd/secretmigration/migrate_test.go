package secretmigration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/housepower/ckman/common"
	"github.com/housepower/ckman/config"
	"github.com/housepower/ckman/log"
	"github.com/housepower/ckman/model"
	"github.com/housepower/ckman/repository/sqlite"
	"github.com/stretchr/testify/require"
)

func TestMigrateSQLiteClusterSecrets(t *testing.T) {
	log.InitLoggerConsole()
	t.Setenv(common.SM4KeyEnv, "0123456789abcdef")
	previous := config.GlobalConfig
	t.Cleanup(func() { config.GlobalConfig = previous; _ = common.Gsypt.SetAlgorithm(common.GosyptAlgorithm) })
	config.GlobalConfig = config.CKManConfig{EncryptAlgorithm: common.GosyptAlgorithm}
	dir := t.TempDir()
	sp := sqlite.NewSQLitePersistent()
	require.NoError(t, sp.Init(sqlite.LocalConfig{ConfigDir: dir, ConfigFile: "clusters"}))
	want := model.CKManClickHouseConfig{
		Cluster: "test", Password: "cluster-pass", SshPassword: "ssh-pass",
		TlsConf:   &model.TlsConf{Key: "-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----\n"},
		UsersConf: model.UsersConf{Users: []model.User{{Name: "alice", Password: "user-pass"}}},
		Storage:   &model.Storage{Disks: []model.Disk{{Type: "s3", DiskS3: &model.DiskS3{AccessKeyID: "access", SecretAccessKey: "secret"}}}},
	}
	require.NoError(t, sp.CreateCluster(want))
	var before sqlite.TblCluster
	require.NoError(t, sp.Client.Where("cluster_name = ?", want.Cluster).First(&before).Error)
	require.NotContains(t, before.Config, common.SM4SecretPrefix)

	confPath := filepath.Join(dir, "ckman.json")
	conf := fmt.Sprintf(`{"encrypt_algorithm":"SM4","server":{"persistent_policy":"local"},"persistent_config":{"local":{"config_dir":%q,"config_file":"clusters"}}}`, dir)
	require.NoError(t, os.WriteFile(confPath, []byte(conf), 0600))
	count, err := Run(confPath, true)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	var dry sqlite.TblCluster
	require.NoError(t, sp.Client.Where("cluster_name = ?", want.Cluster).First(&dry).Error)
	require.Equal(t, before.Config, dry.Config)

	count, err = Run(confPath, false)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	var after sqlite.TblCluster
	require.NoError(t, sp.Client.Where("cluster_name = ?", want.Cluster).First(&after).Error)
	var stored model.CKManClickHouseConfig
	require.NoError(t, json.Unmarshal([]byte(after.Config), &stored))
	for _, secret := range []string{stored.Password, stored.SshPassword, stored.UsersConf.Users[0].Password,
		stored.Storage.Disks[0].DiskS3.AccessKeyID, stored.Storage.Disks[0].DiskS3.SecretAccessKey, stored.TlsConf.Key} {
		require.True(t, strings.HasPrefix(secret, common.SM4SecretPrefix), "ciphertext should be SM4")
	}
	got, err := sp.GetClusterbyName(want.Cluster)
	require.NoError(t, err)
	require.Equal(t, want.Password, got.Password)
	require.Equal(t, want.SshPassword, got.SshPassword)
	require.Equal(t, want.UsersConf.Users[0].Password, got.UsersConf.Users[0].Password)
	require.Equal(t, want.Storage.Disks[0].DiskS3.AccessKeyID, got.Storage.Disks[0].DiskS3.AccessKeyID)
	require.Equal(t, want.Storage.Disks[0].DiskS3.SecretAccessKey, got.Storage.Disks[0].DiskS3.SecretAccessKey)
	require.Equal(t, want.TlsConf.Key, got.TlsConf.Key)

	t.Setenv(common.SM4KeyEnv, "fedcba9876543210")
	_, err = Run(confPath, true)
	require.Error(t, err)
}
