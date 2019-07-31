package jwt_util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	JwtTool.SetSecretKey("HELLO")
	fmt.Println(JwtTool.GenerateJWT(map[string]interface{}{
		"user_id": int(1),
		"version": 1,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}))
}

func TestValidateJWT(t *testing.T) {
	JwtTool.SetSecretKey("HELLO")
	token, msg := JwtTool.ValidateJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0MTcsInVzZXJfbmFtZSI6IuW8oOWuh-mUiyIsImdyb3VwX2lkIjoyLCJpc19zdXBlcnVzZXIiOmZhbHNlLCJleHAiOjE1NTY1OTYxOTV9.imW2r7aReVJgUsfGPTYA4jtK9imMwBWucjABRHdjq90")

	if !token.Valid {
		fmt.Println("valid fail", msg)
		return
	}
	fmt.Println(token.Claims.(jwt.MapClaims)["user_id"])
	fmt.Println(reflect.TypeOf(token.Claims.(jwt.MapClaims)["user_id"]))

	var r map[string]interface{}
	r = token.Claims.(jwt.MapClaims)

	fmt.Println(r)
	fmt.Println(r["user_id"])
	fmt.Println(reflect.TypeOf(r["user_id"]))
}
