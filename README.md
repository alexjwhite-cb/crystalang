# Jet

## Introduction

Jet is a dynamically-typed interpretted object-_orchestration_ language inspired by Go, PHP, and Python.

## Language Features

### Entrypoint

Jet uses the `main` method as its default entrypoint.

`.jet` is the default file extension for Jet files.

### Types

Jet supports the following types:

* `int`
* `float`
* `boolean`
* `string`
* `map`

Jet treats all arrays as maps in order to have a unified function set:

* `(map, value)->add` - adds the value to the map with the lowest available int starting at `0`
* `(map, key)->remove` - removes the specified key and associated value from the map
* `map->unique` - returns a copy of the map with any duplicate **values** removed
* `map->ksort` - sorts by the key in ascending order
* `map->vsort` - sorts by the value in ascending order
* `map->reverse` - reverses the current order

### Returning

Jet does not have the return keyword, but instead uses the passthrough syntax: `->`

Methods in Jet do not have explicit expectations in regard to return values, so numerous arguments can be returned by encapsulating them within parenthesis like so: `(x, y)->`

### Methods

#### 1. Declaration

Methods are declared with the `meth` keyword and code blocks are defined with braces `{}`. Arguments are declared after a colon (`:`) and are comma seperated. When an unknown number of arguments are required, the `*` suffix can be used. Additional arguments will be compiled into a map and can be accessed via the argument name that precedes the `*` token. 

Methods can be declared in any of the following formats.

```
meth myFunction {}

meth myFunction: arg1, arg2 {}

meth myFunction: arg1, arg* {}

myClosure = meth: x, y { (x + y)-> }

myClosure = meth {}
```

#### 2. Calling Methods

If a method has no parameters, parenthesis `()` can be omitted.
```
myMethod
myObject.MyMethod
myMethod(x, y)
myObject.MyMethod()
```
```
a, _ = myMethod  // Returns 
_, b = myMethod
c = myMethod
```

#### 3. Non-Declarative Argument Parsing

Just as `->` is used to return, values can be passed directly into functions to create function chains as follows:
```
meth Foo: array { (a, b, c)-> }
meth Bar: args* { (d, e)-> }
meth Baz: arg1, arg2 { (string)-> }

myString  = Foo(myArray)->Bar->Baz

a, b, c   = Foo(myArray)
d, e      = Bar(a, b, c)
myString2 = Baz(d, e)

if myString == myString2 {
    // Evaluates to true.
} 
```

### Descriptors

Descriptors are Jet's response to Classes. A descriptor is used to **describe** the functionality of a given object.

Both descriptors and their properties should be capitalised. Properties not declared as a descriptor argument can be labelled as constant.

```
describe Vehicle: Seats {
    const Material = "Metal"
    Wheels = 4
}

describe Jet: Name, TopSpeed {
    meth speedBoost { (TopSpeed * 2)-> }

    meth canFly { (true)-> }
}
```

### Objects

Objects are orchestrated from descriptors. Arguments are inherited from the descriptors in the order they are assigned to the object.

Inherited methods can be overloaded by the object. New methods can also be added to the object, allowing utilisation of properties from across descriptors.

Properties not included as a descriptor argument can be updated later. Accessing properties within the object requires the `Descriptor.Property` format.

New properties cannot be added in runtime code, and constants cannot be updated. Attempting to do so will cause Jet to panic.

```
object FighterJet: Vehicle, Jet {
    overload speedBoost { (Jet.TopSpeed * 4)-> }
    
    meth describeFighterJet {
        ("{Jet.Name} is made of {Vehicle.Material} has {Vehicle.Wheels}")->print
    }
}

myFighter = FighterJet(2, "Falcon", 100)
3->myFighter.Wheels
```

## Language Objectives

* [ ] Jet uses a common entrypoint; `main` will always be used to initialise a program.
* [ ] Inheritence is "shallow". Object-types can be defined and their default methods defined and implemented, however one object-type cannot inherit from another. They must be orchestrated together.
* [ ] Attributes and child methods are accessed via `.` syntax
* [ ] Values are piped into and out of functions with the "spoon" syntax: `()->`

## Resources

[High Level Principles / Tokenisation](https://www.freecodecamp.org/news/the-programming-language-pipeline-91d3f449c919/)

[Abstract Syntax Trees](https://www.twilio.com/blog/abstract-syntax-trees)

[Golang AST Package](https://tech.ingrid.com/introduction-ast-golang/)