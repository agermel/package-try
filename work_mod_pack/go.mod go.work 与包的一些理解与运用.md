## go.mod go.work 与包的一些理解与运用

#### go.mod

是一个定义模块(Module)的文件
```
module package1 //创建名为package1的模块
```

cd到相关文件夹下用go mod init package1创建

一个模块下不能有其他模块

#### 包（package)

函数的集成体

**在模块建本地包版本**

一个在package1模块下名为greetings的包//

/package1
│  go.mod // `module package1`
│
└─greetings
        hello.go 
        ISAM.go

hello.go

```
package greetings //代表greetings包的一部分

import "fmt"

func Morning() {
	fmt.Println("Good Morning")
}

func Evening() {
	fmt.Println(("Good Evening"))
}

```

**在本地pkg建本地包版本**

Coming soon

**包内函数开头要大写才能被被包外程序使用！！！**

#### go.work **//使用在模块内建本地包的例子**

如上文 greeting包在package1模块里，而运行程序在Try模块里
` import "package1/greetings" `会出现问题 ==go中不能跨模块引入包==
此时go.work可以出手了（实现多模块在一个“模块下”运用？）

设置 `go.work` 文件

`go work init ./package1 ./Try`

实现
├── go.work            # 工作区文件
├── package1
│   ├── go.mod         # 包含 "package1" 模块的 Go 配置
│   └── greetings
│       ├── hello.go
│       └── ISAM.go
└── Try
    ├── go.mod         # 包含 "Try" 模块的 Go 配置
    └── main.go        # 主程序文件

#### 有趣的东西

1. 包级别的全局变量，存储最后输入的值!!! 不被函数所束缚(函数执行完毕它还在，等Get函数来返回)

```
package fibo

// 包级别的全局变量，存储最后输入的值!!! 不被函数所束缚(函数执行完毕它还在，等Get函数来返回)
var LastInput int

// 计算斐波那契数列的函数
func Fibonacci(n int) int {
    LastInput = n
    if n <= 0 {
        return 0
    } else if n == 1 {
        return 1
    } else {
        return Fibonacci(n-1) + Fibonacci(n-2)
    }
}

// 输出最后输入的值
func GetLastInput() int {
    return LastInput
}
```

