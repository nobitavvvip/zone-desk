package model

type DesktopSettings struct {
	Wallpaper string  `json:"wallpaper"`
	Blur      float64 `json:"blur"`
	Mask      float64 `json:"mask"`
	Accent    string  `json:"accent"`
}
