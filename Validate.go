package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"
)

type Val struct {
	Phone string    //电话号码
	Code  string    // 验证码
	Dup   time.Time //短期重复发送限制
	Exp   time.Time //验证码过期时间
	LimT  int       //发送限制日期
	Lim   int       //同一天发送次数
}

func NewVal() *Val {
	return &Val{
		Phone: "",
		Code:  "",
		Dup:   time.Time{},
		Exp:   time.Time{},
		LimT:  time.Now().Day(),
		Lim:   0,
	}
}

func main() {
	Val := NewVal()
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
	for {
		fmt.Println("1:请输入验证码登录 2：获取验证码")
		var num int
		_, err := fmt.Scanln(&num)
		if err != nil {
			panic(err)
		}
		switch num {
		case 1:
			fmt.Println("请输入验证码：")
			var code string
			_, err := fmt.Scanln(&code)
			if err != nil {
				panic(err)
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

		case 2:
			if time.Now().Before(Val.Dup) {
				fmt.Println("一分钟内已获取验证码无法重复获取")
			} else {
				if Val.Lim > 5 && time.Now().Day() == Val.LimT {
					fmt.Println("今日获取验证码已达上限，明天再来吧！")
				} else {
					Val.Code = GenerateCode()
					if time.Now().Day() != Val.LimT {
						Val.Lim = 0
						Val.LimT = time.Now().Day()
					}
					Val.Lim += 1
					fmt.Println("验证码已发送，请注意查收！")
					fmt.Print("验证码为：")
					fmt.Println(Val.Code)
					Val.Dup = time.Now().Add(60 * time.Second)
					Val.Exp = time.Now().Add(300 * time.Second)
					_, err := os.OpenFile("TimeOut.json", os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						panic(err)
					}
					jsonData, err := json.Marshal(Val)
					if err != nil {
						panic(err)
					}
					fmt.Println(string(jsonData))
					err = os.WriteFile("TimeOut.json", jsonData, 0644)
					if err != nil {
						panic(err)
					}
				}
			}
		default:
			fmt.Println("输入错误！")
		}
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
