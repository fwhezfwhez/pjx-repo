## 1. Pjx-repo
Here is my personal repo for pjx.There are some README.md files for each packaget.

By the way, this repo is an example of using [pjx](https://github.com/fwhezfwhez/pjx).To merge packages here into local pjx repo, you might need to `pjx clone https://github.com/fwhezfwhez/pjx-repo.git global -u`

## 2. Dep
https://github.com/fwhezfwhez/pjx

`go get -u github.com/fwhezfwhez/pjx`

## 3. Clone
`pjx clone https://github.com/fwhezfwhez/pjx-repo.git global -u` it will clone repo into %pjx_path%/global for pjx use.

## 4. Repo-component
Details of existing component.

#### 4.1. db
It provide gorm db package.Including db auto reconnection.It might need a tiny modify if you want to use your local config.More about gorm refer to https://github.com/jinzhu/gorm

`cd %your_project%/dep` cd where you want to put db
`pjx use db` add db package into your project

Then in your code, you can use it like:
```go
package main
import "dep/db"

func main() {
    db.DB.Exec("insert into ...")
}
```

#### 4.2 errorReport
It provide error report component.It will handle error depending on mode.More refer to https://github.com/fwhezfwhez/errorx

`cd %your_project%/dep` cd where you want to put db
`pjx use errorReport` add db package into your project

```go
package main
import "dep/errorReport"

func main() {
    errorReport.Er.SaveError(fmt.Errorf("nil return"), nil)
}
```

#### 4.3 jwt-util
It provide jwt validate.

```go
package main
import "dep/jwt-util"

func main() {
    	jwt_util.JwtTool.SetSecretKey("HELLO")
    	token,e:=JwtTool.GenerateJWT(map[string]interface{}{
    		"user_id": int(1),
    		"version": 1,
    		"exp":     time.Now().Add(2 * time.Hour).Unix(),
    	})
    	fmt.Println(token, e)

    	token, msg := JwtTool.ValidateJWT(token)

        if !token.Valid {
            fmt.Println("valid fail", msg)
            return
        }
        var r map[string]interface{}
        r = token.Claims.(jwt.MapClaims)

        fmt.Println(r)
        fmt.Println(r["user_id"])
}
```

#### 4.4 redistool
It provide redis client with pool.

```go
package main
import(
    "dep/redistool"
)
func main() {
    	pool := redistool.GetRedis("redis://localhost:6379")
    	c := pool.Get()
    	defer c.Close()

    	_, err := c.Do("MSET", "user_name", "kkk")
    	if err != nil {
    		fmt.Println("mset error :", err.Error())
    		return
    	}
    	username, err := redis.String(c.Do("GET", "user_name"))

    	if err != nil && err != redis.ErrNil{
    		fmt.Println("get error:", err.Error())
    		return
    	} else {
    		fmt.Println("get user_name", username)
    	}
}
```

#### 4.5 vx
It provides wechat api.This lib will be developing.Now it including:

- CheckSessionKey: Check session key valid or not
- MidasPay: Decrease diamond number
- GetBalance: Get diamond balance
- GenerateSigAndMpSig: Generate sig and mp_sig
- MidasPresent: present diamond
