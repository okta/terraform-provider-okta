package sdk

type ProtocolSettings struct {
	NameFormat            string `json:"nameFormat,omitempty"`
	HonorPersistentNameId bool   `json:"honorPersistentNameId"`
	ParticipateSLO        bool   `json:"participateSlo,omitempty"`
}
