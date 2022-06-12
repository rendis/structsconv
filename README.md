# Structs Conversion (structsconv)

[![coverage](https://api.codiga.io/project/32663/score/svg)](https://app.codiga.io/public/project/32663/structsconv/dashboard)
[![quality](https://api.codiga.io/project/32663/status/svg)](https://app.codiga.io/public/project/32663/structsconv/dashboard)
[![GoDoc](https://pkg.go.dev/badge/rendis/structsconv?status.svg)](https://pkg.go.dev/github.com/rendis/structsconv?tab=doc)
---

The **structsconv** package makes it easy to convert between structs of different types. This development is inspired by
the java [mapstruct](https://mapstruct.org/) library.

## Motivation
The main purpose of this library is that, within a `hexagonal architecture`, the conversion can be done without 
the need to _dirty_ the domain layer with `tags`.

For more information on how to use, see the [examples](https://github.com/rendis/structsconv/tree/master/example).

## Release Compatibility

- v1.0.0: is compatible with Golang versions:
    - 1.17+

## Installation

1. Install the package:

```bash
go get -u github.com/rendis/structsconv
```

2. Import the package:

```go
import "github.com/rendis/structsconv"
```

## Quick start

### Default Mapping
Fields with the `same name` and `same type` are mapped by default.
<br/>
<br/>

#### simple struct &rarr; simple struct

Structs
```go
// Source struct
type SiteSource struct {
  Name string
  URL  string
}
```

```go
// Target struct
type SiteTarget struct {
  Name string
  URL  string
}
```

Mapping
```go
// Mapping struct
source := SiteSource{
  Name: "Golang",
  URL:  "https://golang.org",
}
target := SiteTarget{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target) // Output: {Name:Golang URL:https://golang.org }
```
---
<br/>

#### nested struct &rarr; nested struct

Structs
```go
// Source struct
type SiteSource struct {
  Name string
  URL  string
  Meta MetaInfoSource
}

type MetaInfoSource struct {
    Description string
}
```

```go
// Target struct
type SiteTarget struct {
  Name string
  URL  string
  Meta MetaInfoTarget
}

type MetaInfoTarget struct {
    Description string
}
```

Mapping
```go
// Mapping struct
source := SiteSource{
  Name: "Golang",
  URL:  "https://golang.org",
  Meta: MetaInfoSource{
    Description: "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
  },
}
target := SiteTarget{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target) // Output: {Name:Golang URL:https://golang.org Meta:{Description:Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.}}
```
---
<br/>

#### slice &rarr; slice

Structs
```go
// Source struct
type ProductDto struct {
  Name     string
  Price    int
  Keywords []string
  Locals   []LocalDto
}

type LocalDto struct {
  Id    int
  Stock int
}
```

```go
// Target struct
type ProductDomain struct {
  Name     string
  Price    int
  Keywords []string
  Locals   []LocalDomain
}

type LocalDomain struct {
  Id    int
  Stock int
}
```

Mapping
```go
// Mapping struct
source := ProductDto{
  Name:     "go shoes",
  Price:    100,
  Keywords: []string{"shoes", "sneakers"},
  Locals: []LocalDto{
    {Id: 1, Stock: 10},
    {Id: 2, Stock: 20},
  },
}
target := ProductDomain{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target) // Output: {Name:go shoes Price:100 Keywords:[shoes sneakers] Locals:[{Id:1 Stock:10} {Id:2 Stock:20}]}
```
---
<br/>

#### map &rarr; map
>  - Source map key type must be equal to target map key type
>  - Source map value type must be equal to target map value type

> Everything directly associated with the map type **must be exported**.

Structs
```go
// Source struct
type SchoolDto struct {
    Students map[string]StudentDto
}

type StudentDto struct {
  Id  string
  Name string
  Age  int
}
```

```go
// Target struct
type SchoolDomain struct {
    Students map[string]StudentDomain
}

type StudentDomain struct {
  Id  string
  Name string
  Age  int
}
```

Mapping
```go
// Mapping struct
source := SchoolDto{
  Students: map[string]StudentDto{
    "stu124": {
      Id:    "stu124",
      Name: "John",
      Age:  20,
    },
    "stu125": {
      Id:    "stu125",
      Name: "Jane",
      Age:  21,
    },
  },
}
target := SchoolDomain{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target) // Output: {Students:{stu124:{Id:stu124 Name:John Age:20} stu125:{Id:stu125 Name:Jane Age:21}}}
```
---
<br/>

### Rules
Rules are used to customize the assignment. Use when desired:
* Map fields with _different names_ or _different types_ or both.
* Define a _constant_.
* Ignore a field.
* A custom mapping.

<br/>

#### Rules for fields with different names
Fields with the `different name` and `same type` are mapped using `rules`.

Structs
```go
// Source struct
type SiteSource struct {
  Name string
  URL  string
}
```

```go
// Target struct
type SiteTarget struct {
  SiteName string
  URL  string
}
```

Create the rules
```go
var rules = structsconv.RulesSet{}   // Create a new rules set
rules["SiteName"] = "Name"           // Rule to map 'Name' to 'SiteName'
siteRulesDefinition := structsconv.RulesDefinition{
    Rules:  rules,
    Source: SiteSource{},            // To get the source struct type
    Target: SiteTarget{},            // To get the target struct type
}
```
> The rules are created based on the **names on the target struct**.
<br/>

Register the rules
```go
structsconv.RegisterRulesDefinitions(siteRulesDefinition)
```

Mapping
```go
source := SiteSource{
    Name: "Golang",
    URL:  "https://golang.org",
}
target := SiteTarget{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target)
```
---
<br/>


#### Rules for constants
Structs
```go
//Source struct
type ItemDto struct {
	Id   int
	Name string
}
```

```go
// Target struct
type ItemDomain struct {
	Id            int
	ItemName      string
	ConstantValue string  // Field does not exist in the source struct
}
```

Create the rules
```go
var rules = structsconv.RulesSet{}

rules["ItemName"] = "Name"

rules["ConstantValue"] = func() string {
    return "Constant value"
}

rulesDefinition := structsconv.RulesDefinition{
    Rules:  rules,
    Source: ItemDto{},
    Target: ItemDomain{},
}
```
Register the rules
```go
structsconv.RegisterRulesDefinitions(rulesDefinition)
```
Mapping
```go
source := ItemDto{
    Id:   1,
    Name: "Toy",
}
target := ItemDomain{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target) // Output: {Id:1 ItemName:Toy ConstantValue:Constant value}
```
---
<br/>


#### Rules for ignoring fields
To ignore a field, just set the `rules` to `nil`.

Structs
```go
//Source struct
type ItemDto struct {
	Id   int
	Name string
}
```

```go
// Target struct
type ItemDomain struct {
	Id            int
	ItemName      string
	IgnorableValue string  // Field does not exist in the source struct
}
```

Create the rules
```go
var rules = structsconv.RulesSet{}

rules["ItemName"] = "Name"

rules["IgnorableValue"] = nil

rulesDefinition := structsconv.RulesDefinition{
    Rules:  rules,
    Source: ItemDto{},
    Target: ItemDomain{},
}
```

Register the rules
```go
structsconv.RegisterRulesDefinitions(rulesDefinition)
```

Mapping
```go
source := ItemDto{
    Id:   1,
    Name: "Toy",
}
target := ItemDomain{}

structsconv.Map(&source, &target)
fmt.Printf("%+v\n", target) // Output: { Id:1 ItemName:Toy IgnorableValue:"" }
```
---
<br/>


### Rules with arguments

#### Root Struct & Current Struct
You will always be able to request the `root origin struct` and the `current origin struct`, to use them in custom mapping, you just need to request them. 

Let see an example with 2 levels of nested structs:
```go
// Source struct
type ParentRootDto struct {
  ParentField string
  Child1      Child1Dto
}

type Child1Dto struct {
  Level1Field string
  Child       Child2Dto
}

type Child2Dto struct {
    Level2Field string
}
```

```go
// Target struct
type ParentRootDomain struct {
  ParentField string       // Mapped directly
  Child1      Child1Domain // Mapped directly
}

type Child1Domain struct {
  Level1Field string       // Mapped directly
  Child       Child2Domain // Mapped directly
  
  ParentField string       // Custom mapping
  Description string       // Custom mapping
}

type Child2Domain struct {
  Level2Field string       // Mapped directly
  
  ParentField string       // Custom mapping
  Level1Field string       // Custom mapping
  Description string       // Custom mapping
}
```

Rules for `Child1Domain`
```go
var child1Rules = structsconv.RulesSet{}

child1Rules["ParentField"] = func(parent ParentRootDto) string {
    return parent.ParentField
}

child1Rules["Description"] = func(current Child1Dto, parent ParentRootDto) string {
    return fmt.Sprintf("root parent = %s,  current = %s", parent.ParentField, current.Level1Field)
}

child1Definition := structsconv.RulesDefinition{
    Rules:  child1Rules,
    Source: Child1Dto{},
    Target: Child1Domain{},
}
```

Rules for `Child2Domain`
```go
var child2Rules = structsconv.RulesSet{}

// Note that parent is ParentRootDto and not Child1Dto
child2Rules["ParentField"] = func(parent ParentRootDto) string {
    return parent.ParentField
}

// Note that parent is ParentRootDto and not Child1Dto
child2Rules["Level1Field"] = func(parent ParentRootDto) string {
    return parent.Child1.Level1Field
}

// Note that parent is ParentRootDto and not Child1Dto
child2Rules["Description"] = func(parent ParentRootDto, current Child2Dto) string {
    return fmt.Sprintf(
        "root parent = %s, parent = %s, current = %s",
        parent.ParentField, parent.Child1.Level1Field, current.Level2Field,
    )
}

child2Definition := structsconv.RulesDefinition{
    Rules:  child2Rules,
    Source: Child2Dto{},
    Target: Child2Domain{},
}
```
> **Note that parent is `ParentRootDto` and not `Child1Dto`**.

Register the rules
```go
structsconv.RegisterRulesDefinitions(child1Definition)
structsconv.RegisterRulesDefinitions(child2Definition)
```

Mapping
```go
source := ParentRootDto{
    ParentField: "ParentField (root)",
    Child1: Child1Dto{
        Level1Field: "level1 (Child of rut, parent of child2)",
        Child: Child2Dto{
            Level2Field: "level2 (last level, Child of Child1)",
        },
    },
}

target := ParentRootDomain{}

structsconv.Map(&source, &target) // To long to print but works :)
fmt.Printf("%+v\n", target)
```
---
<br/>

#### More and more arguments
You can also get more than just the root struct and the current struct as an argument.

By passing arguments to the `map()` function, they will be available for request in all rules.

```go
...
str1 = "str1"
str2 = "str2"
int1 = 1
int2 = 2
struct1 = struct{}{}
struct2 = struct{}{}
structsconv.Map(&source, &target, str1, str2, int1, int2, struct1, struct2)
...
```
To use them in the rules, you only need to request the ones you need.
* Request all the arguments (Ordered)
```go 
rules["targetField1"] = func(str1, str2 string, int1, int2 int, struct1 struct{}, struct2 struct{}) string {
    return fmt.Sprintf("%s %s %d %d %+v %+v", str1, str2, int1, int2, struct1, struct2)
}
```

* Request all argument "Unordered" (1)
```go
rules["targetField1"] = func(struct1 struct{}, struct2 struct{}, int1, int2 int, str1, str2 string) string {
    return fmt.Sprintf("%s %s %d %d %+v %+v", str1, str2, int1, int2, struct1, struct2)
} 
```

* Request all argument "Unordered" (2)
```go
rules["targetField1"] = func(struct1 struct{}, int1 int, str1 string, struct2 struct{}, str2 string, int2 int) string {
    return fmt.Sprintf("%s %s %d %d %+v %+v", str1, str2, int1, int2, struct1, struct2)
} 
```

Notice the double quotes around the word "Unordered", this is because arguments can be requested out of order relative to their type, arguments of the same type must be requested in the same order as they were passed.

Let's see some examples:
```go
str1 = "str1"
str2 = "str2"
structsconv.Map(&source, &target, str1, str2)
```
```go
rules["targetField1"] = func(struct1 struct{}, arg1 string, arg2 string) string {
    return fmt.Sprintf("%s %s", str1, str2)
}
```
```go
rules["targetField1"] = func(struct1 struct{}, arg2 string, arg1 string) string {
    return fmt.Sprintf("%s %s", str1, str2)
}
```
In both cases, the first argument is `str1` and the second is `str2`.

In summary, the rules are:
* Order does not matter on different types.
* Order matters on the same types.

Another way to use argument passing is through `slices`.
```go
str1 = "hello"
str2 = "world"
int1 = 1
int2 = 2
args := []interface{}{str1, str2, int1, int2}
structsconv.Map(&source, &target, args)
```
```go
rules["targetField1"] = func(args []interfaces{}) string {
    return fmt.Sprintf("first args = %s, last args = %d", args[0], str2[2]) // first args = hello, last args = 2
}
```
