package token

import (
	"coolcar/secret"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)
 

func TestGenerateToken(t *testing.T) {
	token, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret.PrivateKey))
	if err != nil {
		//测试会中止
		t.Fatalf("Failed to parse RSA private key: %v", err)
	}
	t.Errorf("%v", token)
	g:=NewJWTToken("coolcar/auth",token)
	g.nowFunc=func() time.Time{
		return time.Unix(1516239022,0)
	}
	tkn,err:=g.GenerateToken("623442b31628572a726be0b6",2*time.Hour)
	if err!=nil{
		//进行下去 测试不会终止
		t.Errorf("Failed to generate token: %v", err)
	}
	if tkn!=secret.Jwttoken{
		t.Errorf("Failed to generate token,\n want:\n%q\ngot:\n%q",secret.Jwttoken, tkn)
	}

}
