package constants

import "errors"

// User roles
const (
	RoleAdmin = "ADMIN"
	RoleStaff = "STAFF"
	RoleUser  = "USER"
)

// BookStock status
const (
	StatusAvailable = "AVAILABLE"
	StatusBorrowed  = "BORROWED"
	StatusDamaged   = "DAMAGED"
	StatusLost      = "LOST"
)

// Journal status
const (
	StatusActive   = "BORROWED"
	StatusReturned = "RETURNED"
	StatusOverdue  = "OVERDUE"
)

// Error messages
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrBookNotFound       = errors.New("book not found")
	ErrBookstockNotFound  = errors.New("book stock not found")
	ErrJournalNotFound    = errors.New("journal not found")
	ErrMediaNotFound      = errors.New("media not found")
	ErrChargeNotFound     = errors.New("charge not found")
	ErrBookNotAvailable   = errors.New("book is not available")
	ErrInternalServer     = errors.New("internal server error")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden access")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredential  = errors.New("invalid credential")
)

// Success messages
const (
	MsgLoginSuccess    = "Login successful"
	MsgRegisterSuccess = "Registration successful"
	MsgDeleteSuccess   = "Successfully deleted"
	MsgUpdateSuccess   = "Successfully updated"
	MsgCreateSuccess   = "Successfully created"
	MsgBorrowSuccess   = "Book successfully borrowed"
	MsgReturnSuccess   = "Book successfully returned"
)

// Default values
const (
	DefaultDailyLateFee = 1000.0 // Default late fee per day
)
