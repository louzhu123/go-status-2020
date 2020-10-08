package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/louzhu123/gcrawl"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Data struct {
}

var (
	mysql_user = "root"
	mysql_pwd  = "root"
	mysql_db   = "gcrawl"
	conditions = map[string]interface{}{ //  https://github.com/louzhu123/gcrawl
		"position": "后端开发",
	}
)

func main() {
	var page = 1

	for {
		if page == 2000 {
			break
		}

		fmt.Println("正在爬取" + strconv.Itoa(page) + "页")

		// hrefs
		gcrawl := gcrawl.New51job()
		conditions["page"] = page
		res := gcrawl.Where(conditions).Get()

		// moreInfo
		items := gjson.Get(res, "engine_search_result").Array()
		result := []map[string]interface{}{}
		for _, item := range items {
			m, _ := gjson.Parse(item.String()).Value().(map[string]interface{})
			href := gjson.Get(item.String(), "job_href").String()
			moreInfo := GetMoreInfo(href)
			m["positionDescribe"] = moreInfo["positionDescribe"]
			m["keyword"] = moreInfo["keyword"]

			// deal with the []interface{} in the map result, otherwise it will be become () in sql by gorm result in a sql bug.
			for key, value := range m {
				if "[]interface {}" == reflect.TypeOf(value).String() {
					jsonStr, _ := json.Marshal(value)
					m[key] = string((jsonStr))
				}
			}
			result = append(result, m)
		}

		// storage the result
		dsn := mysql_user + ":" + mysql_pwd + "@tcp(127.0.0.1:3306)/" + mysql_db + "?charset=utf8mb4&parseTime=True&loc=Local"
		DB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		DB.Model(&Data{}).Create(result)

		page++
		time.Sleep(5 * time.Second)
	}
}

// 职位的详细信息
func GetMoreInfo(url string) map[string]string {
	var (
		result = map[string]string{}
	)

	res, _ := http.Get(url)
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	job_msg := doc.Find(".job_msg")
	positionDescribe := job_msg.Find("p").Text()
	keyword := job_msg.Find(".mt10").Find("a").Text()

	result["positionDescribe"] = Gbk2Utf8(positionDescribe)
	result["keyword"] = Gbk2Utf8(keyword)

	return result
}

// 51job是gbk编码，需要转成utf-8
func Gbk2Utf8(gbkStr string) string {
	enc := mahonia.NewDecoder("gbk")
	utf8Str := enc.ConvertString(gbkStr)
	return utf8Str
}
