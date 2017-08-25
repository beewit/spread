package main

import (
	"time"
)

func main() {
	test_1()
	println("我这执行了test_2")
}

func test_1() {
	time.Sleep(10 * time.Second)
	println("我这执行了test_1")
}
