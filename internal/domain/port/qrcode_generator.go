package port

type QRCodeGenerator interface {
	GeneratePNG(content string) ([]byte, error)
}
