package http_helper

type Helper struct {
	sentryDSN string
}

func New(sentryDSN string) *Helper {
	return &Helper{
		sentryDSN: sentryDSN,
	}
}
