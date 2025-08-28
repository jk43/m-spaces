package main

import (
	"context"
	"crypto/md5"
	"file/models"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (serv *grpcServer) DeleteFile(ctx context.Context, in *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := serv.app.DeleteFile(in.OrgID, in.UserID, in.UuID)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteFileResponse{}, nil
}

func (app *application) DeleteFile(orgID, userID, uuID string) error {
	file := &models.FileInfo{
		UUID:           uuID,
		UserID:         userID,
		OrganizationID: orgID,
	}
	err := app.DB.GetFile(file)
	if err != nil {
		return err
	}
	s3Path := file.S3Key
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, awsToken,
		)),
	)
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Path),
	})
	if err != nil {
		return err
	}
	err = app.DB.DeleteFile(file)
	if err != nil {
		return err
	}
	utils.TermDebugging(`s3Path`, s3Path)
	return nil
}

// pull rules from service
func getRules(userID, orgID, service, serviceCtx string) (*pb.FileRulesResponse, error) {
	serviceAddr := os.Getenv(strings.ToUpper(service)+"_ADDR") + ":5001"
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		utils.TermDebugging(`err`, err)
		return nil, err
	}
	c := pb.NewFileClientClient(conn)
	req := &pb.FileRulesRequest{
		UserID:     userID,
		OrgID:      orgID,
		ServiceCtx: serviceCtx,
	}
	return c.GetFileRules(context.Background(), req)
}

// send file data to service
func saveFileData(service string, fileInfos []*models.FileInfo) (*pb.FileSaveDataResponse, error) {
	conn, err := grpc.Dial(os.Getenv(strings.ToUpper(service)+"_ADDR")+":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		utils.TermDebugging(`err`, err)
		return nil, err
	}
	c := pb.NewFileClientClient(conn)
	req := &pb.FileSaveDataRequest{}
	req.BatchID = fileInfos[0].BatchID
	req.UserID = fileInfos[0].UserID
	req.OrgID = fileInfos[0].OrganizationID
	req.FormID = fileInfos[0].FormID
	req.FileInfos = []*pb.FileInfo{}
	for _, fileInfo := range fileInfos {
		req.FileInfos = append(req.FileInfos, &pb.FileInfo{
			OriginalFileName: fileInfo.OriginalFileName,
			FileName:         fileInfo.FileName,
			S3Key:            fileInfo.S3Key,
			ContentType:      fileInfo.ContentType,
			Size:             int64(fileInfo.Size),
			FormKey:          fileInfo.FormKey,
			UuID:             fileInfo.UUID,
		})
	}
	return c.SaveFileData(context.Background(), req)
}

func validateFile(rules *pb.FileRulesResponse, fileInfo *models.FileInfo) error {
	if rules.MaxFileSize < fileInfo.Size {
		return fmt.Errorf(fmt.Sprintf("The maximum allowed file size is %dMB. However, the size of the %v is %dMB.", rules.MaxFileSize/(1024*1024), fileInfo.OriginalFileName, fileInfo.Size/(1024*1024)))
	}
	inputFileType := strings.ToLower(strings.Split(fileInfo.ContentType, "/")[0])
	var err error
	err = nil
	for _, fileType := range rules.AllowedFileTypes {
		if fileType == "*" {
			return nil
		}
		if fileType == inputFileType {
			return nil
		}
		err = fmt.Errorf("The file type %v is not allowed", inputFileType)
	}
	if err != nil {
		return err
	}
	for _, contentType := range rules.AllowedContentTypes {
		if contentType == "*" {
			return nil
		}
		if contentType == fileInfo.ContentType {
			return nil
		}
		err = fmt.Errorf("The file type %v is not allowed", fileInfo.ContentType)
	}
	return err
}

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ping")
}

func (app *application) SaveFiles(w http.ResponseWriter, r *http.Request) {
	var maxFileSize int64 = 10 << 20
	fileDir := "/tmp"
	batchID := getBatchID()
	claims := utils.GetClaims(r)
	userID := claims.Subject
	orgIDObjectID, _ := service.GetOrgID(r)
	orgID := orgIDObjectID.Hex()
	// If a call occurs between servers, the orgID cannot be found.
	if orgID == "000000000000000000000000" {
		orgID = claims.OrganizationID
	}
	err := r.ParseMultipartForm(1 << 20) // Maximum memory 10MB
	//delete temp files
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFile, w, err)
		return
	}
	svc := r.FormValue("service")
	svcCtx := r.FormValue("serviceCtx")
	formID := r.FormValue("formId")
	rules, err := getRules(userID, orgID, svc, svcCtx)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err)
		return
	}
	fileInfos := []*models.FileInfo{}
	for k, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			fileInfo := &models.FileInfo{}
			fileInfos = append(fileInfos, fileInfo)
			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFile, w, err)
				return
			}
			defer file.Close()
			fileInfo.BatchID = batchID
			fileInfo.UUID = uuid.New().String()
			fileInfo.OriginalFileName = fileHeader.Filename
			fileInfo.FileName = hashFileName(fileInfo.OriginalFileName)
			fileInfo.ContentType = fileHeader.Header.Get("Content-Type")
			fileInfo.Size = fileHeader.Size
			fileInfo.Service = svc
			fileInfo.ServiceCtx = svcCtx
			fileInfo.FormID = formID
			fileInfo.FormKey = k
			fileInfo.UserID = userID
			fileInfo.OrganizationID = orgID
			fileInfo.IP = r.RemoteAddr
			fileInfo.S3Key = rules.S3Dir + fileInfo.FileName
			// Create a new file in the uploads directory
			dst, err := os.Create(fileDir + "/" + fileInfo.FileName)
			fileInfo.OSFile = dst
			if err != nil {
				cleanUpFiles(fileInfos)
				utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFile, w, err)
				return
			}
			err = validateFile(rules, fileInfo)
			if err != nil {
				utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFile, w, err)
				return
			}
			defer dst.Close()
			// Copy the uploaded file to the filesystem at the specified destination
			byteWritten, err := io.Copy(dst, file)
			if err != nil {
				utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFile, w, err)
				return
			}
			// Validate file size
			if byteWritten > maxFileSize {
				cleanUpFiles(fileInfos)
				utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFile, w, err)
				return
			}
		}
	}
	err = app.uploadToS3(fileInfos)
	if err != nil {
		cleanUpFiles(fileInfos)
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeS3, w, err)
		return
	}
	// Save file info to database
	for _, fileInfo := range fileInfos {
		cleanUpFiles(fileInfos)
		_, err = app.DB.InsertFile(fileInfo)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
			return
		}
	}
	cleanUpFiles(fileInfos)
	res, err := saveFileData(svc, fileInfos)
	utils.TermDebugging(`Res from grpc `, res.Output)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: fileInfos})
}

func hashFileName(name string) string {
	timestamp := time.Now().UnixNano() / int64(time.Microsecond)
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(timestamp, 10)+name)
	nameParts := strings.Split(name, ".")
	ext := nameParts[len(nameParts)-1]
	return fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
}

func getBatchID() string {
	timestamp := time.Now().UnixNano() / int64(time.Microsecond)
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(timestamp, 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func cleanUpFiles(filesInfos []*models.FileInfo) {
	for _, file := range filesInfos {
		file.OSFile.Close()
		os.Remove(file.OSFile.Name())
	}
}

func (app *application) uploadToS3(fileInfos []*models.FileInfo) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, awsToken,
		)),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	var wg sync.WaitGroup
	s3Error := false

	for _, fileInfo := range fileInfos {
		wg.Add(1)
		go func(fileInfo *models.FileInfo) {
			defer wg.Done()
			_, err = fileInfo.OSFile.Seek(0, 0)
			if err != nil {
				fileInfo.Error = true
				s3Error = true
				return
			}

			_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
				Bucket:      aws.String(s3Bucket),
				Key:         aws.String(fileInfo.S3Key),
				ContentType: aws.String(fileInfo.ContentType),
				Body:        fileInfo.OSFile,
			})
			if err != nil {
				fileInfo.Error = true
				fileInfo.ErrorMessage = err.Error()
				s3Error = true
			}
		}(fileInfo)
	}

	wg.Wait()

	if s3Error {
		for _, fileInfo := range fileInfos {
			_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: aws.String(s3Bucket),
				Key:    aws.String(fileInfo.S3Key),
			})
			if err != nil {
				app.Logger.Error().Err(err).Msg("Failed to delete object from S3")
			}
		}
		return err
	}

	return nil
}
