package model

import (
	"github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"regexp"
	"time"
)

var tokenKey = []byte("lowcode")

type Token struct {
	UserId    int
	UserName  string
	UserEmail string
	jwt.RegisteredClaims
}

func UnixToDate(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

func GetUnix() int64 {
	return time.Now().Unix()
}

func BcryptPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CheckBcryptPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// CheckUserName 用户名只能包含中文、英文、数字、下划线，长度为3-15个字符
func CheckUserName(userName string) bool {
	regex := `^[\p{Han}a-zA-Z0-9_]{3,15}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(userName)
}

// CheckPassword 6-50个字符，至少一个大写字母，一个小写字母和一个数字，其他可以是任意字符（除了空白符）
func CheckPassword(password string) bool {
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^\s]{6,50}$`
	re := regexp2.MustCompile(regex, 0)
	match, _ := re.MatchString(password)
	return match
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func SendEmail(toEmail string, subject string, body string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "yuhui_z@foxmail.com", "医学低代码")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.qq.com", 587, "yuhui_z@foxmail.com", "mjkouqozyfjpcahh")

	go func() {
		_ = d.DialAndSend(m)
	}()
}

func SignToken(userId int, userName string, userEmail string) (string, error) {
	claims := Token{
		UserId:    userId,
		UserName:  userName,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tokenKey)
}

func ParseToken(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return tokenKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Token), nil
}
