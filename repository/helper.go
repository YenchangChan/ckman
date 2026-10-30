package repository

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/housepower/ckman/common"
	"github.com/housepower/ckman/config"
	"github.com/housepower/ckman/log"
	"github.com/housepower/ckman/model"
)

var (
	ErrRecordNotFound   = fmt.Errorf("record not found")
	ErrRecordExists     = fmt.Errorf("record is exists already")
	ErrTransActionBegin = fmt.Errorf("transaction already begin")
	ErrTransActionEnd   = fmt.Errorf("transaction already commit or rollback")
)

// EncodePasswd encrypts every secret in a cluster before persistence.
func EncodePasswd(conf *model.CKManClickHouseConfig) error {
	encode := func(field string, value *string) error {
		ciphertext, err := common.EncryptSecret(*value, config.GlobalConfig.EncryptAlgorithm)
		if err != nil {
			return fmt.Errorf("cluster %s %s: %w", conf.Cluster, field, err)
		}
		*value = ciphertext
		return nil
	}
	if err := encode("password", &conf.Password); err != nil {
		return err
	}
	if err := encode("sshPassword", &conf.SshPassword); err != nil {
		return err
	}
	conf.UsersConf.Users = append([]model.User(nil), conf.UsersConf.Users...)
	for i := range conf.UsersConf.Users {
		if err := encode(fmt.Sprintf("user[%d].password", i), &conf.UsersConf.Users[i].Password); err != nil {
			return err
		}
	}
	if conf.Storage != nil {
		storage := *conf.Storage
		storage.Disks = append([]model.Disk(nil), conf.Storage.Disks...)
		for i := range storage.Disks {
			disk := &storage.Disks[i]
			if disk.Type == "s3" && disk.DiskS3 != nil {
				s3 := *disk.DiskS3
				disk.DiskS3 = &s3
				if err := encode(fmt.Sprintf("disk[%d].accessKeyID", i), &s3.AccessKeyID); err != nil {
					return err
				}
				if err := encode(fmt.Sprintf("disk[%d].secretAccessKey", i), &s3.SecretAccessKey); err != nil {
					return err
				}
			}
		}
		conf.Storage = &storage
	}
	if conf.TlsConf != nil {
		// A fresh object prevents shared in-memory TLS pointers from being changed.
		tls := *conf.TlsConf
		if tls.Key != "" {
			tls.Key = base64.StdEncoding.EncodeToString([]byte(tls.Key))
			if err := encode("tlsConf.key", &tls.Key); err != nil {
				return err
			}
		}
		conf.TlsConf = &tls
	}
	return nil
}

// DecodePasswd rejects corrupt or undecryptable ciphertext before it reaches callers.
func DecodePasswd(conf *model.CKManClickHouseConfig) error {
	decode := func(field string, value *string) error {
		plain, err := common.DecryptSecret(*value)
		if err != nil {
			return fmt.Errorf("cluster %s %s: %w", conf.Cluster, field, err)
		}
		*value = plain
		return nil
	}
	if err := decode("password", &conf.Password); err != nil {
		return err
	}
	if err := decode("sshPassword", &conf.SshPassword); err != nil {
		return err
	}
	users := append([]model.User(nil), conf.UsersConf.Users...)
	for i := range users {
		if err := decode(fmt.Sprintf("user[%d].password", i), &users[i].Password); err != nil {
			return err
		}
	}
	conf.UsersConf.Users = users
	if conf.Storage != nil {
		storage := *conf.Storage
		storage.Disks = append([]model.Disk(nil), conf.Storage.Disks...)
		for i := range storage.Disks {
			disk := &storage.Disks[i]
			if disk.Type == "s3" && disk.DiskS3 != nil {
				s3 := *disk.DiskS3
				disk.DiskS3 = &s3
				if err := decode(fmt.Sprintf("disk[%d].accessKeyID", i), &s3.AccessKeyID); err != nil {
					return err
				}
				if err := decode(fmt.Sprintf("disk[%d].secretAccessKey", i), &s3.SecretAccessKey); err != nil {
					return err
				}
			}
		}
		conf.Storage = &storage
	}
	if conf.TlsConf != nil {
		tls := *conf.TlsConf
		if err := decode("tlsConf.key", &tls.Key); err != nil {
			return err
		}
		if tls.Key != "" {
			if raw, err := base64.StdEncoding.DecodeString(tls.Key); err == nil && strings.Contains(string(raw), "PRIVATE KEY") {
				tls.Key = string(raw)
			} else {
				log.Logger.Warnf("cluster[%s] TLS private key decrypt sanity check failed; keeping raw decrypted value", conf.Cluster)
			}
		}
		conf.TlsConf = &tls
	}
	return nil
}
