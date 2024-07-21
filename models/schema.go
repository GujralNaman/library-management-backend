package models

import "time"

type Library struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Users struct {
	ID            uint    `json:"id" gorm:"primary_key"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	ContactNumber string  `json:"contactNumber"`
	Role          string  `json:"role"`
	LibID         uint    `json:"libid"`
	Password      string  `json:"password"`
	Library       Library `gorm:"foreignKey:LibID;references:ID"`
}

type BookInventory struct {
	ISBN            uint    `json:"isbn" gorm:"primary_key"`
	LibID           uint    `json:"libID"`
	Title           string  `json:"title"`
	Authors         string  `json:"authors"`
	Publisher       string  `json:"publisher"`
	Version         uint    `json:"version"`
	TotalCopies     uint    `json:"totalCopies"`
	AvailableCopies uint    `json:"availableCopies"`
	Library         Library `gorm:"foreignKey:LibID;references:ID"`
}

type RequestEvents struct {
	ReqID         uint          `json:"reqID" gorm:"primary_key"`
	BookID        uint          `json:"bookID"`
	ReaderID      uint          `json:"readerID"`
	RequestDate   time.Time     `json:"requestdate"`
	ApprovalDate  *time.Time    `json:"approvalDate"`
	ApproverID    *uint         `json:"approverID"`
	RejectDate    *time.Time    `json:"rejectDate"`
	RejectID      *uint         `json:"rejectID"`
	RequestType   string        `json:"requestType"`
	BookInventory BookInventory `gorm:"foreignKey:BookID;references:ISBN"`
	Users         Users         `gorm:"foreignKey:ReaderID;references:ID"`
}

type IssueRegistry struct {
	IssueID            uint          `json:"reqID" gorm:"primary_key"`
	ISBN               uint          `json:"isbn"`
	BookInventory      BookInventory `gorm:"foreignKey:ISBN;references:ISBN"`
	ReaderID           uint          `json:"readerID"`
	IssueApproveID     uint          `json:"issueApproveID"`
	IssueStatus        string        `json:"issueStatus"`
	IssueDate          time.Time     `json:"issueDate"`
	ExpectedReturnDate time.Time     `json:"expectedReturnDate"`
	ReturnDate         time.Time    `json:"returnDate"`
	ReturnApproverID   uint         `json:"returnApproverID"`
	Users              Users         `gorm:"foreignKey:ReaderID;references:ID"`
}

type DeleteData struct {
	ID uint `json:"id"`
}

type UpdateData struct {
	ISBN            uint   `json:"isbn"`
	Title           string `json:"title"`
	Authors         string `json:"authors"`
	Publisher       string `json:"publisher"`
	Version         uint   `json:"version"`
	TotalCopies     uint   `json:"totalCopies"`
	AvailableCopies uint   `json:"availableCopies"`
	LibID           uint   `json:"lib_id"`
}

type EventsData struct {
	BookID   uint `json:"bookID"`
	ReaderID uint `json:"readerID"`
}

type ApproveData struct {
	ReqID uint `json:"reqID"`
	ID    uint `json:"id"`
}

type RejectData struct {
	ReqID uint `json:"reqID"`
	ID    uint `json:"id"`
}

type SearchBook struct {
	Query string `json:"query"`
}

type AuthInput struct {
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Name          string `json:"name"`
	ContactNumber string `json:"contactNumber"`
	Role          string `json:"role"`
}
