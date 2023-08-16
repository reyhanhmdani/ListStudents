package request

type MahasiswaCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
	Birthdate   string `json:"birthdate"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email" binding:"required"`
	UserID      int64  `json:"user_id"`
}

type MahasiswaUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Age         int    `json:"age"`
	Address     string `json:"address"`
	Birthdate   string `json:"birthdate"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email" binding:"required"`
}

func (r *MahasiswaUpdateRequest) ReqMhs() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.Name != "" {
		updates["name"] = r.Name
	}
	updates["email"] = r.Email

	return updates
}
