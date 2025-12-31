package messageprocessor

func (m *MessageProcessor) EscapeMarkdown(s string) string {
	return m.sender.EscapeMarkdown(s)
}
