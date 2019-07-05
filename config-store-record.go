package cloud66

type ConfigStoreRecords struct {
	Records []*ConfigStoreRecord `json:"records" yaml:"records"`
}

type ConfigStoreRecord struct {
	Key      string            `json:"key" yaml:"key"`
	RawValue string            `json:"raw_value" yaml:"raw_value"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
	Ttl      int               `json:"ttl" yaml:"ttl"`
}
