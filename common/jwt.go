package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MySecret 定义一个密钥，实际项目建议放在配置文件里
var MySecret = []byte("这是我的加密密钥")

// MyClaims 自定义声明结构体，也就是你要存入 Token 的信息
type MyClaims struct {
	UserID               uint   `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 内置的标准声明，包含过期时间等
}

// GenerateToken 生成 JWT
func GenerateToken(userID uint, username string) (string, error) {
	// 1. 创建我们要存的信息
	claims := MyClaims{
		userID,
		username,
		jwt.RegisteredClaims{
			// 设置 24 小时后过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// 签发人
			Issuer: "my-gin-project",
		},
	}
	// 2. 使用指定的签名方法创建 Token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 3. 使用密钥签名并获得完整的字符串 Token
	return token.SignedString(MySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析并校验 Token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 这行代码最核心：告诉解析器，用哪个密钥去验证签名
		return MySecret, nil
	})

	if err != nil {
		// 这里可以细分错误类型
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("Token 已过期")
		}
		return nil, errors.New("无效的 Token")
	}

	// 将解析出来的 Claims 强转回我们自定义的结构体
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 Token")
}
