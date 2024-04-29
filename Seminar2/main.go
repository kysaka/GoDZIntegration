package main

// задача 1

/*import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	file, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	fmt.Println("Введите сообщения (для завершения введите 'exit'):")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			fmt.Println("Программа завершена.")
			return
		}
		// Получаем текущее время в формате строки
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		// Записываем данные в файл
		_, err := fmt.Fprintf(file, "%s %s\n", currentTime, input)
		if err != nil {
			fmt.Println("Ошибка записи в файл:", err)
			return
		}
		fmt.Println("Введите следующее сообщение:")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения ввода:", err)
	}
}*/

// задача 2

/*import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("log.txt")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Ошибка получения информации о файле:", err)
		return
	}

	// Проверяем размер файла, чтобы убедиться, что он не пустой
	if stat.Size() == 0 {
		fmt.Println("Файл пуст.")
		return
	}

	fmt.Println("Содержимое файла:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла:", err)
	}
}*/

// Задача 3

/*import (
	"fmt"
	"os"
)

func main() {
	// Создаем файл с правами доступа только для чтения (0444)
	file, err := os.Create("readonly.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	// Изменяем права доступа к файлу
	if err := os.Chmod("readonly.txt", 0444); err != nil {
		fmt.Println("Ошибка изменения прав доступа к файлу:", err)
		return
	}

	// Пытаемся записать данные в файл
	_, err = fmt.Fprintf(file, "Это сообщение не должно быть записано в файл")
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return
	}

	fmt.Println("Успешно записали в файл.")
}*/

// Задача 4

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	fileName := "log.txt"

	fmt.Println("Введите сообщения (для завершения введите 'exit'):")
	for {
		var input string
		fmt.Scanln(&input)
		if input == "exit" {
			fmt.Println("Программа завершена.")
			return
		}

		currentTime := time.Now().Format("2006-01-02 15:04:05")
		message := fmt.Sprintf("%s %s\n", currentTime, input)

		err := ioutil.WriteFile(fileName, []byte(message), os.ModeAppend)
		if err != nil {
			fmt.Println("Ошибка записи в файл:", err)
			return
		}

		fmt.Println("Введите следующее сообщение:")
	}
}


/*import (
	"fmt"
	"io/ioutil"
)

func main() {
	fileName := "log.txt"

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	if len(content) == 0 {
		fmt.Println("Файл пуст.")
		return
	}

	fmt.Println("Содержимое файла:")
	fmt.Println(string(content))
}*/
