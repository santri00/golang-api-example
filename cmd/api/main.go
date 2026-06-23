package main

func main() {
	server, rdb := initialize()
	defer rdb.Close()
	server.Run()
}
