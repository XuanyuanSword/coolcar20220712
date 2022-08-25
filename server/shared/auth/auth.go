package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"coolcar/shared/id"
	"strings"

	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)
func Interceptor(publicKeyFile string)(grpc.UnaryServerInterceptor,error){
	
	
	f,err:=os.Open(publicKeyFile)
	if err!=nil{
		return nil,fmt.Errorf("cannot open public key file:%v",err)
	}
	b,err:=ioutil.ReadAll(f)
	
	if err!=nil{
		return nil,fmt.Errorf("cannot read public key file:%v",err)
	}
	publicKey,err:=jwt.ParseRSAPublicKeyFromPEM(b)
	if err!=nil{
		return nil,fmt.Errorf("cannot read public key file:%v",err)
	}

	i:=&interceptor{verifier:&token.JWTTokenVerifier{
		PublicKey:publicKey,
	}}
	return i.HandlerReq,nil
}
type tokenVerify interface{
	Verify(token string) (string,error)
}

type  interceptor struct {
	//  publicKey *rsa.PublicKey
	 verifier tokenVerify
}

func (i *interceptor) HandlerReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
	 tkn,err:=tokenFromText(ctx)
	 if err!=nil{
		 return nil,status.Error(codes.Unauthenticated,"")
	 }
	 aid,err:=i.verifier.Verify(tkn)
	 
	 if err!=nil{
		 return nil,status.Error(codes.Unauthenticated,"token not vaild")
	 }
	return handler(ContextWithAccountId(ctx,id.AccountIDs(aid)), req)
}
const (
	
	authorization="authorization"
	tknPrefix="token "
)

func tokenFromText(ctx context.Context) (string,error){
	
        m,ok:=metadata.FromIncomingContext(ctx)
		if !ok{
			return "",fmt.Errorf("no metadata")
		}
		tkn:=""
		for _,v:=range m[authorization]{
			if strings.HasPrefix(v,tknPrefix){

				tkn=v[len(tknPrefix):]
				
			}

			
		}
		if tkn==""{
			return "",fmt.Errorf("no token")
		}
		
		return tkn,nil
}
type accountIDKEY struct{}

func ContextWithAccountId(c context.Context,aid id.AccountIDs) context.Context{
	
	return  context.WithValue(c,accountIDKEY{},aid)
}

func AccountID(c context.Context)(id.AccountIDs,error){
	v := c.Value(accountIDKEY{})
	aid, ok := v.(id.AccountIDs)


	if !ok{

		return "", status.Error(codes.Unauthenticated, "")
	}
	return aid,nil
}