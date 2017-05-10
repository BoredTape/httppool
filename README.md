# httppool
http连接池

## 安装

```shell
go get -u github.com/BoredTape/httppool
```

## 导入

```go
import "github.com/BoredTape/httppool"
```

## 快速开始

```go
func GET() {
	p := httppool.NewPools(&apool.Options{})

	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	}

	var request = apool.Request{
		Header: header,
		Url:    "www.google.com",
	}
	result := p.Open(request).Resault()
	fmt.Println(string(result.Body))
}

func POST(){
     p := httppool.NewPools(&apool.Options{})

	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	}
    form := map[string]string{
        "key":"value",
    }
	var request = apool.Request{
		Header: header,
		Url:    "www.google.com",
		Method: "POST",
		Form:   form,
	}
	result := p.Open(request).Resault()
	fmt.Println(string(result.Body))
}

func With_Proxy() {
	p := httppool.NewPools(&apool.Options{})

	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	}

	var request = apool.Request{
		Header: header,
		Url:    "www.google.com",
		Proxy:  "123.123.123.123",
	}
	result := p.Open(request).Resault()
	fmt.Println(string(result.Body))
}
```