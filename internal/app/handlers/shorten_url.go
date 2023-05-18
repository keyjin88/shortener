package handlers

func (h *Handler) ShortenURL(c RequestContext) {
	if c.FullPath() == "/" {
		h.shortenURLText(c)
	}
	if c.FullPath() == "/api/shorten" {
		h.shortenURLJSON(c)
	}
}
