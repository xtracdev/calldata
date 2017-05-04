package apitimings

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const sumoTimingStruct = `"-9223372036261114666","/vc2coma2078845n/log/xtrac-api/api_rest_transformation.log","vc2coma2078845n","/xapi/DEV/NONPROD","1467400978000","1467400988935","117209063","103886658","201","plain:atp:o:6:l:25:p:yyyy-MM-dd'T'HH:mm:ssZZZZ","","UTF8","time=""2016-07-01T15:22:58-04:00"" level=info msg=""request for /xtrac/api/v1/work-items/search with method POST"" {""Name"":""xtracApi-POST-workItems-search"",""Tags"":{""aud"":""a79fcb28-2621-4973-8a1e-c09a2ab30f79"",""jti"":""02fe5ff1-c242-4d18-ac7e-73166de395df"",""sub"":""XWHRon""},""Duration"":143114959,""time"":""2016-07-01T15:22:58.411471789-04:00"",""TxnId"":""56bfe3d5-5204-ec24-077d-6a8213fdf8a5"",""Contributors"":[{""Name"":""JWT Authentication plugin"",""Duration"":143016086,""Error"":"""",""ServiceCalls"":null},{""Name"":""Whitelist plugin"",""Duration"":142842257,""Error"":"""",""ServiceCalls"":null},{""Name"":""Session Management plugin"",""Duration"":142833386,""Error"":"""",""ServiceCalls"":null},{""Name"":""REST plugin"",""Duration"":138057652,""Error"":"""",""ServiceCalls"":null},{""Name"":""workflow-backend"",""Duration"":133996676,""Error"":"""",""ServiceCalls"":[{""Name"":""Core-WorkItem-Search"",""Endpoint"":""vc2coma2078845n.fmr.com:11000"",""Duration"":133775560,""Error"":""""}]}],""ErrorFree"":true,""Error"":""""}"`

func TestParseSumoJSON(t *testing.T) {
	timingRec, err := NewAPITimingRec(sumoTimingStruct)
	if assert.Nil(t, err) {
		fmt.Printf("%s\n", timingRec.Message)

		var callRecord endToEndTimer

		err := json.Unmarshal([]byte(timingRec.Message), &callRecord)
		assert.Nil(t, err)
		assert.Equal(t, "xtracApi-POST-workItems-search", callRecord.Name)
	}
}

const sumoTimingStructNonAPI = `"-9223372036261114666","/vc2coma2078845n/log/xtrac-api/api_rest_transformation.log","vc2coma2078845n","/xapi/DEV/NONPROD","1467400978000","1467400988935","117209063","103886658","201","plain:atp:o:6:l:25:p:yyyy-MM-dd'T'HH:mm:ssZZZZ","","UTF8","time=""2016-07-01T15:22:58-04:00"" level=info msg=""request for /xtrac/api/v1/work-items/search with method POST"" {}"`

func TestParseNonTimingMessage(t *testing.T) {
	timingRec, err := NewAPITimingRec(sumoTimingStructNonAPI)
	assert.Nil(t, err)

	_, err = timingRec.CallRecord()
	assert.NotNil(t, err)
}
