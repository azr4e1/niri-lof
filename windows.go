package nirilof

type Size struct {
	X int
	Y int
}

type Position struct {
	Column int
	Tile   int
}

type FloatingPosition struct {
	X int
	Y int
}

type Layout struct {
	TileSize              Size
	ScrollingPos          Position
	WindowSize            Size
	WindowOffsetTile      Size
	WorkspaceViewPosition FloatingPosition
}

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
