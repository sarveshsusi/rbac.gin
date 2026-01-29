package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"rbac/config"
)

type ImageKitUploader struct {
	publicKey  string
	privateKey string
	endpoint   string
}

func NewImageKitUploader(cfg *config.Config) ImageUploader {
	return &ImageKitUploader{
		publicKey:  cfg.ImageKit.PublicKey,
		privateKey: cfg.ImageKit.PrivateKey,
		endpoint:   cfg.ImageKit.Endpoint,
	}
}

func (u *ImageKitUploader) Upload(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, src); err != nil {
		return "", err
	}

	_ = writer.WriteField("fileName",
		fmt.Sprintf("tickets/%d%s",
			time.Now().UnixNano(),
			filepath.Ext(file.Filename),
		),
	)

	_ = writer.WriteField("useUniqueFileName", "true")
	_ = writer.WriteField("folder", "/tickets")

	writer.Close()

	req, err := http.NewRequest(
		"POST",
		"https://upload.imagekit.io/api/v1/files/upload",
		body,
	)
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(u.publicKey + ":" + u.privateKey),
	)

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return "", errors.New(string(raw))
	}

	var res struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.URL, nil
}
