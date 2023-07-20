package main

import "fmt"


func ordinalNumber(n int) string {
  if n >= 11 && n <= 13 {
    return fmt.Sprintf("%dth", n)
  }

  switch n % 10 {
  case 1: return fmt.Sprintf("%dst", n)
  case 2: return fmt.Sprintf("%dnd", n)
  case 3: return fmt.Sprintf("%drd", n)
  default: return fmt.Sprintf("%dth", n)
  }
}

func max(a, b int) int {
  if a > b { return a }
  return b
}
