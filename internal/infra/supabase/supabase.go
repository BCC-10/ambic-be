package supabase

import (
	"ambic/internal/domain/env"
	"fmt"
	storageGo "github.com/supabase-community/storage-go"
	"mime/multipart"
)

type SupabaseIf interface {
	UploadFile(bucket string, filePath string, contentType string, file multipart.File) (string, error)
}

type Supabase struct {
	Client storageGo.Client
}

func New(env *env.Env) SupabaseIf {
	storageClient := storageGo.NewClient(fmt.Sprintf("%s/storage/v1", env.SupabaseURL), env.SupabaseSecret, nil)

	return &Supabase{
		Client: *storageClient,
	}
}

func (s *Supabase) UploadFile(bucket string, filePath string, contentType string, file multipart.File) (string, error) {
	_, err := s.Client.UploadFile(bucket, filePath, file, storageGo.FileOptions{
		ContentType: &contentType})
	if err != nil {
		return "", err
	}

	publicURL := s.Client.GetPublicUrl(bucket, filePath).SignedURL
	return publicURL, nil
}
