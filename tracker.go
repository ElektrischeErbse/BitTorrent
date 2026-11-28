package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Peer struct {
	PeerID string
	IP     string
	Port   int
}

type Tracker struct {
	mutex sync.Mutex
	// key: hash_info, value: peers
	Torrents map[string][]Peer
}

func NewTracker() *Tracker {
	return &Tracker{
		Torrents: make(map[string][]Peer),
	}
}

func (t *Tracker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 解析URL
	query := r.URL.Query()
	infoHash := query.Get("info_hash")
	peerID := query.Get("peer_id")
	portStr := query.Get("port")
	uploadedStr := query.Get("uploaded")
	downloadedStr := query.Get("downloaded")
	leftStr := query.Get("left")

	// 获取客户端IP
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Get client IP failed", http.StatusBadRequest)
		return
	}
	if infoHash == "" || peerID == "" || portStr == "" || uploadedStr == "" || downloadedStr == "" || leftStr == "" {
		http.Error(w, "Missing info_hash, peer_id, port, uploaded, downloaded or left", http.StatusBadRequest)
		return
	}

	// 添加或更新peer信息
	port, _ := strconv.Atoi(portStr)
	uploaded, _ := strconv.Atoi(uploadedStr)
	downloade, _ := strconv.Atoi(downloadedStr)
	left, _ := strconv.Atoi(leftStr)
	fmt.Fprintf(os.Stdout, "info_hash: %s, peer_id: %s, port: %d, uploaded: %d, downloaded: %d, left: %d\n",
		infoHash, peerID, port, uploaded, downloade, left)

	t.mutex.Lock()
	peers := t.Torrents[infoHash]
	var newPeers []Peer
	for _, peer := range peers {
		if peerID != peer.PeerID {
			newPeers = append(newPeers, peer)
		}
	}
	newPeers = append(newPeers, Peer{PeerID: peerID, IP: host, Port: port})
	t.Torrents[infoHash] = newPeers
	t.mutex.Unlock()
	fmt.Fprintf(os.Stdout, "peers: %v\n", t.Torrents[infoHash])

	// 构建bencode响应
	// interval: 客户端下次请求的间隔时间（秒）
	// peers: 紧凑格式（Compact）的节点列表，每个节点 6 字节（4字节IP + 2字节Port）
	var peerBuffer []byte
	t.mutex.Lock()
	peers = t.Torrents[infoHash]
	for _, peer := range peers {
		if peerID == peer.PeerID {
			continue
		}
		ipBytes := net.ParseIP(peer.IP).To4()
		peerBuffer = append(peerBuffer, ipBytes...)
		// 拼接 2 字节 Port（大端序）
		peerBuffer = append(peerBuffer, byte(peer.Port>>8), byte(peer.Port))
	}
	t.mutex.Unlock()
	fmt.Fprintf(os.Stdout, "peerBuffer: %s\n", hex.EncodeToString(peerBuffer))
	var data = map[string]any{
		"interval": 10,
		"peers":    string(peerBuffer),
	}
	bencode, _ := BuildBencode(data)
	fmt.Fprintf(os.Stdout, "bencode: %s\n", bencode)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(bencode))
}

func (t *Tracker) StartTrackerService() {
	http.Handle("/announce", t)
	fmt.Println("Simple Tracker 启动，监听端口 :8080...")
	fmt.Println("Announce URL: http://127.0.0.1:8080/announce")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Start tracker error: %v", err)
	}
}
