package database

import (
	"context"
	"errors"
	"fmt"
	"ginDatabaseMhs/model/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// adaptop pattern
type MahasiswaRepository struct {
	DB       *gorm.DB
	S3Bucket *s3.Client
}

func NewMahasiswaRepository(DB *gorm.DB, s3Bucket *s3.Client) *MahasiswaRepository {
	return &MahasiswaRepository{
		DB:       DB,
		S3Bucket: s3Bucket,
	}
}

//func (t MahasiswaRepository) GetAll() ([]entity.Mahasiswa, error) {
//	var data []entity.Mahasiswa
//
//	result := t.DB.Preload("Attachments").Preload("User").Find(&data)
//	return data, result.Error
//}

func (t *MahasiswaRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	if err := t.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (t *MahasiswaRepository) DeleteUserByIDAndRole(userID int64, role string) error {
	// Cek apakah pengguna dengan ID dan peran tertentu ada dalam database
	user := &entity.User{}
	err := t.DB.Where("id = ? AND role = ?", userID, role).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found or does not have the specified role")
		}
		return err
	}

	// Hapus pengguna dari database
	if err := t.DB.Delete(user).Error; err != nil {
		return err
	}

	return nil
}

func (t MahasiswaRepository) GetAllUserByID(UserID int64) ([]entity.User_data, error) {
	var data []entity.User_data

	// Ambil semua data berdasarkan user_id
	result := t.DB.Preload("Attachments").Where("user_id = ?", UserID).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (t MahasiswaRepository) GetByID(mhsID, userID int64) (*entity.User_data, error) {
	var data entity.User_data
	result := t.DB.Preload("Attachments").Where("id = ? AND user_id = ?", mhsID, userID).First(&data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &data, result.Error
}

func (t MahasiswaRepository) Create(mahasiswa *entity.User_data) (*entity.User_data, error) {
	result := t.DB.Create(mahasiswa)
	return mahasiswa, result.Error

}

func (t MahasiswaRepository) GetMahasiswaByNameAndEmail(name, email string) (*entity.User_data, error) {
	var data entity.User_data
	result := t.DB.Where("name = ? AND email = ?", name, email).First(&data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &data, nil
}

func (t MahasiswaRepository) Update(mhsID, userID int64, updates map[string]interface{}) (*entity.User_data, error) {
	var data entity.User_data
	result := t.DB.Model(&data).Where("id = ? AND user_id = ?", mhsID, userID).Updates(updates)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &data, result.Error
}

func (t MahasiswaRepository) UpdatetoAtch(mhs *entity.User_data) error {
	err := t.DB.Save(mhs).Error
	return err
}

func (t MahasiswaRepository) Delete(mhsID, userID int64) (int64, error) {
	data := entity.User_data{}

	// Fetch the data by ID and user_id
	if err := t.DB.Where("id = ? AND user_id = ?", mhsID, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If data not found, return 0 RowsAffected
			return 0, nil
		}
		return 0, err
	}

	// Delete the fetched data
	result := t.DB.Delete(&data)
	return result.RowsAffected, result.Error
}
func (t MahasiswaRepository) CreateAdmin(admin *entity.Admin) error {
	result := t.DB.Create(admin)
	return result.Error
}

func (t MahasiswaRepository) CreateUser(user *entity.User) error {
	if err := t.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (t MahasiswaRepository) GetUserByUsernameOrEmail(username, email string) (*entity.User, error) {
	var user entity.User
	result := t.DB.Where("username = ? OR email = ?", username, email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil

}

/////////////////////////////////////////

func (t *MahasiswaRepository) CreateAttachment(mhsID int64, path string, order int64) (*entity.Attachment, error) {
	attachment := &entity.Attachment{
		UserID:          mhsID,
		Path:            path,
		AttachmentOrder: order,
	}
	if err := t.DB.Create(attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}

func (t *MahasiswaRepository) UploadFileS3Atch(file *multipart.FileHeader, mhsID, userID int64) (*entity.Attachment, error) {
	//Mengambil data berdasarkan ID dan user_id

	dataMhs := &entity.User_data{}
	if err := t.DB.Where("id = ? AND user_id = ?", mhsID, userID).First(dataMhs).Error; err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer src.Close()

	// bikin nama file yang uniq untuk menghindari konflik
	uniqueFilename := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(file.Filename))

	// Upload the file to S3
	bucketName := "bucketwithrey"
	objectKey := uniqueFilename
	_, err = t.S3Bucket.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   src,
		//ACL:    types.ObjectCannedACLPublicRead, // Optional: Mengatur ACL agar file yang diunggah dapat diakses oleh publik
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Return the public URL of the uploaded file
	publicURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey)

	// Create an attachment record in the database
	var attachmentOrder int64 = 1 // Set the initial attachment_order
	// Get the count of existing attachments for the dataMhs
	existingAttachmentCount := int64(0)
	t.DB.Model(&entity.Attachment{}).Where("user_id = ?", mhsID).Count(&existingAttachmentCount)
	attachmentOrder = existingAttachmentCount + 1 // Set attachment_order dynamically

	// Create an attachment record in the database
	attachment := &entity.Attachment{
		UserID:          mhsID,
		Path:            publicURL,
		AttachmentOrder: attachmentOrder, // atur order
		Timestamp:       time.Now(),
	}
	err = t.DB.Create(attachment).Error
	if err != nil {
		return nil, err
	}

	return attachment, nil
}
func (t *MahasiswaRepository) UpdateWithAttachments(mhs *entity.User_data) error {
	return t.DB.Transaction(func(tx *gorm.DB) error {
		// Pertama, hapus semua lampiran yang ada yang terkait dengan data
		if err := t.DB.Where("user_id = ?", mhs.ID).Delete(&entity.Attachment{}).Error; err != nil {
			return err
		}

		// Next, create new attachment records
		for i := range mhs.Attachments {
			attachment := &mhs.Attachments[i]
			if err := t.DB.Create(attachment).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (t *MahasiswaRepository) UploadFileS3Buckets(file io.Reader, fileName string) (*string, error) {
	bucketName := "bucketwithrey"
	objectKey := fileName

	_, err := t.S3Bucket.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Return the public URL of the uploaded file
	publicURL := aws.String(fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectKey))

	return publicURL, nil
}

func (t *MahasiswaRepository) UploadFileLocalAtch(file *multipart.FileHeader, mhsID, userID int64) (*entity.Attachment, error) {
	// Fetch the data by ID and user_id
	data := &entity.User_data{}
	if err := t.DB.Where("id = ? AND user_id = ?", mhsID, userID).First(data).Error; err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// bikin nama file yang uniq untuk menghindari konflik
	uniqueFilename := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(file.Filename))

	// Upload the file to Local
	uploadDir := "uploads"

	// Buat direktori unggahan jika belum ada
	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return nil, err
	}

	// Create the destination file
	dest, err := os.Create(filepath.Join(uploadDir, uniqueFilename))
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	// Copy file nya ke file tujuan
	_, err = io.Copy(dest, src)
	if err != nil {
		return nil, err
	}
	// Return the local file path
	localFilePath := filepath.Join(uploadDir, uniqueFilename)

	// Create an attachment record in the database
	var attachmentOrder int64 = 1 // Set the initial attachment_order
	// Get the count of existing attachments for the dataMhs
	existingAttachmentCount := int64(0)
	t.DB.Model(&entity.Attachment{}).Where("user_id = ?", mhsID).Count(&existingAttachmentCount)
	attachmentOrder = existingAttachmentCount + 1 // Set attachment_order dynamically

	// Create an attachment record in the database
	attachment := &entity.Attachment{
		UserID:          mhsID,
		Path:            localFilePath,
		AttachmentOrder: attachmentOrder, // atur order
		Timestamp:       time.Now(),
	}
	err = t.DB.Create(attachment).Error
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (t *MahasiswaRepository) SearchMahasiswaByUser(userID int64, search string, page, perPage int) ([]entity.User_data, int64, error) {
	var dataMhs []entity.User_data

	// Menghitung total data
	var total int64
	t.DB.Model(&entity.User_data{}).Where("user_id = ? AND title LIKE ?", userID, "%"+search+"%").Count(&total)

	// Mengambil data dengan paginasi
	offset := (page - 1) * perPage
	err := t.DB.Where("user_id = ? AND title LIKE ?", userID, "%"+search+"%").
		Offset(offset).Limit(perPage).
		Preload("Attachments").Find(&dataMhs).Error

	return dataMhs, total, err
}

// //////////////////////////////////////////////////////
// GetRoleByName mengambil data peran berdasarkan nama peran
func (t *MahasiswaRepository) GetRoleByName(roleName string) (*entity.Roles, error) {
	var role entity.Roles
	if err := t.DB.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
