package dto

type User struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ChangePasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
type UpdateAddressRequest struct {
	Email      string `json:"email" binding:"required"`
	NewAddress string `json:"new_address" binding:"required"`
}
