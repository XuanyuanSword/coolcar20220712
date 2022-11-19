package id
type AccountID string
func (a AccountID) String() string{
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

// BlobID defines blob id object.
type BlobID string

func (i BlobID) String() string {
	return string(i)
}
