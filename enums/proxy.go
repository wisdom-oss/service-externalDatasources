package enums

type ProxyMode uint8

const (
	PROXY_SIMPLE ProxyMode = 1 << iota
	PROXY_TRANSFORM_REQUEST
	PROXY_TRANSFORM_RESPONSE
)
