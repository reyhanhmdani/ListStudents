package repository

import (
	"ginDatabaseMhs/model/entity"
	"io"
	"mime/multipart"
)

type MahasiswaRepository interface {
	GetAllUsers() ([]entity.User, error)
	DeleteUserByIDAndRole(userID int64, role string) error

	// ALL //////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	GetAllUserByID(UserID int64) ([]entity.User_data, error)
	GetByID(mhsID, userID int64) (*entity.User_data, error)
	Create(mahasiswa *entity.User_data) (*entity.User_data, error)
	GetMahasiswaByNameAndEmail(name, email string) (*entity.User_data, error)
	Update(mhsID, userID int64, updates map[string]interface{}) (*entity.User_data, error)
	UpdatetoAtch(todo *entity.User_data) error
	CreateAdmin(admin *entity.Admin) error
	Delete(mhsID, userID int64) (int64, error)
	CreateUser(user *entity.User) error
	GetUserByUsernameOrEmail(username, email string) (*entity.User, error)
	//UploadTodoFileS3(file *multipart.FileHeader, url string) error
	//UploadTodoFileLocal(file *multipart.FileHeader, url string) error
	/////////////////////
	CreateAttachment(mhsID int64, path string, order int64) (*entity.Attachment, error)
	UploadFileS3Atch(file *multipart.FileHeader, mhsID, userID int64) (*entity.Attachment, error)
	UpdateWithAttachments(mhs *entity.User_data) error
	UploadFileS3Buckets(file io.Reader, fileName string) (*string, error)
	UploadFileLocalAtch(file *multipart.FileHeader, mhsID, userID int64) (*entity.Attachment, error)
	SearchMahasiswaByUser(userID int64, search string, page, perPage int) ([]entity.User_data, int64, error)
	//
	GetRoleByName(roleName string) (*entity.Roles, error)
}
