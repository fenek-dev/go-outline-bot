package outline_client

func ConvertMBtoBytes(mb float32) int {
	return int(mb * 1000 * 1000)
}

func ConvertBytesToMB(bytes int) int {
	return bytes / 1000 / 1000
}

func ConvertGBtoBytes(gb float32) int {
	return int(gb * 1000 * 1000 * 1000)
}

func ConvertBytesToGB(bytes int) int {
	return bytes / 1000 / 1000 / 1000
}
