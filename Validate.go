package main

//注释
import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"
)

type Validate struct {
	Phone string    //电话号码
	Code  string    // 验证码
	Dup   time.Time //短期重复发送限制
	Exp   time.Time //验证码过期时间
	LimT  time.Time //发送限制日期
	Lim   int       //同一天发送次数
	Validater
}

type Validater interface {
	ValidatePhone(phone string) bool
	GenerateCode() string
	encrypt(code string, key []byte) string
	decrypt(encode string, key []byte) string
}

func NewVal() *Validate {
	now := time.Now()
	return &Validate{
		Phone: "",
		Code:  "",
		Dup:   time.Time{},
		Exp:   time.Time{},
		LimT:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
		Lim:   0,
	}
}

// ValidatePhone 验证电话号码格式
func ValidatePhone(phone string) bool {
	reg := regexp.MustCompile(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`)
	return reg.MatchString(phone) //MatchString匹配包含正则字符串的，但^和$必须要匹配整个字符串，所以重复输入的不会匹配
}

// GenerateCode 生成验证码
func GenerateCode() string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdfghijklmnopqrstuvwxyz"
	code := ""
	for i := 0; i < 6; i++ {
		code += string(str[rand.Intn(len(str))]) //Intn返回[0,n)的随机数
	}
	return code
}

// 加密
func encrypt(code string, key []byte) string {
	plaintext := []byte(code)        //明文
	block, err := aes.NewCipher(key) //创建AES密码块
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))   //密文,长度为明文长度+AES块大小(16字节)
	iv := ciphertext[:aes.BlockSize]                           //初始向量，初始向量是一个随机生成的固定长度的字节数组，用于增加加密数据的随机性和安全性。在使用分组密码模式加密数据时，每个块都需要使用前一个块的密文和当前块的明文进行加密。但是，在加密第一个块时，没有前一个块的密文，因此需要使用初始向量来代替前一个块的密文。
	stream := cipher.NewCFBEncrypter(block, iv)                //CFB 加密模式是一种分组密码模式，它将前一个密文块作为输入，生成一个伪随机数流，然后将该流与明文块进行异或运算，生成密文块。
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext) //异或加密
	Encode := base64.StdEncoding.EncodeToString(ciphertext)    //base64编码
	return Encode
}

// 解密
func decrypt(encode string, key []byte) string {
	ciphertext, err := base64.StdEncoding.DecodeString(encode)
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext[aes.BlockSize:]))
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	code := string(plaintext)
	return code
}

func main() {
	key := []byte("1767938490412345")
	Val := NewVal()
	now := time.Now()
	Date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	flag := false
	flagVal := NewVal() //用于存储已存在的电话号码的结构体
	for {
		fmt.Println("请输入电话号码：(输入quit退出)")
		_, err := fmt.Scanln(&Val.Phone)
		if err != nil {
			panic(err)
		}
		if Val.Phone == "quit" {
			return
		}
		if Val.ValidatePhone(Val.Phone) {
			break
		} else {
			fmt.Println("电话号码格式错误！")
		}
	} //验证电话号码格式循环

	f, err := os.OpenFile("TimeOut.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	var jsonArray []Validate
	decoder := json.NewDecoder(f)    //解码器
	err = decoder.Decode(&jsonArray) //解码
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	//检查是否已存在该电话号码
	for i, v := range jsonArray {
		if v.Phone == Val.Phone {
			flagVal = &jsonArray[i]
			Val.Dup = v.Dup
			Val.Exp = v.Exp
			Val.LimT = v.LimT
			Val.Lim = v.Lim
			flag = true
			break
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1:请输入验证码登录 2：获取验证码")
		num := ""
		_, err := fmt.Scanln(&num)
		if err != nil {
			panic(err)
		}
		switch num {
		case "1":
			fmt.Println("请输入验证码：")
			scanner.Scan()
			code := scanner.Text()
			if flag {
				Val.Code = Val.decrypt(flagVal.Code, key)
			} else {
				Val.Code = Val.decrypt(Val.Code, key)
			}
			if code == Val.Code {
				if time.Now().Before(Val.Exp) {
					fmt.Println("登录成功！")
					return
				} else if time.Now().After(Val.Exp) {
					fmt.Println("验证码已过期！")
				}
			} else {
				fmt.Println("验证码错误！")
			}

		case "2":
			if time.Now().Before(Val.Dup) {
				fmt.Println("一分钟内已获取验证码无法重复获取")
			} else {
				if Val.Lim >= 5 && Date == Val.LimT {
					fmt.Println("今日获取验证码已达上限，明天再来吧！")
					return
				} else {
					Val.Code = Val.GenerateCode()
					if Date != Val.LimT {
						Val.Lim = 0
						Val.LimT = Date
					}
					Val.Lim += 1
					fmt.Println("验证码已发送，请注意查收！")
					fmt.Print("验证码为：")
					fmt.Println(Val.Code)
					Val.Dup = time.Now().Add(60 * time.Second)
					Val.Exp = time.Now().Add(300 * time.Second)
					Val.Code = Val.encrypt(Val.Code, key)
					if flag {
						flagVal.Code = Val.Code
						flagVal.Dup = Val.Dup
						flagVal.Exp = Val.Exp
						flagVal.LimT = Val.LimT
						flagVal.Lim += 1
					} else {
						jsonArray = append(jsonArray, *Val)
					}
					encoder := json.NewEncoder(f)
					encoder.SetIndent("", "    ")
					f, err = os.OpenFile("TimeOut.json", os.O_TRUNC|os.O_WRONLY, 0766)
					if err != nil {
						panic(err)
					}
					defer func(f *os.File) {
						err := f.Close()
						if err != nil {
							panic(err)
						}
					}(f)
					err = encoder.Encode(jsonArray)
					if err != nil {
						panic(err)
					}
					//以下五行为测试用，可删
					jsonData, err := json.Marshal(Val)
					if err != nil {
						panic(err)
					}
					fmt.Println(string(jsonData))
				}
			}
		default:
			fmt.Println("输入错误！请重新输入。")
		}
	}
}
