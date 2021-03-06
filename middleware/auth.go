package middleware

import (
	"Etpmls-Admin-Server/core"
	"Etpmls-Admin-Server/database"
	"Etpmls-Admin-Server/library"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)


// Only check if the token exists
// 只检查token是否存在
func BasicCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get Token
		// 获取token
		token, err := core.GetToken(c)
		if err != nil {
			core.JsonError(c, http.StatusPaymentRequired, core.ERROR_Code, core.Translate(c, "ERROR_MESSAGE_GetToken"), nil, err)
			c.Abort()
			return
		}

		// Get Claims
		// 获取Claims
		_, err = library.Jwt_Token.ParseToken(token)
		if err != nil {
			core.JsonError(c, http.StatusPaymentRequired, core.ERROR_Code, core.Translate(c, "ERROR_MESSAGE_TokenVerificationFailed"), nil, err)
			c.Abort()
			return
		}

		c.Next()
		return
	}
}

// Check whether the token exists, check whether the user's role has permissions
// 检查token是否存在，检查用户所在角色是否拥有权限
func RoleCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get Token
		// 获取token
		token, err := core.GetToken(c)
		if err != nil {
			core.JsonError(c, http.StatusPaymentRequired, core.ERROR_Code, core.Translate(c, "ERROR_MESSAGE_GetToken"), nil, err)
			c.Abort()
			return
		}

		// Get Claims
		// 获取Claims
		tmp, err := library.Jwt_Token.ParseToken(token)
		tk, ok := tmp.(*jwt.Token)
		if !ok || err != nil {
			core.JsonError(c, http.StatusPaymentRequired, core.ERROR_Code, core.Translate(c, "ERROR_MESSAGE_TokenVerificationFailed"), nil, err)
			c.Abort()
			return
		}

		// 判断所属角色是否有相应的权限
		if claims,ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
			if userId, ok := claims["jti"].(string); ok {
				b, err := PermissionCheck(c, userId)
				if err == nil && b {
					c.Next()
					return
				}
			}
		}
		// 没权限就是401
		core.JsonError(c, http.StatusUnauthorized, core.ERROR_Code, core.Translate(c, "ERROR_MESSAGE_PermissionDenied"), nil, err)
		c.Abort()
		return
	}
}


func PermissionCheck(c *gin.Context, idStr string) (b bool, err error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return b, err
	}

	// 1.获取用户ID
	var u database.User
	database.DB.Preload("Roles").First(&u, id)
	var ids []uint
	for _, v := range u.Roles {
		// 如果为管理员组
		if v.ID == 1 {
			return true, nil
		}
		ids = append(ids, v.ID)
	}
	// 获取角色相关权限
	var r []database.Role
	database.DB.Preload("Permissions").Where(ids).Find(&r)

	// 获取当前URL Path
	tmpUri, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		return b, err
	}
	uri := tmpUri.Path

	// Determine whether there is a request permission
	// 判断是否有请求权限
	for _, v := range r {
		for _, subv := range v.Permissions {

			// define an empty slice
			// 定义一个空切片
			var mtd = []string{}
			mtd = strings.Split(subv.Method, ",")

			// Path comparison
			// 路径对比
			b, _ := filepath.Match(subv.Path, uri)
			if b {

				// Method comparison
				// 方法对比
				for _, mtdv := range mtd {
					// If it is ALL, return the permission verification success directly
					// 如果是ALL直接返回权限验证成功
					if mtdv == "ALL" {
						return true, nil
					}
					// If the method is the same as the current request, return the verification success
					// 如果与当前请求方法相同，返回验证成功
					if mtdv == c.Request.Method {
						return true, nil
					}
				}

			}
		}
	}

	return false, err
}