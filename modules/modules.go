package modules

type User struct {
	ID               int    `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	Password         string `json:"password"`
	RoleID           int    `json:"role_id"`
	Phone            string `json:"phone"`
	DateRegistration string `json:"date_registration"`
}

type Ad struct {
	ID          int
	NameOfGoods string
	Overview    string
	Phone       string
	UserID      int
	Images      []string
	PostTime    string
	Location    string
	Category    string
	Name        string
	Price       string
	IsFavorite  bool
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Role struct {
	ID        int
	Role      string
	Privilege int
}

type Favorites struct {
	ID     int
	UserID int
	AdID   int
}

type Error struct {
	Message string
	Path    string
}

type ViewIndexData struct {
	User       User
	Ads        []Ad
	Categories []Category
}

type ViewPlaceAd struct {
	User       User
	Categories []Category
}

type ViewEditAd struct {
	User       User
	Ad         Ad
	Categories []Category
}

type ViewUsers struct {
	User  User
	Users []User
}

type ViewEditUser struct {
	AuthUser User
	EditUser User
}

type ViewSettings struct {
	User User
}

type ViewAllAds struct {
	User User
	Ads  []Ad
}

type ViewFavorites struct {
	User       User
	Ads        []Ad
	Categories []Category
}

type ViewAd struct {
	User User
	Ad   Ad
}

type ViewUser struct {
	User        User
	ProfileUser User
	Categories  []Category
	Ads         []Ad
}

type ViewCategory struct {
	User         User
	Categories   []Category
	Ads          []Ad
	CategoryName string
}

type Chat struct {
	ID           int
	Name         string
	FirstUserID  int
	SecondUserID int
}

type Message struct {
	ID      int
	ChatID  int
	UserID  int
	Content string
	Date    string
	Time    string
}

type MessageStatus struct {
	ID        int
	MessageID int
	UserID    int
	IsRead    bool
}

type ViewChat struct {
	User           User
	Chats          []Chat
	Messages       []Message
	SelectedChatID int
}
