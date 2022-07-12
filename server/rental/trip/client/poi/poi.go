package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"github.com/golang/protobuf/proto"
	"hash/fnv"
)

var poi=[]string{
	"中关村",
	"武当山",
	"陆家嘴",
	"迪士尼",
	"天河体育中心",
	"广州塔",
}

type Manager struct{

}

func (*Manager) Resolve(c context.Context,loc *rentalpb.Location)(string,error){
	b,err:=proto.Marshal(loc)
	if err!=nil{
		return "",err
	}
	h:=fnv.New32()
	h.Write(b)

    return poi[int(h.Sum32())%len(poi)],nil
}