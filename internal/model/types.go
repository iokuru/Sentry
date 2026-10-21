package model

import "time"

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type Category string

const (
	CategoryBank           Category = "bank"
	CategoryDelivery       Category = "delivery"
	CategoryUtility        Category = "utility"
	CategoryPrize          Category = "prize"
	CategoryTechSupport    Category = "tech support"
	CategoryRelativeInNeed Category = "relative in need"
	CategoryRomance        Category = "romance"
	CategoryTax            Category = "tax"
	CategoryOther          Category = "other"
)

type IncidentStatus string

const (
	StatusPending      IncidentStatus = "pending"
	StatusAnalyzed     IncidentStatus = "analyzed"
	StatusAcknowledged IncidentStatus = "acknowledged"
	StatusResolved     IncidentStatus = "resolved"
	StatusFailed       IncidentStatus = "failed"
)

type SignalMatch struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}

type HotlineInfo struct {
	Country string `json:"country"`
	Name    string `json:"name"`
	Number  string `json:"number"`
	TelURI  string `json:"tel_uri"`
	URL     string `json:"url"`
}

type SeniorOutput struct {
	HoldOff         string      `json:"hold_off"`
	SafeAlternative string      `json:"safe_alternative"`
	Hotline         HotlineInfo `json:"hotline"`
}

type Incident struct {
	ID            string         `json:"id" dynamodbav:"id"`
	SeniorID      string         `json:"senior_id,omitempty" dynamodbav:"senior_id,omitempty"`
	ImageKey      string         `json:"image_key,omitempty" dynamodbav:"image_key,omitempty"`
	ExtractedText string         `json:"extracted_text,omitempty" dynamodbav:"extracted_text,omitempty"`
	UserContext   string         `json:"user_context,omitempty" dynamodbav:"user_context,omitempty"`
	Score         int            `json:"score" dynamodbav:"score"`
	RiskLevel     RiskLevel      `json:"risk_level" dynamodbav:"risk_level"`
	Category      Category       `json:"category" dynamodbav:"category"`
	Signals       []SignalMatch  `json:"signals" dynamodbav:"signals"`
	Output        SeniorOutput   `json:"output" dynamodbav:"output"`
	Status        IncidentStatus `json:"status" dynamodbav:"status"`
	Country       string         `json:"country" dynamodbav:"country"`
	CreatedAt     string         `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt     string         `json:"updated_at" dynamodbav:"updated_at"`
	TTL           int64          `json:"ttl,omitempty" dynamodbav:"ttl,omitempty"`
}

type CaregiverLink struct {
	CaregiverID string    `json:"caregiver_id" dynamodbav:"caregiver_id"`
	SeniorID    string    `json:"senior_id" dynamodbav:"senior_id"`
	SeniorName  string    `json:"senior_name" dynamodbav:"senior_name"`
	LinkedAt    time.Time `json:"linked_at" dynamodbav:"linked_at"`
}
