package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"net/url"
	"log"
	"os"
	"github.com/go-chi/chi"
)

const (
	usersFileName = "users.json"
)

// SaveUsersToFile сохраняет пользователей в файл в формате JSON
func (us *UserService) SaveUsersToFile() error {
	file, err := os.Create(usersFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(us.users); err != nil {
		return err
	}

	return nil
}

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Friends  []int    `json:"friends"`
}

type UserService struct {
	users map[int]*User
	mu    sync.RWMutex
}

func main() {
	userService := UserService{
		users: make(map[int]*User),
	}

	r := chi.NewRouter()
	r.Post("/create", userService.CreateUser)
	r.Post("/make_friends", userService.MakeFriends)
	r.Delete("/user/{id}", userService.DeleteUser)
	r.Get("/friends/{id}", userService.GetFriends)
	r.Put("/user/{id}", userService.UpdateAge)

	server := httptest.NewServer(r)
	defer server.Close()

	fmt.Println("Сервер запущен:", server.URL)

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		log.Fatalf("Ошибка при разборе URL сервера: %v", err)
	}

	port := serverURL.Port()

	// Отправка запроса на создание пользователя
	createUserURL := fmt.Sprintf("http://127.0.0.1:%s/create", port)
	fmt.Println("Отправка запроса на создание пользователя...")
	userID, err := sendCreateUserRequest(createUserURL, &User{Name: "Test User", Age: 30, Friends: []int{}})
	if err != nil {
		log.Fatalf("Ошибка при создании пользователя: %v", err)
	}
	fmt.Println("ID созданного пользователя:", userID)

	// Отправка запроса на создание дружбы
	makeFriendsURL := fmt.Sprintf("http://127.0.0.1:%s/make_friends", port)
	fmt.Println("Отправка запроса на создание дружбы...")
	err = sendMakeFriendsRequest(makeFriendsURL, 1, userID)
	if err != nil {
		log.Fatalf("Ошибка при создании дружбы: %v", err)
	}
	fmt.Println("Дружба успешно создана")

	// Отправка запроса на получение друзей пользователя
	getFriendsURL := fmt.Sprintf("http://127.0.0.1:%s/friends/%d", port, userID)
	fmt.Println("Отправка запроса на получение друзей пользователя...")
	friends, err := sendGetFriendsRequest(getFriendsURL)
	if err != nil {
		log.Fatalf("Ошибка при получении друзей пользователя: %v", err)
	}
	fmt.Println("Друзья пользователя:", friends)

	// Отправка запроса на обновление возраста пользователя
	updateAgeURL := fmt.Sprintf("http://127.0.0.1:%s/user/%d", port, userID)
	fmt.Println("Отправка запроса на обновление возраста пользователя...")
	err = sendUpdateAgeRequest(updateAgeURL, 35)
	if err != nil {
		log.Fatalf("Ошибка при обновлении возраста пользователя: %v", err)
	}
	fmt.Println("Возраст пользователя успешно обновлен")
}

func sendCreateUserRequest(url string, user *User) (int, error) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		ID int `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

func sendMakeFriendsRequest(url string, sourceID, targetID int) error {
	data := struct {
		SourceID int `json:"source_id"`
		TargetID int `json:"target_id"`
	}{sourceID, targetID}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func sendGetFriendsRequest(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var friends []string
	err = json.NewDecoder(resp.Body).Decode(&friends)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func sendUpdateAgeRequest(url string, newAge int) error {
	data := struct {
		NewAge int `json:"new_age"`
	}{newAge}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	id := len(us.users) + 1
	newUser.ID = id
	us.users[id] = &newUser

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, id)
}

func (us *UserService) MakeFriends(w http.ResponseWriter, r *http.Request) {
	var data struct {
		SourceID int `json:"source_id"`
		TargetID int `json:"target_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	source, ok1 := us.users[data.SourceID]
	target, ok2 := us.users[data.TargetID]
	if !ok1 || !ok2 {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	source.Friends = append(source.Friends, data.TargetID)
	target.Friends = append(target.Friends, data.SourceID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "%s и %s теперь друзья"}`, source.Name, target.Name)
}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	user, ok := us.users[userID]
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	delete(us.users, userID)

	for _, friendID := range user.Friends {
		friend := us.users[friendID]
		if friend != nil {
			for i, id := range friend.Friends {
				if id == userID {
					friend.Friends = append(friend.Friends[:i], friend.Friends[i+1:]...)
					break
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Пользователь %s удалён"}`, user.Name)
}

func (us *UserService) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us.mu.RLock()
	defer us.mu.RUnlock()

	user, ok := us.users[userID]
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	friendNames := make([]string, 0, len(user.Friends))
	for _, friendID := range user.Friends {
		friend := us.users[friendID]
		if friend != nil {
			friendNames = append(friendNames, friend.Name)
		}
	}

	json.NewEncoder(w).Encode(friendNames)
}

func (us *UserService) UpdateAge(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newData struct {
		NewAge int `json:"new_age"`
	}
	err = json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	user, ok := us.users[userID]
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	user.Age = newData.NewAge

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Возраст пользователя успешно обновлён"}`)
}