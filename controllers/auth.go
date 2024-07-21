// controllers/auth.go

package controllers

import (
	"library/task/models"
	"library/task/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterDate struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	ContactNumber string `json:"contactNumber"`
	Role          string `json:"role"`
	LibraryName   string `json:"libraryName"`
}

func Login(c *gin.Context) {
	var authInput models.AuthInput
	// var emailvalidation models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// find := models.DB.Where("email = ?", emailvalidation.Email).First(&emailvalidation)
	// if find.Error == nil {
	// 	// user found
	// 	c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
	// 	return
	// }

	user, err := models.GetUserByEmail(authInput.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email doesn't exist"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
		return
	}

	// If validated - generating token ....
	token, err := utils.GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role, "user": user})
}

// registration

func Register(c *gin.Context) {
	var phone models.User
	var duplicate models.User
	var user RegisterDate
	var Library models.Library
	var NewUser models.User
	var NewLibrary models.Library

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// dupliacte library
	lib := models.DB.Where("name = ?", user.LibraryName).First(&Library)
	if lib.Error == nil {
		// library found
		c.JSON(http.StatusConflict, gin.H{"message": "library with the same name already exists"})
		return
	}

	// duplicate user
	find := models.DB.Where("email = ?", user.Email).First(&duplicate)
	if find.Error == nil {
		// user found
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	// duplicate contact
	contact := models.DB.Where("contact_number = ?", user.ContactNumber).First(&phone)
	if contact.Error == nil {
		// contact found
		c.JSON(http.StatusConflict, gin.H{"error": "Contact Number already in use"})
		return
	}

	// create library
	NewLibrary = models.Library{Name: user.LibraryName}
	newLib := models.DB.Create(&NewLibrary)
	if newLib.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create library"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	NewUser = models.User{Name: user.Name, Email: user.Email, Password: string(hashedPassword), Role: user.Role, ContactNumber: user.ContactNumber, LibID: &NewLibrary.ID}

	user.Password = string(hashedPassword)

	if res := models.DB.Create(&NewUser); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// If validated - generating token ....
	token, err := utils.GenerateToken(&NewUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "role": user.Role, "user": user})
}

func OnboardAdmin(c *gin.Context) {
	var user models.Users
	var check models.Users
	var phone models.Users

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something not in favor"})
		return
	}

	find := models.DB.Where("email = ?", user.Email).First(&check)
	if find.Error == nil {
		// user found
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	contact := models.DB.Where("contact_number = ?", user.ContactNumber).First(&phone)
	if contact.Error == nil {
		// contact found
		c.JSON(http.StatusConflict, gin.H{"error": "Contact Number already in use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	users := models.DB.Create(&models.Users{Email: user.Email, Name: user.Name, ContactNumber: user.ContactNumber, Role: user.Role, LibID: user.LibID,
		Password: string(hashedPassword)})

	if users.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "admin has been created"})
}

func OnboardReader(c *gin.Context) {
	var user models.Users
	var check models.Users
	var phone models.Users

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something not in favor"})
		return
	}

	find := models.DB.Where("email = ?", user.Email).First(&check)
	if find.Error == nil {
		// user found
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	contact := models.DB.Where("contact_number = ?", user.ContactNumber).First(&phone)
	if contact.Error == nil {
		// contact found
		c.JSON(http.StatusConflict, gin.H{"error": "Contact Number already in use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	users := models.DB.Create(&models.Users{Email: user.Email, Name: user.Name, ContactNumber: user.ContactNumber, Role: user.Role, LibID: user.LibID,
		Password: string(hashedPassword)})

	if users.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Reader has been created"})
}
