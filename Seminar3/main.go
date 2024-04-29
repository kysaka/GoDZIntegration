/*package main

import (
    "fmt"
    "strconv"
    "sync"
)

func main() {
    // Создаем каналы для передачи данных между горутинами
    inputChan := make(chan string)
    squareChan := make(chan int)

    // Создаем WaitGroup для синхронизации горутин
    var wg sync.WaitGroup
    wg.Add(2)

    // Горутина для обработки квадрата числа
    go func() {
        defer wg.Done()
        for {
            numStr := <-inputChan
            if numStr == "стоп" {
                return
            }
            num, err := strconv.Atoi(numStr)
            if err != nil {
                fmt.Println("Ошибка:", err)
                continue
            }
            square := num * num
            fmt.Println("Квадрат:", square)
            squareChan <- square
        }
    }()

    // Горутина для обработки произведения
    go func() {
        defer wg.Done()
        for {
            square := <-squareChan
            if square == 0 {
                return
            }
            product := square * 2
            fmt.Println("Произведение:", product)
        }
    }()

    // Чтение чисел из стандартного ввода
    for {
        var input string
        fmt.Scanln(&input)
        if input == "стоп" {
            close(inputChan)
            break
        }
        inputChan <- input
    }

    // Ждем завершения всех горутин
    wg.Wait()
    fmt.Println("Программа завершена.")
}*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Создаем канал для обработки сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Канал для остановки горутины вывода квадратов чисел
	stopChan := make(chan struct{})

	// Горутина для вывода квадратов натуральных чисел
	go func() {
		num := 1
		for {
			select {
			case <-sigChan:
				fmt.Println("Получен сигнал завершения ^C")
				close(stopChan)
				return
			case <-stopChan:
				fmt.Println("Горутина завершена")
				return
			default:
				fmt.Println("Квадрат числа:", num*num)
				num++
			}
		}
	}()

	// Ожидаем сигнала завершения
	<-sigChan
}