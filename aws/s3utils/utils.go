package s3utils

// Common file sizes
const (
	Byte     = 1.0
	Kilobyte = 1024 * Byte
	Megabyte = 1024 * Kilobyte
	Gigabyte = 1024 * Megabyte
	Terabyte = 1024 * Gigabyte
)

func derefStringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
