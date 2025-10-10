package tray

import _ "embed"

// Icon assets embedded at compile time
// These are placeholder icons - replace with actual ZeroTrace logo icons

//go:embed assets/icon_green.png
var iconGreenData []byte

//go:embed assets/icon_gray.png
var iconGrayData []byte

//go:embed assets/icon_red.png
var iconRedData []byte

//go:embed assets/icon_default.png
var iconDefaultData []byte

// GetGreenIcon returns the green (connected) icon
func GetGreenIcon() []byte {
	if len(iconGreenData) > 0 {
		return iconGreenData
	}
	return getDefaultIcon()
}

// GetGrayIcon returns the gray (disconnected) icon
func GetGrayIcon() []byte {
	if len(iconGrayData) > 0 {
		return iconGrayData
	}
	return getDefaultIcon()
}

// GetRedIcon returns the red (error) icon
func GetRedIcon() []byte {
	if len(iconRedData) > 0 {
		return iconRedData
	}
	return getDefaultIcon()
}

// GetDefaultIcon returns the default icon
func GetDefaultIcon() []byte {
	if len(iconDefaultData) > 0 {
		return iconDefaultData
	}
	return getDefaultIcon()
}

// getDefaultIcon returns a minimal fallback PNG icon (16x16 transparent)
func getDefaultIcon() []byte {
	// This is a minimal 16x16 transparent PNG
	// Replace with actual ZeroTrace logo
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, 0x61, 0x00, 0x00, 0x00,
		0x19, 0x49, 0x44, 0x41, 0x54, 0x38, 0x8D, 0x63, 0x64, 0xC0, 0x00, 0x8C,
		0x0C, 0x0C, 0x0C, 0x8C, 0x6C, 0x0C, 0x0C, 0x0C, 0x0C, 0x40, 0x12, 0x00,
		0x00, 0x3A, 0x80, 0x02, 0x7E, 0x0B, 0x0D, 0xFE, 0xB4, 0x00, 0x00, 0x00,
		0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}
