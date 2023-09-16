package main

import "myzinx/znet"

func main() {
	s := znet.NewServer("[zinx V0.1]")
	s.Serve()
}
