package utils

func TimeOnly(dt string) string {
	// dt format: "2006-01-02 15:04"
	if len(dt) >= 16 {
		return dt[11:16]
	}
	return dt
}
