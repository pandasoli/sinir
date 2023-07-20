package main

import (
  "fmt"
  "strings"

  "github.com/pandasoli/goterm"
)


func makeBorder(rect Rect, focus string) {
  color := func(side string, start bool) string {
    if side == focus {
      if start {
        if focus == "center" { return "\033[41m" }
        return "\033[1;31;49m"
      } else {
        if focus == "center" { return "\033[49m" }
        return "\033[0;39;49m"
      }
    }

    return ""
  }

  goterm.GoToXY(rect.X, rect.Y)
  fmt.Print(color("top left", true), "╭", color("top left", false))
  fmt.Print(color("top", true), strings.Repeat("─", rect.W - 2), color("top", false))
  fmt.Print(color("top right", true), "╮", color("top right", false))

  for i := range make([]int, rect.H - 2) {
    goterm.GoToXY(rect.X, rect.Y + 1 + i)
    fmt.Print(
      color("left", true), "│", color("left", false),
      color("center", true), strings.Repeat(" ", rect.W - 2), color("center", false),
      color("right", true), "│", color("right", false),
    )
  }

  goterm.GoToXY(rect.X, rect.Y + rect.H - 1)
  fmt.Print(color("bottom left", true), "╰", color("bottom left", false))
  fmt.Print(color("bottom", true), strings.Repeat("─", rect.W - 2), color("bottom", false))
  fmt.Print(color("bottom right", true), "╯", color("bottom right", false))
}
