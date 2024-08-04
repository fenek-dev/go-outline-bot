package utils

func ConvertMBtoBytes(mb float32) uint64 {
	return uint64(mb * 1000 * 1000)
}

func ConvertBytesToMB(bytes uint64) uint64 {
	return bytes / 1000 / 1000
}

func ConvertGBtoBytes(gb float32) uint64 {
	return uint64(gb * 1000 * 1000 * 1000)
}

func ConvertBytesToGB(bytes uint64) uint64 {
	return bytes / 1000 / 1000 / 1000
}
