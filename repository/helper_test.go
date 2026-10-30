package repository

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/housepower/ckman/common"
	"github.com/housepower/ckman/log"
	"github.com/housepower/ckman/model"
	"github.com/stretchr/testify/assert"
)

// TestMain 初始化 log.Logger:DecodePasswd 在私钥解密健全性检查失败时会打一条
// warning,若 log.Logger 未初始化（nil *zap.SugaredLogger）会直接 panic。
func TestMain(m *testing.M) {
	log.InitLoggerConsole()
	os.Exit(m.Run())
}

func TestEncodeDecodeTlsKey(t *testing.T) {
	plain := "-----BEGIN PRIVATE KEY-----\nMIIabc\n-----END PRIVATE KEY-----\n"
	conf := model.CKManClickHouseConfig{
		Password: "123456",
		Secure:   true,
		TlsConf:  &model.TlsConf{Cert: "CERT-PEM", Key: plain, SelfSigned: true},
	}

	assert.NoError(t, EncodePasswd(&conf))
	assert.NotEqual(t, plain, conf.TlsConf.Key)
	// 证书是公开信息，不加密
	assert.Equal(t, "CERT-PEM", conf.TlsConf.Cert)

	assert.NoError(t, DecodePasswd(&conf))
	assert.Equal(t, plain, conf.TlsConf.Key)
	assert.True(t, conf.TlsConf.SelfSigned)
}

func TestDecodeTlsKeyNotMutateShared(t *testing.T) {
	plain := "REAL-KEY"
	shared := &model.TlsConf{Cert: "CERT-PEM", Key: plain}
	conf := model.CKManClickHouseConfig{TlsConf: shared}

	assert.NoError(t, EncodePasswd(&conf))
	// EncodePasswd 不得改到调用方持有的那个对象
	assert.Equal(t, plain, shared.Key)
}

func TestEncodeDecodeTlsConfNil(t *testing.T) {
	conf := model.CKManClickHouseConfig{Password: "123456"}
	assert.NoError(t, EncodePasswd(&conf))
	assert.Nil(t, conf.TlsConf)
	assert.NoError(t, DecodePasswd(&conf))
	assert.Nil(t, conf.TlsConf)
}

// 损坏的密文必须返回错误，避免原始密文被误用为私钥。
func TestDecodeTlsKeyCorruptCiphertext_NoPanic(t *testing.T) {
	garbage := "not-a-valid-ciphertext!!"
	conf := model.CKManClickHouseConfig{
		TlsConf: &model.TlsConf{Cert: "CERT-PEM", Key: garbage},
	}

	assert.Error(t, DecodePasswd(&conf))
	assert.Equal(t, garbage, conf.TlsConf.Key)
}

// TestDecodeTlsKeyInvalidContent_KeepsRawForDiagnostics：密文本身合法（能被
// AesDecryptECB 解密、能被 base64 还原），但还原出的内容不是 PEM 私钥
// （不含 "PRIVATE KEY"）——例如误把别的内容加密后存进了 TlsConf.Key。
//
// 这是唯一能真正验证「健全性检查」是否生效的场景：在加检查之前的旧实现下，
// base64 解码本身不会报错，会把 Key 静默替换成 notAKey，这条断言会 FAIL；
// 加了 "PRIVATE KEY" 内容检查之后，Key 应保留 AesDecryptECB 的直接输出。
func TestDecodeTlsKeyInvalidContent_KeepsRawForDiagnostics(t *testing.T) {
	notAKey := "not a real private key"
	cipher := common.AesEncryptECB(base64.StdEncoding.EncodeToString([]byte(notAKey)))
	conf := model.CKManClickHouseConfig{
		TlsConf: &model.TlsConf{Cert: "CERT-PEM", Key: cipher},
	}

	assert.NoError(t, DecodePasswd(&conf))
	assert.Equal(t, common.AesDecryptECB(cipher), conf.TlsConf.Key)
	assert.NotEqual(t, notAKey, conf.TlsConf.Key)
}
