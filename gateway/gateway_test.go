package gateway

import (
	"io/ioutil"
	"net/http"
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestSelfRegisterNoUpstreamNoAPI(t *testing.T) {
	client := &http.Client{}

	defer gock.Off()

	gock.New("http://kong:8001").
		Get("/upstreams/user.api.jormugandr.org").
		Reply(404).
		JSON(map[string]string{"message": "Not Found"})

	gock.New("http://kong:8001").
		Get("/apis/user-microservice").
		Reply(404).
		JSON(map[string]string{"message": "Not Found"})

	gock.New("http://kong:8001").
		Post("/upstreams/").
		MatchParam("name", "user.api.jormugandr.org").
		MatchParam("slots", "10").
		Reply(200).
		JSON(map[string]interface{}{
			"id":   "13611da7-703f-44f8-b790-fc1e7bf51b3e",
			"name": "service.v1.xyz",
			"orderlist": []int{
				1,
				2,
				7,
				9,
				6,
				4,
				5,
				10,
				3,
				8,
			},
			"slots":      10,
			"created_at": 1485521710265,
		})

	gock.InterceptClient(client)

	config := &MicroserviceConfig{
		MicroserviceName: "user-microservice",
		MicroservicePort: 8080,
		ServicesMaxSlots: 10,
		UpstreamURL:      "http://localhost:8080",
		VirtualHost:      "user.api.jormugandr.org",
		Weight:           10,
		Hosts:            []string{"localhost", "user.api.jormugandr.org"},
	}
	gateway := NewKongGateway("http://kong:8001", client, config)

	err := gateway.SelfRegister()
	if err != nil {
		panic(err)
	}

	all := gock.GetAll()

	for mock := range all {
		t.Log("MOCK >>>", mock)
	}

	// resp, err := client.Get("http://kong:8080/apis/test")
	// if err != nil {
	// 	panic(err)
	// }
	// printResp(resp, t)
	// resp, err = client.Get("http://kong:8080/apis/user-microservice")
	// if err != nil {
	// 	panic(err)
	// }
	// printResp(resp, t)
}

func printResp(resp *http.Response, t *testing.T) {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	t.Log("STATUS:", resp.Status)
	for k, v := range resp.Header {
		t.Logf("%s: %s", k, v)
	}
	t.Log(string(bytes))
}
