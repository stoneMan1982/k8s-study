#  SHELL 脚本学习

## 元字符（特殊字符）

### 引号类

- 单引号 '
完全保留字符的字面含义，禁止所有扩展。

```sh
#!/bin/bash

name="World"

echo 'Hello $name'       # output: Hello $name
echo 'Today is `date`'   # output: Today is `date`
echo 'Line1\nLine2'      # output: Line1\nLine2

```

- 双引号 "

部分保留，允许变量扩展和命令替换，但禁止通配符扩展。  

```sh
#!/bin/bash

name="World"

echo "Hello $name"          # output: Hello World
echo "Today is `date`"      # output: Today is 'current date'
echo "File: *.txt"          # output: File: *.txt
```

- 反引号 ` 和 $()  

命令替换，执行命令并返回输出

```sh

#!/bin/bash

# 旧语法
files=`ls`
echo "files: $files"


current_dir=$(pwd)
echo "current directory: $current_dir"

files1=$(ls)

echo "files1: $files1"

```

### 变量相关 $

1. 变量引用

```sh
#!/bin/bash

name="Alice"
echo $name         # simple variable
echo ${name}       # clear boundary of variable
echo ${#name}      # the length of a variable
echo ${name:0:2}   # the substr of a variable
```

2. 特殊变量  

- $0 命令行参数的第一个参数，即脚本的名称，如: ./script.sh  
- $1-$9 命令行参数位置，对应第二个到第九个参数
- $# 参数个数
- $@ 所有参数，每个参数独立存在： "$1" "$2"  "$2"  
- $* 所有参数，一个整体："$1 $2 $3"  
- $? 上一个命令的退出状态，0表示成功  
- $! 最后一个后台进程的PID  
- $$ 当前shell的PID  

3. 通配符  

星号 * 匹配任意字符串

4. 问号 ?  

匹配任意单个字符。  

5. 方括号 []  

匹配指定范围内的任意单个字符。  

### 重定向类  

1. 输入/输出重定向  

```sh
#!/bin/bash

# 输出重定向  
command > file   ## 标准输出到文件（覆盖）
command >> file  ## 标准输出到文件（追加）
command 2> error.log ## 标准错误到文件
command &> output.log ## 所有输出到文件

# 输入重定向
command < file  # 从文件中读取
command << EOF  # 文件尾  
line1
line2
EOF

# 文件描述符
command 2>&1 ## 错误输出重定向到标准输出

```

### 括号类  

1. 圆括号  

在子shell中执行命令  

```sh
#!/bin/bash
(cd /tmp; ls -al)   ## 在子shell中切换目录
pwd                 ## 仍在原目录
```

2. 花括号  

- 在当前shell中执行

```sh
{ echo "start"; ls; echo "end"; } > output.log
```

- 变量边界:

```sh
#!/bin/bash
name="Alice"
echo "Hello ${name}"
```

- 扩展用法：  

```sh
echo file{1,2,3}.txt   ## file1.txt file2.txt file3.txt
echo {a..z}            ## a b c ... z

mkdir -p project/{src, bin, docs} ## 创建多个目录
```
