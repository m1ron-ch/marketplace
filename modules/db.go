package modules

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

var database *sql.DB

func ConnectDB() sql.DB {
	db, err := sql.Open("mysql", "root:Qweaz123@/Website2")

	database = db
	if err != nil {
		panic(err)
	}

	return *database
}

func GetUsersFromDB() []User {
	rows, err := database.Query("SELECT * FROM Users")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	users := []User{}

	for rows.Next() {
		user := User{}
		var phone sql.NullString
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &phone, &user.DateRegistration)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if phone.String != "" {
			user.Phone = phone.String
		} else {
			user.Phone = ""
		}
		users = append(users, user)
	}

	return users
}

func GetAdsFromDB(userID int) []Ad {
	rows, err := database.Query("SELECT * FROM Ads")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	ads := []Ad{}

	for rows.Next() {
		var images string
		ad := Ad{}
		err := rows.Scan(&ad.ID, &ad.NameOfGoods, &ad.Overview, &ad.Phone, &ad.UserID, &images, &ad.PostTime, &ad.Location, &ad.Category, &ad.Name, &ad.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var arr []string
		_ = json.Unmarshal([]byte(images), &arr)
		ad.Images = nil
		for i := 0; i < len(arr); i++ {
			ad.Images = append(ad.Images, arr[i])
		}

		ad.PostTime = AdFormatDate(ad.PostTime)
		ad.IsFavorite = IsFavoriteAd(userID, ad.ID)

		ads = append(ads, ad)
	}

	for i, j := 0, len(ads)-1; i < j; i, j = i+1, j-1 {
		ads[i], ads[j] = ads[j], ads[i]
	}

	return ads
}

func GetAdsFromDBByCategoryID(categoryID, userID int) []Ad {
	rows, err := database.Query("SELECT * FROM Ads WHERE Category = ?", categoryID)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	ads := []Ad{}

	for rows.Next() {
		var images string
		ad := Ad{}
		err := rows.Scan(&ad.ID, &ad.NameOfGoods, &ad.Overview, &ad.Phone, &ad.UserID, &images, &ad.PostTime, &ad.Location, &ad.Category, &ad.Name, &ad.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var arr []string
		_ = json.Unmarshal([]byte(images), &arr)
		ad.Images = nil
		for i := 0; i < len(arr); i++ {
			ad.Images = append(ad.Images, arr[i])
		}

		ad.PostTime = AdFormatDate(ad.PostTime)
		ad.IsFavorite = IsFavoriteAd(userID, ad.ID)

		ads = append(ads, ad)
	}

	for i, j := 0, len(ads)-1; i < j; i, j = i+1, j-1 {
		ads[i], ads[j] = ads[j], ads[i]
	}

	return ads
}

func GetUserFromDBByEmail(email string) User {
	if email == "" {
		return User{}
	}

	var user User
	var phone sql.NullString
	query := "SELECT * FROM Users WHERE Email = ?"

	row := database.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &phone, &user.DateRegistration)

	if err != nil {
		panic(err)
	}

	if phone.String != "" {
		user.Phone = phone.String
	} else {
		user.Phone = ""
	}

	return user
}

func GetUserFromDBByID(id int) User {
	var user User
	var phone sql.NullString

	query := "SELECT * FROM Users WHERE ID = ?"
	row := database.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &phone, &user.DateRegistration)
	if err != nil {
		panic(err)
	}

	if phone.String != "" {
		user.Phone = phone.String
	} else {
		user.Phone = ""
	}

	return user
}

func GetCategories() []Category {
	rows, err := database.Query("SELECT * FROM Categories")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	categories := []Category{}

	for rows.Next() {
		category := Category{}
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		categories = append(categories, category)
	}

	return categories
}

func GetAdsFromDBByUserID(userID int) []Ad {
	rows, err := database.Query("SELECT * FROM Ads WHERE UserID = ?", userID)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	ads := []Ad{}

	for rows.Next() {
		var images string
		ad := Ad{}
		err := rows.Scan(&ad.ID, &ad.NameOfGoods, &ad.Overview, &ad.Phone, &ad.UserID, &images, &ad.PostTime, &ad.Location, &ad.Category, &ad.Name, &ad.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var arr []string
		_ = json.Unmarshal([]byte(images), &arr)
		ad.Images = nil
		for i := 0; i < len(arr); i++ {
			ad.Images = append(ad.Images, arr[i])
		}

		ad.PostTime = AdFormatDate(ad.PostTime)
		ads = append(ads, ad)
	}

	for i, j := 0, len(ads)-1; i < j; i, j = i+1, j-1 {
		ads[i], ads[j] = ads[j], ads[i]
	}

	return ads
}

func GetAdFromDBByID(adID int) Ad {
	var ad Ad
	var images string
	query := "SELECT * FROM Ads WHERE ID = ?"

	row := database.QueryRow(query, adID)
	err := row.Scan(&ad.ID, &ad.NameOfGoods, &ad.Overview, &ad.Phone, &ad.UserID, &images, &ad.PostTime, &ad.Location, &ad.Category, &ad.Name, &ad.Price)
	if err != nil {
		panic(err)
	}

	var arr []string
	_ = json.Unmarshal([]byte(images), &arr)
	ad.Images = nil
	for i := 0; i < len(arr); i++ {
		ad.Images = append(ad.Images, arr[i])
	}

	return ad
}

func GetCategoryByID(id int) string {
	var category Category
	query := "SELECT Name FROM Categories WHERE ID = ?"
	row := database.QueryRow(query, id)
	err := row.Scan(&category.Name)
	if err != nil {
		panic(err)
	}

	return category.Name
}

func GetFavoritesAds(userId int) []Ad {
	query := "SELECT * FROM Ads WHERE ID IN (SELECT AdID FROM Favorites WHERE UserID IN (?))"
	rows, err := database.Query(query, userId)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	ads := []Ad{}

	for rows.Next() {
		var images string
		ad := Ad{}
		err := rows.Scan(&ad.ID, &ad.NameOfGoods, &ad.Overview, &ad.Phone, &ad.UserID, &images, &ad.PostTime, &ad.Location, &ad.Category, &ad.Name, &ad.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var arr []string
		_ = json.Unmarshal([]byte(images), &arr)
		ad.Images = nil
		for i := 0; i < len(arr); i++ {
			ad.Images = append(ad.Images, arr[i])
		}

		ad.PostTime = AdFormatDate(ad.PostTime)
		ads = append(ads, ad)
	}

	for i, j := 0, len(ads)-1; i < j; i, j = i+1, j-1 {
		ads[i], ads[j] = ads[j], ads[i]
	}

	return ads

}

func GetChatsByFirstUser(userID int) []Chat {
	query := "SELECT * FROM Chat WHERE FirstUserID = ?"
	rows, err := database.Query(query, userID)
	if err != nil {
		log.Println(err.Error())
	}

	defer rows.Close()
	chats := []Chat{}

	for rows.Next() {
		chat := Chat{}
		err := rows.Scan(&chat.ID, &chat.Name, &chat.FirstUserID, &chat.SecondUserID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		chats = append(chats, chat)
	}

	if len(chats) == 0 {
		return GetChatsBySecondUser(userID)
	}

	for i, j := 0, len(chats)-1; i < j; i, j = i+1, j-1 {
		chats[i], chats[j] = chats[j], chats[i]
	}

	return chats
}

func GetChatsBySecondUser(userID int) []Chat {
	query := "SELECT * FROM Chat WHERE SecondUserID = ?"
	rows, err := database.Query(query, userID)
	if err != nil {
		log.Println(err.Error())
	}

	defer rows.Close()
	chats := []Chat{}

	for rows.Next() {
		chat := Chat{}
		err := rows.Scan(&chat.ID, &chat.Name, &chat.FirstUserID, &chat.SecondUserID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		chats = append(chats, chat)
	}

	for i, j := 0, len(chats)-1; i < j; i, j = i+1, j-1 {
		chats[i], chats[j] = chats[j], chats[i]
	}

	return chats
}

func GetMessages(chatID int) []Message {
	query := "SELECT * FROM Messages WHERE ChatID = ?"
	rows, err := database.Query(query, chatID)
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	messages := []Message{}

	for rows.Next() {
		message := Message{}
		err := rows.Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.Date, &message.Time)
		if err != nil {
			fmt.Println(err)
			continue
		}

		message.Time = GetHM(message.Time)
		messages = append(messages, message)
	}

	return messages
}

func GetLastChatIDByFirstUserID(userID int) int {
	var id int
	query := "SELECT MAX(ID) FROM Chat WHERE FirstUserID = ?"

	row := database.QueryRow(query, userID)
	err := row.Scan(&id)
	if err != nil {
		log.Println(err.Error())
		return GetLastChatIDBySecondUserID(userID)
	}

	return id
}

func GetLastChatIDBySecondUserID(userID int) int {
	var id int
	query := "SELECT MAX(ID) FROM Chat WHERE SecondUserID = ?"

	row := database.QueryRow(query, userID)
	err := row.Scan(&id)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return id
}

func IsFavoriteAd(userID, adID int) bool {
	var id int
	query := "SELECT ID FROM Favorites WHERE UserID = ? AND AdID = ?"
	row := database.QueryRow(query, userID, adID)
	err := row.Scan(&id)

	return err == nil
}

func IsHaveUserInDB(email, password string) (User, bool) {
	var user User
	query := "SELECT * FROM Users WHERE Email = ? AND Password = ?"

	row := database.QueryRow(query, email, password)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.Phone, &user.DateRegistration)
	if err != nil {
		return User{}, false
	}

	fmt.Println(user.Email)
	return user, true
}

func IsHaveEmailInDB(email string) bool {
	var user User
	query := "SELECT * FROM Users WHERE Email = ? "

	row := database.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.Phone, &user.DateRegistration)

	return err == nil
}

func IsHisAd(userID, adID int) bool {
	var ad Ad
	var images string
	query := "SELECT * FROM Ads WHERE ID = ? AND UserID = ?"

	row := database.QueryRow(query, adID, userID)
	err := row.Scan(&ad.ID, &ad.NameOfGoods, &ad.Overview, &ad.Phone, &ad.UserID, &images, &ad.PostTime, &ad.Location, &ad.Category, &ad.Name, &ad.Price)
	ad.Images = append(ad.Images, images)

	return err == nil
}

func GetChatIDByFirstUserIDSecondByUserID(name string, firstUserID, secondUserID int) int {
	var id int
	query := "SELECT ID FROM Chat WHERE Name = ? AND FirstUserID = ? AND SecondUserID = ?"

	row := database.QueryRow(query, name, firstUserID, secondUserID)
	err := row.Scan(&id)

	if err != nil {
		return GetChatIDBySecondUserIDFirstByUserID(name, secondUserID, firstUserID)
	}

	return id
}

func GetChatIDBySecondUserIDFirstByUserID(name string, secondUserID, firstUserID int) int {
	var id int
	query := "SELECT ID FROM Chat WHERE Name = ? AND SecondUserID = ? AND FirstUserID = ?"

	row := database.QueryRow(query, name, secondUserID, firstUserID)
	err := row.Scan(&id)

	if err != nil {
		return -1
	}

	return id
}

func AddAdToDB(ad Ad) {
	query := "INSERT INTO Ads (NameOfGoods, Overview, Phone, UserID, Images, PostTime, Location, Category, Name, Price) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := database.Exec(query, ad.NameOfGoods, ad.Overview, ad.Phone, ad.UserID, ad.Images[0], ad.PostTime, ad.Location, ad.Category, ad.Name, ad.Price)
	if err != nil {
		panic(err)
	}
}

func AddUserToDB(user User) {
	query := "INSERT INTO Users (Name, Email, Password, DateRegistration) VALUES (?, ?, ?, ?)"
	_, err := database.Exec(query, user.Name, user.Email, user.Password, user.DateRegistration)
	if err != nil {
		panic(err)
	}
}

func AddFavorite(userID, adID int) {
	query := "INSERT INTO Favorites (UserID, AdID) VALUES (?, ?)"
	_, err := database.Exec(query, userID, adID)
	if err != nil {
		panic(err)
	}
}

func AddChat(name string, firstUserID, secondUserID int) {
	query := "INSERT INTO Chat (Name, FirstUserID, SecondUserID) VALUES (?, ?, ?)"
	_, err := database.Exec(query, name, firstUserID, secondUserID)
	if err != nil {
		panic(err)
	}
}

func AddMessageToChat(message Message) {
	query := "INSERT INTO Messages (ChatID, UserID, Content, Date, Time) VALUES (?, ?, ?, ?, ?)"
	_, err := database.Exec(query, message.ChatID, message.UserID, message.Content, message.Date, message.Time)
	if err != nil {
		panic(err)
	}
}

func UpdateUser(user User) {
	query := "UPDATE Users SET Name = ?, Email = ?, Password = ?, Phone = ? WHERE ID = ? "
	_, err := database.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.ID)
	if err != nil {
		panic(err)
	}
}

func UpdateAd(ad Ad) {
	query := "UPDATE Ads SET NameOfGoods = ?, Overview = ?, Phone = ?, Location = ?, Category = ?, Name = ?, Price = ? WHERE ID = ?"
	_, err := database.Exec(query, ad.NameOfGoods, ad.Overview, ad.Phone, ad.Location, ad.Category, ad.Name, ad.Price, ad.ID)
	if err != nil {
		panic(err)
	}
}

func DeleteUserFromDBByID(id int) error {
	query := "DELETE FROM Users WHERE ID = ?"
	_, err := database.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAdFromDBByID(id int) {
	query := "DELETE FROM Ads WHERE ID = ?"
	_, err := database.Exec(query, id)
	if err != nil {
		panic(err)
	}
}

func DeleteFavoriteFromDB(userID, adID int) {
	query := "DELETE FROM Favorites WHERE UserID = ? AND AdID = ?"
	_, err := database.Exec(query, userID, adID)
	if err != nil {
		panic(err)
	}
}

func GetRolesFromDB() []Role {
	rows, err := database.Query("SELECT * FROM Roles")
	if err != nil {
		log.Println(err)
	}

	defer rows.Close()
	roles := []Role{}

	for rows.Next() {
		role := Role{}
		err := rows.Scan(&role.ID, &role.Role, &role.Privilege)
		if err != nil {
			fmt.Println(err)
			continue
		}
		roles = append(roles, role)
	}

	log.Println(roles)

	return roles
}
