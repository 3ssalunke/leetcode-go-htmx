package util

func SanitizeDockerLogs(logSlice []byte, sequence []byte) []byte {
	var updatedSlice []byte

	for {
		index := -1
		for i := 0; i < len(logSlice); i++ {
			if logSlice[i] == sequence[0] {
				match := true
				for j := 0; j < len(sequence); j++ {
					if i+j >= len(logSlice) || logSlice[i+j] != sequence[j] {
						match = false
						break
					}
				}
				if match {
					index = i
					break
				}
			}
		}
		if index == -1 {
			break
		}
		before := logSlice[:index]
		after := logSlice[index+len(sequence):]
		updatedSlice = append(updatedSlice, before...)
		logSlice = after
	}

	updatedSlice = append(updatedSlice, logSlice...)
	return updatedSlice
}
