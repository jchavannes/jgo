package db

type Session struct {
	Id       uint `gorm:"primary_key"`
	CookieId string `gorm:"unique_index"`
	UserId   uint
	StartTs  uint
}

func (s *Session) Save() error {
	db, err := getDb()
	if err != nil {
		return err
	}
	result := db.Save(s)
	return result.Error
}

func GetSession(cookieId string) (*Session, error) {
	session := &Session{
		CookieId: cookieId,
	}
	db, err := getDb()
	if err != nil {
		return session, err
	}
	result := db.Find(session, session)
	if result.Error != nil && result.Error.Error() == "record not found" {
		result = db.Create(session)
	}
	return session, result.Error
}
