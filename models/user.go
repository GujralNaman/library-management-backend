// models/user.go

package models

// User represents a user model
type User struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	ContactNumber string `json:"contactNumber"`
	Role          string `json:"role"`
	LibID         *uint   `json:"libid"`
	Password      string `json:"password"`
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(userID uint) (*User, error) {
	var user User
	if err := DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) error {
	return DB.Create(user).Error
}
