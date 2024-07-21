package controllers

import (
	"fmt"
	"library/task/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// library and user creation

func CreateLibrary(c *gin.Context) {
	var library models.Library

	err := c.ShouldBindJSON(&library)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	lib := models.DB.Create(&models.Library{Name: library.Name})
	if lib.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "library has been created"})
}

func CreateUser(c *gin.Context) {
	var user models.Users

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something not in favor"})
		return
	}

	users := models.DB.Create(&models.Users{Email: user.Email, Name: user.Name, ContactNumber: user.ContactNumber, Role: user.Role, LibID: user.LibID, Password: user.Password})
	if users.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "user has been created"})
}

// admin panel
func CreateBookInventory(c *gin.Context) {

	var book models.BookInventory
	err := c.ShouldBindJSON(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something happened wrong", "err": err.Error()})
		return
	}

	books := models.DB.Create(&book)
	if books.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something happen wrong"})
		return
	}

	qrContent := fmt.Sprintf("Title:-%s, Publisher:-%s, AvailableCopies:-%d, TotalCopies:-%d, Authors:-%s", book.Title, book.Publisher, book.AvailableCopies, book.TotalCopies, book.Authors)
	qrCodePath := fmt.Sprintf("qrcodes/book_%d.png", book.ISBN)
	err = qrcode.WriteFile(qrContent, qrcode.Medium, 256, qrCodePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error generating QR code", "error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book Inventory has been created"})
}

func CreateRequestEvents(c *gin.Context) {
	var req models.RequestEvents

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	requests := models.DB.Create(&models.RequestEvents{ReqID: req.ReqID, BookID: req.BookID, ReaderID: req.ReaderID, RequestDate: req.RequestDate, ApprovalDate: req.ApprovalDate, ApproverID: req.ApproverID, RequestType: "issue"})

	if requests.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Request Events has been created"})
}

func CreateIssueRequests(c *gin.Context) {
	var issue models.RequestEvents

	err := c.ShouldBindJSON(&issue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	reqissue := models.DB.Create(&models.RequestEvents{ReqID: issue.ReqID, BookID: issue.BookID, ReaderID: issue.ReaderID, RequestDate: time.Now(), RequestType: issue.RequestType})

	if reqissue.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Issue Requests has been created"})
}

func ReturnRequests(c *gin.Context) {
	var issue models.RequestEvents

	err := c.ShouldBindJSON(&issue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	reqissue := models.DB.Create(&models.RequestEvents{ReqID: issue.ReqID, BookID: issue.BookID, ReaderID: issue.ReaderID, RequestDate: issue.RequestDate, RequestType: "return"})

	if reqissue.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book has been returned and updated in registry"})
}

// GET all event requests
func Requests(c *gin.Context) {
	libid := c.Param("libid")
	fmt.Println("libid:", libid)
	var requests []models.RequestEvents
	res := models.DB.Where("reader_id IN (SELECT id FROM users WHERE lib_id = ? )", libid).Find(&requests)
	if res.RowsAffected != 0 {
		c.JSON(200, requests)

	} else {
		c.AbortWithStatusJSON(404, "not found")

	}
}

// show all issued books to users
func Issued(c *gin.Context) {
	user := c.Param("user")
	var issuing []models.IssueRegistry
	res := models.DB.Where("reader_id = ? ", user).Find(&issuing)
	if res.RowsAffected != 0 {
		c.JSON(200, issuing)

	} else {
		c.AbortWithStatusJSON(404, "not found")

	}
}

// GET all books in the inventory
func FetchAllBooks(c *gin.Context) {
	libid := c.Param("libid")
	var books []models.BookInventory
	if err := models.DB.Where("lib_id = ?", libid).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func DeleteBook(c *gin.Context) {
	var Book models.BookInventory
	id := c.Param("id")

	// bind
	res := models.DB.Where("isbn = ?", id).First(&Book)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book is not available"})
		return
	}

	if Book.AvailableCopies > 0 {
		Book.AvailableCopies -= 1
		Book.TotalCopies -= 1
		save := models.DB.Save(&Book)
		if save.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error updating the book"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Book Deleted"})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No copies are available to delete"})
	}
}

func UpdateBook(c *gin.Context) {

	var data models.UpdateData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	res := models.DB.Model(&models.BookInventory{}).
		Where("isbn = ?", data.ISBN).
		Updates(map[string]interface{}{
			"title":            data.Title,
			"authors":          data.Authors,
			"publisher":        data.Publisher,
			"version":          data.Version,
			"total_copies":     data.TotalCopies,
			"available_copies": data.AvailableCopies,
			"lib_id":           data.LibID,
		})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update book details", "error": res.Error.Error()})
		return
	}

	var updatedBook models.BookInventory
	if err := models.DB.First(&updatedBook, "isbn = ?", data.ISBN).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch updated book details", "error": err.Error()})
		return
	}

	// Generate QR
	qrContent := fmt.Sprintf("Title: %s, Publisher: %s, AvailableCopies: %d, TotalCopies: %d, Authors: %s", updatedBook.Title, updatedBook.Publisher, updatedBook.AvailableCopies, updatedBook.TotalCopies, updatedBook.Authors)
	qrCodePath := fmt.Sprintf("qrcodes/book_%d.png", updatedBook.ISBN)
	if err := qrcode.WriteFile(qrContent, qrcode.Medium, 256, qrCodePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate QR code", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

func ApproveDisapprove(c *gin.Context) {
	var data models.ApproveData
	var Event models.RequestEvents
	var Book models.BookInventory

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something wrong"})
		return
	}
	fmt.Println(data.ReqID)
	res := models.DB.First(&Event, "req_id = ?", data.ReqID)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	// update
	Event.ApproverID = &data.ID
	Event.ApprovalDate = &[]time.Time{time.Now()}[0]
	save := models.DB.Save(&Event)

	if save.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something save wrong"})
		return
	}

	if Event.RequestType == "issue" {
		// update the availability
		book := models.DB.Where("isbn = ?", Event.BookID).First(&Book)
		if book.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find book"})
			return
		}

		if Book.AvailableCopies > 0 {
			Book.AvailableCopies = Book.AvailableCopies - 1
		}

		saveBook := models.DB.Save(&Book)

		if saveBook.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving the book"})
			return
		}
		reg := models.IssueRegistry{
			ISBN:           Book.ISBN,
			ReaderID:       Event.ReaderID,
			IssueApproveID: data.ID,
			IssueStatus:    "issued",
			IssueDate:      time.Now(),
		}

		reg.ExpectedReturnDate = time.Now().Add(time.Hour * 24 * 14)
		err = models.DB.Create(&reg).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error saving the book"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Issue request approved"})

	} else {
		// update the availability
		book := models.DB.Where("isbn = ?", Event.BookID).First(&Book)
		if book.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find book"})
			return
		}

		Book.AvailableCopies += 1

		saveBook := models.DB.Save(&Book)
		if saveBook.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving the book"})
			return
		}

		// updating to registry
		reg := models.IssueRegistry{
			ISBN:           Book.ISBN,
			ReaderID:       Event.ReaderID,
			IssueApproveID: data.ID,
			IssueStatus:    "return",
			IssueDate:      time.Now(),

			ReturnDate:       time.Now(),
			ReturnApproverID: data.ID,
		}

		err = models.DB.Create(&reg).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error saving the book"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Issue request approved"})

	}

}

func Disapprove(c *gin.Context) {
	var data models.RejectData
	var Event models.RequestEvents

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something wrong"})
		return
	}
	fmt.Println(data.ReqID)
	res := models.DB.First(&Event, "req_id = ?", data.ReqID)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	Event.RejectID = &data.ID
	Event.RejectDate = &[]time.Time{time.Now()}[0]
	save := models.DB.Save(&Event)

	if save.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something save wrong"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Request Rejected"})

}

// reader panel
func CreateRequest(c *gin.Context) {
	var data models.EventsData
	var check models.BookInventory
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	res := models.DB.Create(&models.RequestEvents{ReaderID: data.ReaderID, BookID: data.BookID, RequestDate: time.Now(), RequestType: "issue"})

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}

	if check.AvailableCopies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No copies are available to request for"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Issue request created"})

}

func SearchBookBy(c *gin.Context) {
	var data models.SearchBook
	var book models.BookInventory
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "something not in favor"})
		return
	}

	// Search for the book in the database
	res := models.DB.Where("title = ? OR publisher = ? OR authors = ?", data.Query, data.Query, data.Query).First(&book)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}

	// Generate QR code URL
	qrCodeURL := fmt.Sprintf("qrcodes/book_%d.png", book.ISBN)

	// Merge book and qrCodeURL into a single object
	responseObject := struct {
		Message   string               `json:"message"`
		Book      models.BookInventory `json:"book"`
		QRCodeURL string               `json:"qr_code_url"`
	}{
		Message:   "Book found",
		Book:      book,
		QRCodeURL: qrCodeURL,
	}

	c.IndentedJSON(http.StatusOK, responseObject)
}
