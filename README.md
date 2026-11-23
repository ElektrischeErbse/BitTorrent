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
