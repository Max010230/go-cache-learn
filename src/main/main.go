package main

import (
	"flag"
	"fmt"
	mircool_cache "go-cache-learn/mircool-cache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *mircool_cache.Group {
	return mircool_cache.NewGroup("score", 2<<10, mircool_cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key ", key)
		if value, ok := db[key]; ok {
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, group *mircool_cache.Group) {
	peers := mircool_cache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("mircool cache is run at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *mircool_cache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("fontend server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	//vaule1 := int64(1)
	//vaule2 := int64(2)
	//vaule3 := int64(3)
	//a := A{
	//	Id:    &vaule1,
	//	Level: &vaule1,
	//}
	//b := A{
	//	Id:    &vaule2,
	//	Level: &vaule2,
	//}
	//c := A{
	//	Id:    &vaule3,
	//	Level: &vaule2,
	//}
	//a.Child= []*A{&b,&c}
	//d := A{
	//	Id:    &vaule3,
	//	Level: &vaule3,
	//}
	//b.Child= []*A{&d}
	//if bytes, err := json.Marshal(a);err!=nil{
	//	log.Fatal(err.Error())
	//}else {
	//	log.Println(string(bytes))
	//}

	//mircool_cache.NewGroup("score", 2<<10, mircool_cache.GetterFunc(func(key string) ([]byte, error) {
	//	log.Println("[SlowDB] search key ", key)
	//	if value, ok := db[key]; ok {
	//		return []byte(value), nil
	//	}
	//	return nil, fmt.Errorf("%s not exist", key)
	//}))
	//addr := "localhost:9999"
	//peers := mircool_cache.NewHTTPPool(addr)
	//log.Println("mircool cache is run at ", addr)
	//log.Fatal(http.ListenAndServe(addr, peers))

	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}

type A struct {
	Id    *int64
	Level *int64
	Child []*A
}
