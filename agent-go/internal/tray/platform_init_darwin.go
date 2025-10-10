// +build darwin

package tray

// GetPlatformOperations returns the macOS-specific implementation
func GetPlatformOperations() PlatformOperations {
	return &DarwinOperations{}
}
