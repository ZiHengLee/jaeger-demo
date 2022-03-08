package main

import (
	"jaeger_demo/tracing"
	"net/http"
	"fmt"
	"time"
)

func main()  {
	fmt.Println("===service B start===")
	tracing.InitTracer("service-B", "127.0.0.1:6831")

	ListenHTTP()
}

func ListenHTTP()  {
	http.HandleFunc("/serviceB/api/test", func(w http.ResponseWriter, r *http.Request) {
		span, traceId, _:= tracing.StartSpan(r.RequestURI, r.Header.Get("traceid"), false)

		go callServiceC(traceId)

		time.Sleep(300*time.Millisecond)
		tracing.FinishSpan(span)
		w.Write([]byte("serviceB done"))
	})
	fmt.Println(http.ListenAndServe(":9992", nil))
}

func callServiceC(traceId string)  {
	req, _ := http.NewRequest("GET", "http://localhost:9993/serviceC/api/test?traceid="+traceId, nil)
	http.DefaultClient.Do(req)
}