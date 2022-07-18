package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/littlefish12345/FishBot"
)

func main() {
	device := &FishBot.DeviceInfo{}
	_, err := os.Stat("device.json")
	if os.IsNotExist(err) {
		device = FishBot.NewDevice()
		os.WriteFile("device.json", device.ToJson(), 0644)
	} else {
		deviceData, _ := os.ReadFile("device.json")
		device.FromJson(deviceData)
	}

	var uinString string
	var password string
	fmt.Print("Uin: ")
	fmt.Scanln(&uinString)
	fmt.Print("Password: ")
	passwordByte, _ := gopass.GetPasswd()
	//fmt.Scanln(&password)
	password = string(passwordByte)
	uin, _ := strconv.ParseInt(uinString, 10, 64)
	client, err := FishBot.NewClient(uin, md5.Sum([]byte(password)), device)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("login...")
	loginResponse, _ := client.Login(FishBot.LoginMethodPassword)
	handleLogin(client, loginResponse)
	fmt.Println("login done")
	fmt.Println("friends:")
	friendList := client.GetFriendList()
	for _, friendInfo := range friendList {
		fmt.Println(friendInfo)
	}
	fmt.Println("troops:")
	troopList := client.GetGroupList()
	for _, troop := range troopList {
		fmt.Println(troop)
	}
	ch := make(chan int)
	<-ch
}

func handleLogin(client *FishBot.QQClient, loginResponse *FishBot.LoginResponse) {
	for !loginResponse.Success {
		if loginResponse.Error == FishBot.LoginResponseNeedSlider {
			fmt.Println("need slider")
			urlList, _ := FishBot.StartSliderCaptchaServer()
			var replacedUrlList []string
			for _, url := range urlList {
				replacedUrlList = append(replacedUrlList, strings.Replace(loginResponse.SliderVerifyUrl, "https://ssl.captcha.qq.com/template/wireless_mqq_captcha.html", url, 1))
			}
			fmt.Println(replacedUrlList)
			ticket := FishBot.GetSliderTicket()
			loginResponse = client.SubmitSliderTicket(ticket)
			continue
		} else if loginResponse.Error == FishBot.LoginResponseNeedSMS {
			client.RequestSMSCode()
			var smsCode string
			fmt.Printf("SMS Code: ")
			fmt.Scanf("%s", &smsCode)
			loginResponse, _ = client.SubmitSMSCode(smsCode)
		} else if loginResponse.Error == FishBot.LoginResponseOtherError {
			panic(loginResponse.Message)
		} else {
			fmt.Println(loginResponse)
			panic("unknow type")
		}
	}
}
