package id
type AccountIDs string
func (a AccountIDs) String() string{
	return string(a)
}

type TripID string
func (t TripID) String() string{
	return string(t)
}

type IdentityID string
func (i IdentityID) String() string{
	return string(i)
}

type CarID string
func (c CarID) String() string{
	return string(c)
}


