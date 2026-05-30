package framework

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/db_param"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/helper/goJose"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/go-jose/go-jose/v4"
)

var methodOfEncryptBody = map[string]bool{
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

// 访问控制中间件服务层
type AccessMiddleware struct {
}

func (s *AccessMiddleware) DealSignAndEncrypt(ctx *beecontext.Context) error {
	var err error
	moduleName := ctx.Input.GetData(constant.MODULE_NAME)
	contentType := ctx.Request.Header.Get("Content-Type")
	reqMethod := helper.GetReqMethod(ctx)
	if _, ok := methodOfEncryptBody[reqMethod]; !ok {
		return nil
	}
	if !strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "multipart/form-data") {
		return nil
	}

	// 代码配置文件是否开启api安全配置
	encryptEnable, err1 := s.checkApiSecurityConfig("apiSecurity.encrypt")
	signatureEnable, err2 := s.checkApiSecurityConfig("apiSecurity.signature")
	if err1 != nil || err2 != nil {
		err := helper.Ternary(err1 != nil, err1, err2)
		return err
	}

	// 请求：body加密数据
	bodyEncryptDto, err := helper.GetReqBody[dto.BodyEncryptDto](ctx)
	if err != nil {
		return err
	}
	if bodyEncryptDto.EncryptedData == "" {
		tmpBody := "{\"error\":\"DealSignAndEncrypt(): body encrypted data is empty\"}"
		s.doResetBody(ctx, []byte(tmpBody))
		return nil
	}

	// 请求: header头携带的签名字符串
	signatureStr, err := s.getHeaderSignature(ctx)
	if err != nil {
		return err
	} else if signatureEnable && signatureStr == "" {
		global.Log.Error("未知签名")
		return errors.New("未知签名")
	}

	// 数据库保存rsa配置名称
	configSignName := ""
	configEncryptName := ""
	switch moduleName {
	case constant.ADMIN_MCHNT, constant.ADMIN_PLAT:
		configSignName = constant.ADMIN_API_SIGN_KEY
		configEncryptName = constant.ADMIN_API_ENCRYPT_KEY
	case constant.API:
		configSignName = constant.OPEN_API_SIGN_KEY
		configEncryptName = constant.OPEN_API_ENCRYPT_KEY
	default:
		err = errors.New("未知模块")
	}
	if err != nil {
		return err
	}

	// 从数据库获取rsa密钥
	configModelMap := make(map[string][]models.Config)
	if encryptEnable || signatureEnable {
		configModelMap, err1 = s.getDatabaseRsaData(ctx)
	}
	if err1 != nil {
		return err1
	}

	// 解密body, 并重新覆盖请求body
	if encryptEnable {
		configModel := configModelMap[configEncryptName][0]
		bodyEncryptStr := bodyEncryptDto.EncryptedData
		err = s.decryptBodyAndReset(ctx, configModel, bodyEncryptStr)
	}
	if err != nil {
		return err
	}

	// 验证签名
	if signatureEnable {
		configModel := configModelMap[configSignName][0]
		err = s.verifyBodySign(ctx, configModel, signatureStr)
	}
	return err
}

// 解密请求体并覆盖
func (s *AccessMiddleware) decryptBodyAndReset(ctx *beecontext.Context, configModel models.Config, bodyEncryptStr string) error {
	rsaConfigData, err := s.parseSignatureRsaKey(configModel)
	if err != nil {
		return err
	}
	privateKey, err := helper.RsaPemToPrivateKey(rsaConfigData.Private)
	if err != nil {
		return err
	}
	bodyDecryptByte, err := goJose.JweDecrypt(bodyEncryptStr, privateKey, jose.A256CBC_HS512, jose.RSA_OAEP_256)
	if err != nil {
		return err
	}
	// 解密成功，重置body
	s.doResetBody(ctx, bodyDecryptByte)

	return nil
}
func (s *AccessMiddleware) doResetBody(ctx *beecontext.Context, newBody []byte) {
	ctx.Input.RequestBody = newBody
}

// 验证请求体签名
func (s *AccessMiddleware) verifyBodySign(ctx *beecontext.Context, configModel models.Config, signatureStr string) error {
	rsaConfigData, err := s.parseSignatureRsaKey(configModel)
	if err != nil {
		return err
	}
	publicKey, err := helper.RsaPemToPublicKey(rsaConfigData.Public)
	if err != nil {
		return err
	}

	bodyMap := make(map[string]interface{})
	json.Unmarshal(ctx.Input.RequestBody, &bodyMap)
	payload := goJose.SerializeForPayload(bodyMap)

	err = goJose.JwsSignVerify([]byte(payload), signatureStr, publicKey, jose.RS512)
	if err != nil {
		logErrStr := fmt.Sprintf("签名验证失败：%s\npayload=%s", err.Error(), payload)
		global.Log.Error(logErrStr)

		logErrStr = helper.Ternary(helper.IsDevRunMode(), logErrStr, "签名验证失败")
		return errors.New(logErrStr)
	}
	return nil
}

// 是否开启接口安全配置
func (s *AccessMiddleware) checkApiSecurityConfig(configKey string) (bool, error) {
	configItf, err := global.GetConfigDiy(configKey)
	if err != nil || configItf == nil {
		msg := "配置apiSecurity.signature错误："
		global.Log.Error(msg + err.Error())
		return false, errors.New(msg)
	}
	return configItf.(bool), nil
}

// 获取请求头签名
func (s *AccessMiddleware) getHeaderSignature(ctx *beecontext.Context) (string, error) {
	signatureStr := helper.GetReqSignature(ctx)
	signatureStr = strings.TrimSpace(signatureStr)
	return signatureStr, nil
}

// 获取数据库签名密钥
func (s *AccessMiddleware) getDatabaseRsaData(ctx *beecontext.Context) (configModelMap map[string][]models.Config, err error) {
	moduleName := ctx.Input.GetData(constant.MODULE_NAME)
	configModelMap = make(map[string][]models.Config, 0)
	configName := make([]string, 0)

	switch moduleName {
	case constant.ADMIN_MCHNT, constant.ADMIN_PLAT:
		configName = append(configName, constant.ADMIN_API_SIGN_KEY)
		configName = append(configName, constant.ADMIN_API_ENCRYPT_KEY)

	case constant.API:
		configName = append(configName, constant.OPEN_API_SIGN_KEY)
		configName = append(configName, constant.OPEN_API_ENCRYPT_KEY)
	default:
		err = errors.New("未知模块")
	}
	if err != nil {
		return
	}

	configModels, err := (&models_ar.ConfigAr{}).GetByNames(configName)
	if err != nil {
		return configModelMap, err
	}
	for k, v := range configModels {
		if _, ok := configModelMap[v.Name]; !ok {
			configModelMap[v.Name] = make([]models.Config, 0)
		}
		configModelMap[v.Name] = append(configModelMap[v.Name], configModels[k])
	}

	return configModelMap, nil
}

// 获取数据库签名密钥 - 读取签名密钥
func (s *AccessMiddleware) parseSignatureRsaKey(configModel models.Config) (rsaConfigData *db_param.Db_config_rsa, err error) {
	tmpValue := configModel.Value
	tmpValue = strings.ReplaceAll(tmpValue, "\r", "")
	tmpValue = strings.ReplaceAll(tmpValue, "\n", "\\n")

	rsaConfigData = &db_param.Db_config_rsa{}
	err = json.Unmarshal([]byte(tmpValue), &rsaConfigData)
	return rsaConfigData, err
}
