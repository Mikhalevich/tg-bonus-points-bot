package port

type VerificationCodeGenerator interface {
	Generate() string
}
