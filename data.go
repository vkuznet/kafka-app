package main

type Metadata struct {
	Hostname   string  `json:"hostname"`
	TypePrefix string  `json:"type_prefix"`
	Producer   string  `json:"producer"`
	ID         string  `json:"_id"`
	Type       string  `json:"type"`
	Version    string  `json:"version"`
	Timestamp  float64 `json:"timestamp"`
}

type CMSSWRecord struct {
	Metadata Metadata  `json:"metadata"`
	Data     CMSSW_POP `json:"data"`
}

type CMSSW_POP struct {
	SiteName               string  `json:"site_name"`
	Fallback               bool    `json:"fallback"`
	UserDN                 string  `json:"user_dn"`
	ClientHost             string  `json:"client_host"`
	ClientDomain           string  `json:"client_domain"`
	ServerHost             string  `json:"server_host"`
	ServerDomain           string  `json:"server_domain"`
	UniqueID               string  `json:"unique_id"`
	FileLFN                string  `json:"file_lfn"`
	AppInfo                string  `json:"app_info"`
	FileSize               float64 `json:""file_size`
	ReadSingleSigma        float64 `json:"read_single_sigma"`
	ReadSingleAverage      float64 `json:"read_single_average"`
	ReadVectorSigma        float64 `json:"read_vector_sigma"`
	ReadVectorAverage      float64 `json:"read_vector_average"`
	ReadVectorCountSigma   float64 `json:"read_vector_count_sigma"`
	ReadVectorCountAverage float64 `json:"read_vector_count_average"`
	ReadBytes              float64 `json:"read_bytes"`
	ReadBytesAtClose       float64 `json:"read_bytes_at_close"`
	ReadSingleOperations   float64 `json:"read_single_operations"`
	ReadSingleBytes        float64 `json:"read_single_bytes"`
	ReadVectorOperations   float64 `json:"read_vector_operations"`
	ReadVectorBytes        float64 `json:"read_vector_bytes"`
	StartTime              float64 `json:"start_time"`
	EndTime                float64 `json:"end_time"`
}
