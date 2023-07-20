package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pandasoli/goterm"
)


func main() {
  argv := os.Args

  if len(argv) == 1 {
    fmt.Print(help_page)
    return
  }

  switch argv[1] {
  case "cli":
    if len(argv) > 2 {
      fmt.Print(help_page)
      return
    }

    termios, err := goterm.SetRawMode()
    if err != nil {
      panic(fmt.Errorf("Could not set terminal to raw mode: %s", err))
    }

    defer func() {
      err := goterm.RestoreMode(termios)
      if err != nil {
        panic(fmt.Errorf("Could not restore terminal mode: %s", err))
      }
    }()

    fmt.Print("\033[?1003h\033[?1015h\033[?1006h")
    fmt.Print("\033[?1002h")

    defer fmt.Print("\033[?1000l")
    defer fmt.Print("\033[?1003l\033[?1015l\033[?1006l")

    fmt.Print("\033[?47h")
    defer fmt.Print("\033[?47l")

    goterm.HideCursor()
    defer goterm.ShowCursor()

    rect := Rect { 40, 6, 4, 2 }
    var click_pos Pos
    selected_side := "none"

    var x, y int

    for {
      fmt.Print("\033[2J")
      show_border(rect, selected_side)

      goterm.GoToXY(50, 2)
      fmt.Printf("\033[100m 🌦  \033[0m { \033[31m%d\033[0m, \033[31m%d\033[0m, \033[31m%d\033[0m, \033[31m%d \033[0m}", rect.W, rect.H, rect.X, rect.Y)
      goterm.GoToXY(50, 3)
      fmt.Print(selected_side)

      key, err := goterm.Getch()
      if err != nil {
        panic(fmt.Errorf("Could not get character: %s", err))
      }

      if key == "q" || key == "\n" { break }

      if strings.HasPrefix(key, "\033[<") {
        list := strings.Split(key[3:], ";")

        ev := list[0]
        x, _ = strconv.Atoi(list[1])
        y, _ = strconv.Atoi(list[2][:len(list[2]) - 1])
        kind := list[2][len(list[2]) - 1]

        x--
        y--

        switch ev {
        case "0":
          if kind == 'M' {
            sides := map[string]bool {
              "center": (rect.X < x && x < rect.X + rect.W - 1) && (rect.Y < y && y < rect.Y + rect.H - 1),

              "top": y == rect.Y && (rect.X < x && x < rect.X + rect.W - 1),
              "left": x == rect.X && (rect.Y < y && y < rect.Y + rect.H - 1),
              "bottom": y == rect.Y + rect.H - 1 && (rect.X < x && x < rect.X + rect.W - 1),
              "right": x == rect.X + rect.W - 1 && (rect.Y < y && y < rect.Y + rect.H - 1),

              "top left": y == rect.Y && x == rect.X,
              "top right": y == rect.Y && x == rect.X + rect.W - 1,
              "bottom left": y == rect.Y + rect.H - 1 && x == rect.X,
              "bottom right": y == rect.Y + rect.H - 1 && x == rect.X + rect.W - 1,
            }

            selected_side = "none"

            for side, ok := range sides {
              if ok { selected_side = side; break }
            }

            click_pos.X = x - rect.X
            click_pos.Y = y - rect.Y
          } else if kind == 'm' {
            selected_side = "none"
          }
        case "32":
          var funcs map[string]func(x, y int)
          funcs = map[string]func(x, y int) {
            "center": func(x, y int) {
              new_x := (x + rect.X) - click_pos.X
              new_y := (y + rect.Y) - click_pos.Y

              if new_x >= 0 { rect.X = new_x }
              if new_y >= 0 { rect.Y = new_y }
            },
            "top": func(x, y int) {
              new_y := rect.Y + y
              new_h := rect.H - y

              if new_y >= 0 { rect.Y = new_y }
              if new_h >= 2 { rect.H = new_h }
            },
            "bottom": func(x, y int) {
              new_h := y + 1

              if new_h >= 2 { rect.H = new_h }
            },
            "left": func(x, y int) {
              new_x := rect.X + x
              new_w := rect.W - x

              if new_x >= 0 { rect.X = new_x }
              if new_w >= 2 { rect.W = new_w }
            },
            "right": func(x, y int) {
              new_w := x + 1

              if new_w >= 2 { rect.W = new_w }
            },
            "top left": func(x, y int) {
              new_h := rect.H - y
              new_y := rect.Y + y
              new_w := rect.W - x
              new_x := rect.X + x

              if new_h >= 2 { rect.H = new_h }
              if new_y >= 0 { rect.Y = new_y }
              if new_w >= 2 { rect.W = new_w }
              if new_x >= 0 { rect.X = new_x }
            },
            "top right": func(x, y int) {
              new_h := rect.H - y
              new_y := rect.Y + y
              new_w := x + 1

              if new_h >= 2 { rect.H = new_h }
              if new_y >= 0 { rect.Y = new_y }
              if new_w >= 2 { rect.W = new_w }
            },
            "bottom left": func(x, y int) {
              new_w := rect.W - x
              new_x := (x + rect.X)
              new_h := (y + rect.Y) - 1

              if new_w >= 2 { rect.W = new_w }
              if new_x >= 0 { rect.X = new_x }
              if new_h >= 2 { rect.H = new_h }
            },
            "bottom right": func(x, y int) {
              new_h := y + 1
              new_w := x + 1

              if new_h >= 2 { rect.H = new_h }
              if new_w >= 2 { rect.W = new_w }
            },
          }

          if fn, ok := funcs[selected_side]; ok {
            fn(x - rect.X, y - rect.Y)
          }
        }
      }
    }

  case "show":
    /*
      0: program name
      1: "show"
      2: <width>
      3: <height>
      4: <x>
      5: <y>
    */

    if len(argv) != 6 {
      fmt.Print(help_page)
      return
    }

    w, w_err := strconv.Atoi(argv[2])
    h, h_err := strconv.Atoi(argv[3])
    x, x_err := strconv.Atoi(argv[4])
    y, y_err := strconv.Atoi(argv[5])

    icon_err := "\033[30;43m 🟡️ \033[39;49m"

    if w_err != nil {
      fmt.Printf(" %s Dönüştürülemedi `%s` hedef int (beklenen width).\n\n", icon_err, argv[2])
      return
    }

    if h_err != nil {
      fmt.Printf(" %s Dönüştürülemedi `%s` hedef int (beklenen heigth).\n\n", icon_err, argv[3])
      return
    }

    if x_err != nil {
      fmt.Printf(" %s Dönüştürülemedi `%s` hedef int (beklenen x).\n\n", icon_err, argv[4])
      return
    }

    if y_err != nil {
      fmt.Printf(" %s Dönüştürülemedi `%s` hedef int (beklenen y).\n\n", icon_err, argv[5])
      return
    }

    show_border(Rect { w, h, x, y }, "")
    fmt.Println()
  default:
    fmt.Printf(" \033[41m 🍒️ \033[49m Anlayamadım `%s` seçenek.\n\n", argv[1])
  }
}
