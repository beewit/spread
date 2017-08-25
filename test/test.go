package main

import (
	"fmt"

	"net/url"

)

func main() {
	u, _ := url.Parse("http://www.baidu.com/index.html")
	fmt.Println(u.Host)

}
