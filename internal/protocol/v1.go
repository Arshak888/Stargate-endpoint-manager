package protocol

const Version = 1

type Request struct {
	ProtocolVersion int         `json:"protocol_version"`
	RequestID       string      `json:"request_id"`
	EndpointID      string      `json:"endpoint_id"`
	Operation       string      `json:"operation"`
	SentAt          string      `json:"sent_at"`
	Payload         interface{} `json:"payload"`
}

type Response struct {
	ProtocolVersion int         `json:"protocol_version"`
	RequestID       string      `json:"request_id"`
	OperationID     string      `json:"operation_id,omitempty"`
	Accepted        bool        `json:"accepted"`
	Status          string      `json:"status"`
	ObservedAt      string      `json:"observed_at"`
	Result          interface{} `json:"result,omitempty"`
	Error           *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}
