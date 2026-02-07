package suite

// Suite represents the top-level eval suite structure
type Suite struct {
	Name       string            `yaml:"suite_name"`
	Capability string            `yaml:"capability"`
	Docs       []Doc             `yaml:"docs"`
	Schemas    map[string]Schema `yaml:"schemas"`
	Tools      []Tool            `yaml:"tools"`
	Grounding  GroundingConfig   `yaml:"grounding"`
	Cases      []Case            `yaml:"cases"`
}

// Doc represents a documentation document with chunks
type Doc struct {
	ID     string  `yaml:"id"`
	Title  string  `yaml:"title"`
	Chunks []Chunk `yaml:"chunks"`
}

// Chunk represents a documentation chunk
type Chunk struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
}

// Schema represents a JSON schema definition
type Schema map[string]interface{}

// Tool represents a tool definition
type Tool struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Args        map[string]interface{} `yaml:"args"`
	Returns     map[string]interface{} `yaml:"returns"`
}

// GroundingConfig represents grounding configuration
type GroundingConfig struct {
	CitationFormat  string   `yaml:"citation_format"`
	ValidDocIDs     []string `yaml:"valid_doc_ids"`
	ValidChunkIDs   []string `yaml:"valid_chunk_ids"`
	CitationPattern string   `yaml:"citation_pattern"`
}

// Case represents a test case
type Case struct {
	ID         string      `yaml:"id"`
	Prompt     string      `yaml:"prompt"`
	Assertions []Assertion `yaml:"assertions"`
}

// Assertion represents a test assertion
type Assertion struct {
	Type     string      `yaml:"type"`
	Expected interface{} `yaml:"expected"`
	Weight   float64     `yaml:"weight,omitempty"`
}

// ValidatorType represents the type of validator to use
type ValidatorType string

const (
	ValidatorContains   ValidatorType = "contains"
	ValidatorRegex      ValidatorType = "regex"
	ValidatorJSONSchema ValidatorType = "json_schema"
	ValidatorGrounding  ValidatorType = "grounding"
)

// CaseType represents the category of a test case
type CaseType string

const (
	CaseTypeGrounding CaseType = "grounding"
	CaseTypeSchema    CaseType = "schema"
	CaseTypeTool      CaseType = "tool"
)

// GetCaseType determines the case type from the case ID prefix
func GetCaseType(caseID string) CaseType {
	if len(caseID) >= 9 && caseID[:9] == "grounding" {
		return CaseTypeGrounding
	}
	if len(caseID) >= 6 && caseID[:6] == "schema" {
		return CaseTypeSchema
	}
	if len(caseID) >= 4 && caseID[:4] == "tool" {
		return CaseTypeTool
	}
	return CaseTypeGrounding // default
}

// Result represents the result of running a test case
type Result struct {
	ID             string   `json:"id"`
	Status         Status   `json:"status"`
	Validator      string   `json:"validator"`
	DurationMS     int64    `json:"duration_ms"`
	FailureReasons []string `json:"failure_reasons,omitempty"`
	Metrics        *Metrics `json:"metrics,omitempty"`
}

// Status represents the status of a test case
type Status string

const (
	StatusPass  Status = "PASS"
	StatusFail  Status = "FAIL"
	StatusError Status = "ERROR"
	StatusSkip  Status = "SKIP"
)

// Metrics represents optional performance metrics
type Metrics struct {
	Tokens  *int `json:"tokens,omitempty"`
	Latency *int `json:"latency,omitempty"`
	Cost    *int `json:"cost,omitempty"`
}
