package auth

import "context"

// contextKey 是用于context的键类型，防止键冲突
type contextKey string

const (
	// UserOpenIDKey 用户OpenID的context键
	UserOpenIDKey contextKey = "user_open_id"
	// UserNameKey 用户姓名的context键
	UserNameKey contextKey = "user_name"
)

// GetUserOpenID 从context中获取用户OpenID
func GetUserOpenID(ctx context.Context) (string, bool) {
	openID, ok := ctx.Value(UserOpenIDKey).(string)
	return openID, ok
}

// GetUserName 从context中获取用户姓名
func GetUserName(ctx context.Context) (string, bool) {
	userName, ok := ctx.Value(UserNameKey).(string)
	return userName, ok
}
