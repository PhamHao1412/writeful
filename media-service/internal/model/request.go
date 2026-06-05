package model

type TransformRequest struct {
	ID        string `json:"id"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	X         int    `json:"x,omitempty"` // crop offset
	Y         int    `json:"y,omitempty"`
	Angle     int    `json:"angle,omitempty"`     // rotation
	FlipAxis  string `json:"flip_axis,omitempty"` // "horizontal" / "vertical"
	Watermark string `json:"watermark,omitempty"` // overlay image public_id
	Quality   string `json:"quality,omitempty"`   // "auto", "60", ...
	Filter    string `json:"filter,omitempty"`
	Format    string `json:"format,omitempty"`
}

type UploadParams struct {
	Id     string
	Type   string
	Folder string
}

type GetListRequest struct {
	Limit      int    `form:"limit" json:"limit,omitempty"`
	Offset     int    `form:"offset" json:"offset,omitempty"`
	Sort       string `form:"sort" json:"sort,omitempty"`
	Type       string `form:"type" json:"type,omitempty"`
	Provider   string `form:"provider" json:"provider,omitempty"`
	UploadedBy string `form:"uploaded_by" json:"uploaded_by,omitempty"`
}
