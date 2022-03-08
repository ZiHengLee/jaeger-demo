package main

import (
	"fmt"
	"jaeger_demo/tracing"
	"net/http"
	"time"
)

func main()  {
	fmt.Println("===service A start===")
	_ = tracing.InitTracer("service-A", "127.0.0.1:6831")

	ListenHTTPB()
}

func ListenHTTPB() {
	http.HandleFunc("/serviceA/api/test", func(w http.ResponseWriter, r *http.Request) {
		span, traceId, _ := tracing.StartSpan(r.RequestURI, "", true)
		tracing.SpanSetTag(span, "jaeger-demo", 1)
		print(traceId)
		// 查询 mysql db 数据
		mysqlQuerySpan, _, _ := tracing.StartSpan("mysql查询", traceId, false)
		time.Sleep(123*time.Millisecond)
		tracing.FinishSpan(mysqlQuerySpan)
		// 查询 mongo db 数据
		mongoQuerySpan, _, _ := tracing.StartSpan("Mongo查询", traceId, false)
		time.Sleep(345*time.Millisecond)
		tracing.FinishSpan(mongoQuerySpan)

		// 请求服务 B
		callServiceBSpan, callServiceBSpanId, _ := tracing.StartSpan("HTTP GET : serviceB", traceId, false)
		req, _ := http.NewRequest("GET", "http://localhost:9992/serviceB/api/test", nil)
		req.Header.Add("traceid", callServiceBSpanId) // 把 traceid 通过 header 传给服务B
		http.DefaultClient.Do(req)
		tracing.FinishSpan(callServiceBSpan)

		tracing.FinishSpan(span)
		w.Write([]byte("serviceA done"))
	})
	fmt.Println(http.ListenAndServe(":9991", nil))
}