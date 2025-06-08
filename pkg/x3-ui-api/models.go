package x3uiapi

type WrapAPI[T interface{}] struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
	Object  T      `json:"obj"`
}

type (
	ClientStats struct {
		ID         int32  `json:"id"`
		InboundID  int64  `json:"inboundId"`
		Enable     bool   `json:"enable"`
		Email      string `json:"email"`
		Up         int64  `json:"up"`
		Down       int64  `json:"down"`
		ExpiryTime int64  `json:"expiryTime"`
		Total      int64  `json:"total"`
		Reset      int64  `json:"reset"`
	}

	Inbound struct {
		ID             int32         `json:"id"`
		Up             int64         `json:"up"`
		Down           int64         `json:"down"`
		Total          int32         `json:"total"`
		Remark         string        `json:"remark"`
		Enable         bool          `json:"enable"`
		ExpiryTime     int32         `json:"expiryTime"`
		ClientsStats   []ClientStats `json:"clientStats"`
		Listen         string        `json:"listen"`
		Port           int32         `json:"port"`
		Protocol       string        `json:"protocol"`
		Settings       string        `json:"settings"`
		StreamSettings string        `json:"streamSettings"`
		Tag            string        `json:"tag"`
		Sniffing       string        `json:"sniffing"`
		Allocate       string        `json:"allocate"`
	}
)

type Online []string

type (
	ClientTraffic struct {
		ID         uint64 `json:"id"`
		InboundID  uint64 `json:"inboundId"`
		Enable     bool   `json:"enable"`
		Email      string `json:"email"`
		Up         uint64 `json:"up"`
		Down       uint64 `json:"down"`
		ExpiryTime uint64 `json:"expiryTime"`
		Total      uint64 `json:"total"`
		Reset      uint64 `json:"reset"`
	}

	InboundTraffic struct {
		IsInbound  bool   `json:"IsInbound"`
		IsOutbound bool   `json:"IsOutbound"`
		Tag        string `json:"Tag"`
		Up         uint64 `json:"Up"`
		Down       uint64 `json:"Down"`
	}

	TrafficUpdates struct {
		Client  []ClientTraffic  `json:"clientTraffics"`
		Inbound []InboundTraffic `json:"inboundTraffics"`
	}
)

func (ctf ClientTraffic) DownTraffic() float64  { return float64(ctf.Down) }
func (ctf ClientTraffic) UpTraffic() float64    { return float64(ctf.Up) }
func (ctf ClientTraffic) TotalTraffic() float64 { return float64(ctf.Total) }
func (ctf ClientTraffic) EmailString() string   { return ctf.Email }

func (itf InboundTraffic) DownTraffic() float64 { return float64(itf.Down) }
func (itf InboundTraffic) UpTraffic() float64   { return float64(itf.Up) }
func (itf InboundTraffic) TagString() string    { return itf.Tag }
