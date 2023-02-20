package constants

type HintCode int

const (
	BadRequestError     HintCode = 1
	InternalServerError HintCode = 2
)

var ErrorMessages = map[HintCode]string{
	BadRequestError:     "Could not process payload events",
	InternalServerError: "System malfunction",
}
