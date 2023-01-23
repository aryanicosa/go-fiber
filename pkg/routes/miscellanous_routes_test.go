package routes

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	test := struct {
		route          string // input route
		expectedCode   int
		expectedResult string
	}{
		route:          "/v1/misc/base64encode",
		expectedCode:   201,
		expectedResult: "\"YWRtaW46c2VjcmV0\"",
	}

	reqBodyStr := `admin:secret`

	req := httptest.NewRequest("POST", test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	// Perform the request plain with the AppTest.
	resp, err := AppTest.Test(req, -1) // the -1 disables request latency
	if err != nil {
		log.Fatal("fail to create request")
	}

	var encodedResponse string
	responseBodyBytes, _ := ioutil.ReadAll(resp.Body)
	encodedResponse = string(responseBodyBytes)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, test.expectedResult, encodedResponse)
}
