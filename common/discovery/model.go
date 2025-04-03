package discovery

import (
	"encoding/json"
)

// etcd的元数据结构
type EndpointInfo struct {
	IP       string                 `json:"ip"`
	Port     string                 `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	ed := &EndpointInfo{}
	err := json.Unmarshal(data, ed)
	if err != nil {
		return nil, err
	}
	return ed, nil
}
func (edi *EndpointInfo) Marshal() string {
	data, err := json.Marshal(edi)
	if err != nil {
		panic(err)
	}
	return string(data)
}
