package utils

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
)

var port string = ":80"
var fileDir string = os.Getenv("FILE_DIR")
var s3Bucket string = os.Getenv("AWS_S3_BUCKET")
var s3Region string = os.Getenv("AWS_S3_REGION")

var awsAccessKey string = os.Getenv("AWS_ACCESS_KEY_ID")
var awsSecretKey string = os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsToken string = os.Getenv("AWS_TOKEN")

var awsCloudFrontKeyPairID string = os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")
var awsCloudFrontPrivateKey []byte = []byte(os.Getenv("AWS_CLOUDFRONT_PRIVATE_KEY"))
var awsCloudFrontDomain string = os.Getenv("AWS_CLOUDFRONT_DOMAIN")

var s3Client *s3.Client

func GetS3Client() (*s3.Client, error) {
	if s3Client != nil {
		return s3Client, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, awsToken,
		)),
	)
	if err != nil {
		return nil, err
	}
	s3Client = s3.NewFromConfig(cfg)
	return s3Client, nil
}

func GetPresignedUrl(key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, awsToken,
		)),
	)
	if err != nil {
		return "", err
	}

	client := s3.NewPresignClient(s3.NewFromConfig(cfg))

	presignedReq, err := client.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: &s3Bucket,
			Key:    &key,
		},
		s3.WithPresignExpires(15*time.Minute),
	)
	if err != nil {
		return "", err
	}

	return presignedReq.URL, nil
}

func GetCloudFrontPresignedURL(objectPath string) (string, error) {
	block, _ := pem.Decode(awsCloudFrontPrivateKey)
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey
	var err error

	if parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		switch key := parsedKey.(type) {
		case *rsa.PrivateKey:
			privateKey = key
		default:
			return "", fmt.Errorf("not an RSA private key")
		}
	} else {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key: %v", err)
		}
	}

	signer := sign.NewURLSigner(awsCloudFrontKeyPairID, privateKey)

	resourceURL := awsCloudFrontDomain + "/" + objectPath
	expires := time.Now().Add(1 * time.Minute)

	fmt.Println("resourceURL", resourceURL)
	signedURL, err := signer.Sign(resourceURL, expires)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func GetCloudFrontSignedCookie() ([]*http.Cookie, error) {
	block, _ := pem.Decode(awsCloudFrontPrivateKey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey
	var err error

	if parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		switch key := parsedKey.(type) {
		case *rsa.PrivateKey:
			privateKey = key
		default:
			return nil, fmt.Errorf("not an RSA private key")
		}
	} else {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
	}

	signer := sign.NewCookieSigner(awsCloudFrontKeyPairID, privateKey)
	expires := time.Now().Add(60 * time.Minute)
	policy := &sign.Policy{
		Statements: []sign.Statement{
			{

				Resource: "*",
				Condition: sign.Condition{
					DateLessThan: &sign.AWSEpochTime{Time: expires},
				},
			},
		},
	}
	cookies, err := signer.SignWithPolicy(policy)
	if err != nil {
		return nil, err
	}
	return cookies, nil
}
