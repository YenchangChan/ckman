package ckconfig

import (
	"github.com/housepower/ckman/common"
	"github.com/housepower/ckman/model"
	"github.com/housepower/ckman/service/zookeeper"
)

func GenerateMetrikaXML(filename string, conf *model.CKManClickHouseConfig) (string, error) {
	xml := common.NewXmlFile(filename)
	rootTag := GetRootTag(conf.Version)
	xml.Begin(rootTag)
	xml.Append(GenZookeeperMetrika(xml.GetIndent(), conf))
	xml.Begin("remote_servers")
	xml.Append(GenLocalMetrika(xml.GetIndent(), conf))
	xml.End("remote_servers")
	xml.End(rootTag)
	err := xml.Dump()
	if err != nil {
		return "", err
	}
	return filename, nil
}

func GenerateMetrikaXMLwithLogic(filename string, conf *model.CKManClickHouseConfig, logicMrtrika string) (string, error) {
	xml := common.NewXmlFile(filename)
	rootTag := GetRootTag(conf.Version)
	xml.Begin(rootTag)
	xml.Append(GenZookeeperMetrika(xml.GetIndent(), conf))
	xml.Begin("remote_servers")
	xml.Append(GenLocalMetrika(xml.GetIndent(), conf))
	xml.Append(logicMrtrika)
	xml.End("remote_servers")
	xml.End(rootTag)
	err := xml.Dump()
	if err != nil {
		return "", err
	}
	return filename, nil
}

func GenZookeeperMetrika(indent int, conf *model.CKManClickHouseConfig) string {
	xml := common.NewXmlFile("")
	xml.SetIndent(indent)
	xml.Begin("zookeeper")
	nodes, port := zookeeper.GetZkInfo(conf)
	for index, zk := range nodes {
		xml.BeginwithAttr("node", []common.XMLAttr{{Key: "index", Value: index + 1}})
		xml.Write("host", zk)
		xml.Write("port", port)
		xml.End("node")
	}
	xml.Write("operation_timeout_ms", 30000)
	xml.Write("session_timeout_ms", 60000)
	xml.End("zookeeper")
	return xml.GetContext()
}

func GenLocalMetrika(indent int, conf *model.CKManClickHouseConfig) string {
	xml := common.NewXmlFile("")
	xml.SetIndent(indent)
	xml.Begin(conf.Cluster)
	secret := true
	if common.CompareClickHouseVersion(conf.Version, "20.10.3.30") < 0 {
		secret = false
	}
	if secret {
		xml.Write("secret", "foo")
	}
	for _, shard := range conf.Shards {
		xml.Begin("shard")
		xml.Write("internal_replication", conf.IsReplica)
		for _, replica := range shard.Replicas {
			xml.Begin("replica")
			xml.Write("host", replica.Ip)
			xml.Write("port", conf.Port)
			if !secret {
				xml.Write("user", conf.User)
				xml.Write("password", conf.Password)
			}
			xml.End("replica")
		}
		xml.End("shard")
	}
	xml.End(conf.Cluster)
	return xml.GetContext()
}

func GenLogicMetrika(logicName string, clusters []model.CKManClickHouseConfig, secret bool) string {
	xml := common.NewXmlFile("")
	xml.SetIndent(2)
	xml.Begin(logicName)
	if secret {
		xml.Write("secret", "foo")
	}
	for _, conf := range clusters {
		for _, shard := range conf.Shards {
			xml.Begin("shard")
			xml.Write("internal_replication", conf.IsReplica)
			for _, replica := range shard.Replicas {
				xml.Begin("replica")
				xml.Write("host", replica.Ip)
				xml.Write("port", conf.Port)
				if !secret {
					xml.Write("user", conf.User)
					xml.Write("password", conf.Password)
				}
				xml.End("replica")
			}
			xml.End("shard")
		}
	}
	xml.End(logicName)
	return xml.GetContext()
}
