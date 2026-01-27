package am

const (
	ReplyHeaderPrefix  = "REPLY_"
	ReplyNameHeader    = ReplyHeaderPrefix + "NAME"
	ReplyOutcomeHeader = ReplyHeaderPrefix + "OUTCOME"

	FailureReply = "am.Failure"
	SuccessReply = "am.Success"

	OutcomeSuccess = "SUCCESS"
	OutcomeFailure = "FAILURE"
)
