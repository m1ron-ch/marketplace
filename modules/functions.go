package modules

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetCookieUser(w http.ResponseWriter, r *http.Request) User {
	cookieUserEmail, err := ReadCookie(r, "User")
	if err != nil {
		WriteCookie(w, "User", "")
		return User{}
	}

	return GetUserFromDBByEmail(cookieUserEmail)
}

func IsHaveAccess(user User) (bool, error) {
	return user.RoleID >= 1, errors.New("Not access")
}

func IsPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`(?:^|[^0-9])(1[34578][0-9]{9})(?:$|[^0-9])`)
	submatch := re.FindStringSubmatch(phone)
	if len(submatch) < 2 {
		return false
	}
	return true
}

func AdFormatDate(dateTime string) string {
	var date string
	var filedsDateTime = strings.Fields(dateTime)
	if filedsDateTime[0] == GetNowDate() {
		date += "Today"
	} else if filedsDateTime[0] == time.Now().AddDate(0, 0, -1).Format("2006-01-02") {
		date += "Yesterday"
	} else {
		date += filedsDateTime[0]
	}

	date += " " + filedsDateTime[1][:5]
	return date
}

func GetHM(time string) string {
	return time[:5]
}

func GetNowDateTime() string {

	dt := time.Now()
	return dt.Format("2006-01-02 15:04:05")
}

func GetNowDate() string {

	dt := time.Now()
	return dt.Format("2006-01-02")
}

func GetNowTimeHoureMinute() string {

	dt := time.Now()
	return dt.Format("15:04")
}
