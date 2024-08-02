package payment_service

type TransactionResult struct {
	Success    bool                    `json:"success"`
	Data       TransactionResultData   `json:"data"`
	Error      *TransactionResultError `json:"error"`
	ExternalID *string                 `json:"external_id"`
}

type TransactionResultData struct {
	Type        TransactionResultType `json:"type"`
	Message     *string               `json:"message"`
	RedirectURL *string               `json:"redirect_url"`
	Response    interface{}           `json:"response"`
}

type TransactionResultError struct {
	Message string `json:"message"`
}

type TransactionResultType string

const (
	ResultTypeSuccess  TransactionResultType = "success"
	ResultTypeRedirect TransactionResultType = "redirect"
	ResultTypeMessage  TransactionResultType = "message"
	ResultTypeFail     TransactionResultType = "fail"
	ResultTypeCancel   TransactionResultType = "cancel"
)
