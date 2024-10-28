package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Файл для хранения состояния
const dbFile = "db.json"

// Хранилище для списков
var store = make(map[string][]int)

// LoadFromDisk загружает состояние из файла на диск
func LoadFromDisk() error {
	file, err := os.Open(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файл не существует, загружать нечего
		}
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &store)
}

// SaveToDisk сохраняет состояние на диск
func SaveToDisk() error {
	bytes, err := json.Marshal(store)
	if err != nil {
		return err
	}

	return os.WriteFile(dbFile, bytes, 0644)
}

// LPUSH добавляет элементы в начало списка
func LPUSH(key string, elements ...int) int {
	if _, exists := store[key]; !exists {
		store[key] = []int{}
	}
	store[key] = append(elements, store[key]...)
	return len(store[key])
}

// RPUSH добавляет элементы в конец списка
func RPUSH(key string, elements ...int) int {
	if _, exists := store[key]; !exists {
		store[key] = []int{}
	}
	store[key] = append(store[key], elements...)
	return len(store[key])
}

// RADDTOSET добавляет элементы в конец списка, если их еще нет в списке
func RADDTOSET(key string, elements ...int) int {
	if _, exists := store[key]; !exists {
		store[key] = []int{}
	}
	existingElements := make(map[int]bool)
	for _, el := range store[key] {
		existingElements[el] = true
	}
	for _, el := range elements {
		if !existingElements[el] {
			store[key] = append(store[key], el)
			existingElements[el] = true
		}
	}
	return len(store[key])
}

// LPOP удаляет и возвращает элементы с начала списка
func LPOP(key string, count ...int) ([]int, error) {
	if _, exists := store[key]; !exists || len(store[key]) == 0 {
		return nil, errors.New("list is empty or does not exist")
	}
	start, end := parseCountIndexes(count, len(store[key]))

	removedElements := store[key][start:end]
	store[key] = store[key][end:]
	return removedElements, nil
}

// RPOP удаляет и возвращает элементы с конца списка
func RPOP(key string, count ...int) ([]int, error) {
	if _, exists := store[key]; !exists || len(store[key]) == 0 {
		return nil, errors.New("list is empty or does not exist")
	}
	start, end := parseCountIndexes(count, len(store[key]))

	removedElements := store[key][len(store[key])-end : len(store[key])-start]
	store[key] = store[key][:len(store[key])-end]
	return removedElements, nil
}

// LSET устанавливает значение элемента с индексом
func LSET(key string, index int, element int) error {
	if _, exists := store[key]; !exists || index >= len(store[key]) || index < -len(store[key]) {
		return errors.New("index out of range")
	}
	if index < 0 {
		index = len(store[key]) + index
	}
	store[key][index] = element
	return nil
}

// LGET получает значение элемента с индексом
func LGET(key string, index int) (int, error) {
	if _, exists := store[key]; !exists || index >= len(store[key]) || index < -len(store[key]) {
		return 0, errors.New("index out of range")
	}
	if index < 0 {
		index = len(store[key]) + index
	}
	return store[key][index], nil
}

// Вспомогательная функция для обработки индексов count в командах LPOP и RPOP
func parseCountIndexes(count []int, listLength int) (int, int) {
	if len(count) == 1 {
		if count[0] > listLength {
			return 0, listLength
		}
		return 0, count[0]
	}
	start, end := count[0], count[1]
	if start < 0 {
		start = listLength + start
	}
	if end < 0 {
		end = listLength + end
	}
	if start < 0 {
		start = 0
	}
	if end > listLength {
		end = listLength
	}
	return start, end
}

// Пример использования с загрузкой и сохранением состояния
func main() {
	// Загружаем состояние из файла
	err := LoadFromDisk()
	if err != nil {
		fmt.Println("Ошибка при загрузке состояния:", err)
		return
	}

	// Операции с базой данных
	fmt.Println("LPUSH:", LPUSH("list1", 1, 2, 3))      // LPUSH list1 1 2 3
	fmt.Println("RPUSH:", RPUSH("list1", 4, 5, 6))      // RPUSH list1 4 5 6
	fmt.Println("RADDTOSET:", RADDTOSET("list1", 3, 7)) // RADDTOSET list1 3 7

	removedElements, err := LPOP("list1", 2) // LPOP list1 2
	if err != nil {
		fmt.Println("Ошибка LPOP:", err)
	} else {
		fmt.Println("LPOP:", removedElements)
	}

	removedElements, err = RPOP("list1", 1) // RPOP list1 1
	if err != nil {
		fmt.Println("Ошибка RPOP:", err)
	} else {
		fmt.Println("RPOP:", removedElements)
	}

	err = LSET("list1", 1, 99) // LSET list1 1 99
	if err != nil {
		fmt.Println("Ошибка LSET:", err)
	} else {
		fmt.Println("LSET OK")
	}

	element, err := LGET("list1", 1) // LGET list1 1
	if err != nil {
		fmt.Println("Ошибка LGET:", err)
	} else {
		fmt.Println("LGET:", element)
	}

	// Сохраняем состояние на диск
	err = SaveToDisk()
	if err != nil {
		fmt.Println("Ошибка при сохранении состояния:", err)
	}
}
