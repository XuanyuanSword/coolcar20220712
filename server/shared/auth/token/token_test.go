package token

import (
	"coolcar/secret"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestVerify(t *testing.T){
	token, err := jwt.ParseRSAPublicKeyFromPEM([]byte(secret.PublicKey))
	if err != nil {
		//测试会中止
		t.Fatalf("Failed to parse RSA private key: %v", err)
	}
	v:=&JWTTokenVerifier{
		PublicKey:token,
	}
	cases :=[]struct {
		tkn string
		now time.Time
		want string
		wantErr bool
		name string
	}{
		{
			name: "有效的vaild_token",
			tkn:secret.Jwttoken,
			now:time.Unix(1516239222,0),
			want:"623442b31628572a726be0b6",
			
		},
		{
			name: "过期的token",
			tkn:secret.Jwttoken,
			now:time.Unix(1517239222,0),
			want:"623442b31628572a726be0b6",
			wantErr:false,
		},
		{
			name: "无效的token",
			tkn:"Jwttoken",
			now:time.Unix(1516239222,0),
			want:"623442b31628572a726be0b6",
			wantErr:true,
		},
		{
			name: "错误的签名",
			tkn:secret.Jwttoken,
			now:time.Unix(1516239222,0),
			want:"623442b31628572a726be0b6",
			wantErr:true,
		},
	}
	for _,c:=range cases{
		t.Run(c.name,func(t *testing.T){
			jwt.TimeFunc=func() time.Time{
				return c.now
			}
			got,err:=v.Verify(c.tkn)
			if err!=nil{
				if c.wantErr{
					return
				}
				t.Fatalf("\n验证token失败:\n%v", err)
			}
			if got!=c.want{
				t.Errorf("\n验证accountID失败:\nwant:\n %q\ngot:\n%q",c.want, got)
			}
		})
	}
}