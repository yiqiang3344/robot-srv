package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/yiqiang3344/go-lib/helper"
	"log"
	helper2 "robot-srv/helper"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	robotSrv "github.com/yiqiang3344/go-lib/proto/robot-srv"
)

type RobotSrv struct{}
type RobotAppConfig struct {
	Id   int    `db:"id"`
	Type string `db:"type"`
	Cfg  string `db:"cfg"`
}
type RobotAppConfigCfg struct {
	SecretKey   string `json:"secret_key"`
	AccessToken string `json:"access_token"`
}
type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}
type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type Text struct {
	Content string `json:"content"`
}
type MarkdownData struct {
	Msgtype  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
	At       At       `json:"at"`
}
type TextData struct {
	Msgtype string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

type TestData struct {
	Test string `json:"test"`
}

var Db *sqlx.DB

// Call is a single request handler called via client.Call or the generated client code
func (e *RobotSrv) SendMsg(ctx context.Context, req *robotSrv.Request, rsp *robotSrv.Response) error {
	sp := helper.NewInnerSpan(helper.RunFuncName(), ctx)
	if sp != nil {
		defer sp.Finish()
	}

	go sendRobotMsg(ctx, req)

	rsp.Status = "1"
	rsp.Msg = "success"

	return nil
}

func sendRobotMsg(ctx context.Context, req *robotSrv.Request) {
	redisInstance := helper.DefaultRedis()
	defer redisInstance.Close()
	//频率限制
	_r, tips := helper2.CheckFreq("robotMsg", req.Title, 60)
	if _r == false {
		return
	}

	//redis读取机器人配置
	key := helper.GenRedisKey("robotAppConfig:" + req.BizType)
	r, err := redis.String(redisInstance.Do("get", key))
	if r == "" {
		//从mysql查询biz_type配置
		Db, err := helper.DefaultDB()
		if err != nil {
			return
		}

		var robotAppConfig []RobotAppConfig
		err = Db.Select(&robotAppConfig, "select id, type, cfg from robot_app_config where biz_type=? order by id asc limit 1", req.BizType)
		if err != nil {
			helper.ErrorLog("select robotAppConfig["+req.BizType+"] failed:"+err.Error(), helper.RunFuncName())
			return
		}

		if len(robotAppConfig) == 0 {
			helper.ErrorLog("配置不存在 robotAppConfig["+req.BizType+"]", helper.RunFuncName())
			return
		}

		_, err = redisInstance.Do("set", key, robotAppConfig[0].Cfg)
		if err != nil {
			helper.ErrorLog("redis set "+key+" failed:"+err.Error(), helper.RunFuncName())
			return
		}
		_, err = redisInstance.Do("expire", key, 3600) //缓存1小时
		if err != nil {
			helper.ErrorLog("redis expire "+key+" failed:"+err.Error(), helper.RunFuncName())
			return
		}

		r = robotAppConfig[0].Cfg
	}
	//json配置解析为结构体
	var cfg RobotAppConfigCfg
	err = json.Unmarshal([]byte(r), &cfg)
	if err != nil {
		helper.ErrorLog("json decode failed:"+r+" --- "+err.Error(), helper.RunFuncName())
		return
	}

	var data interface{}
	//根据消息类型拼装消息体
	if req.MsgType == "markdown" {
		data = MarkdownData{
			Msgtype: "markdown",
			Markdown: Markdown{
				Title: req.Title,
				Text:  tips + req.Content,
			},
			At: At{
				AtMobiles: req.AtMobiles,
				IsAtAll:   req.AtAll,
			},
		}
	} else {
		data = TextData{
			Msgtype: "text",
			Text: Text{
				Content: tips + req.Content,
			},
			At: At{
				AtMobiles: req.AtMobiles,
				IsAtAll:   req.AtAll,
			},
		}
	}

	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	signStr := fmt.Sprintf("%s\n%s", timestamp, cfg.SecretKey)
	sign := helper.ComputeHmacSha256(signStr, cfg.SecretKey)
	_url := "https://oapi.dingtalk.com/robot/send?access_token=" + cfg.AccessToken + "&timestamp=" + timestamp + "&sign=" + sign
	ret, statusCode, retStr := helper.PostJson(ctx, _url, data, 5*time.Second)
	if !ret {
		helper.BizLog("钉钉消息发送失败["+strconv.Itoa(statusCode)+"]:"+retStr, "")
		return
	}
	helper.BizLog("钉钉消息发送成功:"+retStr, "")
}

// Call is a single request handler called via client.Call or the generated client code
func (e *RobotSrv) Test(ctx context.Context, req *robotSrv.TestRequest, rsp *robotSrv.Response) error {
	sp := helper.NewInnerSpan(helper.RunFuncName(), ctx)
	if sp != nil {
		defer sp.Finish()
	}

	log.Println("start sleep")
	time.Sleep(3 * time.Second)
	log.Println("end sleep")
	//go test1(ctx)

	rsp.Status = "1"
	rsp.Msg = "success"

	return nil
}

func test(ctx context.Context) {
	data := TestData{
		Test: "test",
	}
	_url := "http://localhost:8080/robot/test"
	ret, statusCode, retStr := helper.PostJson(ctx, _url, data, 5*time.Second)
	if !ret {
		helper.DebugLog("测试消息发送失败["+strconv.Itoa(statusCode)+"]:"+retStr, "")
		return
	}
	helper.DebugLog("测试消息发送成功:"+retStr, "")
}

func test1(ctx context.Context) {
	data := TestData{
		Test: "test1",
	}
	_url := "http://localhost:8080/robot/test1"
	ret, statusCode, retStr := helper.PostJson(ctx, _url, data, 5*time.Second)
	if !ret {
		helper.DebugLog("测试消息发送失败["+strconv.Itoa(statusCode)+"]:"+retStr, "")
		return
	}
	helper.DebugLog("测试消息发送成功:"+retStr, "")
}
