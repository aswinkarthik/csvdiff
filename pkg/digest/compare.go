package digest

// Compare compares two Digest maps and returns the additions and modification
// keys as arrays.
func Compare(baseDigest, newDigest map[uint64]uint64) (additions []uint64, modifications []uint64) {
	maxSize := len(newDigest)
	additions = make([]uint64, maxSize)
	modifications = make([]uint64, maxSize)

	additionCounter := 0
	modificationCounter := 0
	for k, newVal := range newDigest {
		if oldVal, present := baseDigest[k]; present {
			if newVal != oldVal {
				//Modifications
				modifications[modificationCounter] = k
				modificationCounter++
			}
		} else {
			//Additions
			additions[additionCounter] = k
			additionCounter++
		}
	}
	return additions[:additionCounter], modifications[:modificationCounter]
}
