package main

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
func ValidatePhone(phone string) bool {
	reg := regexp.MustCompile(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`)
	return reg.MatchString(phone)
}

func GenerateCode() string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdfghijklmnopqrstuvwxyz"
	code := ""
	for i := 0; i < 6; i++ {
		code += string(str[rand.Intn(len(str))])
	}
	return code
}
func encrypt(code string, key []byte) string {
	plaintext := []byte(code)        //明文
	block, err := aes.NewCipher(key) //加密
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))   //密文
	iv := ciphertext[:aes.BlockSize]                           //初始化向量
	stream := cipher.NewCFBEncrypter(block, iv)                //加密
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext) //加密
	Encode := base64.StdEncoding.EncodeToString(ciphertext)    //base64编码
	return Encode
}

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
	flagVal := NewVal()
	for {
		fmt.Println("请输入电话号码：(输入quit退出)")
		_, err := fmt.Scanln(&Val.Phone)
		if err != nil {
			panic(err)
		}
		if Val.Phone == "quit" {
			return
		}
		if ValidatePhone(Val.Phone) {
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
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&jsonArray)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
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
				Val.Code = decrypt(flagVal.Code, key)
			} else {
				Val.Code = decrypt(Val.Code, key)
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
					Val.Code = GenerateCode()
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
					if flag {
						flagVal.Code = encrypt(Val.Code, key)
						flagVal.Dup = Val.Dup
						flagVal.Exp = Val.Exp
						flagVal.LimT = Val.LimT
						flagVal.Lim += 1
					}
					SaveVal := Val
					SaveVal.Code = encrypt(Val.Code, key)
					if !flag {
						jsonArray = append(jsonArray, *SaveVal)
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
