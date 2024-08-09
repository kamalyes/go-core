/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 11:02:00
 * @FilePath: \go-core\jwt\jwt.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamalyes/go-core/global"
	"gorm.io/gorm"
)

// 定义一些常量
var (
	TokenExpired     error = errors.New("Token 已经过期")
	TokenNotValidYet error = errors.New("Token 尚未激活")
	TokenMalformed   error = errors.New("Token 格式错误")
	TokenInvalid     error = errors.New("Token 无法解析")
	jwtSignKey             = "82011FC650590620FEFAC6500ADAB0F77" // 默认签名用的key
)

// JWT jwt签名结构
type JWT struct {
	SigningKey []byte
}

// SetJWTSignKey 动态设置JWT签名密钥
func SetJWTSignKey(key string) {
	jwtSignKey = key
}

// GetJWTSignKey 获取JWT签名密钥
func GetJWTSignKey() string {
	return jwtSignKey
}

// NewJWT 新建一个 jwt 实例
func NewJWT() *JWT {
	return &JWT{[]byte(GetJWTSignKey())}
}

// RegisteredClaims expiresAt 过期时间单位秒
func RegisteredClaims(issuer string, expiresAt int64) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
	}
}

// CreateToken 生成 token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 判断多点登录拦截是否开启
	if global.CONFIG.JWT.UseMultipoint {
		// 拦截
		if global.REDIS != nil {
			// 优先存入到 redis
			jsonData, _ := json.Marshal(claims)
			toJson := string(jsonData)
			// 此处过期时间等于jwt过期时间
			timer := time.Duration(global.CONFIG.JWT.ExpiresTime) * time.Second
			err := global.REDIS.Set(context.Background(), claims.UserId, toJson, timer).Err()
			if err != nil {
				return "", err
			}
			return token.SignedString(j.SigningKey)
		}
		// 没有redis存入到 数据库
		err := global.DB.Save(&claims).Error
		if err != nil {
			return "", err
		} else {
			return token.SignedString(j.SigningKey)
		}
	}
	// 不拦截
	return token.SignedString(j.SigningKey)
}

// DeleteToken 强制删除Token记录，用途--用户账号被盗后，强制下线
func DeleteToken(userId string) (err error) {
	if global.REDIS != nil {
		err = global.REDIS.Del(context.Background(), userId).Err()
		return err
	}
	err = global.DB.Where("user_id = ?", userId).Delete(&CustomClaims{}).Error
	return err
}

// ResolveToken 解析token
func (j *JWT) ResolveToken(tokenString string) (*CustomClaims, error) {
	token, parseErr := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if parseErr != nil {
		return handleTokenParseError(parseErr)
	}

	if token != nil && token.Valid {
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			return nil, TokenInvalid
		}

		if !j.isMultipointAuthEnabled(claims) {
			return claims, nil
		}

		if err := j.checkMultipointAuth(claims); err != nil {
			return nil, err
		}

		return claims, nil
	}

	return nil, TokenInvalid
}

// handleTokenParseError 处理token解析错误
func handleTokenParseError(err error) (*CustomClaims, error) {
	if ve, ok := err.(*jwt.ValidationError); ok {
		switch {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, TokenMalformed
		case ve.Errors&jwt.ValidationErrorExpired != 0:
			return nil, TokenExpired
		case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	return nil, err
}

// isMultipointAuthEnabled 检查是否启用多点登录拦截
func (j *JWT) isMultipointAuthEnabled(claims *CustomClaims) bool {
	return global.CONFIG.JWT.UseMultipoint
}

// checkMultipointAuth 检查多点登录验证
func (j *JWT) checkMultipointAuth(claims *CustomClaims) error {
	if global.REDIS != nil {
		if jsonStr, err := global.REDIS.Get(context.Background(), claims.UserId).Result(); err == redis.Nil || jsonStr == "" {
			return nil
		} else {
			return j.checkRedisMultipointAuth(claims, jsonStr)
		}
	}
	return j.checkDBMultipointAuth(claims)
}

// checkRedisMultipointAuth 检查Redis中的多点登录验证
func (j *JWT) checkRedisMultipointAuth(claims *CustomClaims, jsonStr string) error {
	var clis CustomClaims
	if err := json.Unmarshal([]byte(jsonStr), &clis); err != nil {
		return errors.New("解析Redis中的用户token时出错: " + err.Error())
	}

	if clis.TokenId != "" && claims.TokenId != clis.TokenId {
		return errors.New("账号已在其他地方登录，您已被迫下线")
	}

	return nil
}

// checkDBMultipointAuth 检查数据库中的多点登录验证
func (j *JWT) checkDBMultipointAuth(claims *CustomClaims) error {
	var clis CustomClaims
	if err := global.DB.Where("user_id = ?", claims.UserId).First(&clis).Error; err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("从数据库获取用户token异常：" + err.Error())
	}

	if claims.TokenId != clis.TokenId {
		return errors.New("账号已在其他地方登录，您已被迫下线")
	}

	return nil
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(time.Now().Unix()+global.CONFIG.JWT.ExpiresTime, 0))
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}

// GetClaims 获取Claims
func GetClaims(c *gin.Context) (*CustomClaims, error) {
	if claims, exists := c.Get("claims"); !exists {
		global.LOG.Error("从Gin的Context中获取从jwt解析出来的用户claims失败, 请检查路由是否使用jwt中间件")
		return nil, errors.New("获取用户用户claims失败")
	} else {
		token := claims.(*CustomClaims)
		return token, nil
	}
}

// ClaimHandlerFunc 定义处理声明的函数
type ClaimHandlerFunc func(*CustomClaims) interface{}

// ClaimHandlers 存储不同类型Claim的处理函数
var ClaimHandlers = map[string]ClaimHandlerFunc{
	"TokenId":      func(claims *CustomClaims) interface{} { return claims.TokenId },
	"UserId":       func(claims *CustomClaims) interface{} { return claims.UserId },
	"UserName":     func(claims *CustomClaims) interface{} { return claims.UserName },
	"UserType":     func(claims *CustomClaims) interface{} { return claims.UserType },
	"NickName":     func(claims *CustomClaims) interface{} { return claims.NickName },
	"PhoneNumber":  func(claims *CustomClaims) interface{} { return claims.PhoneNumber },
	"MerchantNo":   func(claims *CustomClaims) interface{} { return claims.MerchantNo },
	"AuthorityId":  func(claims *CustomClaims) interface{} { return claims.AuthorityId },
	"AppProductId": func(claims *CustomClaims) interface{} { return claims.AppProductId },
	"PlatformType": func(claims *CustomClaims) interface{} { return claims.PlatformType },
	"BufferTime":   func(claims *CustomClaims) interface{} { return claims.BufferTime },
	"Extend":       func(claims *CustomClaims) interface{} { return claims.Extend },
}

// GetClaimValue 从Gin的Context中获取特定类型的Claim值，通过ClaimHandlers映射来获取
func GetClaimValue(c *gin.Context, key string) interface{} {
	claims, exists := c.Get("claims")
	if !exists {
		return nil
	}

	customClaims, ok := claims.(*CustomClaims)
	if !ok {
		return nil
	}

	handler, found := ClaimHandlers[key]
	if !found {
		return nil
	}

	return handler(customClaims)
}

// GetStringClaimValue 从Gin的Context中获取字符串类型的Claim值
func GetStringClaimValue(c *gin.Context, key string) string {
	value := GetClaimValue(c, key)
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return ""
}

// GetInt32ClaimValue 从Gin的Context中获取Int32类型的Claim值
func GetInt32ClaimValue(c *gin.Context, key string) int32 {
	value := GetClaimValue(c, key)
	if intValue, ok := value.(int32); ok {
		return intValue
	}
	return 0
}

// GetInt64ClaimValue 从Gin的Context中获取Int64类型的Claim值
func GetInt64ClaimValue(c *gin.Context, key string) int64 {
	value := GetClaimValue(c, key)
	if intValue, ok := value.(int64); ok {
		return intValue
	}
	return 0
}

// GetTokenId 获取Token Id
func GetTokenId(c *gin.Context) string {
	return GetStringClaimValue(c, "TokenId")
}

// GetUserId 获取用户Id
func GetUserId(c *gin.Context) string {
	return GetStringClaimValue(c, "UserId")
}

// GetUserName 获取用户名
func GetUserName(c *gin.Context) string {
	return GetStringClaimValue(c, "UserName")
}

// GetUserType 获取用户类型
func GetUserType(c *gin.Context) string {
	return GetStringClaimValue(c, "UserType")
}

// GetNickName 获取用户昵称
func GetNickName(c *gin.Context) string {
	return GetStringClaimValue(c, "NickName")
}

// GetPhoneNumber 获取用户手机号
func GetPhoneNumber(c *gin.Context) string {
	return GetStringClaimValue(c, "PhoneNumber")
}

// GetMerchantNo 获取商户号
func GetMerchantNo(c *gin.Context) string {
	return GetStringClaimValue(c, "MerchantNo")
}

// GetUserAuthorityId 获取用户角色Id
func GetUserAuthorityId(c *gin.Context) string {
	return GetStringClaimValue(c, "AuthorityId")
}

// GetAppProductId 获取AppProduct Id
func GetAppProductId(c *gin.Context) int32 {
	return GetInt32ClaimValue(c, "AppProductId")
}

// GetPlatformType 获取Platform Type
func GetPlatformType(c *gin.Context) int32 {
	return GetInt32ClaimValue(c, "PlatformType")
}

// GetBufferTime 获取BufferTime
func GetBufferTime(c *gin.Context) int64 {
	return GetInt64ClaimValue(c, "BufferTime")
}

// GetExtend 获取Extend
func GetExtend(c *gin.Context) string {
	return GetStringClaimValue(c, "Extend")
}
