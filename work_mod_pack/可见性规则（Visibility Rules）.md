在 Go 中，**可见性规则**（Visibility Rules）决定了标识符（如变量、函数、字段、方法等）是否能在包外被访问。通过使用这些规则，可以限制或禁止某些功能的使用，例如禁止使用 `new` 函数。

### 可见性规则（Visibility Rules）
- 在 Go 中，**首字母大写**的标识符是导出的（public），可以被其他包访问。
- **首字母小写**的标识符是未导出的（private），只能在定义它的包内使用，外部包不能访问。

### `new` 函数简介
`new` 函数是 Go 内置的一个函数，用于分配内存并返回一个指向类型的指针。它返回的是类型的零值指针：

```go
p := new(Type)  // p 是 *Type 类型的指针，指向 Type 的零值
```

### 如何通过可见性规则禁止使用 `new` 函数

通过可见性规则，限制 `new` 函数的使用并不是直接禁止它，而是通过将某些结构体字段、方法或类型的可见性设置为未导出（即首字母小写），这样就可以防止外部包通过 `new` 函数创建这些类型的实例。

### 例子：禁止通过 `new` 函数创建类型的实例

假设你希望禁止外部包使用 `new` 函数来创建某个类型的实例。可以通过将该类型的构造函数字段设置为未导出，并且提供一个只能在包内使用的初始化函数。

#### 例子 1: 禁止外部包使用 `new` 创建实例

```go
package rectangle

type Rectangle struct {
    width  float64  // 未导出的字段
    height float64 // 未导出的字段
}

// 在包内提供一个构造函数，用于创建实例
func NewRectangle(w, h float64) *Rectangle {
    return &Rectangle{width: w, height: h}
}
```

在上面的例子中，`Rectangle` 结构体的字段是未导出的（首字母小写），因此外部包不能直接访问 `width` 和 `height` 字段，也不能直接使用 `new(Rectangle)` 来创建该类型的实例。相反，外部包需要使用提供的 `NewRectangle` 构造函数来创建 `Rectangle` 类型的实例。

#### 例子 2: 尝试禁止 `new` 的使用

假设你在 `main` 包中尝试通过 `new` 函数创建 `Rectangle` 实例：

```go
package main

import (
    "fmt"
    "your-module-path/rectangle" // 导入 rectangle 包
)

func main() {
    // 这行会导致编译错误，因为 Rectangle 的字段是未导出的，不能直接使用 new(Rectangle)
    r := new(rectangle.Rectangle)
    fmt.Println(r)
}
```

由于 `Rectangle` 的字段是未导出的，外部包无法通过 `new` 来访问和设置这些字段，因此会引发错误。外部包必须通过 `NewRectangle` 函数来创建该类型的实例。

### 总结

通过应用 Go 的可见性规则，你可以控制哪些字段、方法或类型是可访问的，从而间接禁止外部包使用 `new` 函数来创建某些类型的实例。常见的做法是将类型的字段和构造函数设为未导出，只提供包内的初始化函数，让外部包不能直接访问类型的内部字段。