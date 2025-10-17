package services

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/model/constant"
	"crypsis-backend/internal/repository"
	"fmt"
	"log/slog"
	"mime/multipart"

	"gorm.io/gorm"
)

type FileService struct {
	cryptoService         CryptographicInterface
	storageService        StorageInterface
	kmsService            KMSInterface
	fileRepository        repository.FileRepository
	fileLogsRepository    repository.FileLogsRepository
	applicationRepository repository.ApplicationRepository
	adminRepository       repository.AdminRepository
	db                    *gorm.DB
	keyConfig             *model.KeyConfig
	bucketName            string
	hashMethod            string
	hashEncryptedFile     bool
	encryptionMethod      string
}

func NewFileService(params FileServiceParams) FileInterface {
	return &FileService{
		cryptoService:         params.CryptoService,
		storageService:        params.StorageService,
		kmsService:            params.KMSService,
		fileRepository:        params.FileRepository,
		fileLogsRepository:    params.FileLogsRepository,
		applicationRepository: params.ApplicationRepository,
		adminRepository:       params.AdminRepository,
		db:                    params.DB,
		keyConfig:             params.KeyConfig,
		bucketName:            params.BucketName,
		hashMethod:            params.HashMethod,
		hashEncryptedFile:     params.HashEncryptedFile,
		encryptionMethod:      params.EncryptionMethod,
	}
}

func (c *FileService) ListFiles(ctx context.Context, clientID string, limit, offset int, sortBy, order string) (int64, *[]model.FileResponse, error) {
	if clientID == "" {
		return 0, nil, model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return 0, nil, err
	}

	fmt.Println("Validated App ID:", validatedAppID)
	fmt.Println("Offset:", offset, "Limit:", limit, "SortBy:", sortBy, "Order:", order)

	count, files, err := c.fileRepository.GetListFiles(ctx, validatedAppID, offset, limit, sortBy, order)
	if err != nil {
		return 0, nil, err
	}

	var fileResponse []model.FileResponse
	for _, file := range files {
		fileResponse = append(fileResponse, model.FileResponse{
			ID:        file.ID,
			Name:      file.Name,
			Size:      file.Size,
			MimeType:  file.MimeType,
			UpdatedAt: file.UpdatedAt.String(),
		})
	}

	return count, &fileResponse, nil
}

func (c *FileService) ListFilesForAdmin(ctx context.Context, adminID, appID string, limit, offset int, sortBy, order string) (int64, *[]model.FileResponse, error) {
	if adminID == "" {
		return 0, nil, model.ErrInvalidInput
	}

	var count int64
	var files []entity.Files
	var err error

	if appID == "" {
		// Admin can list all files
		count, files, err = c.fileRepository.GetListFilesForAdmin(ctx, offset, limit, sortBy, order)
		if err != nil {
			return 0, nil, err
		}
	} else {
		// Admin can list files for a specific client
		count, files, err = c.fileRepository.GetListFiles(ctx, appID, offset, limit, sortBy, order)
		if err != nil {
			return 0, nil, err
		}
	}

	if c.adminRepository.IsAdmin(ctx, adminID) == false {
		return 0, nil, model.AdminErrNotFound
	}

	var fileResponse []model.FileResponse
	for _, file := range files {
		fileResponse = append(fileResponse, model.FileResponse{
			ID:        file.ID,
			Name:      file.Name,
			OwnerID:   file.AppID,
			Size:      file.Size,
			MimeType:  file.MimeType,
			UpdatedAt: file.UpdatedAt.String(),
			Deleted:   file.DeletedAt.Valid,
		})
	}

	return count, &fileResponse, nil
}

func (c *FileService) UploadFile(ctx context.Context, clientID, fileName string, input multipart.File) (fileUID string, err error) {
	// Check Client ID
	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return "", err
	}

	// Generate file UID
	fileUID = helper.GenerateCustomUUID().String()

	// Generate Key and Encrypt file
	encryptedFile, metaDataDTO, err := c.encryptFile(ctx, "", fileUID, input)
	// encryptedFile, metaDataDTO, err := c.encryptFile(ctx, "", fileUID, input)
	if err != nil {
		return "", err
	}

	// Wrap key
	wrappedKey, err := c.cryptoService.EncryptString(c.keyConfig.KEK, metaDataDTO.Key)
	if err != nil {
		return "", err
	}

	// File to be saved to db
	fileToBeSaved := &entity.Files{
		ID:       fileUID,
		Name:     fileName,
		AppID:    validatedAppID,
		UserID:   "Not Available", // TO BE ADDED
		Size:     metaDataDTO.Size,
		MimeType: metaDataDTO.MimeType,
	}

	metadataToBeSaved := &entity.Metadata{
		ID:      helper.GenerateCustomUUID().String(),
		FileID:  fileToBeSaved.ID,
		Hash:    metaDataDTO.Hash,
		EncHash: metaDataDTO.EncryptedFileHash,
		KeyUID:  metaDataDTO.KeyUID,
		EncKey:  wrappedKey,
		KeyAlgo: c.encryptionMethod,
	}

	// Create multipart file
	toBeUploadedFile, tobeUploadSize, err := helper.CreateMultipartFileFromBytes(encryptedFile, fileName)
	if err != nil {
		return "", err
	}

	// Upload encrypted file
	transcationResponse, err := c.storageService.UploadFile(ctx, c.bucketName, createFileName(fileToBeSaved.ID), toBeUploadedFile, tobeUploadSize)
	if err != nil {
		return "", err
	}

	metadataToBeSaved.VersionID = transcationResponse.VersionID
	fileToBeSaved.Location = transcationResponse.Location

	err = c.fileRepository.CreateFileWithMetadata(ctx, fileToBeSaved, metadataToBeSaved)
	if err != nil {
		return "", err
	}

	// Save to log
	c.saveFileLog(ctx, validatedAppID, fileToBeSaved.ID, constant.ActorTypeClient, string(constant.ActionTypeUpload), fileName)
	return fileToBeSaved.ID, nil
}

func (c *FileService) DownloadFile(ctx context.Context, clientID, fileUID string) ([]byte, string, error) {
	if fileUID == "" {
		return nil, "", model.ErrInvalidInput
	}

	// Check Client ID
	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return nil, "", err
	}

	//check file existence
	fileMetaData, err := c.fileRepository.GetMetadataByAppIDAndFileID(ctx, validatedAppID, fileUID)
	if err != nil {
		return nil, "", err
	}

	isExist, _, err := c.storageService.Exists(ctx, c.bucketName, createFileName(fileMetaData.FileID))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", model.ErrFileNotFound
	}

	//save to log
	c.saveFileLog(ctx, validatedAppID, fileMetaData.FileID, constant.ActorTypeClient, string(constant.ActionTypeDownload), fileMetaData.File.Name)

	//unwrap key
	key, err := c.cryptoService.DecryptString(c.keyConfig.KEK, fileMetaData.EncKey)
	if err != nil {
		return nil, "", err
	}

	//download file
	encryptedFile, err := c.storageService.DownloadFile(ctx, c.bucketName, createFileName(fileMetaData.FileID))
	if err != nil {
		return nil, "", err
	}

	//decrypt file
	decryptedFile, err := c.decryptFile(key, fileMetaData.Hash, encryptedFile)
	if err != nil {
		return nil, "", err
	}
	return decryptedFile, fileMetaData.File.Name, nil
}

func (c *FileService) EncryptFile(ctx context.Context, clientID, fileName string, input multipart.File) ([]byte, string, error) {
	if fileName == "" || input == nil {
		return nil, "", model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return nil, "", err
	}

	//save to log
	c.saveFileLog(ctx, validatedAppID, fileName, constant.ActorTypeClient, string(constant.ActionTypeEncrypt), fileName)

	// Generate file UID
	fileUID := helper.GenerateCustomUUID().String()

	encryptedFile, metadataDTO, err := c.encryptFile(ctx, "", fileUID, input)
	if err != nil {
		return nil, "", err
	}

	//key wrapping
	wrappedKey, err := c.cryptoService.EncryptString(c.keyConfig.KEK, metadataDTO.Key)
	if err != nil {
		return nil, "", err
	}

	// File to be saved to db
	fileToBeSaved := &entity.Files{
		ID:       fileUID,
		Name:     fileName,
		AppID:    validatedAppID,
		UserID:   "Not Available", // TO BE ADDED
		Size:     metadataDTO.Size,
		MimeType: metadataDTO.MimeType,
	}

	metadataToBeSaved := &entity.Metadata{
		ID:      helper.GenerateCustomUUID().String(),
		FileID:  fileToBeSaved.ID,
		Hash:    metadataDTO.Hash,
		EncHash: metadataDTO.EncryptedFileHash,
		KeyUID:  metadataDTO.KeyUID,
		EncKey:  wrappedKey,
		KeyAlgo: c.encryptionMethod,
	}

	err = c.fileRepository.CreateFileWithMetadata(ctx, fileToBeSaved, metadataToBeSaved)
	if err != nil {
		return nil, "", err
	}

	return encryptedFile, fileUID, nil
}

func (c *FileService) DecryptFile(ctx context.Context, clientID, fileUID string, input multipart.File) ([]byte, error) {
	// Input validation
	if clientID == "" || fileUID == "" || input == nil {
		return nil, model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return nil, err
	}

	// Check file existence
	fileMetaData, err := c.fileRepository.GetMetadataByAppIDAndFileID(ctx, validatedAppID, fileUID)
	if err != nil {
		return nil, err
	}

	//save to log
	c.saveFileLog(ctx, validatedAppID, fileMetaData.FileID, constant.ActorTypeClient, string(constant.ActionTypeDecrypt), fileMetaData.File.Name)

	//unwrap key
	key, err := c.cryptoService.DecryptString(c.keyConfig.KEK, fileMetaData.EncKey)
	if err != nil {
		return nil, err
	}

	//read file
	encryptedFile, _, _, err := helper.GetFileBytesFromMultipart(input)
	if err != nil {
		return nil, err
	}

	//decrypt file
	decryptedFile, err := c.decryptFile(key, fileMetaData.Hash, encryptedFile)
	if err != nil {
		return nil, err
	}
	return decryptedFile, nil
}

func (c *FileService) GetFileMetadata(ctx context.Context, clientID, fileUID string) (*model.FileMetadataResponse, error) {
	// Input validation
	if clientID == "" || fileUID == "" {
		return nil, model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return nil, err
	}

	fmt.Println("Validated App ID:", validatedAppID)

	// check file existence
	result, err := c.fileRepository.GetMetadataByAppIDAndFileID(ctx, validatedAppID, fileUID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, model.ErrFileNotFound
	}

	return &model.FileMetadataResponse{
		ID:         result.ID,
		Name:       result.File.Name,
		Size:       result.File.Size,
		MimeType:   result.File.MimeType,
		VersionID:  result.VersionID,
		Hash:       result.Hash,
		BucketName: result.File.BucketName,
		Location:   result.File.Location,
		CreatedAt:  result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil

}

func (c *FileService) UpdateFile(ctx context.Context, clientID, fileUID, fileName string, input multipart.File) (string, error) {
	// Input validation
	if clientID == "" || fileUID == "" || fileName == "" || input == nil {
		return "", model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return "", err
	}

	// check file existence
	fileMetaData, err := c.fileRepository.GetMetadataByAppIDAndFileID(ctx, validatedAppID, fileUID)
	if err != nil {
		return "", err
	}

	// Unwrap Key
	key, err := c.cryptoService.DecryptString(c.keyConfig.KEK, fileMetaData.EncKey)
	if err != nil {
		return "", err
	}

	// Encrypt File
	encryptedFile, metaDataDTO, err := c.encryptFile(ctx, key, fileMetaData.FileID, input)
	if err != nil {
		return "", err
	}

	fileToBeUpdated := &entity.Files{
		ID:       fileMetaData.FileID,
		Name:     fileName,
		AppID:    validatedAppID,  // TO BE ADDED
		UserID:   "Not Available", // TO BE ADDED
		Size:     metaDataDTO.Size,
		MimeType: metaDataDTO.MimeType,
	}

	// Create multipart file
	toBeUploadedFile, tobeUploadSize, err := helper.CreateMultipartFileFromBytes(encryptedFile, fileName)
	if err != nil {
		return "", err
	}

	// Update File to Storage
	resp, err := c.storageService.UpdateFile(ctx, c.bucketName, createFileName(fileMetaData.FileID), toBeUploadedFile, tobeUploadSize)
	if err != nil {
		return "", err
	}

	// UPDATING NEW METADADATA
	fileMetaData.Hash = metaDataDTO.Hash
	fileMetaData.EncHash = metaDataDTO.EncryptedFileHash
	fileMetaData.VersionID = resp.VersionID
	fileToBeUpdated.Location = resp.Location

	// Update data IN DB
	err = c.fileRepository.UpdateFileAndMetadata(ctx, fileToBeUpdated, fileMetaData)
	if err != nil {
		return "", err
	}

	// save to log
	c.saveFileLog(ctx, validatedAppID, fileMetaData.FileID, constant.ActorTypeClient, string(constant.ActionTypeUpdate), fileMetaData.File.Name)

	return fmt.Sprintf("File %s updated successfully", fileMetaData.File.Name), nil
}

func (c *FileService) DeleteFile(ctx context.Context, clientID, fileUID string) error {
	if fileUID == "" || clientID == "" {
		return model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return err
	}

	// check file existence
	result, err := c.fileRepository.GetMetadataByAppIDAndFileID(ctx, validatedAppID, fileUID)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("file %s does not exist", fileUID)
	}

	// delete file in storage
	err = c.storageService.DeleteFile(ctx, c.bucketName, createFileName(result.FileID))
	if err != nil {
		return err
	}

	// soft delete file in DB
	err = c.fileRepository.Delete(ctx, result.FileID)
	if err != nil {
		return err
	}

	c.saveFileLog(ctx, validatedAppID, fileUID, constant.ActorTypeClient, string(constant.ActionTypeDelete), result.File.Name)
	return nil
}

func (c *FileService) RecoverFile(ctx context.Context, clientID, fileUID string) (string, error) {
	if fileUID == "" || clientID == "" {
		return "", model.ErrInvalidInput
	}

	validatedAppID, err := c.checkClientID(ctx, clientID)
	if validatedAppID == "" {
		return "", err
	}

	file, err := c.fileRepository.GetDeletedMetadataByAppIDAndFileID(ctx, validatedAppID, fileUID)
	if err != nil {
		return "", err
	} else if file == nil {
		return "", model.ErrFileNotFound
	}

	if !file.DeletedAt.Valid {
		return "", model.ErrFileAlreadyExists
	}

	err = c.storageService.RestoreFile(ctx, c.bucketName, createFileName(file.FileID), file.VersionID)
	if err != nil {
		return "", err
	}
	err = c.fileRepository.RestoreFile(ctx, file.FileID)
	if err != nil {
		return "", err
	}
	c.saveFileLog(ctx, validatedAppID, fileUID, constant.ActorTypeClient, string(constant.ActionTypeRecover), fileUID)
	return fmt.Sprintf("File %s recovered successfully", file.File.Name), nil
}

// ADMIN ONLY
func (c *FileService) ReKey(ctx context.Context, appID, keyUID string) (string, error) {
	if keyUID == "" {
		return "", model.ErrInvalidInput
	}

	if !c.keyConfig.KMSEnable {
		return "", fmt.Errorf("KMS is not enabled")
	}

	_, err := c.kmsService.ReKey(ctx, keyUID)
	if err != nil {
		return "", err
	}

	toBeUpdates := map[string]string{}

	c.saveFileLog(ctx, appID, constant.ActorTypeAdmin, "REKEY", string(constant.ActionTypeReKey), "")

	keyUIDs, err := c.fileRepository.GetAllKeyUIDs(ctx)
	for _, fileKeyUID := range keyUIDs {
		key, err := c.kmsService.ExportKey(ctx, fileKeyUID)
		if err != nil {
			continue
		}
		//key wrapping
		wrappedKey, err := c.cryptoService.EncryptString(c.keyConfig.KEK, key)
		if err != nil {
			continue
		}

		toBeUpdates[fileKeyUID] = wrappedKey
	}

	err = c.fileRepository.BatchUpdateEncKeys(ctx, toBeUpdates)
	if err != nil {
		return "", fmt.Errorf("failed to update keys: %w", err)
	}
	return "", nil
}

// ADMIN ONLY
func (c *FileService) ListLogs(ctx context.Context, limit, offset int, sortBy, order string) (int64, *[]model.FileLogResponse, error) {

	count, result, err := c.fileLogsRepository.List(ctx, offset, limit, sortBy, order)
	if err != nil {
		return 0, nil, err
	}

	var fileLogResponse []model.FileLogResponse
	for _, fileLog := range *result {
		fileLogResponse = append(fileLogResponse, model.FileLogResponse{
			ID:        fileLog.FileID,
			ActorID:   fileLog.ActorID,
			ActorType: fileLog.ActorType,
			Action:    fileLog.Action,
			Timestamp: fileLog.Timestamp,
			IP:        fileLog.IP,
			UserAgent: fileLog.UserAgent,
			Metadata:  fileLog.Metadata,
		})
	}

	return count, &fileLogResponse, nil
}

func (c *FileService) encryptFile(ctx context.Context, fileKey, fileUID string, file multipart.File) ([]byte, *model.MetaDataDTO, error) {
	var key string
	var keyUID string

	// Read file bytes
	fileBytes, fileSize, mimeType, err := helper.GetFileBytesFromMultipart(file)
	if err != nil {
		slog.Error("Failed to read file", slog.Any("error", err))
		return nil, nil, model.ErrFailedToReadFile
	}
	if fileSize == 0 {
		slog.Error("File is empty")
		return nil, nil, model.ErrFileIsEmpty
	}

	// Calculate file hash
	hashValue, err := c.cryptoService.HashFile(c.hashMethod, fileBytes)
	if err != nil {
		slog.Error("Failed to hash file", slog.Any("error", err))
		return nil, nil, model.ErrHashCalculationFailed
	}

	//	Use provided key
	if fileKey != "" {
		key = fileKey
		slog.Debug("Using provided key for encryption")
	} else if fileUID != "" { // Generate encryption key form KMS
		key, keyUID, err = c.getEncryptionKey(ctx, fileUID)
		if err != nil {
			slog.Error("Failed to generate key", slog.Any("error", err))
			return nil, nil, model.ErrKeyGenerationFailed
		}
		slog.Debug("Generated new key for encryption", slog.Int("key_length", len(key)))
	} else {
		return nil, nil, model.ErrFileUidOrKeyInvalid
	}

	// Encrypt file
	encryptedFile, err := c.cryptoService.EncryptFile(key, fileBytes)
	if err != nil {
		slog.Error("Failed to encrypt file", slog.Any("error", err))
		return nil, nil, model.ErrFileEncryptionFailed
	}

	// Prepare metadata
	metadata := c.createMetadataDTO(keyUID, key, mimeType, fileSize, hashValue, encryptedFile)
	return encryptedFile, metadata, nil
}

// getEncryptionKey generates or retrieves an encryption key
func (c *FileService) getEncryptionKey(ctx context.Context, fileUID string) (key, keyUID string, err error) {
	if c.keyConfig.KMSEnable {
		slog.Info("KMS is enabled, generating key from KMS")
		keyUID, err = c.kmsService.GenerateSymetricKey(ctx, fileUID)
		if err != nil {
			return "", "", model.ErrFailedToGenerateKeyFromKMS
		}

		keyHex, err := c.kmsService.ExportKey(ctx, keyUID)
		if err != nil {
			return "", "", model.ErrFailedToImportKeyFromKMS
		}

		// Convert hex string to bytes
		keyBytes, err := helper.HexToBytes(keyHex)
		if err != nil {
			return "", "", fmt.Errorf("failed to decode hex key: %w", err)
		}

		// Convert raw key bytes to Tink keyset format
		key, err = c.cryptoService.ImportRawKeyAsBase64(keyBytes)
		if err != nil {
			return "", "", fmt.Errorf("failed to convert raw key to Tink keyset: %w", err)
		}
	} else {
		slog.Info("KMS is not enabled, generating local key")
		key, err = c.cryptoService.GenerateKey()
		if err != nil {
			return "", "", model.ErrKeyGenerationFailed
		}
	}
	return key, keyUID, nil
}

// createMetadataDTO constructs metadata DTO
func (c *FileService) createMetadataDTO(keyUID, key, mimeType string, size int64, hash string, encryptedFile []byte) *model.MetaDataDTO {
	metadata := &model.MetaDataDTO{
		KeyUID:   keyUID,
		Key:      key,
		MimeType: mimeType,
		Size:     size,
		Hash:     hash,
	}

	if c.hashEncryptedFile {
		if encryptedFileHash, err := c.cryptoService.HashFile(c.hashMethod, encryptedFile); err == nil {
			metadata.EncryptedFileHash = encryptedFileHash
		} else {
			slog.Warn("Failed to calculate encrypted file hash", slog.Any("error", err))
		}
	}

	return metadata
}

func (c *FileService) decryptFile(key, hashValue string, encryptedFile []byte) ([]byte, error) {
	if encryptedFile == nil && len(encryptedFile) == 0 && len(hashValue) == 0 && len(key) == 0 {
		return nil, model.ErrFileIsEmpty
	}

	//decrypt file
	decryptedFile, err := c.cryptoService.DecryptFile(key, encryptedFile)
	if err != nil {
		return nil, err
	}
	//compare hash
	if !c.cryptoService.CompareHashFile(c.hashMethod, decryptedFile, hashValue) {
		return nil, model.ErrHashNotMatch
	}
	return decryptedFile, nil
}

func createFileName(fileName string) string {
	return fileName + ".enc"
}

func (c *FileService) saveFileLog(ctx context.Context, appID, fileID, actorType, action string, fileName string) error {
	log := &entity.FileLogs{
		FileID:    fileID,
		ActorID:   appID,
		ActorType: actorType,
		Action:    action,
		IP:        helper.GetClientIP(ctx),
		UserAgent: helper.GetUserAgent(ctx),
		Metadata: map[string]interface{}{
			"file_name": fileName,
		},
	}
	return c.fileLogsRepository.Create(ctx, log)
}

func (c *FileService) checkClientID(ctx context.Context, clientID string) (string, error) {
	if clientID == "" {
		return "", model.ErrInvalidInput
	}
	app, err := c.applicationRepository.GetByClientIDLimited(ctx, clientID)
	if err != nil {
		return "", err
	}
	if app == nil {
		return "", model.ErrAppNotFound
	} else if app.IsActive == false {
		return "", model.ErrAppNotActive
	}
	return app.ID, nil
}

type FileServiceParams struct {
	CryptoService         CryptographicInterface
	StorageService        StorageInterface
	KMSService            KMSInterface
	FileRepository        repository.FileRepository
	FileLogsRepository    repository.FileLogsRepository
	ApplicationRepository repository.ApplicationRepository
	AdminRepository       repository.AdminRepository
	DB                    *gorm.DB
	KeyConfig             *model.KeyConfig
	BucketName            string
	HashMethod            string
	HashEncryptedFile     bool
	EncryptionMethod      string
}
