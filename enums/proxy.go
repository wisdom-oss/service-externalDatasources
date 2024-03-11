package enums

type ProxyMode uint8

const (
	SimpleProxy ProxyMode = 1 << iota
	TransformRequest
	TransformResponse
)
