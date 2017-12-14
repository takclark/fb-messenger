package messenger

import (
	"encoding/json"
	"fmt"
	"testing"
)

var sampleJson string = `
{"object":"page","entry":[{"id":"323004561521128","time":1513204114307,"messaging":[{"recipient":{"id":"323004561521128"},"timestamp":1513204114307,"sender":{"id":"1764799940231555"},"postback":{"payload":"postback_payload","title":"Push Me"}}]}]}
`

func TestJsonUnmarshalling(t *testing.T) {
	fbEvent := new(IncomingFacebookEvent)
	err := json.Unmarshal([]byte(sampleJson), fbEvent)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%+v\n", fbEvent)
	}

}
