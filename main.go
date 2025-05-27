// примитивная змейка на GOlang. Для тренировки.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type Point struct {
	x             int
	y             int
	directionFrom rune
	directionTo   rune
	symbol        string
}

var Snake []Point
var Food = Point{0, 0, 'c', 'c', "○ "}
var n, m int = 23, 23
var area [23][23]string
var score, speed int
var wg sync.WaitGroup
var stop bool
var direction, directionNext, directionPast rune = 'r', 'r', 'l'
var gameText string = " - use the wasd keys to control, and press esc to exit."
var endGameText string = "- THE GAME IS OVER. I hope you enjoyed it."

func clearConsole() {
	switch runtime.GOOS {
	case "linux", "darwin": // Linux и MacOS
		fmt.Print("\033[H\033[2J")
	case "windows": // Windows (не всегда поддерживается в стандартном терминале)
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Очистка консоли не поддерживается на этой ОС")
	}
	return
}

func spawnFood() (int, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	y := r.Intn(n-3) + 1
	x := r.Intn(m-3) + 1
	return y, x
}

func firstDrawArea() {
	Snake = append(Snake, Point{7, 7, 'l', 'r', "──"})
	Snake = append(Snake, Point{8, 7, 'l', 'r', "──"})
	Snake = append(Snake, Point{9, 7, 'l', 'r', "▷ "})
	for i := 0; i < n; i++ {
		area[i][0], area[i][m-1] = "██ ", "██"
	}
	for j := 1; j < m-1; j++ {
		area[0][j], area[n-1][j] = "██", "██"
	}
	area[0][0], area[n-1][0] = "███", "███"
	for i := 1; i < n-1; i++ {
		for j := 1; j < m-1; j++ {
			area[i][j] = "  "
		}
	}

	for _, v := range Snake {
		area[v.y][v.x] = v.symbol
	}

	for true {
		y, x := spawnFood()
		if area[y][x] == "  " {
			area[y][x] = Food.symbol
			Food.x = x
			Food.y = y
			break
		}
	}
	return
}

func drawArea() {
	clearConsole()
	fmt.Println()
	fmt.Println("SNAKE GAME | made by Zigrik")
	fmt.Println("Score :", score, gameText)
	for i := 0; i < n; i++ {
		fmt.Println()
		for j := 0; j < m; j++ {
			fmt.Print(area[i][j])
		}
	}
	fmt.Println()
	return
}

func steps() {
	defer wg.Done()
	var xNext, yNext int
	for !stop {
		pause := time.Duration(100_000/(250+(speed*5))) * time.Millisecond
		time.Sleep(pause)
		if directionNext == 'u' {
			yNext, xNext, directionPast = Snake[len(Snake)-1].y-1, Snake[len(Snake)-1].x, 'd'
		} else if directionNext == 'd' {
			yNext, xNext, directionPast = Snake[len(Snake)-1].y+1, Snake[len(Snake)-1].x, 'u'
		} else if directionNext == 'l' {
			yNext, xNext, directionPast = Snake[len(Snake)-1].y, Snake[len(Snake)-1].x-1, 'r'
		} else if directionNext == 'r' {
			yNext, xNext, directionPast = Snake[len(Snake)-1].y, Snake[len(Snake)-1].x+1, 'l'
		}
		direction = directionNext
		Snake[len(Snake)-1].directionTo = directionNext

		//убираем хвост, если в этот ход змея не ела
		if area[yNext][xNext] == "  " {
			area[Snake[0].y][Snake[0].x] = "  "
			Snake = Snake[1:]
		}

		//меняем отрисовку шеи
		k := len(Snake) - 1
		if (Snake[k].directionFrom == 'u' && Snake[k].directionTo == 'd') || (Snake[k].directionFrom == 'd' && Snake[k].directionTo == 'u') {
			Snake[k].symbol = "│ "
		} else if (Snake[k].directionFrom == 'l' && Snake[k].directionTo == 'r') || (Snake[k].directionFrom == 'r' && Snake[k].directionTo == 'l') {
			Snake[k].symbol = "──"
		} else if (Snake[k].directionFrom == 'u' && Snake[k].directionTo == 'r') || (Snake[k].directionFrom == 'r' && Snake[k].directionTo == 'u') {
			Snake[k].symbol = "└─"
		} else if (Snake[k].directionFrom == 'd' && Snake[k].directionTo == 'r') || (Snake[k].directionFrom == 'r' && Snake[k].directionTo == 'd') {
			Snake[k].symbol = "┌─"
		} else if (Snake[k].directionFrom == 'l' && Snake[k].directionTo == 'u') || (Snake[k].directionFrom == 'u' && Snake[k].directionTo == 'l') {
			Snake[k].symbol = "┘ "
		} else if (Snake[k].directionFrom == 'l' && Snake[k].directionTo == 'd') || (Snake[k].directionFrom == 'd' && Snake[k].directionTo == 'l') {
			Snake[k].symbol = "┐ "
		}
		area[Snake[k].y][Snake[k].x] = Snake[k].symbol

		//направление головы
		if area[yNext][xNext] == "  " || area[yNext][xNext] == Food.symbol {
			Snake = append(Snake, Point{xNext, yNext, directionPast, directionNext, "  "})
			k = len(Snake) - 1
			if directionNext == 'r' {
				Snake[k].symbol = "▢ "
			} else if directionNext == 'l' {
				Snake[k].symbol = "▢─"
			} else if directionNext == 'u' {
				Snake[k].symbol = "▢ "
			} else if directionNext == 'd' {
				Snake[k].symbol = "▢ "
			}
			if area[yNext][xNext] == Food.symbol {
				score++
				speed++
				for true {
					y, x := spawnFood()
					if area[y][x] == "  " {
						area[y][x] = Food.symbol
						Food.x = x
						Food.y = y
						break
					}
				}
			}
			area[Snake[k].y][Snake[k].x] = Snake[k].symbol
		} else {
			gameText = endGameText
			drawArea()
			stop = true
			return
		}
		drawArea()
	}
	return
}

func clicks() {
	defer wg.Done()
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}
		if key == keyboard.KeyEsc {
			stop = true
			return
		}
		switch {
		case char == 'w' || char == 'W' || char == 'Ц' || char == 'ц':
			if direction != 'd' {
				directionNext = 'u'
			}
		case char == 'a' || char == 'A' || char == 'Ф' || char == 'ф':
			if direction != 'r' {
				directionNext = 'l'
			}
		case char == 's' || char == 'S' || char == 'Ы' || char == 'ы':
			if direction != 'u' {
				directionNext = 'd'
			}
		case char == 'd' || char == 'D' || char == 'В' || char == 'в':
			if direction != 'l' {
				directionNext = 'r'
			}
		}
	}
	return
}

func main() {
	firstDrawArea()
	wg.Add(2)
	go clicks()
	go steps()
	wg.Wait()
}
