package main

import "fmt"

func BubbleSort(array[] int)[]int {
   for i:=0; i< len(array)-1; i++ {
      for j:=0; j < len(array)-i-1; j++ {
         if (array[j] > array[j+1]) {
            array[j], array[j+1] = array[j+1], array[j]
         }
      }
   }
   return array
}

func SelectionSort(array[] int, size int) []int {
   var min_index int
   var temp int
   for i := 0; i < size - 1; i++ {
      min_index = i
      for j := i + 1; j < size; j++ {
         if array[j] < array[min_index] {
            min_index = j
         }
      }
      temp = array[i]
      array[i] = array[min_index]
      array[min_index] = temp
   }
   return array
}

func InsertionSort(arr []int) []int {
    for i := 1; i < len(arr); i++ {
        key := arr[i]
        j := i - 1
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j = j - 1
        }
        arr[j+1] = key
    }
    return arr
}

