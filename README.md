# tech-detecter
> 不喜欢httpx json格式的指纹,自己二开了一个Web指纹识别功能,替换它原来的指纹识别模块.
欢迎大家pr共享指纹

## Usage
匹配的信息
```azure
title
header
server
cert
```

```azure
package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func main(t *testing.T) {
	resp, err := http.DefaultClient.Get("https://stg-data-in.ads.heytapmobi.com/")
	if err != nil {
		log.Fatal(err)
	}
	tech := TechDetecter{}
	err = tech.Init("/Users/wing/RedTeamWing/Wing/02-WingCoding/GoWing/tech-detecter/rules/")
	if err != nil {
		log.Fatal(err)
	}
	result, err := tech.Detect(resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
```
## 测试结果
![img.png](img/img.png)

## 参考项目
- https://github.com/jweny/pocassist
- https://github.com/hakuQAQ/Holmes