# BitTorrent

### bencode格式
整数 (Integer)：以 i 开头，e 结尾1112。例如 i42e 代表 42

字符串 (String)：<长度>:<字符串内容>。例如 4:spam 代表 "spam"

列表 (List)：以 l 开头，e 结尾。例如 l4:spam4:eggse 代表 ["spam", "eggs"]

字典 (Dictionary)：以 d 开头，e 结尾，键值对交替出现，且键必须按字典序升序排列

### torrent文件结构
announce: Tracker 服务器的 URL

info: 最核心的部分，包含文件的元数据。它的 SHA-1 哈希值就是 Info Hash

announce-list: 包含多个 Tracker 服务器的 URL 列表，用于备用

creation date: 创建时间，Unix 时间戳格式

comment: 备注

### info内部结构
#### 单文件
- name：文件名（字符串）
- length：文件总大小（整数，单位字节）
- piece length：每个数据块的大小（整数，通常是 256KB 等 2 的幂次方）
- pieces：所有数据块的 SHA-1 哈希值拼接而成的二进制串（用于校验完整性）
#### 多文件
- name：顶层目录名（字符串）
- piece length：每个数据块的大小
- pieces：同上
- files：一个列表，里面包含多个字典，每个字典描述一个文件：
    - length：该文件的大小
    - path：列表结构，包含子目录名和最终的文件名（例如 l4:test8:test.txte）

### info hash
info hash 是 torrent 文件中 info 字典的 SHA-1 哈希值，用于唯一标识一个 torrent 文件。
它在 BitTorrent 协议中起着关键作用，尤其是在 peer 之间交换数据时，用于验证文件的完整性和一致性。
```go
info, ok := getValueByKey[map[string]any](result, "info")
if !ok {
    fmt.Fprintln(os.Stderr, "Failed to get info")
}
// 构建 bencode 编码的 info 字典
info_content, err := BuildBencode(info)
if err != nil {
    fmt.Fprintf(os.Stderr, "BuildBencode error: %v\n", err)
}
// 计算 info hash
hash := sha1.Sum([]byte(info_content))
torrent.InfoHash = hex.EncodeToString(hash[:])
```

### BitTorrent Tracker 的 GET 请求
请求的 Base URL 来源于 .torrent 文件中定义的 announce 字段。参数通过标准的 CGI 方式（即 ? 后跟 & 分隔的键值对）追加到 URL 后面
```
GET /announce?info_hash=%124Vx%9A%BC%DE%F1%23Eg%89%AB%CD%EF%124Vx%9A&peer_id=-LT3000-123456789ABC&port=6881&uploaded=0&downloaded=0&left=1048576&compact=1&event=started HTTP/1.1
Host: tracker.example.com
```
Tracker 请求包含以下参数，其中所有二进制数据（特别是 info_hash 和 peer_id）必须进行严格的 URL 转义（不在 0-9, a-z, A-Z, ., -, _, ~ 集合中的字节必须以 %nn 的十六进制形式编码）
1. 必填字段
- info_hash 20 字节的 SHA-1 哈希值，即 .torrent 文件中 info 字典的哈希值
- peer_id 20 字节的唯一标识符，用于标识客户端
- port 客户端监听的端口号
- uploaded 已上传的字节数
- downloaded 已下载的字节数
- left 剩余需要下载的字节数

2. 可选字段
- compact 是否请求紧凑的 peer 列表（1 表示是，0 表示否）
- event 事件类型（started、stopped、completed 等）
- ip 可选的 IP 地址，通常由 Tracker 服务器自动获取
- numwant 请求的 peer 数量，默认值为 50
- no_peer_id 是否请求 peer_id（1 表示不请求，0 表示请求）

### BitTorrent Tracker 的响应
Tracker 响应是一个 bencode 编码的字典，包含以下字段：
1. 必填字段
- interval 客户端应等待的秒数，之后再次向 Tracker 发送请求
- peers peer 列表，格式取决于 compact 参数的值：
    - compact=1：peers 是一个二进制字符串，每 6 个字节表示一个 peer（前 4 字节为 IP 地址，后 2 字节为端口号）
    - compact=0：peers 是一个列表，每个元素是一个字典，包含 peer_id、ip 和 port 字段

