package profile

import (
	"context"
	"coolcar/shared/id"
)

type Manager struct{

}

func (p *Manager)Verify(context.Context,id.AccountIDs)(id.IdentityID,error){
	return id.IdentityID("identity1"), nil
}