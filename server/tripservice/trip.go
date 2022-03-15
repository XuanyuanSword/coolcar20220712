package trip

import (
	"context"
	trippb "coolcar/proto/gen/go"
)
type Service struct{

}
func (s *Service) GetTrip(ctx context.Context, req *trippb.GetTripRequest) (*trippb.GetTripResponse,error) {
	return &trippb.GetTripResponse{
		Id:req.Id,
		Trip: &trippb.Trip{
			Start:"abc",
			End:"def",
			DurationSec:3600,
			FeeCent:10000,
			StartPos: &trippb.Location{
				Latitude:1.1,
				Longitude:2.2,
			},
			EndPos: &trippb.Location{
				Latitude:3.3,
				Longitude:4.4,
			},
			Status: trippb.TripStatus_IN_PROGRESS,
		},
		 
	},nil
}