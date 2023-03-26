package main

type ThreatfoxIOC struct {
	ID             string `json:"id,omitempty" validate:"required"`
	Value          string `json:"ioc,omitempty" validate:"required"`
	Type           string `json:"ioc_type,omitempty" validate:"required"`
	ThreatType     string `json:"threat_type,omitempty" validate:"required"`
	ThreatTypeDesc string `json:"threat_type_desc,omitempty" validate:"required"`

	Malware          string   `json:"malware,omitempty" validate:"required"`
	MalwarePrintable string   `json:"malware_printable,omitempty" validate:"required"`
	MalwareAlias     string   `json:"malware_alias,omitempty" validate:"required"`
	MalwareMalpedia  string   `json:"malware_malpedia,omitempty" validate:"required"`
	ConfidenceLevel  uint8    `json:"confidence_level,omitempty" validate:"numeric,gte=0,lte=100"`
	FirstSeen        string   `json:"first_seen,omitempty" validate:"required"`
	LastSeen         string   `json:"last_seen,omitempty" validate:"required"`
	Reference        string   `json:"reference,omitempty" validate:"required"`
	Reporter         string   `json:"reporter,omitempty" validate:"required"`
	Tags             []string `json:"tags,omitempty" validate:"required"`
}

type ThreatFoxResponse struct {
	Query  string         `json:"query_status"`
	Result []ThreatfoxIOC `json:"data"`
}

type Edge struct {
	From string `json:"_from,omitempty"`
	To   string `json:"_to,omitempty"`
}

type IOC struct {
	ID              string `json:"_id"`
	Key             string `json:"_key"`
	Value           string //key
	Type            string
	ThreatType      string
	ThreatTypeDesc  string
	ConfidenceLevel uint8
	FirstSeen       string
	LastSeen        string
}

type Malware struct {
	ID               string `json:"_id"`
	Key              string `json:"_key"`
	Malware          string //key
	MalwarePrintable string
	MalwareAlias     string
	MalwareMalpedia  string
	FirstSeen        string
	LastSeen         string
	Reference        string
	Tags             []string
}

type Reporter struct {
	ID       string `json:"_id"`
	Key      string `json:"_key"`
	Reporter string //key
}

type ThreatSource struct {
	ID    string `json:"_id"`
	Key   string `json:"_key"`
	Value string
}

// type ThreatActor struct {
// }

// type AttackFlow struct {
// }
