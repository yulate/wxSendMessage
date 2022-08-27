package utils

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"wxSendMessage/defs/config"
)

/**
两个日期的自然天之差
*/
func GetTimeSubDay(bef, now string) int {
	var day int
	t1, _ := time.Parse("2006-01-02 15:04:05", bef)
	t2, _ := time.Parse("2006-01-02 15:04:05", now)
	swap := false
	if t1.Unix() > t2.Unix() {
		t1, t2 = t2, t1
		swap = true
	}

	t1_ := t1.Add(time.Duration(t2.Sub(t1).Milliseconds()%86400000) * time.Millisecond)
	day = int(t2.Sub(t1).Hours() / 24)
	// 计算在t1+两个时间的余数之后天数是否有变化
	if t1_.Day() != t1.Day() {
		day += 1
	}

	if swap {
		day = -day
	}

	return day
}

// 读取配置文件
func ReadConfigurationFile() {
	configPath := "./config.txt"
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("读取配置文件错误:%s", err)
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, "AppId:") {
			config.AppId = strings.Trim(strings.Split(line, ":")[1], "\n")
		} else if strings.Contains(line, "AppSecret:") {
			config.AppSecret = strings.Trim(strings.Split(line, ":")[1], "\n")
		} else if strings.Contains(line, "ToUserId:") {
			config.ToUserId = strings.Trim(strings.Split(line, ":")[1], "\n")
		} else if strings.Contains(line, "TemplateId:") {
			config.TemplateId = strings.Trim(strings.Split(line, ":")[1], "\n")
		} else if strings.Contains(line, "City:") {
			config.City = strings.Trim(strings.Split(line, ":")[1], "\n")
		} else if strings.Contains(line, "TianApi:") {
			config.TianApi = strings.Trim(strings.Split(line, ":")[1], "\n")
		} else if strings.Contains(line, "LoveTime:") {
			config.LoveTime = strings.Trim(strings.Split(line, ":")[1], "\n")
		}
	}

	defer file.Close()
}

// 随机生成16进制颜色
func RandomlyGenerateHexadecimalColor() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("0123456789ABCDEF")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return "#" + string(b)
}