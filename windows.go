package nirilof

type NumericalPair []float64

// Layout represents the layout of a window in Niri
type Layout struct {
	TileSize              NumericalPair `json:"tile_size"`
	ScrollingPos          NumericalPair `json:"pos_in_scrolling_layout"`
	WindowSize            NumericalPair `json:"window_size"`
	WindowOffsetTile      NumericalPair `json:"window_offset_in_tile"`
	WorkspaceViewPosition NumericalPair `json:"tile_pos_in_workspace_view"`
}

// Window represents the properties of a window in Niri
type Window struct {
	Focused     bool   `json:"is_focused"`
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AppID       string `json:"app_id"`
	IsFloating  bool   `json:"is_floating"`
	PID         int    `json:"pid"`
	WorkspaceID int    `json:"workspace_id"`
	Layout      Layout `json:"layout"`
}
