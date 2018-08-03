package main

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/minio/minio-go"
)

type DonationSummary struct {
	UpdatedAt      time.Time `json:"updated_at"`
	UsdDonations   float32   `json:"usd_donations"`
	EurDonations   float32   `json:"eur_donations"`
	EurToUsdRate   float32   `json:"eur_to_usd_rate"`
	TotalDonations float32   `json:"total_donations"`
}

const minioHost = "cdn.haiku-os.org"

func UploadJson(summary *DonationSummary) error {
	accessKeyID := config.Minio.AccessKeyId
	secretAccessKey := config.Minio.SecretAccessKey
	minioClient, err := minio.New(minioHost, accessKeyID, secretAccessKey, true)
	if err != nil {
		return err
	}

	b, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	_, err = minioClient.PutObject(
		"haiku-inc",
		"donations.json",
		bytes.NewBuffer(b),
		int64(len(b)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	if err != nil {
		return err
	}

	return nil
}
