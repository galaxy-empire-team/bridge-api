package gameconfighandlers

type ErrorResponse struct {
	Err string `json:"err"`
}

type VersionRequest struct {
	Version uint16 `json:"version,omitempty"`
}

type GameConfig struct {
	Version uint16 `json:"version"`
	Config  []byte `json:"config,omitempty"`
}
