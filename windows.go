package nirilof

// Size is used for window sizes
type Size struct {
	X int
	Y int
}

// Position is used for window position
// within column and workspaces
type Position struct {
	Column int
	Tile   int
}

// FloatingPosition is used for position
// within a workspace of a floating window
type FloatingPosition struct {
	X int
	Y int
}

// Layout represents the layout of a window in Niri
type Layout struct {
	TileSize              Size
	ScrollingPos          Position
	WindowSize            Size
	WindowOffsetTile      Size
	WorkspaceViewPosition FloatingPosition
}

// Window represents the properties of a window in Niri
type Window struct {
	Focused     bool
	ID          int
	Title       string
	AppID       string
	IsFloating  bool
	PID         int
	WorkspaceID int
	Layout      Layout
}
