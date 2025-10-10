// +build !windows,!linux,!darwin

package tray

// GetPlatformOperations returns the fallback implementation for unsupported platforms
func GetPlatformOperations() PlatformOperations {
	return &FallbackOperations{}
}
