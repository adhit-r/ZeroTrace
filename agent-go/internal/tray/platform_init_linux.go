// +build linux

package tray

// GetPlatformOperations returns the Linux-specific implementation
func GetPlatformOperations() PlatformOperations {
	return &LinuxOperations{}
}
