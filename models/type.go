package models

type NewsType struct {
	Title string `bson:"title"`
	Link  string `bson:"link"`
}

type User struct {
	ChatID   int64  `bson:"chatid"`
	Username string `bson:"username"`
}
