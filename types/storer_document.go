package types

//go:generate msgp

//固态存储中存储的struct
type StorerDocument struct{
	DocId uint64 `msg:"id"`
	Latitude float64 `msg:"lat"`
	Longitude float64 `msg:"long"`
	Fields interface{} `msg:"f"`
}

//marshal
func(z *StorerDocument) MarshalBinary() (data []byte, err error){
	return z.MarshalMsg(nil)
}

//unmarshal
func(z *StorerDocument) UnmarshalBinary(data []byte) error{
	_,err := z.UnmarshalMsg(data)
	return err
}
