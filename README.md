# glog

自用的golang日志库

```go
import "github.com/AkvicorEdwards/glog"
```

# 类型

使用`SetMask(m int)`设置每种消息是否输出

- `Unknown(format string, values ...any)`: 写入`os.Stdout`
- `Debug(format string, values ...any)`: 写入`os.Stdout`
- `Trace(format string, values ...any)`: 写入`os.Stdout`
- `Info(format string, values ...any)`: 写入`os.Stdout`
- `Warning(format string, values ...any)`: 写入`os.Stdout`
- `Error(format string, values ...any)`: 写入`os.Stderr`
- `Fatal(format string, values ...any)`: 写入`os.Stderr`

# 前缀

支持添加时间、类型、调用位置，使用`SetFlag(f int)`设置前缀。

# 文件

支持在显示在控制台的同时写入文件，通过`SetLogFile(path string) error`设置文件，设置后会自动写入此文件。

通过`CloseFile()`关闭文件，关闭后不再写入文件

