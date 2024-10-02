package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type RPCResponse struct {
	Version string `json:"jsonrpc"`
	Result  struct {
		Gateways []struct {
			ID    string `json:"ID"`
			IP    string `json:"IP"`
			State string `json:"State"`
		} `json:"Gateways"`
	} `json:"result"`
	Id string `json:"id"`
}

func parse_dr_gw_status(body_response []byte) {
	// Parse the JSON response
	var rpc_response RPCResponse
	err := json.Unmarshal(body_response, &rpc_response)
	if err != nil {
		fmt.Printf("client: could not parse response rpc json: %s\n", err)
		os.Exit(1)
	}

	// Extract all Gateways
	gateways := rpc_response.Result.Gateways

  fmt.Println("# TYPE opensips_dr_gw_status gauge")
	
  for _, element := range gateways {
		var state int = 0
		if element.State == "Active" {
			state = 1
		}
		fmt.Printf("opensips_dr_gw_status{id=\"%s\",ip=\"%s\"} %d\n", element.ID, element.IP, state)
	}
}

func rpc_get(request_url string, method string) []byte {
	// Build request ID
	id := uuid.New()

	// Build request body
	rpc_request := fmt.Sprintf(`{"jsonrpc":"2.0","method":"%s","id":"%s"}`, method, id.String())
	jsonBody := []byte(rpc_request)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(request_url, "application/json", bodyReader)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	return resBody
}

func main() {
	opensips_mi_address := os.Args[1]
	rpc_method := os.Args[2]

	// Perform the RPC request
	resBody := rpc_get(fmt.Sprintf("http://%s:8888/mi", opensips_mi_address), rpc_method)

	if rpc_method == "dr_gw_status" {
		parse_dr_gw_status(resBody)
	} else {
		fmt.Printf("this rpc method is not yet implemented\n")
		os.Exit(1)
	}

}
