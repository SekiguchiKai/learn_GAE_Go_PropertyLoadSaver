package model

import (
	"google.golang.org/appengine/datastore"
)

// 変更前の構造体
type User struct {
	Name string
	Address string
	Age int
}

// 返納後の構造体
//type Person struct {
//	Name string
//	Age int
//	From string
//}
//

func (u *User) Load(ps []datastore.Property) error {
	// LoadStruct loads the properties from p to dst.
	// dst must be a struct pointer.
	// LoadStructは、第二引数のdatastore.Propertyのslice(property)から
	// dst(構造体のポインタにロードする)
	err := datastore.LoadStruct(u, ps)
	if fmerr, ok := err.(*datastore.ErrFieldMismatch); ok && fmerr != nil && fmerr.Reason == "no such struct field" {
	} else if err != nil {
		return err
	}


	return nil
}

func (u *User) Save() ([]datastore.Property, error) {
	// 第一引数のPropertyをSliceにして、それを返す
	// 引数は、構造体のポインタでないといけない
	pr, err := datastore.SaveStruct(u)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
