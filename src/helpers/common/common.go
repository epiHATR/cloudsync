package common

// Sort the list of string by its length
func SortStrings(arr []string) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// Check if input arrays is the same or not.
// Return true if arr1 = arr2.
// Return false if arr1 is not equal arr2
func IsSameArray(arr1, arr2 []string) bool {
	// Check if the lengths of the arrays are equal
	if len(arr1) != len(arr2) {
		return false
	}

	// Sort both arrays to ensure consistent order
	// before comparing the elements
	sortedArr1 := make([]string, len(arr1))
	copy(sortedArr1, arr1)
	sortedArr2 := make([]string, len(arr2))
	copy(sortedArr2, arr2)

	SortStrings(sortedArr1)
	SortStrings(sortedArr2)

	// Compare each element of the sorted arrays
	for i := range sortedArr1 {
		if sortedArr1[i] != sortedArr2[i] {
			return false
		}
	}
	return true
}

// Find the shortest array includes the input array.
// Return whole result or exclude input from result
func GetShortestArray(input []string, excludeInput bool, arrays ...[]string) []string {
	var matchingArray []string
	matchingScore := -1

	for _, arr := range arrays {
		score := 0
		for _, item := range input {
			if IsItemExistInArray(arr, item) {
				score++
			}
		}

		if score > matchingScore {
			matchingArray = arr
			matchingScore = score
		} else if score == matchingScore && len(arr) < len(matchingArray) {
			matchingArray = arr
			matchingScore = score
		}
	}

	if excludeInput {
		// Filter out input elements from the matchingArray
		filteredArray := make([]string, 0, len(matchingArray))
		for _, item := range matchingArray {
			if !IsItemExistInArray(input, item) {
				filteredArray = append(filteredArray, item)
			}
		}
		return filteredArray
	}

	return matchingArray
}

// Check if a string exists in arr of strings
func IsItemExistInArray(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}
