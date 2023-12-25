package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("read content failed:%v", err)
	}

	doc.Find("div.ant-col h2").Each(func(i int, selection *goquery.Selection) {
		//获取匹配标签中的文本
		title := selection.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}

// 获取网页内容，检测网页字符编码并将文本统一转换为utf8格式
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%d", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)
	uft8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(uft8Reader)
}

func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	if err != nil {
		fmt.Println("fetch error:%v", err)
		return unicode.UTF8
	}

	//检测最多前1024字节的内容来确定HTML文档的编码
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
