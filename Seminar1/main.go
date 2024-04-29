package main

import "fmt"

/*func mergeSortedArrays(arr1 []int, arr2 []int) []int {
    mergedArray := make([]int, 0, len(arr1)+len(arr2))
    i, j := 0, 0
    for i < len(arr1) && j < len(arr2) {
        if arr1[i] < arr2[j] {
            mergedArray = append(mergedArray, arr1[i])
            i++
        } else {
            mergedArray = append(mergedArray, arr2[j])
            j++
        }
    }
    mergedArray = append(mergedArray, arr1[i:]...)
    mergedArray = append(mergedArray, arr2[j:]...)
    return mergedArray
}

func main() {
    arr1 := []int{1, 3, 5, 7}
    arr2 := []int{2, 4, 6, 8, 9}
    merged := mergeSortedArrays(arr1, arr2)
    fmt.Println(merged) // Output: [1 2 3 4 5 6 7 8 9]
}*/

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func main() {
	arr := []int{6, 3, 9, 5, 2, 8}
	fmt.Println("Before sorting:", arr)
	bubbleSort(arr)
	fmt.Println("After sorting:", arr)
}
