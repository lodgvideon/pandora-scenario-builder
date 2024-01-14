package scenario

type HttpRequestBuilder struct {
	method string
	uri    string
}

func (b *HttpRequestBuilder) SetURI(uri string) *HttpRequestBuilder {
	b.uri = uri
	return b
}

func HTTP(method string) *HttpRequestBuilder {
	return &HttpRequestBuilder{method: method}
}

func GET() *HttpRequestBuilder {
	return HTTP("GET")
}

func POST() *HttpRequestBuilder {
	return HTTP("POST")
}

func PUT() *HttpRequestBuilder {
	return HTTP("PUT")
}

func DELETE() *HttpRequestBuilder {
	return HTTP("DELETE")
}

func PATCH() *HttpRequestBuilder {
	return HTTP("PATCH")
}
