![GitHub release (latest by date)](https://img.shields.io/github/v/release/mauserzjeh/null?style=flat-square)

# Null

Null is a package that provides a generic nullable variable that can be used for both `JSON` and `SQL` operations.

## Features
- Uses a generic type for the underlying value
- Provides an easy way to test if certain struct fields were set after unmarshalling a `JSON` message
- Makes it easier to update only those fields in the database which were set previously. Especially in `NoSQL` databases using `JSON` documents
- Compatible with the standard `database/sql` package
    - If the underlying type is a custom type that the standard `database/sql` package can't handle, then it is possible to implement the `sql.Scanner` and `driver.Valuer` interfaces for that given type to be compatible
- Possible to control how the underlying value is handled during json marshalling and unmarshalling by implementing the `json.Marshaler` and `json.Unmarshaler` interfaces for the underlying type
- Helper functions provide methods to easily filter unset fields in structs and maps

## Installation
```
go get -u github.com/mauserzjeh/null
```

## Usage & Examples
Below you can see how to use the package in general and also in a more complex scenario. There are a few examples in `test_null.go` as well.

## General usage
```go
var nullableStr null.Var[string]

// IsSet returns if the value was set
nullableStr.IsSet() // false

// Valid returns if the value is NULL
nullableStr.Valid() // false

// Val returns the value
nullableStr.Val() // ""

// ------------------------------------------------------------------

nullableStr.Set("foo")
nullableStr.IsSet() // true
nullableStr.Valid() // true
nullableStr.Val() // "foo"

// ------------------------------------------------------------------

nullableStr.SetNil()
nullableStr.IsSet() // true
nullableStr.Valid() // false
nullableStr.Val() // ""

// ------------------------------------------------------------------

nullableStr.Unset()
nullableStr.IsSet() // false
nullableStr.Valid() // false
nullableStr.Val() // ""
```

## Complex usage

### 1. Let's have the following types
```go
// custom integer type
type SiblingType int64

const (
    Sister SiblingType = iota
    Brother
)

// custom slice type
type CustomSlice []string

// custom generic map
type CustomMap[K, V comparable] map[K]V

// custom struct
type Sibling struct {
    Name null.Var[string]      `json:"name"`
    Age  null.Var[int64]       `json:"age"`
    Type null.Var[SiblingType] `json:"type"`
}

// another custom struct
type Person struct {
    Name       null.Var[string]                   `json:"name"`
    Age        null.Var[int64]                    `json:"age"`
    BirthDate  null.Var[time.Time]                `json:"birth_date"`
    Books      null.Var[CustomSlice]              `json:"books"`
    ExamScores null.Var[CustomMap[string, int64]] `json:"exam_scores"`
    Sibling    null.Var[Sibling]                  `json:"sibling"`
}
```

### 2. Default `JSON` marshaling
```go
var p Person
log.Printf("%+v", p)
// {
//     Name:       {set:false valid:false value:} 
//     Age:        {set:false valid:false value:0} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:false valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }

j, _ := json.Marshal(p)
log.Printf("%s", j)
// {
// 	"name":null,
// 	"age":null,
// 	"birth_date":null,
// 	"books":null,
// 	"exam_scores":null,
// 	"sibling":null
// }

p.Name.Set("Peter")
p.Age.Set(25)

var s Sibling
log.Printf("%+v", s)
// {
//     Name:   {set:false valid:false value:} 
//     Age:    {set:false valid:false value:0} 
//     Type:   {set:false valid:false value:0}
// }

s.Name.Set("Anna")
s.Age.Set(20)
s.Type.Set(Sister)

p.Sibling.Set(s)

j, _ = json.Marshal(p)
log.Printf("%s", j)
// {
// 	"name":"Peter",
// 	"age":25,
// 	"birth_date":null,
// 	"books":null,
// 	"exam_scores":null,
// 	"sibling":{
// 		"name":"Anna",
// 		"age":20,
// 		"type":0
// 	}
// }

p.BirthDate.Set(time.Unix(852073200, 0))
p.Books.Set([]string{"George Orwell - 1984", "Stephen E. Ambrose - Band of Brothers"})
p.ExamScores.Set(CustomMap[string, int64]{
    "math":    80,
    "physics": 90,
})

j, _ = json.Marshal(p)
log.Printf("%s", j)
// {
// 	"name":"Peter",
// 	"age":25,
// 	"birth_date":"1997-01-01T00:00:00+01:00",
// 	"books":[
// 		"George Orwell - 1984",
// 		"Stephen E. Ambrose - Band of Brothers"
// 	],
// 	"exam_scores":{
// 		"math":80,
// 		"physics":90
// 	},
// 	"sibling":{
// 		"name":"Anna",
// 		"age":20,
// 		"type":0
// 	}
// }
```

### 3. Implement `JSON` marshaler for `Person` and `Sibling`
```go
// MarshalJSON implements the json.Marshaler interface
func (s Sibling) MarshalJSON() ([]byte, error) {
    m := map[string]any{}
    if s.Name.IsSet() {
        m["name"] = s.Name
    }

    if s.Age.IsSet() {
        m["age"] = s.Age
    }

    if s.Type.IsSet() {
        m["type"] = s.Type
    }

    return json.Marshal(m)
}

// MarshalJSON implements the json.Marshaler interface
func (p Person) MarshalJSON() ([]byte, error) {
    m := map[string]any{}
    if p.Name.IsSet() {
        m["name"] = p.Name
    }

    if p.Age.IsSet() {
        m["age"] = p.Age
    }

    if p.BirthDate.IsSet() {
        m["birth_date"] = p.BirthDate
    }

    if p.Books.IsSet() {
        m["books"] = p.Books
    }

    if p.ExamScores.IsSet() {
        m["exam_scores"] = p.ExamScores
    }

    if p.Sibling.IsSet() {
        m["sibling"] = p.Sibling
    }

    return json.Marshal(m)
}

var p Person
log.Printf("%+v", p)
// {
//     Name:       {set:false valid:false value:} 
//     Age:        {set:false valid:false value:0} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:false valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }

j, _ := json.Marshal(p)
log.Printf("%s", j)
// {}

p.Name.Set("Peter")
p.Age.Set(25)

var s Sibling
log.Printf("%+v", s)
// {
//     Name:   {set:false valid:false value:} 
//     Age:    {set:false valid:false value:0} 
//     Type:   {set:false valid:false value:0}
// }

s.Name.Set("Anna")
s.Age.Set(20)
s.Type.Set(Sister)

p.Sibling.Set(s)

j, _ = json.Marshal(p)
log.Printf("%s", j)
// {
//     "age":25,
//     "name":"Peter",
//     "sibling":{
//         "age":20,
//         "name":"Anna",
//         "type":0
//     }
// }

p.BirthDate.Set(time.Unix(852073200, 0))
p.Books.Set([]string{"George Orwell - 1984", "Stephen E. Ambrose - Band of Brothers"})
p.ExamScores.Set(CustomMap[string, int64]{
    "math":    80,
    "physics": 90,
})

j, _ = json.Marshal(p)
log.Printf("%s", j)
// {
//     "age":25,
//     "birth_date":"1997-01-01T00:00:00+01:00",
//     "books":[
//         "George Orwell - 1984",
//         "Stephen E. Ambrose - Band of Brothers"
//     ],
//     "exam_scores":{
//         "math":80,
//         "physics":90
//     },
//     "name":"Peter",
//     "sibling":{
//         "age":20,
//         "name":"Anna",
//         "type":0
//     }
// }

p.Books.SetNil()     // will be null
p.ExamScores.Unset() // will be unset
p.Sibling.Unset()    // will be unset

j, _ = json.Marshal(p)
log.Printf("%s", j)
// {
//     "age":25,
//     "birth_date":"1997-01-01T00:00:00+01:00",
//     "books":null,
//     "name":"Peter"
// }
```

### 4. Use `FilterStruct` and `FilterMap`

`FilterStruct` and `FilterMap` are helper functions that can filter either a struct or a map from unset nullable variables. They provide an easy way to implement `json.Marshaller` interface without having to check each field in a struct. These helper functions both return a map without the unset fields.

However there are a few requirements:
- All fields should be tagged with either a `json` or custom tag
- To be able to recursively filter custom structs the `Filterable` interface should be embedded into the given structs
- Custom structs can be embedded without having them to be tagged as long as they have `Filterable` interface embedded inside them and their fields are tagged. In this case the fields of the embedded struct will be on the same level as the struct that embeds it. If a tag is set for this embedded field, then the embedded struct fields will be presented under the given tag.

Similiarly to `FilterStruct`, `FilterMap` can be used to filter maps containing nullable variables. 
It will also recursively filter map keys which are `map[string]any` type.

The `JSON` marshaler example looks like this using the `FilterStruct` helper function. More examples in `filter_test.go`
    
```go
// custom struct
type Sibling struct {
    // Signal for FilterStruct that if this struct type is used
    // for another struct's field type, then it can recursively
    // filter it.
    null.Filterable

    Name null.Var[string]      `json:"name"`
    Age  null.Var[int64]       `json:"age"`
    Type null.Var[SiblingType] `json:"type"`
}

// another custom struct
type Person struct {
    // Signal for FilterStruct that if this struct type is used
    // for another struct's field type, then it can recursively
    // filter it.
    null.Filterable 

    Name       null.Var[string]                   `json:"name"`
    Age        null.Var[int64]                    `json:"age"`
    BirthDate  null.Var[time.Time]                `json:"birth_date"`
    Books      null.Var[CustomSlice]              `json:"books"`
    ExamScores null.Var[CustomMap[string, int64]] `json:"exam_scores"`

    // Sibling is also filterable
    Sibling    Sibling                            `json:"sibling"`
}


// MarshalJSON implements the json.Marshaler interface
func (s Sibling) MarshalJSON() ([]byte, error) {
    m, err := null.FilterStruct(s)
    if err != nil {
        return nil, err
    }

    return json.Marshal(m)
}

// MarshalJSON implements the json.Marshaler interface
func (p Person) MarshalJSON() ([]byte, error) {
    m, err := null.FilterStruct(p)
    if err != nil {
        return nil, err
    }

    return json.Marshal(m)
}
```

`FilterStruct` has an additional option to use custom tags when determining the keys for the filtered map that it will produce. 
By default the `json` tag is used.
```go
m, err := null.FilterStruct(s, null.UseTag("custom_tag"))
```

An example using `FilterMap`.
```go
a := null.Var[string]
b := null.Var[string]
c := null.Var[string]

a.Set("A")
b.SetNil()

m := map[string]any{
    "a": a,
    "b": b,
    "c": c,
}

mf, err := null.FilterMap(m)
if err != nil {
    log.Fatal(err)
}

log.Printf("%+v", mf)
// map[a:A b:<nil>]

```

### 5. Default `JSON` unmarshal
```go
var p Person
log.Printf("%+v", p)
// {
//     Name:       {set:false valid:false value:} 
//     Age:        {set:false valid:false value:0} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:false valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }

jsonStr1 := []byte(`{}`)
_ = json.Unmarshal(jsonStr1, &p)
log.Printf("%+v", p)
// {
//     Name:       {set:false valid:false value:} 
//     Age:        {set:false valid:false value:0} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:false valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }

p = Person{} // reset variale
jsonStr2 := []byte(`{"age":25,"name":"Peter","sibling":{"age":20,"name":"Anna","type":0}}`)
_ = json.Unmarshal(jsonStr2, &p)
log.Printf("%+v", p)
// {
//     Name:       {set:true valid:true value:Peter} 
//     Age:        {set:true valid:true value:25} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:true valid:true value:{
//             Name:   {set:true valid:true value:Anna} 
//             Age:    {set:true valid:true value:20} 
//             Type:   {set:true valid:true value:0}
//         }
//     }
// }


p = Person{} // reset variale
jsonStr3 := []byte(`{"age":25,"birth_date":"1997-01-01T00:00:00+01:00","books":["George Orwell - 1984","Stephen E. Ambrose - Band of Brothers"],"exam_scores":{"math":80,"physics":90},"name":"Peter","sibling":{"age":20,"name":"Anna","type":0}}`)
_ = json.Unmarshal(jsonStr3, &p)
log.Printf("%+v", p)
// {
//     Name:       {set:true valid:true value:Peter} 
//     Age:        {set:true valid:true value:25} 
//     BirthDate:  {set:true valid:true value:{wall:0 ext:62987670000 loc:0x2ca260}} 
//     Books:      {set:true valid:true value:[George Orwell - 1984 Stephen E. Ambrose - Band of Brothers]} 
//     ExamScores: {set:true valid:true value:map[math:80 physics:90]} 
//     Sibling:    {set:true valid:true value:{
//             Name:   {set:true valid:true value:Anna} 
//             Age:    {set:true valid:true value:20} 
//             Type:   {set:true valid:true value:0}
//         }
//     }
// }

p = Person{} // reset variale
jsonStr4 := []byte(`{"age":25,"birth_date":"1997-01-01T00:00:00+01:00","books":null,"name":"Peter"}`)
_ = json.Unmarshal(jsonStr4, &p)
log.Printf("%+v", p)
// {
//     Name:       {set:true valid:true value:Peter} 
//     Age:        {set:true valid:true value:25} 
//     BirthDate:  {set:true valid:true value:{wall:0 ext:62987670000 loc:0x2ca260}} 
//     Books:      {set:true valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:false valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }
```

### 6. Implement `JSON` unmarshaler for `SiblingType`
```go
// UnmarshalJSON implements the json.Unmarshaler interface
func (st *SiblingType) UnmarshalJSON(data []byte) error {
    // null value already handled by null.Var, 
    // so we don't need to check for that here

	str := string(data)
	if str == "brother" {
		*st = Brother
	} else if str == "sister" {
		*st = Sister
	} else {
		return fmt.Errorf("invalid type for %T", st)
	}

	return nil
}

var p Person
log.Printf("%+v", p)
// {
//     Name:       {set:false valid:false value:} 
//     Age:        {set:false valid:false value:0} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:false valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }

jsonStr1 := []byte(`{"age":25,"name":"Peter","sibling":{"age":20,"name":"Anna","type":"sister"}}`)
_ = json.Unmarshal(jsonStr1, &p)
log.Printf("%+v", p)
// {
//     Name:       {set:true valid:true value:Peter} 
//     Age:        {set:true valid:true value:25} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:true valid:true value:{
//             Name:   {set:true valid:true value:Anna} 
//             Age:    {set:true valid:true value:20} 
//             Type:   {set:true valid:true value:0}
//         }
//     }
// }

p = Person{}
jsonStr2 := []byte(`{"age":25,"name":"Peter","sibling":{"age":20,"name":"Anna","type":"foo"}}`)
err := json.Unmarshal(jsonStr2, &p)
log.Printf("%v", err)
// invalid type for *SiblingType

p = Person{}
jsonStr3 := []byte(`{"age":25,"name":"Peter","sibling":null}`)
_ = json.Unmarshal(jsonStr3, &p)
log.Printf("%+v", p)
// {
//     Name:       {set:true valid:true value:Peter} 
//     Age:        {set:true valid:true value:25} 
//     BirthDate:  {set:false valid:false value:{wall:0 ext:0 loc:<nil>}} 
//     Books:      {set:false valid:false value:[]} 
//     ExamScores: {set:false valid:false value:map[]} 
//     Sibling:    {set:true valid:false value:{
//             Name:   {set:false valid:false value:} 
//             Age:    {set:false valid:false value:0} 
//             Type:   {set:false valid:false value:0}
//         }
//     }
// }
```


### 7. Use with the `database/sql` package
```go
// NOTE:
// If the underlying value is not compatible with the database/sql package by default, 
// then implement sql.Scanner and driver.Valuer interfaces for the underlying type

var p Person

// Exec
// If any of the variables are not set or set to NULL, then NULL will be inserted
_ = db.Exec(/* query */, p.Age, p.Name)

// Scan
// If any of the values scanned will be NULL, then the IsValid() will return false. 
// IsSet() will return true for every field used in the scan.
_ = db.QueryRow(/* query */).Scan(&p.Age, &p.Name)
```


