# GAE/GoのDatastoreにおけるEntityとPropertiesとPropertyLoadSaverのお話

## GAE/GoのDatastoreにおける構造体とEntityとProperties
DatastoreのEntityの中身は、構造体のポインタが格納されることが多い。
しかし、PropertyLoadSaver interfaceを実装すればどんなTypeでもEntityの中身になり得る。
構造体のポインタであれば、Datastoreがreflection経由で自動的に変換してくれるので、PropertyLoadSaver interfaceを明示的に実装する必要はない。
ちなみにデフォルトでは、構造体のポインタは潜在的にIndex化されていて、propertyの名前は、構造体のフィールド名と同様である。

構造体のポインタがPropertyLoadSaverを明示的に実装していれば、構造体のポインタのデフォルトの振る舞いより優先的にPropertyLoadSaverのメソッドが使用される。
構造体のポインタは、より強く型付けされていて使用しやすいが、PropertyLoadSaverはより柔軟に使用することができる。

実際のTypeは、 GetとPutで同一のものでなくても良いし、別々のApp Engineのリクエストを横断してもよい。
DatastoreにEntityが格納される際には、Propertyのシーケンス([]Property)として格納される。



[The datastore package  |  App Engine standard environment for Go  |  Google Cloud Platform]
(https://cloud.google.com/appengine/docs/standard/go/datastore/reference#Properties) を非常に参考にさせていただいた。

### EntityとPropertiesの差異の扱い方について
フィールドが欠けていたりする不完全なEntityは、結果として `ErrFieldMismatch` になるけども、
このエラーを致命的(fatal)にするか、回復可能として扱うか、無視するかは呼び出し側による。
これについては、以下の記事がわかりやすい
[DatastoreからGetした時に余計なPropertyがある場合エラーになるが無視してもいい - Qiita]
(https://qiita.com/vvakame/items/e017e7d955f82ddd8af1)

[The datastore package  |  App Engine standard environment for Go  |  Google Cloud Platform]
(https://cloud.google.com/appengine/docs/standard/go/datastore/reference#Properties) を非常に参考にさせていただいた。

## The PropertyLoadSaver Interface
Entityの中身はPropertyLoadSaver interfaceを実装した、いかなる型にも成り得る。
これは構造体のポインタでもそうでなくても良い。
PropertyLoadSaverを実装していれば、Datastore packageはEntityの中身をGetする際には、Loadメソッドを呼びだす。Entityの中身をPutする際にはSaveを呼び出す。

不要になったプロパティの削除やバージョンの切り戻し（バージョンアップでプロパティが増えたがやっぱり切り戻ししたくなった場合）の際に利用すると有効そうだ。

[The datastore package  |  App Engine standard environment for Go  |  Google Cloud Platform]
(https://cloud.google.com/appengine/docs/standard/go/datastore/reference#hdr-The_PropertyLoadSaver_Interface)を非常に参考にさせていただいた。


## 各Typeや関数の詳細

各type、メソッド、関数のシグネチャは以下の公式より引用させていただき、また、説明文は以下の公式の意訳及び、それを参考にさせていただいたものになっている。
[The datastore package  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/datastore/reference)

### type Property

```go
type Property struct {
    Name string
    Value interface{}
    NoIndex bool
    Multiple bool

```

Propertyは、名前と値のペアとその他のメタデータを含む構造体である。
Datastoreのエンティティの中身は、Propertyのシーケンスとしてロードされ、保存される。

### type PropertyList

```go
type PropertyList []Property
```
[]Property(PropertyのSlice)を、PropertyLoadSaverの実装へと変換する。

#### func (*PropertyList) Load

```go
func (l *PropertyList) Load(p []Property) error
```
[]Property(PropertyのSlice)を*PropertyListにロードする。
lを最初に空にするようなことはしない。
Datastore packageがEntityの中身をGetする際には、Loadメソッドを呼びだす。



#### func (*PropertyList) Save
```go
func (l *PropertyList) Save() ([]Property, error)

```
lの全てのPropertyをsliceかPropertiesとして、保存する。
Datastore packageがEntityの中身をPutする際にはSaveを呼び出す。

### type PropertyLoadSaver

```go
type PropertyLoadSaver interface {
    Load([]Property) error
    Save() ([]Property, error)
}
```
PropertyLoadSaverは、[]Property(PropertyのSlice)にもなれるし、[]Property(PropertyのSlice)から変換もできる。

### func LoadStruct

```go
func LoadStruct(dst interface{}, p []Property) error
```
第2引数の[]Property(PropertyのSlice)から第一引数の構造体のポインタにロードする。


### func SaveStruct

```go
func SaveStruct(src interface{}) ([]Property, error)
```
引数の構造体のポインタから抽出した[]Property(PropertyのSlice)を返す。


### 実装
```go
type User struct {
	ID string
	Name    string
	Address string
	Age     int
	UpdatedAt time.Time `json:"updatedAt"`
}

//  Datastore packageがEntityの中身をGetする際には、Loadメソッドを呼びだす。
// []Property(PropertyのSlice)を*PropertyListにロードする。 lを最初に空にするようなことはしない。 
func (u *User) Load(ps []datastore.Property) error {
	// LoadStructは、第二引数のdatastore.Propertyのslice(property)から
	// dst(構造体のポインタ)にロードする
	err := datastore.LoadStruct(u, ps)
	if fmerr, ok := err.(*datastore.ErrFieldMismatch); ok && fmerr != nil && fmerr.Reason == "no such struct field" {
	} else if err != nil {
		return err
	}

	return nil
}


// Datastore packageがEntityの中身をPutする際にはSaveを呼び出す。 
// 全てのPropertyをsliceかPropertiesとして、保存する 
func (u *User) Save() ([]datastore.Property, error) {
	// 第一引数のPropertyをSliceにして、それを返す
	// 引数は、構造体のポインタでないといけない
	pr, err := datastore.SaveStruct(u)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
```
[DatastoreからGetした時に余計なPropertyがある場合エラーになるが無視してもいい - Qiita]
(https://qiita.com/vvakame/items/e017e7d955f82ddd8af1)を参考に実装した


## まとめ
* DatastoreにEntityが格納される際には、Propertyのシーケンス([]Property)として格納される
* Propertyは、名前と値のペアと他のメタデータを持つ構造体である
* 構造体のポインタであればDatastoreがreflection経由で自動的に変換してくれるが、PropertyLoadSaver interfaceを明示的に実装すれば、そうじゃない値も可能だし、 GetとPutで同一のものでなくても良くなる
* PropertyLoadSaverは、[]Property(PropertyのSlice)にもなれるし、[]Property(PropertyのSlice)から変換もできる
*  type PropertyListは、[]Property(PropertyのSlice)を、PropertyLoadSaverの実装へと変換する
*  PropertyLoadSaverを実装していれば、Datastore packageはEntityの中身をGetする際には、Loadメソッドを呼びだす。Entityの中身をPutする際にはSaveを呼び出す。

## 参考にさせていただいた記事
[The datastore package  |  App Engine standard environment for Go  |  Google Cloud Platform](https://cloud.google.com/appengine/docs/standard/go/datastore/reference)

[DatastoreからGetした時に余計なPropertyがある場合エラーになるが無視してもいい - Qiita]
(https://qiita.com/vvakame/items/e017e7d955f82ddd8af1)

