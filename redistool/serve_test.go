package redistool

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"testing"
)

func TestGetRedis(t *testing.T) {
	pool := GetRedis("redis://localhost:6379/1")
	c := pool.Get()
	defer c.Close()
	log.SetFlags(log.Llongfile | log.LstdFlags)
	_, err := c.Do("SADD", "users", "hello")
	if err != nil {
		log.Println( err)
		return
	}
	_, err = c.Do("SADD", "users", "2")
	if err != nil {
		log.Println( err)
		return
	}
	bufs, err := redis.ByteSlices(c.Do("sunion", "users"))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(bufs[0]))
	//fmt.Println(string(buf.([]interface{})[0].([]uint8)[2]))
	_, err = c.Do("MSET", "user_name", "kkk")
	if err != nil {
		fmt.Println("数据设置失败:", err)
	}
	username, err := redis.String(c.Do("GET", "user_name"))

	if err != nil {
		fmt.Println("数据获取失败:", err)
		fmt.Println(err.Error() == "redigo: nil returned")
	} else {
		fmt.Println("2.1.获取user_name", username)
	}

	type Man struct {
		Name string `json:"name"`
	}

	rp, err := c.Do("LRANGE", "test", 0, 10)
	if err != nil {
		fmt.Println("LRANGE失败:" + err.Error())
		return
	}
	res := rp.([]interface{})
	fmt.Println(len(res))
	var man Man
	for i, v := range res {
		man = Man{}
		json.Unmarshal(v.([]byte), &man)
		fmt.Println(fmt.Sprintf("%d:%s", i, man))
	}
}
