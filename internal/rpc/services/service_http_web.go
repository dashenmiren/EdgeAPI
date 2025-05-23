package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dashenmiren/EdgeAPI/internal/db/models"
	"github.com/dashenmiren/EdgeAPI/internal/errors"
	"github.com/dashenmiren/EdgeAPI/internal/utils/domainutils"
	"github.com/dashenmiren/EdgeAPI/internal/utils/regexputils"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
	"github.com/dashenmiren/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/dbs"
	"github.com/iwind/TeaGo/lists"
	"strings"
)

type HTTPWebService struct {
	BaseService
}

// CreateHTTPWeb 创建Web配置
func (this *HTTPWebService) CreateHTTPWeb(ctx context.Context, req *pb.CreateHTTPWebRequest) (*pb.CreateHTTPWebResponse, error) {
	// 校验请求
	adminId, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	webId, err := models.SharedHTTPWebDAO.CreateWeb(tx, adminId, userId, req.RootJSON)
	if err != nil {
		return nil, err
	}

	return &pb.CreateHTTPWebResponse{HttpWebId: webId}, nil
}

// FindEnabledHTTPWeb 查找Web配置
func (this *HTTPWebService) FindEnabledHTTPWeb(ctx context.Context, req *pb.FindEnabledHTTPWebRequest) (*pb.FindEnabledHTTPWebResponse, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	web, err := models.SharedHTTPWebDAO.FindEnabledHTTPWeb(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}

	if web == nil {
		return &pb.FindEnabledHTTPWebResponse{HttpWeb: nil}, nil
	}

	result := &pb.HTTPWeb{}
	result.Id = int64(web.Id)
	result.IsOn = web.IsOn
	return &pb.FindEnabledHTTPWebResponse{HttpWeb: result}, nil
}

// FindEnabledHTTPWebConfig 查找Web配置
func (this *HTTPWebService) FindEnabledHTTPWebConfig(ctx context.Context, req *pb.FindEnabledHTTPWebConfigRequest) (*pb.FindEnabledHTTPWebConfigResponse, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	config, err := models.SharedHTTPWebDAO.ComposeWebConfig(tx, req.HttpWebId, false, false, nil, nil)
	if err != nil {
		return nil, err
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return &pb.FindEnabledHTTPWebConfigResponse{HttpWebJSON: configJSON}, nil
}

// UpdateHTTPWeb 修改Web配置
func (this *HTTPWebService) UpdateHTTPWeb(ctx context.Context, req *pb.UpdateHTTPWebRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}

		req.RootJSON = []byte("{}") // 为了安全
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWeb(tx, req.HttpWebId, req.RootJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebCompression 修改压缩配置
func (this *HTTPWebService) UpdateHTTPWebCompression(ctx context.Context, req *pb.UpdateHTTPWebCompressionRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	if len(req.CompressionJSON) == 0 {
		return nil, errors.New("'compressionJSON' should not be empty")
	}
	var compressionConfig = &serverconfigs.HTTPCompressionConfig{}
	err = json.Unmarshal(req.CompressionJSON, compressionConfig)
	if err != nil {
		return nil, err
	}
	err = compressionConfig.Init()
	if err != nil {
		return nil, err
	}

	err = models.SharedHTTPWebDAO.UpdateWebCompression(tx, req.HttpWebId, compressionConfig)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebOptimization 修改页面优化配置
func (this *HTTPWebService) UpdateHTTPWebOptimization(ctx context.Context, req *pb.UpdateHTTPWebOptimizationRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	if len(req.OptimizationJSON) == 0 {
		return nil, errors.New("invalid 'optimizationJSON'")
	}
	var optimizationConfig = serverconfigs.NewHTTPPageOptimizationConfig()
	err = json.Unmarshal(req.OptimizationJSON, optimizationConfig)
	if err != nil {
		return nil, err
	}

	err = optimizationConfig.Init()
	if err != nil {
		return nil, err
	}

	err = models.SharedHTTPWebDAO.UpdateWebOptimization(tx, req.HttpWebId, optimizationConfig)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebWebP 修改WebP配置
func (this *HTTPWebService) UpdateHTTPWebWebP(ctx context.Context, req *pb.UpdateHTTPWebWebPRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebWebP(tx, req.HttpWebId, req.WebpJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebRemoteAddr 更改RemoteAddr配置
func (this *HTTPWebService) UpdateHTTPWebRemoteAddr(ctx context.Context, req *pb.UpdateHTTPWebRemoteAddrRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()
	err = models.SharedHTTPWebDAO.UpdateWebRemoteAddr(tx, req.HttpWebId, req.RemoteAddrJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebCharset 修改字符集配置
func (this *HTTPWebService) UpdateHTTPWebCharset(ctx context.Context, req *pb.UpdateHTTPWebCharsetRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebCharset(tx, req.HttpWebId, req.CharsetJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebRequestHeader 更改请求Header策略
func (this *HTTPWebService) UpdateHTTPWebRequestHeader(ctx context.Context, req *pb.UpdateHTTPWebRequestHeaderRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebRequestHeaderPolicy(tx, req.HttpWebId, req.HeaderJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebResponseHeader 更改响应Header策略
func (this *HTTPWebService) UpdateHTTPWebResponseHeader(ctx context.Context, req *pb.UpdateHTTPWebResponseHeaderRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebResponseHeaderPolicy(tx, req.HttpWebId, req.HeaderJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebShutdown 更改Shutdown
func (this *HTTPWebService) UpdateHTTPWebShutdown(ctx context.Context, req *pb.UpdateHTTPWebShutdownRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()
	var newShutdownJSON = req.ShutdownJSON
	if len(req.ShutdownJSON) > 0 {
		const maxURLLength = 512
		const maxBodyLength = 32 * 1024

		var shutdownConfig = &serverconfigs.HTTPShutdownConfig{}
		err = json.Unmarshal(req.ShutdownJSON, shutdownConfig)
		if err != nil {
			return nil, err
		}
		err = shutdownConfig.Init()
		if err != nil {
			return nil, errors.New("validate config failed: " + err.Error())
		}

		switch shutdownConfig.BodyType {
		case serverconfigs.HTTPPageBodyTypeURL:
			if len(shutdownConfig.URL) > maxURLLength {
				return nil, errors.New("'url' too long")
			}
			if shutdownConfig.IsOn /** validate when it's on **/ && !regexputils.HTTPProtocol.MatchString(shutdownConfig.URL) {
				return nil, errors.New("invalid 'url' format")
			}

			if len(shutdownConfig.Body) > maxBodyLength { // we keep short body for user experience
				shutdownConfig.Body = ""
			}
		case serverconfigs.HTTPPageBodyTypeRedirectURL:
			if len(shutdownConfig.URL) > maxURLLength {
				return nil, errors.New("'url' too long")
			}
			if shutdownConfig.IsOn /** validate when it's on **/ && !regexputils.HTTPProtocol.MatchString(shutdownConfig.URL) {
				return nil, errors.New("invalid 'url' format")
			}

			if len(shutdownConfig.Body) > maxBodyLength { // we keep short body for user experience
				shutdownConfig.Body = ""
			}
		case serverconfigs.HTTPPageBodyTypeHTML:
			if len(shutdownConfig.Body) > maxBodyLength {
				return nil, errors.New("'body' too long")
			}

			if len(shutdownConfig.URL) > maxURLLength { // we keep short url for user experience
				shutdownConfig.URL = ""
			}

		default:
			return nil, errors.New("invalid 'bodyType': " + shutdownConfig.BodyType)
		}

		newShutdownJSON, err = json.Marshal(shutdownConfig)
		if err != nil {
			return nil, err
		}
	}

	err = models.SharedHTTPWebDAO.UpdateWebShutdown(tx, req.HttpWebId, newShutdownJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebPages 更改Pages
func (this *HTTPWebService) UpdateHTTPWebPages(ctx context.Context, req *pb.UpdateHTTPWebPagesRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	// 检查配置
	var newPages = []*serverconfigs.HTTPPageConfig{}
	if len(req.PagesJSON) > 0 {
		var pages = []*serverconfigs.HTTPPageConfig{}
		err = json.Unmarshal(req.PagesJSON, &pages)
		if err != nil {
			return nil, err
		}

		for _, page := range pages {
			err = page.Init()
			if err != nil {
				return nil, errors.New("validate page failed: " + err.Error())
			}

			// reset not needed fields, keep "id" reference only
			page.URL = ""
			page.Body = ""

			newPages = append(newPages, &serverconfigs.HTTPPageConfig{Id: page.Id})
		}
	}
	newPagesJSON, err := json.Marshal(newPages)
	if err != nil {
		return nil, err
	}

	err = models.SharedHTTPWebDAO.UpdateWebPages(tx, req.HttpWebId, newPagesJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebGlobalPagesEnabled 更改系统自定义页面启用状态
func (this *HTTPWebService) UpdateHTTPWebGlobalPagesEnabled(ctx context.Context, req *pb.UpdateHTTPWebGlobalPagesEnabledRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()
	err = models.SharedHTTPWebDAO.UpdateGlobalPagesEnabled(tx, req.HttpWebId, req.IsEnabled)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebAccessLog 更改访问日志配置
func (this *HTTPWebService) UpdateHTTPWebAccessLog(ctx context.Context, req *pb.UpdateHTTPWebAccessLogRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebAccessLogConfig(tx, req.HttpWebId, req.AccessLogJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebStat 更改统计配置
func (this *HTTPWebService) UpdateHTTPWebStat(ctx context.Context, req *pb.UpdateHTTPWebStatRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebStat(tx, req.HttpWebId, req.StatJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebCache 更改缓存配置
func (this *HTTPWebService) UpdateHTTPWebCache(ctx context.Context, req *pb.UpdateHTTPWebCacheRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	var cacheConfig = &serverconfigs.HTTPCacheConfig{}
	if len(req.CacheJSON) > 0 {
		err = json.Unmarshal(req.CacheJSON, cacheConfig)
		if err != nil {
			return nil, err
		}
		err = cacheConfig.Init()
		if err != nil {
			return nil, fmt.Errorf("validate cache config failed: %w", err)
		}
	}

	// validate host
	if cacheConfig.Key != nil && cacheConfig.Key.IsOn {
		if cacheConfig.Key.Scheme != "http" && cacheConfig.Key.Scheme != "https" {
			return nil, errors.New("key scheme must be 'http' or 'https'")
		}

		cacheConfig.Key.Host = strings.ToLower(cacheConfig.Key.Host)
		if !domainutils.ValidateDomainFormat(cacheConfig.Key.Host) {
			return nil, errors.New("key host must be a valid domain")
		}

		serverId, err := models.SharedHTTPWebDAO.FindWebServerId(tx, req.HttpWebId)
		if err != nil {
			return nil, err
		}
		if serverId > 0 {
			serverNamesJSON, _, _, _, _, err := models.SharedServerDAO.FindServerServerNames(tx, serverId)
			if err != nil {
				return nil, err
			}
			var serverNames = []*serverconfigs.ServerNameConfig{}
			err = json.Unmarshal(serverNamesJSON, &serverNames)
			if err != nil {
				return nil, err
			}
			if !lists.ContainsString(serverconfigs.PlainServerNames(serverNames), cacheConfig.Key.Host) {
				return nil, errors.New("key host '" + cacheConfig.Key.Host + "' not exists in server names")
			}
		}
	}

	err = models.SharedHTTPWebDAO.UpdateWebCache(tx, req.HttpWebId, req.CacheJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebFirewall 更改防火墙设置
func (this *HTTPWebService) UpdateHTTPWebFirewall(ctx context.Context, req *pb.UpdateHTTPWebFirewallRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebFirewall(tx, req.HttpWebId, req.FirewallJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebLocations 更改路由规则设置
func (this *HTTPWebService) UpdateHTTPWebLocations(ctx context.Context, req *pb.UpdateHTTPWebLocationsRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		// 检查用户权限
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebLocations(tx, req.HttpWebId, req.LocationsJSON)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebRedirectToHTTPS 更改跳转到HTTPS设置
func (this *HTTPWebService) UpdateHTTPWebRedirectToHTTPS(ctx context.Context, req *pb.UpdateHTTPWebRedirectToHTTPSRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	// 检查权限
	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebRedirectToHTTPS(tx, req.HttpWebId, req.RedirectToHTTPSJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebWebsocket 更改Websocket设置
func (this *HTTPWebService) UpdateHTTPWebWebsocket(ctx context.Context, req *pb.UpdateHTTPWebWebsocketRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebsocket(tx, req.HttpWebId, req.WebsocketJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebFastcgi 更改Fastcgi设置
func (this *HTTPWebService) UpdateHTTPWebFastcgi(ctx context.Context, req *pb.UpdateHTTPWebFastcgiRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebFastcgi(tx, req.HttpWebId, req.FastcgiJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebRewriteRules 更改重写规则设置
func (this *HTTPWebService) UpdateHTTPWebRewriteRules(ctx context.Context, req *pb.UpdateHTTPWebRewriteRulesRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx = this.NullTx()

	err = models.SharedHTTPWebDAO.UpdateWebRewriteRules(tx, req.HttpWebId, req.RewriteRulesJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebHostRedirects 更改主机跳转设置
func (this *HTTPWebService) UpdateHTTPWebHostRedirects(ctx context.Context, req *pb.UpdateHTTPWebHostRedirectsRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	hostRedirects := []*serverconfigs.HTTPHostRedirectConfig{}
	if len(req.HostRedirectsJSON) == 0 {
		return nil, errors.New("'hostRedirectsJSON' should not be empty")
	}
	err = json.Unmarshal(req.HostRedirectsJSON, &hostRedirects)
	if err != nil {
		return nil, err
	}

	// 校验
	for _, redirect := range hostRedirects {
		err := redirect.Init()
		if err != nil {
			return nil, err
		}
	}

	var tx *dbs.Tx
	err = models.SharedHTTPWebDAO.UpdateWebHostRedirects(tx, req.HttpWebId, hostRedirects)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// FindHTTPWebHostRedirects 查找主机跳转设置
func (this *HTTPWebService) FindHTTPWebHostRedirects(ctx context.Context, req *pb.FindHTTPWebHostRedirectsRequest) (*pb.FindHTTPWebHostRedirectsResponse, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx *dbs.Tx
	redirectsJSON, err := models.SharedHTTPWebDAO.FindWebHostRedirects(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}
	return &pb.FindHTTPWebHostRedirectsResponse{HostRedirectsJSON: redirectsJSON}, nil
}

// UpdateHTTPWebAuth 更改认证设置
func (this *HTTPWebService) UpdateHTTPWebAuth(ctx context.Context, req *pb.UpdateHTTPWebAuthRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(nil, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var tx *dbs.Tx
	err = models.SharedHTTPWebDAO.UpdateWebAuth(tx, req.HttpWebId, req.AuthJSON)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// UpdateHTTPWebCommon 更改通用设置
func (this *HTTPWebService) UpdateHTTPWebCommon(ctx context.Context, req *pb.UpdateHTTPWebCommonRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(tx, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	err = models.SharedHTTPWebDAO.UpdateWebCommon(tx, req.HttpWebId, req.MergeSlashes)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// UpdateHTTPWebRequestLimit 修改请求限制
func (this *HTTPWebService) UpdateHTTPWebRequestLimit(ctx context.Context, req *pb.UpdateHTTPWebRequestLimitRequest) (*pb.RPCSuccess, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(tx, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var config = &serverconfigs.HTTPRequestLimitConfig{}
	err = json.Unmarshal(req.RequestLimitJSON, config)
	if err != nil {
		return nil, err
	}

	err = models.SharedHTTPWebDAO.UpdateWebRequestLimit(tx, req.HttpWebId, config)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// FindHTTPWebRequestLimit 查找请求限制
func (this *HTTPWebService) FindHTTPWebRequestLimit(ctx context.Context, req *pb.FindHTTPWebRequestLimitRequest) (*pb.FindHTTPWebRequestLimitResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	config, err := models.SharedHTTPWebDAO.FindWebRequestLimit(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return &pb.FindHTTPWebRequestLimitResponse{RequestLimitJSON: configJSON}, nil
}

// FindHTTPWebRequestScripts 查找请求脚本
func (this *HTTPWebService) FindHTTPWebRequestScripts(ctx context.Context, req *pb.FindHTTPWebRequestScriptsRequest) (*pb.FindHTTPWebRequestScriptsResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	config, err := models.SharedHTTPWebDAO.FindWebRequestScripts(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return &pb.FindHTTPWebRequestScriptsResponse{
		RequestScriptsJSON: configJSON,
	}, nil
}

// UpdateHTTPWebReferers 修改防盗链设置
func (this *HTTPWebService) UpdateHTTPWebReferers(ctx context.Context, req *pb.UpdateHTTPWebReferersRequest) (*pb.RPCSuccess, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(tx, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var config = &serverconfigs.ReferersConfig{}
	if len(req.ReferersJSON) > 0 {
		err = json.Unmarshal(req.ReferersJSON, config)
		if err != nil {
			return nil, err
		}

		err = config.Init()
		if err != nil {
			return nil, errors.New("validate referers config failed: " + err.Error())
		}
	}

	err = models.SharedHTTPWebDAO.UpdateWebReferers(tx, req.HttpWebId, config)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// FindHTTPWebReferers 查找防盗链设置
func (this *HTTPWebService) FindHTTPWebReferers(ctx context.Context, req *pb.FindHTTPWebReferersRequest) (*pb.FindHTTPWebReferersResponse, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(tx, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	referersJSON, err := models.SharedHTTPWebDAO.FindWebReferers(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}

	return &pb.FindHTTPWebReferersResponse{
		ReferersJSON: referersJSON,
	}, nil
}

// UpdateHTTPWebUserAgent 修改UserAgent设置
func (this *HTTPWebService) UpdateHTTPWebUserAgent(ctx context.Context, req *pb.UpdateHTTPWebUserAgentRequest) (*pb.RPCSuccess, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(tx, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	var config = &serverconfigs.UserAgentConfig{}
	if len(req.UserAgentJSON) > 0 {
		err = json.Unmarshal(req.UserAgentJSON, config)
		if err != nil {
			return nil, err
		}

		err = config.Init()
		if err != nil {
			return nil, errors.New("validate user-agent config failed: " + err.Error())
		}
	}

	err = models.SharedHTTPWebDAO.UpdateWebUserAgent(tx, req.HttpWebId, config)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// FindHTTPWebUserAgent 查找UserAgent设置
func (this *HTTPWebService) FindHTTPWebUserAgent(ctx context.Context, req *pb.FindHTTPWebUserAgentRequest) (*pb.FindHTTPWebUserAgentResponse, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		err = models.SharedHTTPWebDAO.CheckUserWeb(tx, userId, req.HttpWebId)
		if err != nil {
			return nil, err
		}
	}

	userAgentJSON, err := models.SharedHTTPWebDAO.FindWebUserAgent(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}

	return &pb.FindHTTPWebUserAgentResponse{
		UserAgentJSON: userAgentJSON,
	}, nil
}

// FindServerIdWithHTTPWebId 根据WebId查找ServerId
func (this *HTTPWebService) FindServerIdWithHTTPWebId(ctx context.Context, req *pb.FindServerIdWithHTTPWebIdRequest) (*pb.FindServerIdWithHTTPWebIdResponse, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if req.HttpWebId <= 0 {
		return nil, errors.New("invalid 'httpWebId'")
	}

	var tx = this.NullTx()
	serverId, err := models.SharedHTTPWebDAO.FindWebServerId(tx, req.HttpWebId)
	if err != nil {
		return nil, err
	}

	if serverId > 0 && userId > 0 {
		err = models.SharedServerDAO.CheckUserServer(tx, userId, serverId)
		if err != nil {
			return nil, err
		}
	}

	return &pb.FindServerIdWithHTTPWebIdResponse{ServerId: serverId}, nil
}
