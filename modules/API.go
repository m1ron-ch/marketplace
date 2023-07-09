package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var users []User
	if value, err := strconv.Atoi(vars["id"]); err != nil {
		users = GetUsersFromDB()
	} else if value > 0 {
		users = append(users, GetUserFromDBByID(value))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(jsonData))
}

func GetAuthUser(w http.ResponseWriter, r *http.Request) {
	user := GetCookieUser(w, r)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	jsonData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(jsonData))
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	WriteCookie(w, "User", "miron.cherneyko@mail.ru")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	jsonData, err := json.MarshalIndent("", "", "  ")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	fmt.Println("Cookie set")
	fmt.Fprintln(w, string(jsonData))
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryes := GetCategories()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	jsonData, err := json.MarshalIndent(categoryes, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(jsonData))
}

func DeleteUserHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Robit")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	errDeleteUser := DeleteUserFromDBByID(id)
	if errDeleteUser != nil {
		http.Error(w, "User Not Found By ID", http.StatusBadRequest)
		return
	}

	fmt.Println(w, "API. User with ID %d deleted successfully", id)
	w.WriteHeader(http.StatusOK)
}

func GetRoles(w http.ResponseWriter, r *http.Request) {
	roles := GetRolesFromDB()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	jsonData, err := json.MarshalIndent(roles, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(jsonData))
}
