![GitHub release (latest by date)](https://img.shields.io/github/v/release/mauserzjeh/null?style=flat-square)

# Null

Null is a package that provides a generic nullable variable that can be used for both `JSON` and `SQL` operations.

## Features
- Uses a generic type for the underlying value
- Provides an easy way to test if certain struct fields were set after unmarshaling a json message
- Makes it easier to update only those fields in the database which were set previously. Especially in `NoSQL` databases using `JSON` documents
- Compatible with the standard `database/sql` package
    - If the underlying type is a custom type that the standard `database/sql` package can't handle, then it is possible to implement the `sql.Scanner` and `driver.Valuer` interfaces for that given type to be compatible
- Possible to control how the underlying value is handled during json marshalling and unmarshalling by implementing the `json.Marshaler` and `json.Unmarshaler` interfaces for the underlying type

## Installation
```
go get -u github.com/mauserzjeh/null
```

## Usage & Examples
Below you can see how to use the package in general and also in a more complex scenario. There are a few examples in `test_null.go` as well.

#### General usage
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

#### Complex usage

1. Let's have the following types
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

2. Default `JSON` marshaling
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

3. Implement `JSON` marshaler for `Person` and `Sibling`
```go
// NOTE: 
// It is better to implement a generic function that can turn a struct into a map. Then in the MarshalJSON implementations, turn the structure into a map, loop through each item and do a type switch for nullable variables (each type of null.Var used in the original struct should have a case) and just filter out those variables which were not set. To keep this example clean, we will just manually check each field instead.

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

4. Default `JSON` unmarshal
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

5. Implement `JSON` unmarshaler for `SiblingType`
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

6. Use with the `database/sql` package
```go
// NOTE:
// If the underlying value is not compatible with the database/sql package by default, then implement sql.Scanner and driver.Valuer interfaces for the underlying type

var p Person

// Exec
// If any of the variables are not set or set to NULL, then NULL will be inserted
_ = db.Exec(/* query */, p.Age, p.Name)

// Scan
// If any of the values scanned will be NULL, then the IsValid() will return false. IsSet() will return true for every field used in the scan.
_ = db.QueryRow(/* query */, &p.Age, &p.Name)
```


