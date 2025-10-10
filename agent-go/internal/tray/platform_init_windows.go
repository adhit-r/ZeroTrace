// +build windows

package tray

// GetPlatformOperations returns the Windows-specific implementation
func GetPlatformOperations() PlatformOperations {
	return &WindowsOperations{}
}
