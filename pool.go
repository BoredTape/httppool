package httppool

import (
	"net/http"
	"time"
	"errors"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/url"
	"strings"
)

type Pools struct {
	Size    int
	Timeout time.Duration
	Clients chan *http.Client
}

type Options struct {
	Size    int
	Timeout time.Duration
}

type Respond struct {
	Result *http.Response
	Err    error
}

type TReader struct {
	Result   *transform.Reader
	Response *http.Response
	Err      error
}

type Request struct {
	Header map[string]string
	Url    string
	Method string
	Form   map[string]string
	Proxy  string
}

func (opt *Options) init() {
	if opt.Timeout == 0*time.Second {
		opt.Timeout = 5 * time.Second
	}
	if opt.Size == 0 {
		opt.Size = 3000
	}
}
func NewPools(opt *Options) *Pools {
	opt.init()
	pool := new(Pools)
	pool.Size = opt.Size
	pool.Timeout = opt.Timeout
	pool.SetPools()
	return pool
}

func (p *Pools) SetPools() {
	p.Clients = make(chan *http.Client, p.Size)
	for i := 0; i < p.Size; i++ {
		go func() {
			p.Clients <- &http.Client{Timeout: p.Timeout, }
		}()
	}
}

func (p *Pools) Open(args Request) (*Respond) {
	var request *http.Request
	var resp = new(Respond)
	if args.Method == "" {
		resp.Err = errors.New("url not exist")
		return resp
	}
	if args.Url == "" {
		resp.Err = errors.New("url not exist")
		return resp
	}
	if strings.Contains(args.Url, "http://") == false {
		if strings.Contains(args.Url, "https://") == false {
			args.Url = "http://" + args.Url
		}
	}

	if args.Form == nil {
		request, resp.Err = http.NewRequest(args.Method, args.Url, nil)
	} else {
		var form url.Values
		for key, value := range args.Form {
			temp := []string{
				value,
			}
			form[key] = temp
		}
		request, resp.Err = http.NewRequest(args.Method, args.Url, strings.NewReader(form.Encode()))
	}

	if resp.Err != nil {
		return resp
	}
	client, enough := <-p.Clients
	if enough == false {
		resp.Err = errors.New("not enought client")
		return resp
	}

	if args.Proxy != "" {
		if strings.Contains(args.Proxy, "http://") == false {
			if strings.Contains(args.Proxy, "https://") == false {
				args.Proxy = "http://" + args.Proxy
			}
		}
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(args.Proxy)
		}
		transport := &http.Transport{Proxy: proxy, }
		client.Transport = transport
	}

	if args.Header != nil {
		for key, value := range args.Header {
			request.Header.Add(key, value)
		}
	}
	resp.Result, resp.Err = client.Do(request)
	p.Clients <- client
	return resp
}

func (resp *Respond) Transform() (*TReader) {
	reader := new(TReader)
	reader.Err = resp.Err
	if reader.Err == nil {
		reader.Response = resp.Result
		reader.Result = transform.NewReader(resp.Result.Body, simplifiedchinese.GBK.NewDecoder())
	}
	return reader
}

func (resp *Respond) Resault() ([]byte, error) {
	var body []byte
	var err = resp.Err
	if resp.Err == nil {
		body, err = ioutil.ReadAll(resp.Result.Body)
		defer resp.Result.Body.Close()
		return body, err
	}
	return body, err
}

func (reader *TReader) Resault() ([]byte, error) {
	var body []byte
	var err = reader.Err
	if err == nil {
		body, err = ioutil.ReadAll(reader.Result)
		defer reader.Response.Body.Close()
		return body, err
	}
	return body, err
}
