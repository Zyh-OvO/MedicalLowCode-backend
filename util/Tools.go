package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/gomail.v2"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const Device = "cpu"
const Env = "/opt/miniconda3/etc/profile.d/conda.sh"
const TrainPlanFolder = "nnUNetTrainer__nnUNetPlans__3d_fullres"
const TrainFoldFolder = "fold_all"

var tokenKey = []byte("lowcode")

type Token struct {
	UserId    int
	UserName  string
	UserEmail string
	jwt.RegisteredClaims
}

func UnixToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02-15-04-05")
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
	//todo:错误处理
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Token), nil
}

// SetDefault v 需要为一个结构体指针
func SetDefault(v any) {
	value := reflect.ValueOf(v).Elem()
	typ := value.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			fieldValue := value.Field(i)
			switch field.Type.Kind() {
			case reflect.Int:
				defaultValueInt, _ := strconv.Atoi(defaultValue)
				fieldValue.SetInt(int64(defaultValueInt))
			case reflect.String:
				fieldValue.SetString(defaultValue)
			case reflect.Bool:
				defaultValueBool, _ := strconv.ParseBool(defaultValue)
				fieldValue.SetBool(defaultValueBool)
			case reflect.Slice:
				switch field.Type.Elem().Kind() {
				case reflect.Int:
				}
			}
		}
	}
}

// CamelCaseToSnakeCase 驼峰转蛇形
func CamelCaseToSnakeCase(input string) string {
	var buffer bytes.Buffer

	for i, char := range input {
		if unicode.IsUpper(char) && i > 0 {
			buffer.WriteRune('_')
		}
		buffer.WriteRune(unicode.ToLower(char))
	}

	return buffer.String()
}

func GetSSHConfig() (*ssh.ClientConfig, error) {
	key, err := os.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: "wxgroup1",
		Auth: []ssh.AuthMethod{
			//ssh.Password("Buaa123456"),
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 不验证主机密钥
	}
	return config, nil
}

func SliceContains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func GiveStaticToken() (*Token, error) {
	token := &Token{
		UserId:    0,
		UserName:  "nii",
		UserEmail: "nii@test.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return token, nil
}

func CreateFolderIfNotExists(folderPath string) bool {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false
	}
	folderPath = filepath.Join(currentDir, folderPath)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return false
		} else {
			fmt.Println("Folder created successfully:", folderPath)
			return true
		}
	} else {
		fmt.Println("Folder already exists, no action taken:", folderPath)
		return true
	}
}

func Unzip(src, dest string) error {
	// 打开 ZIP 文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	// 解压 ZIP 文件中的每个文件
	for _, f := range r.File {
		// 打开文件
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 跳过奇怪的文件
		if strings.Contains(f.Name, "__MACOSX") {
			continue
		}

		// 创建目标文件
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		// 创建目标文件
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将文件内容拷贝到目标文件
		_, err = io.Copy(file, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindFileNumFolder(folder string) string {
	destFolder := folder
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.Compare(path, folder) != 0 {
			if !strings.HasSuffix(path, "imagesTr") && !strings.HasSuffix(path, "imagesTs") && !strings.HasSuffix(path, "labelsTr") && !strings.HasSuffix(path, "labelsTs") {
				destFolder = path + "/"
				fmt.Println("data in folder:", destFolder)
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历文件夹时发生错误：", err)
	}

	return destFolder
}

func FindFileNumInFolder(folder string, extension string) int {

	count := 0

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, extension) {
			count++
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历文件夹时发生错误：", err)
		return -1
	}

	fmt.Printf("%s文件夹中以%s结尾的文件数量为：%d\n", folder, extension, count)
	return count
}
