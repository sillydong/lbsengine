package types

//go:generate msgp

//固态存储中存储的struct
type IndexedDocument struct{
	DocId uint64 `msg:"id"`
	Latitude float64 `msg:"lat"`
	Longitude float64 `msg:"long"`
	Fields interface{} `msg:"f"`
}

//marshal
func(z *IndexedDocument) MarshalBinary() (data []byte, err error){
	return z.MarshalMsg(nil)
}

//unmarshal
func(z *IndexedDocument) UnmarshalBinary(data []byte) error{
	_,err := z.UnmarshalMsg(data)
	return err
}
