package freq

import (
	"crypto/md5"
	"fmt"
	"github.com/garyburd/redigo/redis"
	cLog "github.com/yiqiang3344/go-lib/utils/log"
	cRedis "github.com/yiqiang3344/go-lib/utils/redis"
	"strconv"
	"strings"
	"time"
)

func CheckFreq(_type string, title string, limitFreq int) (bool, string) {
	redisInstance := cRedis.DefaultRedis()
	defer redisInstance.Close()
	freqKey := cRedis.GenRedisKey("frequency:" + _type + ":" + fmt.Sprintf("%x", md5.Sum([]byte(title))))
	historyKey := cRedis.GenRedisKey("history:" + _type + ":" + fmt.Sprintf("%x", md5.Sum([]byte(title))))
	r, err := redis.Bool(redisInstance.Do("setnx", freqKey, 1))
	if err != nil {
		cLog.ErrorLog("setnx "+freqKey+" error:"+err.Error(), "")
		return false, ""
	}
	if r == false {
		_, _ = redisInstance.Do("lpush", historyKey, fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second()))
		return false, ""
	}
	_, _ = redisInstance.Do("expire", freqKey, limitFreq)
	histores, err := redis.Strings(redisInstance.Do("lrange", historyKey, 0, -1))
	tips := ""
	if len(histores) > 0 {
		tips = histores[0] + "到" + histores[len(histores)-1] + "之间共有" + strconv.Itoa(len(histores)) + "次通知，时间表如下：\n" + strings.Join(histores, "\n") + "\n\n"
		_, _ = redisInstance.Do("del", historyKey)
	}
	return true, tips
}
