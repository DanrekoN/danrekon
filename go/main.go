package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Car struct {
	ID         int64  `json:"id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Mileage    int    `json:"mileage"`
	OwnerCount int    `json:"owner_count"`
}

type Furniture struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Manufacturer string  `json:"manufacturer"`
	Height       float64 `json:"height"`
	Width        float64 `json:"width"`
	Length       float64 `json:"length"`
}

type Flower struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	ArrivalDate string  `json:"arrival_date"`
}

var (
	cars            = []Car{}
	furnitures      = []Furniture{}
	flowers         = []Flower{}
	mutex           sync.Mutex
	nextCarID       int64
	nextFurnitureID int64
	nextFlowerID    int64
)

func loadData() {
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		return
	}
	json.Unmarshal(file, &cars)
	json.Unmarshal(file, &furnitures)
	json.Unmarshal(file, &flowers)

	for _, car := range cars {
		if car.ID >= nextCarID {
			nextCarID = car.ID + 1
		}
	}
	for _, furniture := range furnitures {
		if furniture.ID >= nextFurnitureID {
			nextFurnitureID = furniture.ID + 1
		}
	}
	for _, flower := range flowers {
		if flower.ID >= nextFlowerID {
			nextFlowerID = flower.ID + 1
		}
	}
}

func saveData() {
	data := struct {
		Cars       []Car       `json:"cars"`
		Furnitures []Furniture `json:"furnitures"`
		Flowers    []Flower    `json:"flowers"`
	}{
		Cars:       cars,
		Furnitures: furnitures,
		Flowers:    flowers,
	}
	file, _ := json.MarshalIndent(data, "", " ")
	ioutil.WriteFile("data.json", file, 0644)
}

// Cars Handlers
func createCar(w http.ResponseWriter, r *http.Request) {
	var car Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	car.ID = nextCarID
	nextCarID++
	cars = append(cars, car)
	mutex.Unlock()
	saveData()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

func getCars(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cars)
}

func getCarByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/cars/"):]
	for _, car := range cars {
		if fmt.Sprintf("%d", car.ID) == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(car)
			return
		}
	}
	http.NotFound(w, r)
}

func updateCar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/cars/"):]
	var updatedCar Car
	if err := json.NewDecoder(r.Body).Decode(&updatedCar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	for i, car := range cars {
		if fmt.Sprintf("%d", car.ID) == id {
			cars[i] = updatedCar
			cars[i].ID = car.ID // сохраняем оригинальный ID
			mutex.Unlock()
			saveData()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}
	mutex.Unlock()
	http.NotFound(w, r)
}

func patchCar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/cars/"):]
	var updatedFields Car
	if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	for i, car := range cars {

		if fmt.Sprintf("%d", car.ID) == id {
			if updatedFields.Brand != "" {
				cars[i].Brand = updatedFields.Brand
			}
			if updatedFields.Model != "" {
				cars[i].Model = updatedFields.Model
			}
			if updatedFields.Mileage != 0 {
				cars[i].Mileage = updatedFields.Mileage
			}
			if updatedFields.OwnerCount != 0 {
				cars[i].OwnerCount = updatedFields.OwnerCount
			}
			mutex.Unlock()
			saveData()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	mutex.Unlock()
	http.NotFound(w, r)
}

func deleteCar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/cars/"):]
	mutex.Lock()
	for i, car := range cars {
		if fmt.Sprintf("%d", car.ID) == id {
			cars = append(cars[:i], cars[i+1:]...)
			mutex.Unlock()
			saveData()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	mutex.Unlock()
	http.NotFound(w, r)
}

// Furniture Handlers (аналогично для мебели и цветов)
func createFurniture(w http.ResponseWriter, r *http.Request)  { /*...*/ }
func getFurnitures(w http.ResponseWriter, r *http.Request)    { /*...*/ }
func getFurnitureByID(w http.ResponseWriter, r *http.Request) { /*...*/ }
func updateFurniture(w http.ResponseWriter, r *http.Request)  { /*...*/ }
func patchFurniture(w http.ResponseWriter, r *http.Request)   { /*...*/ }
func deleteFurniture(w http.ResponseWriter, r *http.Request)  { /*...*/ }

// Flower Handlers (аналогично для цветов)
func createFlower(w http.ResponseWriter, r *http.Request)  { /*...*/ }
func getFlowers(w http.ResponseWriter, r *http.Request)    { /*...*/ }
func getFlowerByID(w http.ResponseWriter, r *http.Request) { /*...*/ }
func updateFlower(w http.ResponseWriter, r *http.Request)  { /*...*/ }
func patchFlower(w http.ResponseWriter, r *http.Request)   { /*...*/ }
func deleteFlower(w http.ResponseWriter, r *http.Request)  { /*...*/ }

func main() {
	loadData() // Загружаем данные при старте

	http.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createCar(w, r)
		case http.MethodGet:
			getCars(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/cars/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCarByID(w, r)
		case http.MethodPut:
			updateCar(w, r)
		case http.MethodPatch:
			patchCar(w, r)
		case http.MethodDelete:
			deleteCar(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/furniture", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createFurniture(w, r)
		case http.MethodGet:
			getFurnitures(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/furniture/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getFurnitureByID(w, r)
		case http.MethodPut:
			updateFurniture(w, r)
		case http.MethodPatch:
			patchFurniture(w, r)
		case http.MethodDelete:
			deleteFurniture(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/flowers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createFlower(w, r)
		case http.MethodGet:
			getFlowers(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/flowers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getFlowerByID(w, r)
		case http.MethodPut:
			updateFlower(w, r)
		case http.MethodPatch:
			patchFlower(w, r)
		case http.MethodDelete:
			deleteFlower(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
