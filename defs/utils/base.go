package utils

import (
	"fmt"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"wxSendMessage/defs/config"
)

// 获取access token
func GetAccessToken() string {
	postUrl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + config.AppId + "&secret=" + config.AppSecret
	response, err := http.Get(postUrl)
	if err != nil {
		fmt.Println("ERROR：获取access_toke失败，请检查app_id和app_secret是否正确。 ")
		os.Exit(1)
	}
	body, _ := ioutil.ReadAll(response.Body)
	accessToken := gjson.Get(string(body), "access_token")
	defer response.Body.Close()
	return accessToken.String()
}

// 获取生日倒计时
// TODO 暂时不投入正式使用
func GetBirthday() int {
	var lunarDay string
	var lunarMonth string
	var nowBirth string
	// 将当前日期转换为农历
	c := calendar.BySolar(int64(time.Now().Year()), int64(time.Now().Month()), int64(time.Now().Day()), 0, 0, 0)
	bytes, _ := c.ToJSON()

	// 获取日月
	lunarDay = gjson.Get(string(bytes), "lunar.day").String()
	lunarMonth = gjson.Get(string(bytes), "lunar.month").String()

	intDay, _ := strconv.Atoi(lunarDay)
	intMonth, _ := strconv.Atoi(lunarMonth)

	// 补0操作
	if intDay < 10 {
		lunarDay = "0" + lunarDay
	}
	if intMonth < 10 {
		lunarMonth = "0" + lunarMonth
	}

	nowTime := time.Now().Format("2006") + "-" + lunarMonth + "-" + lunarDay + " 00:00:00"

	//TODO 后面改为从配置中获取
	nowBirth = time.Now().Format("2006") + "-" + "07" + "-" + "27" + " 00:00:00"

	num := GetTimeSubDay(nowTime, nowBirth)
	if num > 0 {
		return num
	} else if num < 0 {
		t := time.Now()
		addData := t.AddDate(1, 0, 0).Year()
		now2Birth := strconv.Itoa(addData) + "-" + "07" + "-" + "25" + " 00:00:00"
		num2 := GetTimeSubDay(nowTime, now2Birth)
		return num2
	} else {
		return num
	}
}

// 获取天气
func GetWeather() (string, string, string, string, string, string, string) {
	postValue := url.Values{"key": {config.TianApi}, "city": {config.City}}
	res, err := http.PostForm("http://api.tianapi.com/tianqi/index", postValue)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(res.Body)

	time.Sleep(3 * time.Second)

	area := gjson.Get(string(body), "newslist.0.area").String()
	date := gjson.Get(string(body), "newslist.0.date").String()
	week := gjson.Get(string(body), "newslist.0.week").String()
	weather := gjson.Get(string(body), "newslist.0.weather").String()
	real := gjson.Get(string(body), "newslist.0.real").String()
	lowest := gjson.Get(string(body), "newslist.0.lowest").String()
	highest := gjson.Get(string(body), "newslist.0.highest").String()
	return area, date, week, weather, real, lowest, highest
}

// 彩虹屁
func Caihongpi() string {
	if config.TianApi != "" {
		postValue := url.Values{"key": {config.TianApi}}
		res, _ := http.PostForm("http://api.tianapi.com/caihongpi/index", postValue)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		content := gjson.Get(string(body), "newslist.0.content").String()
		return content
	}
	return ""
}

// 古代情诗
func GetLovePoems() string {
	postValue := url.Values{"key": {config.TianApi}}
	res, _ := http.PostForm("http://api.tianapi.com/qingshi/index", postValue)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	content := gjson.Get(string(body), "newslist.0.content").String()
	return content
}

func LoveDay() string {
	LoveTime := config.LoveTime + " 00:00:00"
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	return strconv.Itoa(GetTimeSubDay(LoveTime, nowTime))
}

/** template
{{date.DATA}}
城市：{{city.DATA}}
天气：{{weather.DATA}}
当前气温:{{real.DATA}}
最低气温: {{lowest.DATA}}
最高气温: {{highest.DATA}}
今天是我们恋爱的第{{love_day.DATA}}天
爱你哟！
{{lovePoems.DATA}}
*/
func SendMessage() {
	var area string
	var date string
	var week string
	var weather string
	var real string
	var lowest string
	var highest string
	urls := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + GetAccessToken()

	area, date, week, weather, real, lowest, highest = GetWeather()
	var str = `{
        "touser": "` + config.ToUserId + `",
        "template_id": "` + config.TemplateId + `",
        "url": "http://weixin.qq.com/download",
        "topcolor": "#FF0000",
        "data": {
            "date": {
                "value": "` + date + " " + week + `",
				"color": "` + RandomlyGenerateHexadecimalColor() + `"
			},
			"city": {
                "value": "` + area + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            },
            "weather": {
                "value": "` + weather + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            },
            "real": {
                "value": "` + real + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            },
            "lowest": {
                "value": "` + lowest + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            },
            "highest": {
                "value": "` + highest + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            },
            "love_day": {
                "value": "` + LoveDay() + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            },
            "lovePoems": {
                "value": "` + GetLovePoems() + `",
                "color": "` + RandomlyGenerateHexadecimalColor() + `"
            }
        }
    }`
	payload := strings.NewReader(str)
	res, _ := http.Post(urls, "application/json", payload)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	success := gjson.Get(string(body), "errcode").Int()
	if success == 0 {
		fmt.Println("推送成功")
	}
}
