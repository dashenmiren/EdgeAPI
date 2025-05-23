package setup

import (
	"encoding/json"
	"github.com/dashenmiren/EdgeAPI/internal/acme"
	"github.com/dashenmiren/EdgeAPI/internal/db/models"
	"github.com/dashenmiren/EdgeAPI/internal/db/models/stats"
	"github.com/dashenmiren/EdgeAPI/internal/errors"
	"github.com/dashenmiren/EdgeCommon/pkg/serverconfigs"
	"github.com/dashenmiren/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/dashenmiren/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/dashenmiren/EdgeCommon/pkg/systemconfigs"
	"github.com/dashenmiren/EdgeCommon/pkg/userconfigs"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"strconv"
	"strings"
)

type upgradeVersion struct {
	version string
	f       func(db *dbs.DB) error
}

var upgradeFuncs = []*upgradeVersion{
	{
		"0.0.3", upgradeV0_0_3,
	},
	{
		"0.0.5", upgradeV0_0_5,
	},
	{
		"0.0.6", upgradeV0_0_6,
	},
	{
		"0.0.9", upgradeV0_0_9,
	},
	{
		"0.0.10", upgradeV0_0_10,
	},
	{
		"0.2.5", upgradeV0_2_5,
	},
	{
		"0.2.8.1", upgradeV0_2_8_1,
	},
	{
		"0.3.0", upgradeV0_3_0,
	},
	{
		"0.3.1", upgradeV0_3_1,
	},
	{
		"0.3.2", upgradeV0_3_2,
	},
	{
		"0.3.3", upgradeV0_3_3,
	},
	{
		"0.3.7", upgradeV0_3_7,
	},
	{
		"0.4.0", upgradeV0_4_0,
	},
	{
		"0.4.1", upgradeV0_4_1,
	},
	{
		"0.4.5", upgradeV0_4_5,
	},
	{
		"0.4.7", upgradeV0_4_7,
	},
	{
		"0.4.8", upgradeV0_4_8,
	},
	{
		"0.4.9", upgradeV0_4_9,
	},
	{
		"0.4.11", upgradeV0_4_11,
	},
	{
		"0.5.3", upgradeV0_5_3,
	},
	{
		"0.5.6", upgradeV0_5_6,
	},
	{
		"0.5.8", upgradeV0_5_8,
	},
	{
		"1.2.1", upgradeV1_2_1,
	},
	{
		"1.2.9", upgradeV1_2_9,
	},
	{
		"1.2.10", upgradeV1_2_10,
	},
	{
		"1.3.2", upgradeV1_3_2,
	},
	{
		"1.3.4", upgradeV1_3_4,
	},
}

// UpgradeSQLData 升级SQL数据
func UpgradeSQLData(db *dbs.DB) error {
	version, err := db.FindCol(0, "SELECT version FROM edgeVersions")
	if err != nil {
		return err
	}
	var versionString = types.String(version)
	if len(versionString) > 0 {
		for _, f := range upgradeFuncs {
			if CompareVersion(versionString, f.version) >= 0 {
				continue
			}
			err = f.f(db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// v0.0.3
func upgradeV0_0_3(db *dbs.DB) error {
	// 获取第一个管理员
	adminIdCol, err := db.FindCol(0, "SELECT id FROM edgeAdmins ORDER BY id ASC LIMIT 1")
	if err != nil {
		return err
	}
	adminId := types.Int64(adminIdCol)
	if adminId <= 0 {
		return errors.New("'edgeAdmins' table should not be empty")
	}

	// 升级edgeDNSProviders
	_, err = db.Exec("UPDATE edgeDNSProviders SET adminId=? WHERE adminId=0 AND userId=0", adminId)
	if err != nil {
		return err
	}

	// 升级edgeDNSDomains
	_, err = db.Exec("UPDATE edgeDNSDomains SET adminId=? WHERE adminId=0 AND userId=0", adminId)
	if err != nil {
		return err
	}

	// 升级edgeSSLCerts
	_, err = db.Exec("UPDATE edgeSSLCerts SET adminId=? WHERE adminId=0 AND userId=0", adminId)
	if err != nil {
		return err
	}

	// 升级edgeNodeClusters
	_, err = db.Exec("UPDATE edgeNodeClusters SET adminId=? WHERE adminId=0 AND userId=0", adminId)
	if err != nil {
		return err
	}

	// 升级edgeNodes
	_, err = db.Exec("UPDATE edgeNodes SET adminId=? WHERE adminId=0 AND userId=0", adminId)
	if err != nil {
		return err
	}

	// 升级edgeNodeGrants
	_, err = db.Exec("UPDATE edgeNodeGrants SET adminId=? WHERE adminId=0", adminId)
	if err != nil {
		return err
	}

	return nil
}

// v0.0.5
func upgradeV0_0_5(db *dbs.DB) error {
	// 升级edgeACMETasks
	_, err := db.Exec("UPDATE edgeACMETasks SET authType=? WHERE authType IS NULL OR LENGTH(authType)=0", acme.AuthTypeDNS)
	if err != nil {
		return err
	}

	return nil
}

// v0.0.6
func upgradeV0_0_6(db *dbs.DB) error {
	stmt, err := db.Prepare("SELECT COUNT(*) FROM edgeAPITokens WHERE role='user'")
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()
	col, err := stmt.FindCol(0)
	if err != nil {
		return err
	}
	count := types.Int(col)
	if count > 0 {
		return nil
	}

	nodeId := rands.HexString(32)
	secret := rands.String(32)
	_, err = db.Exec("INSERT INTO edgeAPITokens (nodeId, secret, role) VALUES (?, ?, ?)", nodeId, secret, "user")
	if err != nil {
		return err
	}

	return nil
}

// v0.0.9
func upgradeV0_0_9(db *dbs.DB) error {
	// firewall policies
	var tx *dbs.Tx
	dbs.NotifyReady()
	policies, err := models.NewHTTPFirewallPolicyDAO().FindAllEnabledFirewallPolicies(tx)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if policy.ServerId > 0 {
			continue
		}
		policyId := int64(policy.Id)
		webIds, err := models.NewHTTPWebDAO().FindAllWebIdsWithHTTPFirewallPolicyId(tx, policyId)
		if err != nil {
			return err
		}
		serverIds := []int64{}
		for _, webId := range webIds {
			serverId, err := models.NewServerDAO().FindEnabledServerIdWithWebId(tx, webId)
			if err != nil {
				return err
			}
			if serverId > 0 && !lists.ContainsInt64(serverIds, serverId) {
				serverIds = append(serverIds, serverId)
			}
		}
		if len(serverIds) == 1 {
			err = models.NewHTTPFirewallPolicyDAO().UpdateFirewallPolicyServerId(tx, policyId, serverIds[0])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// v0.0.10
func upgradeV0_0_10(db *dbs.DB) error {
	return nil
}

// v0.2.5
func upgradeV0_2_5(db *dbs.DB) error {
	// 更新用户
	_, err := db.Exec("UPDATE edgeUsers SET day=FROM_UNIXTIME(createdAt,'%Y%m%d') WHERE day IS NULL OR LENGTH(day)=0")
	if err != nil {
		return err
	}

	// 更新防火墙规则
	ones, _, err := db.FindOnes("SELECT id, actions, action, actionOptions FROM edgeHTTPFirewallRuleSets WHERE actions IS NULL OR LENGTH(actions)=0")
	if err != nil {
		return err
	}
	for _, one := range ones {
		oneId := one.GetInt64("id")
		action := one.GetString("action")
		options := one.GetString("actionOptions")
		var optionsMap = maps.Map{}
		if len(options) > 0 {
			_ = json.Unmarshal([]byte(options), &optionsMap)
		}
		var actions = []*firewallconfigs.HTTPFirewallActionConfig{
			{
				Code:    action,
				Options: optionsMap,
			},
		}
		actionsJSON, err := json.Marshal(actions)
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE edgeHTTPFirewallRuleSets SET actions=? WHERE id=?", string(actionsJSON), oneId)
		if err != nil {
			return err
		}
	}

	return nil
}

// v0.3.0
func upgradeV0_3_0(db *dbs.DB) error {
	// 升级健康检查
	ones, _, err := db.FindOnes("SELECT id,healthCheck FROM edgeNodeClusters WHERE state=1")
	if err != nil {
		return err
	}
	for _, one := range ones {
		var clusterId = one.GetInt64("id")
		var healthCheck = one.GetString("healthCheck")
		if len(healthCheck) == 0 {
			continue
		}
		var config = &serverconfigs.HealthCheckConfig{}
		err = json.Unmarshal([]byte(healthCheck), config)
		if err != nil {
			continue
		}
		if config.CountDown <= 1 {
			config.CountDown = 3
			configJSON, err := json.Marshal(config)
			if err != nil {
				continue
			}
			_, err = db.Exec("UPDATE edgeNodeClusters SET healthCheck=? WHERE id=?", string(configJSON), clusterId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// v0.3.1
func upgradeV0_3_1(db *dbs.DB) error {
	// 清空域名统计，已使用分表代替
	// 因为可能有权限问题，所以我们忽略错误
	_, _ = db.Exec("TRUNCATE table edgeServerDomainHourlyStats")

	// 升级APIToken
	ones, _, err := db.FindOnes("SELECT uniqueId,secret FROM edgeNodeClusters")
	if err != nil {
		return err
	}
	for _, one := range ones {
		var uniqueId = one.GetString("uniqueId")
		var secret = one.GetString("secret")
		tokenOne, err := db.FindOne("SELECT id FROM edgeAPITokens WHERE nodeId=? LIMIT 1", uniqueId)
		if err != nil {
			return err
		}
		if len(tokenOne) == 0 {
			_, err = db.Exec("INSERT INTO edgeAPITokens (nodeId, secret, role, state) VALUES (?, ?, 'cluster', 1)", uniqueId, secret)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// v0.3.2
func upgradeV0_3_2(db *dbs.DB) error {
	// gzip => compression

	type HTTPGzipRef struct {
		IsPrior bool  `yaml:"isPrior" json:"isPrior"` // 是否覆盖
		IsOn    bool  `yaml:"isOn" json:"isOn"`       // 是否开启
		GzipId  int64 `yaml:"gzipId" json:"gzipId"`   // 使用的配置ID
	}

	webOnes, _, err := db.FindOnes("SELECT id, gzip FROM edgeHTTPWebs WHERE gzip IS NOT NULL AND compression IS NULL")
	if err != nil {
		return err
	}
	for _, webOne := range webOnes {
		var gzipRef = &HTTPGzipRef{}
		err = json.Unmarshal([]byte(webOne.GetString("gzip")), gzipRef)
		if err != nil {
			continue
		}
		if gzipRef == nil || gzipRef.GzipId <= 0 {
			continue
		}
		var webId = webOne.GetInt("id")

		var compressionConfig = &serverconfigs.HTTPCompressionConfig{
			UseDefaultTypes: true,
		}
		compressionConfig.IsPrior = gzipRef.IsPrior
		compressionConfig.IsOn = gzipRef.IsOn

		gzipOne, err := db.FindOne("SELECT * FROM edgeHTTPGzips WHERE id=?", gzipRef.GzipId)
		if err != nil {
			return err
		}
		if len(gzipOne) == 0 {
			continue
		}

		level := gzipOne.GetInt("level")
		if level <= 0 {
			continue
		}
		if level > 0 && level <= 10 {
			compressionConfig.Level = types.Int8(level)
		} else if level > 10 {
			compressionConfig.Level = 10
		}

		var minLengthBytes = []byte(gzipOne.GetString("minLength"))
		if len(minLengthBytes) > 0 {
			var sizeCapacity = &shared.SizeCapacity{}
			err = json.Unmarshal(minLengthBytes, sizeCapacity)
			if err != nil {
				continue
			}
			compressionConfig.MinLength = sizeCapacity
		}

		var maxLengthBytes = []byte(gzipOne.GetString("maxLength"))
		if len(maxLengthBytes) > 0 {
			var sizeCapacity = &shared.SizeCapacity{}
			err = json.Unmarshal(maxLengthBytes, sizeCapacity)
			if err != nil {
				continue
			}
			compressionConfig.MaxLength = sizeCapacity
		}

		var condsBytes = []byte(gzipOne.GetString("conds"))
		if len(condsBytes) > 0 {
			var conds = &shared.HTTPRequestCondsConfig{}
			err = json.Unmarshal(condsBytes, conds)
			if err != nil {
				continue
			}
			compressionConfig.Conds = conds
		}

		configJSON, err := json.Marshal(compressionConfig)
		if err != nil {
			return err
		}
		_, err = db.Exec("UPDATE edgeHTTPWebs SET compression=? WHERE id=?", string(configJSON), webId)
		if err != nil {
			return err
		}
	}

	// 更新服务端口
	var serverDAO = models.NewServerDAO()
	ones, err := serverDAO.Query(nil).
		ResultPk().
		FindAll()
	if err != nil {
		return err
	}
	for _, one := range ones {
		var serverId = int64(one.(*models.Server).Id)
		err = serverDAO.NotifyServerPortsUpdate(nil, serverId)
		if err != nil {
			return err
		}
	}

	return nil
}

// v0.3.3
func upgradeV0_3_3(db *dbs.DB) error {
	// 升级CC请求数Code
	_, err := db.Exec("UPDATE edgeHTTPFirewallRuleSets SET code='8002' WHERE name='CC请求数' AND code='8001'")
	if err != nil {
		return err
	}

	// 清除节点
	// 删除7天以前的info日志
	err = models.NewNodeLogDAO().DeleteExpiredLogsWithLevel(nil, "info", 7)
	if err != nil {
		return err
	}

	return nil
}

// v0.3.7
func upgradeV0_3_7(db *dbs.DB) error {
	// 修改所有edgeNodeGrants中的su为0
	_, err := db.Exec("UPDATE edgeNodeGrants SET su=0 WHERE su=1")
	if err != nil {
		return err
	}

	// WAF预置分组
	_, err = db.Exec("UPDATE edgeHTTPFirewallRuleGroups SET isTemplate=1 WHERE LENGTH(code)>0")
	if err != nil {
		return err
	}

	return nil
}

// v0.4.0
func upgradeV0_4_0(db *dbs.DB) error {
	// 升级SYN Flood配置
	synFloodJSON, err := json.Marshal(firewallconfigs.NewSYNFloodConfig())
	if err == nil {
		_, err := db.Exec("UPDATE edgeHTTPFirewallPolicies SET synFlood=? WHERE synFlood IS NULL AND state=1", string(synFloodJSON))
		if err != nil {
			return err
		}
	}

	return nil
}

// v0.4.1
func upgradeV0_4_1(db *dbs.DB) error {
	// 升级 servers.lastUserPlanId
	_, err := db.Exec("UPDATE edgeServers SET lastUserPlanId=userPlanId WHERE userPlanId>0")
	if err != nil {
		return err
	}

	// 执行域名统计清理
	err = stats.NewServerDomainHourlyStatDAO().CleanDays(nil, 7)
	if err != nil {
		return err
	}

	return nil
}

// v0.4.5
func upgradeV0_4_5(db *dbs.DB) error {
	// 升级访问日志自动分表
	{
		var dao = models.NewSysSettingDAO()
		valueJSON, err := dao.ReadSetting(nil, systemconfigs.SettingCodeAccessLogQueue)
		if err != nil {
			return err
		}
		if len(valueJSON) > 0 {
			var config = &serverconfigs.AccessLogQueueConfig{}
			err = json.Unmarshal(valueJSON, config)
			if err == nil && config.RowsPerTable == 0 {
				config.EnableAutoPartial = true
				config.RowsPerTable = 500_000
				configJSON, err := json.Marshal(config)
				if err == nil {
					err = dao.UpdateSetting(nil, systemconfigs.SettingCodeAccessLogQueue, configJSON)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// 升级一个防SQL注入规则
	{
		ones, _, err := db.FindOnes(`SELECT id FROM edgeHTTPFirewallRules WHERE value=?`, "(updatexml|extractvalue|ascii|ord|char|chr|count|concat|rand|floor|substr|length|len|user|database|benchmark|analyse)\\s*\\(")
		if err != nil {
			return err
		}
		for _, one := range ones {
			var ruleId = one.GetInt64("id")
			_, err = db.Exec(`UPDATE edgeHTTPFirewallRules SET value=? WHERE id=? LIMIT 1`, `\b(updatexml|extractvalue|ascii|ord|char|chr|count|concat|rand|floor|substr|length|len|user|database|benchmark|analyse)\s*\(.*\)`, ruleId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// v0.4.7
func upgradeV0_4_7(db *dbs.DB) error {
	// 升级 edgeServers 中的 plainServerNames
	{
		ones, _, err := db.FindOnes("SELECT id,serverNames FROM edgeServers WHERE state=1")
		if err != nil {
			return err
		}
		for _, one := range ones {
			var serverId = one.GetInt64("id")
			var serverNamesJSON = one.GetBytes("serverNames")
			if len(serverNamesJSON) > 0 {
				var serverNames = []*serverconfigs.ServerNameConfig{}
				err = json.Unmarshal(serverNamesJSON, &serverNames)
				if err != nil {
					return err
				}
				plainServerNamesJSON, err := json.Marshal(serverconfigs.PlainServerNames(serverNames))
				if err != nil {
					return err
				}
				_, err = db.Exec("UPDATE edgeServers SET plainServerNames=? WHERE id=?", plainServerNamesJSON, serverId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// v0.4.8
func upgradeV0_4_8(db *dbs.DB) error {
	// 设置edgeIPLists中的serverId
	{
		firewallPolicyOnes, _, err := db.FindOnes("SELECT inbound,serverId FROM edgeHTTPFirewallPolicies WHERE serverId>0")
		if err != nil {
			return err
		}
		for _, one := range firewallPolicyOnes {
			var inboundBytes = one.GetBytes("inbound")
			var serverId = one.GetInt64("serverId")

			var listIds = []int64{}

			if len(inboundBytes) > 0 {
				var inbound = &firewallconfigs.HTTPFirewallInboundConfig{}
				err = json.Unmarshal(inboundBytes, inbound)
				if err == nil { // we ignore errors
					if inbound.AllowListRef != nil && inbound.AllowListRef.ListId > 0 {
						listIds = append(listIds, inbound.AllowListRef.ListId)
					}
					if inbound.DenyListRef != nil && inbound.DenyListRef.ListId > 0 {
						listIds = append(listIds, inbound.DenyListRef.ListId)
					}
					if inbound.GreyListRef != nil && inbound.GreyListRef.ListId > 0 {
						listIds = append(listIds, inbound.GreyListRef.ListId)
					}
				}
			}

			if len(listIds) == 0 {
				continue
			}
			for _, listId := range listIds {
				isPublicCol, err := db.FindCol(0, "SELECT isPublic FROM edgeIPLists WHERE id=? LIMIT 1", listId)
				if err != nil {
					return err
				}
				var isPublic = types.Bool(isPublicCol)
				if !isPublic {
					_, err = db.Exec("UPDATE edgeIPLists SET serverId=? WHERE id=?", serverId, listId)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// v0.4.11
func upgradeV0_4_11(db *dbs.DB) error {
	// 升级ns端口
	{
		// TCP
		{
			var config = &serverconfigs.TCPProtocolConfig{}
			config.IsOn = true
			config.Listen = []*serverconfigs.NetworkAddressConfig{
				{
					Protocol:  serverconfigs.ProtocolTCP,
					PortRange: "53",
				},
			}
			configJSON, err := json.Marshal(config)
			if err != nil {
				return err
			}
			_, err = db.Exec("UPDATE edgeNSClusters SET tcp=? WHERE tcp IS NULL", configJSON)
			if err != nil {
				return err
			}
		}

		// UDP
		{
			var config = &serverconfigs.UDPProtocolConfig{}
			config.IsOn = true
			config.Listen = []*serverconfigs.NetworkAddressConfig{
				{
					Protocol:  serverconfigs.ProtocolUDP,
					PortRange: "53",
				},
			}
			configJSON, err := json.Marshal(config)
			if err != nil {
				return err
			}
			_, err = db.Exec("UPDATE edgeNSClusters SET udp=? WHERE udp IS NULL", configJSON)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// v1.2.1
func upgradeV1_2_1(db *dbs.DB) error {
	// upgrade generated USER-xxx in old versions
	ones, _, err := db.FindOnes("SELECT id, username, clusterId FROM edgeUsers WHERE username LIKE 'USER-%'")
	if err != nil {
		return err
	}
	for _, one := range ones {
		var userId = one.GetInt64("id")
		var clusterId = one.GetInt64("clusterId")
		var username = one.GetString("username")
		if clusterId <= 0 {
			defaultClusterIdValue, err := db.FindCol(0, "SELECT id FROM edgeNodeClusters WHERE state=1 ORDER BY id ASC LIMIT 1")
			if err != nil {
				return err
			}

			var defaultClusterId = types.Int64(defaultClusterIdValue)
			if defaultClusterId > 0 {
				_, err = db.Exec("UPDATE edgeUsers SET username=?, clusterId=? WHERE id=?", strings.ReplaceAll(username, "-", "_"), defaultClusterId, userId)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// 1.2.10
func upgradeV1_2_10(db *dbs.DB) error {
	{
		type OldGlobalConfig struct {
			// HTTP & HTTPS相关配置
			HTTPAll struct {
				DomainAuditingIsOn   bool   `yaml:"domainAuditingIsOn" json:"domainAuditingIsOn"`     // 域名是否需要审核
				DomainAuditingPrompt string `yaml:"domainAuditingPrompt" json:"domainAuditingPrompt"` // 域名审核的提示
			} `yaml:"httpAll" json:"httpAll"`

			TCPAll struct {
				PortRangeMin int   `yaml:"portRangeMin" json:"portRangeMin"` // 最小端口
				PortRangeMax int   `yaml:"portRangeMax" json:"portRangeMax"` // 最大端口
				DenyPorts    []int `yaml:"denyPorts" json:"denyPorts"`       // 禁止使用的端口
			} `yaml:"tcpAll" json:"tcpAll"`
		}

		globalConfigValue, err := db.FindCol(0, "SELECT value FROM edgeSysSettings WHERE code='serverGlobalConfig'")
		if err != nil {
			return err
		}
		var globalConfigString = types.String(globalConfigValue)
		if len(globalConfigString) > 0 {
			var oldGlobalConfig = &OldGlobalConfig{}
			err = json.Unmarshal([]byte(globalConfigString), oldGlobalConfig)
			if err == nil { // we ignore error
				ones, _, err := db.FindOnes("SELECT id,globalServerConfig FROM edgeNodeClusters")
				if err != nil {
					return err
				}
				for _, one := range ones {
					var id = one.GetInt64("id")
					var globalServerConfigData = []byte(one.GetString("globalServerConfig"))
					if len(globalServerConfigData) > 32 {
						var globalServerConfig = &serverconfigs.GlobalServerConfig{}
						err = json.Unmarshal(globalServerConfigData, globalServerConfig)
						if err != nil {
							return err
						}

						globalServerConfig.HTTPAll.DomainAuditingIsOn = oldGlobalConfig.HTTPAll.DomainAuditingIsOn
						globalServerConfig.HTTPAll.DomainAuditingPrompt = oldGlobalConfig.HTTPAll.DomainAuditingPrompt

						globalServerConfig.TCPAll.DenyPorts = oldGlobalConfig.TCPAll.DenyPorts
						globalServerConfig.TCPAll.PortRangeMin = oldGlobalConfig.TCPAll.PortRangeMin
						globalServerConfig.TCPAll.PortRangeMax = oldGlobalConfig.TCPAll.PortRangeMax

						globalServerConfigJSON, err := json.Marshal(globalServerConfig)
						if err != nil {
							return err
						}
						_, err = db.Exec("UPDATE edgeNodeClusters SET globalServerConfig=? WHERE id=?", globalServerConfigJSON, id)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// 1.3.2
func upgradeV1_3_2(db *dbs.DB) error {
	// waf
	{
		var disableSet = func(setId int64) error {
			_, err := db.Exec("UPDATE edgeHTTPFirewallRuleSets SET state=0 WHERE id=?", setId)
			return err
		}

		var addRuleToGroup = func(groupId int64, setCode string, setName string, actions []*firewallconfigs.HTTPFirewallActionConfig, ruleParam string, ruleOperator string, value string) error {
			actionsJSON, err := json.Marshal(actions)
			if err != nil {
				return err
			}

			// rule
			ruleResult, err := db.Exec("INSERT INTO edgeHTTPFirewallRules (isOn, param, operator, value, isCaseInsensitive, state) VALUES (1, ?, ?, ?, 0, 1)", ruleParam, ruleOperator, value)
			if err != nil {
				return err
			}
			ruleId, err := ruleResult.LastInsertId()
			if err != nil {
				return err
			}
			var ruleRefs = []*firewallconfigs.HTTPFirewallRuleRef{
				{
					IsOn:   true,
					RuleId: ruleId,
				},
			}
			ruleRefsJSON, err := json.Marshal(ruleRefs)
			if err != nil {
				return err
			}

			// set
			setResult, err := db.Exec("INSERT INTO edgeHTTPFirewallRuleSets (isOn, code, name, rules, connector, state, actions) VALUES (1, ?, ?, ?, 'or', 1, ?)", setCode, setName, ruleRefsJSON, actionsJSON)
			if err != nil {
				return err
			}
			setId, err := setResult.LastInsertId()
			if err != nil {
				return err
			}
			var setRefs = []*firewallconfigs.HTTPFirewallRuleSetRef{
				{
					IsOn:  true,
					SetId: setId,
				},
			}
			setRefsJSON, err := json.Marshal(setRefs)
			if err != nil {
				return err
			}

			// group
			_, err = db.Exec("UPDATE edgeHTTPFirewallRuleGroups SET sets=? WHERE id=?", setRefsJSON, groupId)
			if err != nil {
				return err
			}

			return nil
		}

		// sql injection
		{
			ruleGroups, _, err := db.FindOnes("SELECT id, sets FROM edgeHTTPFirewallRuleGroups WHERE code='sqlInjection' AND state=1")
			if err != nil {
				return err
			}
			for _, ruleGroup := range ruleGroups {
				var setsJSON = ruleGroup.GetBytes("sets")
				if len(setsJSON) == 0 {
					continue
				}

				var setRefs = []*firewallconfigs.HTTPFirewallRuleSetRef{}
				err = json.Unmarshal(setsJSON, &setRefs)
				if err != nil {
					continue
				}
				if len(setRefs) != 5 {
					continue
				}

				var isChanged = false
				for setIndex, setRef := range setRefs {
					set, setErr := db.FindOne("SELECT id, rules, isOn, actions FROM edgeHTTPFirewallRuleSets WHERE id=? AND state=1", setRef.SetId)
					if setErr != nil {
						return setErr
					}
					if set == nil {
						isChanged = true
						break
					}
					var rulesJSON = set.GetBytes("rules")
					if len(rulesJSON) == 0 {
						isChanged = true
						break
					}
					var ruleRefs = []*firewallconfigs.HTTPFirewallRuleRef{}
					err = json.Unmarshal(rulesJSON, &ruleRefs)
					if err != nil {
						return err
					}
					if len(ruleRefs) < 1 {
						isChanged = true
						break
					}

					var actionsJSON = set.GetBytes("actions")
					if len(actionsJSON) == 0 {
						isChanged = true
						break
					}
					var actions = []*firewallconfigs.HTTPFirewallActionConfig{}
					err = json.Unmarshal(actionsJSON, &actions)
					if err != nil {
						return err
					}
					if !(len(actions) == 1 && actions[0].Code == firewallconfigs.HTTPFirewallActionBlock) {
						isChanged = true
						break
					}

					var rules = []maps.Map{}
					for _, ruleRef := range ruleRefs {
						rule, ruleErr := db.FindOne("SELECT * FROM edgeHTTPFirewallRules WHERE id=? AND state=1", ruleRef.RuleId)
						if ruleErr != nil {
							return ruleErr
						}
						if rule == nil {
							isChanged = true
							break
						}
						rules = append(rules, rule)
					}
					if isChanged {
						break
					}
					if len(rules) < 1 {
						isChanged = true
						break
					}

					switch setIndex {
					case 0:
						var rule = rules[0]
						if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `union[\s/\*]+select`) {
							isChanged = true
						}
					case 1:
						var rule = rules[0]
						if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `/\*(!|\x00)`) {
							isChanged = true
						}
					case 2:
						if len(rules) != 4 {
							isChanged = true
						} else {
							{
								var rule = rules[0]
								if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `\s(and|or|rlike)\s+(if|updatexml)\s*\(`) {
									isChanged = true
								}
							}

							{
								var rule = rules[1]
								if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `\s+(and|or|rlike)\s+(select|case)\s+`) {
									isChanged = true
								}
							}

							{
								var rule = rules[2]
								if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `\s+(and|or|procedure)\s+[\w\p{L}]+\s*=\s*[\w\p{L}]+(\s|$|--|#)`) {
									isChanged = true
								}
							}

							{
								var rule = rules[3]
								if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `\(\s*case\s+when\s+[\w\p{L}]+\s*=\s*[\w\p{L}]+\s+then\s+`) {
									isChanged = true
								}
							}
						}
					case 3:
						var rule = rules[0]
						if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `\b(updatexml|extractvalue|ascii|ord|char|chr|count|concat|rand|floor|substr|length|len|user|database|benchmark|analyse)\s*\(.*\)`) {
							isChanged = true
						}
					case 4:
						var rule = rules[0]
						if !(rule.GetString("param") == "${requestAll}" && rule.GetString("operator") == "match" && rule.GetString("value") == `;\s*(declare|use|drop|create|exec|delete|update|insert)\s`) {
							isChanged = true
						}
					}
				}
				if isChanged {
					continue
				}

				for _, setRef := range setRefs {
					err = disableSet(setRef.SetId)
					if err != nil {
						return err
					}
				}

				err = addRuleToGroup(ruleGroup.GetInt64("id"), "7010", "SQL注入检测", []*firewallconfigs.HTTPFirewallActionConfig{
					{
						Code:    firewallconfigs.HTTPFirewallActionPage,
						Options: maps.Map{"status": 403, "body": ""},
					},
				}, "${requestAll}", firewallconfigs.HTTPFirewallRuleOperatorContainsSQLInjection, "")
				if err != nil {
					return err
				}
			}
		}

		// xss
		{
			ruleGroups, _, err := db.FindOnes("SELECT id, sets FROM edgeHTTPFirewallRuleGroups WHERE code='xss' AND state=1")
			if err != nil {
				return err
			}
			for _, ruleGroup := range ruleGroups {
				var setsJSON = ruleGroup.GetBytes("sets")
				if len(setsJSON) == 0 {
					continue
				}

				var setRefs = []*firewallconfigs.HTTPFirewallRuleSetRef{}
				err = json.Unmarshal(setsJSON, &setRefs)
				if err != nil {
					continue
				}
				if len(setRefs) != 3 {
					continue
				}

				var isChanged = false
				for setIndex, setRef := range setRefs {
					set, setErr := db.FindOne("SELECT id, rules, isOn, actions FROM edgeHTTPFirewallRuleSets WHERE id=? AND state=1", setRef.SetId)
					if setErr != nil {
						return setErr
					}
					if set == nil {
						isChanged = true
						break
					}
					var rulesJSON = set.GetBytes("rules")
					if len(rulesJSON) == 0 {
						isChanged = true
						break
					}
					var ruleRefs = []*firewallconfigs.HTTPFirewallRuleRef{}
					err = json.Unmarshal(rulesJSON, &ruleRefs)
					if err != nil {
						return err
					}
					if len(ruleRefs) != 1 {
						isChanged = true
						break
					}

					var actionsJSON = set.GetBytes("actions")
					if len(actionsJSON) == 0 {
						isChanged = true
						break
					}
					var actions = []*firewallconfigs.HTTPFirewallActionConfig{}
					err = json.Unmarshal(actionsJSON, &actions)
					if err != nil {
						return err
					}
					if !(len(actions) == 1 && actions[0].Code == firewallconfigs.HTTPFirewallActionBlock) {
						isChanged = true
						break
					}

					rule, ruleErr := db.FindOne("SELECT * FROM edgeHTTPFirewallRules WHERE id=? AND state=1", ruleRefs[0].RuleId)
					if ruleErr != nil {
						return ruleErr
					}
					if rule == nil {
						isChanged = true
						break
					}

					switch setIndex {
					case 0:
						if !(rule.GetString("param") == "${requestURI}" && rule.GetString("operator") == "match" && rule.GetString("value") == `(onmouseover|onmousemove|onmousedown|onmouseup|onerror|onload|onclick|ondblclick|onkeydown|onkeyup|onkeypress)\s*=`) {
							isChanged = true
						}
					case 1:
						if !(rule.GetString("param") == "${requestURI}" && rule.GetString("operator") == "match" && rule.GetString("value") == `(alert|eval|prompt|confirm)\s*\(`) {
							isChanged = true
						}
					case 2:
						if !(rule.GetString("param") == "${requestURI}" && rule.GetString("operator") == "match" && rule.GetString("value") == `<(script|iframe|link)`) {
							isChanged = true
						}
					}
				}
				if isChanged {
					continue
				}

				for _, setRef := range setRefs {
					err = disableSet(setRef.SetId)
					if err != nil {
						return err
					}
				}

				err = addRuleToGroup(ruleGroup.GetInt64("id"), "1010", "XSS攻击检测", []*firewallconfigs.HTTPFirewallActionConfig{
					{
						Code:    firewallconfigs.HTTPFirewallActionPage,
						Options: maps.Map{"status": 403, "body": ""},
					},
				}, "${requestAll}", firewallconfigs.HTTPFirewallRuleOperatorContainsXSS, "")
				if err != nil {
					return err
				}
			}
		}
	}

	// user register config

	var newAddedFeatureCodes = []string{
		userconfigs.UserFeatureCodeServerOptimization,
		userconfigs.UserFeatureCodeServerAuth,
		userconfigs.UserFeatureCodeServerWebsocket,
		userconfigs.UserFeatureCodeServerHTTP3,
		userconfigs.UserFeatureCodeServerCC,
		userconfigs.UserFeatureCodeServerReferers,
		userconfigs.UserFeatureCodeServerUserAgent,
		userconfigs.UserFeatureCodeServerRequestLimit,
		userconfigs.UserFeatureCodeServerCompression,
		userconfigs.UserFeatureCodeServerRewriteRules,
		userconfigs.UserFeatureCodeServerHostRedirects,
		userconfigs.UserFeatureCodeServerHTTPHeaders,
		userconfigs.UserFeatureCodeServerPages,
	}

	{
		value, err := db.FindCol(0, "SELECT value FROM edgeSysSettings WHERE code=?", systemconfigs.SettingCodeUserRegisterConfig)
		if err != nil {
			return err
		}
		if value != nil {
			var valueString = types.String(value)
			if valueString != "null" && len(valueString) > 0 {
				var registerConfig = &userconfigs.UserRegisterConfig{}
				err = json.Unmarshal([]byte(valueString), registerConfig)
				if err != nil {
					return err
				}

				if len(registerConfig.Features) > 0 {
					var newFeatureCodes = registerConfig.Features
					var changed = false
					for _, featureCode := range newAddedFeatureCodes {
						if !lists.ContainsString(newFeatureCodes, featureCode) {
							newFeatureCodes = append(newFeatureCodes, featureCode)
							changed = true
						}
					}

					if changed {
						registerConfig.Features = newFeatureCodes
						registerConfigJSON, err := json.Marshal(registerConfig)
						if err != nil {
							return err
						}
						_, err = db.Exec("UPDATE edgeSysSettings SET value=? WHERE code=?", registerConfigJSON, systemconfigs.SettingCodeUserRegisterConfig)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	// user features
	{
		var sqlPieces []string
		for _, featureCode := range newAddedFeatureCodes {
			if strings.Contains(featureCode, "'") {
				continue
			}
			sqlPieces = append(sqlPieces, "'$', '"+featureCode+"'")
		}

		_, err := db.Exec("UPDATE edgeUsers SET features=JSON_ARRAY_APPEND(features," + strings.Join(sqlPieces, ",") + ") WHERE features IS NOT NULL AND JSON_LENGTH(features)>0 AND NOT JSON_CONTAINS(features, '" + strconv.Quote(newAddedFeatureCodes[0]) + "')")
		if err != nil {
			return err
		}
	}

	return nil
}

// 1.3.4
func upgradeV1_3_4(db *dbs.DB) error {
	_, err := db.Exec("DELETE FROM edgeLoginSessions WHERE adminId>0")
	if err != nil {
		return err
	}

	return nil
}
