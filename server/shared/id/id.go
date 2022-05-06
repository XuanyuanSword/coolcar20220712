package id
type AccountIDs string
func (a AccountIDs) String() string{
	return string(a)
}

type TripID string
func (t TripID) String() string{
	return string(t)
}