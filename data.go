package main

type HDFSRecord struct {
	Metadata map[string]interface{} `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
}
