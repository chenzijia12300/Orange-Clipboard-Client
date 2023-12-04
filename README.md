# orange-clipboard
This is a cross-platform clipboard data synchronization client (Windows/MacOS)

```go
// 克隆代码库
git clone https://github.com/chenzijia12300/orange-clipboard.git

// 移动到文件夹下
cd orange-clipboard

// 拉取项目依赖
go mod tidy

// 编译项目
go build

// 运行服务器
cd server/cmd/
go build -o orange-clipboard-server
./orange-clipboard-server

// 运行客户端
cd client/cmd/
go build -o orange-clipboard-client
./orange-clipboard-client
```

