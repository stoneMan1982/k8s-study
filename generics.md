# Golang Generics

## Prerequisites  

- An installation of Go 1.18 or later  
- A tool to edit your code  
- A command terminal


## Add non-generic functions  

[reference](./generic/non_generics.go)

## Add a generic function to handle multiple types  

## Remove type arguments when calling the generic function 

## Declare a type constraint  

## Conclusion  

## Completed Code

## Underlying Types  

In go, every type has an underlying type. The underlying type is important for determing type compatibility and what operations are allowed on a value.  

### Key Concepts  

- Built-in types  
  Built-in types are their own underlying types

- Named types  
  When you define a new type using type, it creates a named type with an underlying type:

  ``` go

    type MyInt int               // MyInt's underlying type is int
    type MyString string         // MyString's underlying type is string
    type MySlice  []int          // MySlice's underlying type is []int
    type MyMap map[string]int    // MyMap's underlying type is map[string]int
    type MyStruct struct {X int} // MyStruct's underlying type is struct {X int}

  ```

  - Type identity vs. type compatibility  
    - Two named types are identical only if they're the same type
    - But values can be assigned between types if their underlying types are identical

```go
    type A int
    type B int

    var x A = 5
    var y B = x        // Error: cannot assign A to B(different named types)
    
```

- Interface implementation  

A type implements an interface based on its methods, not its underlying type.

- Constants and untyped values  

Untyped constants have underlying types that determine compatibility

```go

    const x = 5                           // untyped int constant
    var a int = x                         // OK
    var b MyInt = MyInt(x)                // requires conversion

```

## Why this matters

- Type safety: Go enforces strict type checking through underlying types  
- Conversion: You can convert between a named type and its underlying  
- Reflection: The reflect package provides ways to inspect underlying types