package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"server/modules"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type chatRoom struct {
	clients       map[*client]bool
	addClient     chan *client
	removeClient  chan *client
	broadcastChan chan []byte
}

type client struct {
	conn *websocket.Conn
	send chan []byte
	room *chatRoom
}

var chat = &chatRoom{
	clients:       make(map[*client]bool),
	addClient:     make(chan *client),
	removeClient:  make(chan *client),
	broadcastChan: make(chan []byte),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var errorMessage modules.Error
var verificationCode string
var user modules.User
var msg modules.Message

func IsValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		var data modules.ViewIndexData
		data.User = user
		data.Ads = modules.GetAdsFromDB(user.ID)
		data.Categories = modules.GetCategories()

		tmpl, _ := template.ParseFiles("./templates/html/index.html")
		tmpl.Execute(w, data)
	default:
		panic("Method don't find")
	}
}

func SignInUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "./templates/html/signIn.html")
	case "POST":
		r.ParseForm()

		var u modules.User
		u.Name = r.PostFormValue("name")
		u.Email = r.PostFormValue("email")
		u.Password = modules.GetMD5(r.PostFormValue("password"))
		confirmPassword := modules.GetMD5(r.PostFormValue("confirmPassword"))
		u.Phone = ""
		u.DateRegistration = modules.GetNowDateTime()

		isHaveUser := modules.IsHaveEmailInDB(u.Email)
		isValidEmail := modules.IsValidEmail(u.Email)
		samePassword := u.Password == confirmPassword

		if isValidEmail && samePassword && !isHaveUser {
			user = u
			code, err := modules.GenerateOTP()
			if err != nil {
				panic(err)
			}

			verificationCode = code
			modules.SendEmailMessage(u.Email, verificationCode)
			fmt.Println(verificationCode)

			http.Redirect(w, r, "user_verification", http.StatusSeeOther)
		} else {
			if isHaveUser {
				errorMessage.Message = "The user is already registered with this mail '" + u.Email + "'"
				errorMessage.Path = "/"
				http.Redirect(w, r, "/error_page", http.StatusSeeOther)
			} else {
				errorMessage.Message = "Password mismatch"
				errorMessage.Path = "/"
				http.Redirect(w, r, "/error_page", http.StatusSeeOther)
			}
		}
	default:
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		email := r.PostFormValue("email")
		password := modules.GetMD5(r.PostFormValue("password"))

		fmt.Println(email)
		fmt.Println(password)

		userDB, isHave := modules.IsHaveUserInDB(email, password)
		if isHave {
			user = userDB
			modules.WriteCookie(w, "User", user.Email)
		} else {
			errorMessage.Message = "Incorrect login or password"
			errorMessage.Path = "/"
			http.Redirect(w, r, "/error_page", http.StatusSeeOther)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
	}
}

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if user.Email == "" {
			return
		}

		tmpl, _ := template.ParseFiles("./templates/html/verification.html")
		tmpl.Execute(w, user)
		return
	case "POST":
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		var inputCode string
		for i := 1; i <= 6; i++ {
			inputCode += r.PostFormValue("num" + strconv.Itoa(i))
		}

		if verificationCode == inputCode {
			modules.WriteCookie(w, "User", user.Email)
			modules.AddUserToDB(user)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			errorMessage.Message = "Incorrect verfiication code"
			errorMessage.Path = "/user_verification"
			http.Redirect(w, r, "/error_page", http.StatusSeeOther)
		}
	default:
	}
}

func MyAdsHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		var data modules.ViewIndexData
		data.Categories = modules.GetCategories()
		data.Ads = modules.GetAdsFromDBByUserID(user.ID)
		data.User = user

		tmpl, _ := template.ParseFiles("./templates/html/my_ads.html")
		tmpl.Execute(w, data)
	case "POST":
	default:
	}
}

func PlaceAdHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		var data modules.ViewPlaceAd
		data.User = user
		data.Categories = modules.GetCategories()

		tmpl, _ := template.ParseFiles("./templates/html/place_ad.html")
		tmpl.Execute(w, data)
	case "POST":
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		r.ParseMultipartForm(0)

		var fileNames []string
		fhs := r.MultipartForm.File["images"]
		for _, fh := range fhs {
			f, err := fh.Open()
			if err != nil {
				panic(err)
			}

			defer f.Close()

			err = os.MkdirAll("./templates/images", os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			imageFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fh.Filename))
			dst, err := os.Create(fmt.Sprintf("./templates/images/%s", imageFileName))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer dst.Close()

			_, err = io.Copy(dst, f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileNames = append(fileNames, imageFileName)
		}

		urlsJson, err := json.Marshal(fileNames)
		if err != nil {
			panic(err)
		}

		var ad modules.Ad
		ad.NameOfGoods = r.PostFormValue("nameOfGoods")
		ad.Overview = r.PostFormValue("overview")
		ad.Phone = r.PostFormValue("phone")
		ad.UserID = user.ID
		ad.Images = append(ad.Images, string(urlsJson))
		ad.PostTime = modules.GetNowDateTime()
		ad.Location = r.FormValue("location")
		ad.Category = r.FormValue("category")
		ad.Name = r.PostFormValue("name")
		ad.Price = r.PostFormValue("price")

		modules.AddAdToDB(ad)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
	}
}

func AllAdsHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)
		if user.RoleID <= 1 {
			return
		}

		var data modules.ViewAllAds
		data.User = user
		data.Ads = modules.GetAdsFromDB(user.ID)

		tmpl, _ := template.ParseFiles("./templates/html/all_ads.html")
		tmpl.Execute(w, data)
	case "POST":
	default:
	}
}

func FavoritesHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		var data modules.ViewFavorites
		data.User = user
		data.Categories = modules.GetCategories()
		data.Ads = modules.GetFavoritesAds(user.ID)

		tmpl, _ := template.ParseFiles("./templates/html/favorites.html")
		tmpl.Execute(w, data)
	case "POST":
		vars := mux.Vars(r)
		adID, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		modules.AddFavorite(user.ID, adID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case "DELETE":
		vars := mux.Vars(r)
		adID, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		modules.DeleteFavoriteFromDB(user.ID, adID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
	}

}

func SettingsHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var data modules.ViewSettings
		data.User = modules.GetCookieUser(w, r)
		tmpl, _ := template.ParseFiles("./templates/html/settings.html")
		tmpl.Execute(w, data)
	case "POST":
		user = modules.GetCookieUser(w, r)
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		var u modules.User
		u.ID = user.ID
		u.Name = r.PostFormValue("name")
		u.Email = r.PostFormValue("email")
		u.Password = r.PostFormValue("password")
		u.Phone = r.PostFormValue("phone")

		if len(u.Password) == 0 {
			u.Password = user.Password
		} else {
			u.Password = modules.GetMD5(u.Password)
		}

		if modules.IsValidEmail(u.Email) {
			var isHave = false
			if u.Email != user.Email {
				if modules.IsHaveEmailInDB(u.Email) {
					isHave = true
				}
			}

			if !isHave {
				modules.UpdateUser(u)
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
			}
		} else {
			fmt.Fprintln(w, "Error save settings")
		}

	default:
	}

}

func LogoutHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		modules.WriteCookie(w, "User", "")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
	}
}

func AdHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var data modules.ViewAd
		data.User = user
		data.Ad = modules.GetAdFromDBByID(id)

		tmpl, _ := template.ParseFiles("./templates/html/ad.html")
		tmpl.Execute(w, data)
	default:
	}
}

func UserHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var data modules.ViewUser
		data.User = user
		data.Categories = modules.GetCategories()
		data.ProfileUser = modules.GetUserFromDBByID(id)
		data.Ads = modules.GetAdsFromDBByUserID(id)

		tmpl, _ := template.ParseFiles("./templates/html/user.html")
		tmpl.Execute(w, data)
	default:
	}
}

func CategoryHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var data modules.ViewCategory
		data.User = user
		data.Categories = modules.GetCategories()
		data.Ads = modules.GetAdsFromDBByCategoryID(id, user.ID)
		data.CategoryName = modules.GetCategoryByID(id)

		tmpl, _ := template.ParseFiles("./templates/html/category.html")
		tmpl.Execute(w, data)
	default:
	}
}

func EditAdHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var data modules.ViewEditAd
		data.User = user
		data.Ad = modules.GetAdFromDBByID(id)
		data.Categories = modules.GetCategories()

		tmpl, _ := template.ParseFiles("./templates/html/edit_ad.html")
		tmpl.Execute(w, data)
	case "POST":
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		r.ParseMultipartForm(0)

		var fileNames []string
		fhs := r.MultipartForm.File["images"]
		for _, fh := range fhs {

			f, err := fh.Open()
			if err != nil {
				panic(err)
			}

			defer f.Close()

			err = os.MkdirAll("./templates/images", os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			imageFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fh.Filename))
			dst, err := os.Create(fmt.Sprintf("./templates/images/%s", imageFileName))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer dst.Close()

			_, err = io.Copy(dst, f)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileNames = append(fileNames, imageFileName)
		}

		// for i := 0; i < 12; i++ {
		// 	var formName = "img" + strconv.Itoa(i)
		// 	fmt.Println("'" + r.FormValue(formName) + "'")
		// 	if r.PostFormValue(formName) != "" {
		// 		fileNames = append(fileNames, r.PostFormValue(formName))
		// 	}
		// }

		// urlsJson, err := json.Marshal(fileNames)
		// if err != nil {
		// 	panic(err)
		// }

		var ad modules.Ad
		id, err := strconv.Atoi(r.PostFormValue("id"))
		if err != nil {
			panic(err)
		}

		ad.ID = id
		ad.NameOfGoods = r.PostFormValue("nameOfGoods")
		ad.Overview = r.PostFormValue("overview")
		ad.Phone = r.PostFormValue("phone")
		ad.UserID = user.ID
		// ad.Images = append(ad.Images, string(urlsJson))
		ad.PostTime = modules.GetNowDateTime()
		ad.Location = r.FormValue("location")
		ad.Category = r.FormValue("category")
		ad.Name = r.PostFormValue("name")
		ad.Price = r.PostFormValue("price")

		modules.UpdateAd(ad)
		http.Redirect(w, r, "/my_ads", http.StatusSeeOther)
	default:
	}
}

func DeleteAdHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	if delete := modules.IsHisAd(user.ID, id); delete {
		modules.DeleteAdFromDBByID(id)
	}
}

func UsersHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)
		if user.RoleID <= 1 {
			fmt.Fprintln(w, "Not access")
			return
		}

		var data modules.ViewUsers
		data.User = user
		data.Users = modules.GetUsersFromDB()
		tmpl, _ := template.ParseFiles("./templates/html/users.html")
		tmpl.Execute(w, data)
	case "POST":
	default:
	}
}

func EditUserHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)
		access, err := modules.IsHaveAccess(user)
		if !access {
			fmt.Fprintln(w, err)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		user = modules.GetCookieUser(w, r)

		var data modules.ViewEditUser
		data.AuthUser = user
		data.EditUser = modules.GetUserFromDBByID(id)
		tmpl, _ := template.ParseFiles("./templates/html/edit_user.html")
		tmpl.Execute(w, data)
	case "POST":
		var u modules.User
		u.ID = user.ID
		u.Name = r.PostFormValue("name")
		u.Email = r.PostFormValue("email")
		u.Password = r.PostFormValue("password")
		u.Phone = r.PostFormValue("phone")

		if len(u.Password) == 0 {
			u.Password = user.Password
		} else {
			u.Password = modules.GetMD5(u.Password)
		}

		if modules.IsValidEmail(u.Email) {
			var isHave = false
			if u.Email != user.Email {
				if modules.IsHaveEmailInDB(u.Email) {
					isHave = true
				}
			}

			if !isHave {
				modules.UpdateUser(u)
				http.Redirect(w, r, "/settings", http.StatusSeeOther)
			}
		} else {
			fmt.Fprintln(w, "Error save settings")
		}
	default:
		panic("Not found")
	}
}

func DeleteUserHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		modules.DeleteUserFromDBByID(id)
		http.Redirect(w, r, "/edit_users", http.StatusSeeOther)
	default:
		panic("Not found")
	}
}

func CheckChatByNameByIDHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)

		path := "/messages/"
		vars := mux.Vars(r)
		adName := vars["name"]
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		chatID := modules.GetChatIDByFirstUserIDSecondByUserID(adName, user.ID, id)
		if chatID > 0 {
			path += strconv.Itoa(chatID)
		} else if chatID == -1 {
			modules.AddChat(adName, user.ID, id)
			path += strconv.Itoa(modules.GetLastChatIDByFirstUserID(user.ID))
		}

		http.Redirect(w, r, path, http.StatusSeeOther)
	default:
		panic("Not found")
	}
}

func AllMessagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user = modules.GetCookieUser(w, r)
		vars := mux.Vars(r)

		var data modules.ViewChat
		data.User = user
		chatByFirstUserID := modules.GetChatsByFirstUser(user.ID)
		chatBySecondUserID := modules.GetChatsBySecondUser(user.ID)
		data.Chats = append(chatByFirstUserID, chatBySecondUserID...)

		if vars["id"] != "" {
			id, err := strconv.Atoi(vars["id"])
			if err != nil {
				panic(err)
			}

			data.SelectedChatID = id
			data.Messages = modules.GetMessages(id)
			msg.ChatID = id
		} else {
			lastID := modules.GetLastChatIDByFirstUserID(user.ID)
			data.Messages = modules.GetMessages(lastID)
			data.SelectedChatID = lastID
			msg.ChatID = lastID
		}

		msg.UserID = user.ID

		tmpl, _ := template.ParseFiles("./templates/html/messages.html")
		tmpl.Execute(w, data)
	default:

	}
}

func ErrorPadgeHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("./templates/html/error_page.html")
		tmpl.Execute(w, errorMessage)
	default:

	}
}

func (c *chatRoom) Run() {
	for {
		select {
		case client := <-c.addClient:
			c.clients[client] = true
		case client := <-c.removeClient:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)
			}
		case message := <-c.broadcastChan:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}
		}
	}
}

func (c *chatRoom) Broadcast(message []byte) {
	c.broadcastChan <- message
}

func (c *client) readPump() {
	defer func() {
		c.room.removeClient <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		msg.Content = string(message)
		msg.Date = modules.GetNowDate()
		msg.Time = modules.GetNowTimeHoureMinute()
		modules.AddMessageToChat(msg)
		c.room.Broadcast(message)
	}
}

func (c *client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func ServeChatRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{
		conn: conn,
		send: make(chan []byte, 256),
		room: chat,
	}

	chat.addClient <- client

	go client.WritePump()
	client.readPump()
}

func addCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешить доступ с любого источника
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешить определенные методы запросов
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// Разрешить определенные заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}

func main() {
	database := modules.ConnectDB()

	go chat.Run()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/sign_in", SignInUserHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/user_verification", VerificationHandler)
	r.HandleFunc("/my_ads", MyAdsHandle)
	r.HandleFunc("/all_ads", AllAdsHandle)
	r.HandleFunc("/place_ad", PlaceAdHandle)
	r.HandleFunc("/favorites", FavoritesHandle)
	r.HandleFunc("/add_favorites/{id:[0-9]+}", FavoritesHandle).Methods("POST")
	r.HandleFunc("/settings", SettingsHandle)
	r.HandleFunc("/logout", LogoutHandle)
	r.HandleFunc("/user/{id:[0-9]+}", UserHandle)
	r.HandleFunc("/ad/{id:[0-9]+}", AdHandle)
	r.HandleFunc("/category/{id:[0-9]+}", CategoryHandle)
	r.HandleFunc("/edit_ad", EditAdHandle).Methods("POST")
	r.HandleFunc("/edit_user", EditUserHandle)
	r.HandleFunc("/edit_ad/{id:[0-9]+}", EditAdHandle)
	r.HandleFunc("/edit_user/{id:[0-9]+}", EditUserHandle)
	r.HandleFunc("/delete_ad/{id:[0-9]+}", DeleteAdHandle).Methods("DELETE")
	r.HandleFunc("/delete_user/{id:[0-9]+}", DeleteUserHandle).Methods("DELETE")
	r.HandleFunc("/delete_favorite/{id:[0-9]+}", FavoritesHandle).Methods("DELETE")
	r.HandleFunc("/users", UsersHandle)
	r.HandleFunc("/check_chat/{name}/{id:[0-9]+}", CheckChatByNameByIDHandle)
	r.HandleFunc("/messages", AllMessagesHandler)
	r.HandleFunc("/messages/{id:[0-9]+}", AllMessagesHandler)
	r.HandleFunc("/ws", ServeChatRoom)
	r.HandleFunc("/error_page", ErrorPadgeHandle)

	r.HandleFunc("/api/users", modules.GetUser).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", modules.GetUser).Methods("GET")
	r.HandleFunc("/api/auth_user", modules.GetAuthUser).Methods("GET")
	r.HandleFunc("/api/categories", modules.GetCategory).Methods("GET")
	r.HandleFunc("/api/delete_user/{id:[0-9]+}", modules.DeleteUserHandle).Methods("DELETE")
	r.HandleFunc("/api/roles", modules.GetRoles).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))

	r.Use(addCorsHeaders)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	defer database.Close()

	fmt.Println(fmt.Sprintf("Server is listening %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
