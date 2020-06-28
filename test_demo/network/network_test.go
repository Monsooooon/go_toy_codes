package network

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handleError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("failed", err)
	}
}

// 先在本地开启http在某个path上的监听（在这个path上注册了要测试的Handler
// 然后通过http.Get发送请求，最后对请求的Response内容进行读取，并判断和预计是否相同
func TestConn(t *testing.T) {
	// 在本地的端口上进行监听
	l, err := net.Listen("tcp", "127.0.0.1:0")
	handleError(t, err)
	defer l.Close()

	// 设置http的监听路径和处理函数
	http.HandleFunc("/hello", helloHandler)
	go http.Serve(l, nil) // 需要传入在随机端口上创建的listner

	// 通过http.Get方法对listener的path发送请求
	resp, err := http.Get("http://" + l.Addr().String() + "/hello") // l.Addr().String() 包含了随机的端口
	handleError(t, err)
	defer resp.Body.Close()

	// 把回复的resp.body（一个io.Reader）内容读取出来，放到body（[]byte）中去
	body, err := ioutil.ReadAll(resp.Body)
	handleError(t, err)

	expect := "Ni Hao!"
	if string(body) != expect {
		t.Fatalf("expect %s, but got %s", expect, string(body))
	}
}

// 通过创建虚假的http.Request和供Handler函数写入的NewRecorder()
// 来判断我们的Handler有没有错误。
func TestConnMock(t *testing.T) {
	// 创建一个虚假的htt req
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	// 创建一个ResponseWrite（虽然实际上是NewRecorder()）得到的io.Writer，并交给handler函数去写
	w := httptest.NewRecorder()
	helloHandler(w, req)
	body, _ := ioutil.ReadAll(w.Result().Body)

	expect := "Ni Hao!"
	if string(body) != expect {
		t.Fatalf("expect %s, but got %s", expect, string(body))
	}
}
