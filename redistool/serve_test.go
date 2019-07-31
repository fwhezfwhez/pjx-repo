package redistool

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestGetRedis(t *testing.T) {
	pool := GetRedis("redis://localhost:6379")
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("MSET", "user_name", "kkk")
	if err != nil {
		fmt.Println("数据设置失败:", err)
	}
	username, err := redis.String(c.Do("GET", "user_nam33e"))

	if err != nil {
		fmt.Println("数据获取失败:", err)
		fmt.Println(err.Error() == "redigo: nil returned")
	} else {
		fmt.Println("2.1.获取user_name", username)
	}

	type Man struct {
		Name string `json:"name"`
	}
	//v, _ := json.Marshal(Man{Name: "ft"})
	//_,err =c.Do("LPUSH", "test", v)
	//if err !=nil {
	//	fmt.Println("LPUSH失败:"+ err.Error())
	//	return
	//}
	//_,err =c.Do("LPUSH", "test", v)
	//if err !=nil {
	//	fmt.Println("LPUSH失败:"+ err.Error())
	//	return
	//}

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
