package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)
type  JWTTokenVerifier struct{
	PublicKey *rsa.PublicKey
}
func (v *JWTTokenVerifier)Verify(tokenString string) (string,error){
	t,err:=jwt.ParseWithClaims(tokenString,&jwt.StandardClaims{

	},func(*jwt.Token)(interface{},error){
		return v.PublicKey,nil
	})
	if err!=nil{
		return "",fmt.Errorf("failed to parse token: %v",err)
	}
	if !t.Valid{
		return "",fmt.Errorf("invalid token")
	}
	clm,ok:=t.Claims.(*jwt.StandardClaims)
	if !ok{
		return "",fmt.Errorf("invalid token")
	}
	if err:=clm.Valid();err!=nil{
		return "",fmt.Errorf("claim not vaild: %v",err)
	}
	return clm.Subject,nil
}