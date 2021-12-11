package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/newrelic/nrjmx/gojmx"
	"os"
	"path/filepath"
)

func main() {

	os.Setenv("NR_JMX_TOOL", filepath.Join("..", "bin", "nrjmx"))
	//os.Setenv("NRIA_NRJMX_DEBUG", "true")

	config := &gojmx.JMXConfig{
		Hostname:        "localhost",
		Port:            7199,
		//RequestTimoutMs: 10000,
	}

	client, err := gojmx.NewClient(context.Background()).Open(config)
	jmxErr, ok := gojmx.IsJMXError(err)
	if ok {
		fmt.Println("here", jmxErr.Message)
	}

	fmt.Println(err)


	result, err := client.GetMBeanAttrs("java.lang:type=Memory", "HeapMemoryUsage")

	j2, ok := gojmx.IsJMXError(err)
	if ok {
		fmt.Println("here", j2.Message)
	} else {
		//panic(err)
	}
	data := map[string]interface{} {"test": result[0].GetValue()}
	b, _ := json.Marshal(data)
	fmt.Println(string(b))
	//time.Sleep(1*time.Hour)
}

//
//func main() {
//	ctx := context.Background()
//
//	client, err := gojmx.NewJMXServiceClient(ctx)
//	if err != nil {
//		panic(err)
//	}
//
//	err = client.Connect(ctx, &nrprotocol.JMXConfig{
//		Hostname: "localhost",
//		Port:     7199,
//		UriPath:  "jmxrmi",
//	})
//	if err != nil {
//		panic(err)
//	}
//	defer client.Disconnect(ctx)
//
//	result, err := client.QueryMbean(ctx, "*:*")
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(result)
//}
