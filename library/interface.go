package library

import "time"

var (
	Jwt_Token = Interface_Jwt(NewJwtGo())
	I18n = Interface_I18n(&Go_i18n{})
	Cache = Interface_Cache(&Redis{})
	Log = Interface_Log(&Logrus{})
)


// Jwt token interface
// JWT 令牌接口
type Interface_Jwt interface {
	CreateToken(interface{}) (string,  error)
	ParseToken(string) (interface{}, error)
	GetIdByToken(string) (uint, error)						// Get user ID
	GetIssuerByToken(string) (issuer string, err error)		// Get Username
}


// i18n interface
// i18n 接口
type Interface_I18n interface {
	Translate (ctx string, language string) string
}


// Cache interface
// 缓存接口
type Interface_Cache interface {
	GetString (key string) (string, error)
	SetString (key string, value string, time time.Duration)
	DeleteString (list ...string)
	GetHash (key string, field string) (string, error)
	SetHash (key string, value map[string]string)
	DeleteHash (key string, list ...string)
	ClearAllCache()
}


// Library_Logrus interface
// 日志接口
type Interface_Log interface {
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}