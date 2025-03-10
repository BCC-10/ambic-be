package supabase

import (
	"ambic/internal/domain/env"
	"ambic/internal/infra/response"
	"fmt"
	storageGo "github.com/nedpals/supabase-go"
	"io"
	"mime/multipart"
)

type SupabaseIf interface {
	UploadFile(bucket string, filePath string, contentType string, file multipart.File) (string, error)
	UploadFileFromIOReader(bucket string, filePath string, contentType string, file io.Reader) (string, error)
	DeleteFile(bucket string, filePath string) error
}

type Supabase struct {
	Client storageGo.Client
}

func New(env *env.Env) SupabaseIf {
	storageClient := storageGo.CreateClient(fmt.Sprintf("%s", env.SupabaseURL), env.SupabaseSecret)

	return &Supabase{
		Client: *storageClient,
	}
}

func (s *Supabase) UploadFile(bucket string, filePath string, contentType string, file multipart.File) (string, error) {
	wrapper := safeWrapper(func() error {
		s.Client.Storage.From(bucket).Upload(filePath, file, &storageGo.FileUploadOptions{ContentType: contentType})
		return nil
	})

	if wrapper != nil {
		return "", wrapper
	}

	publicURL := s.Client.Storage.From(bucket).GetPublicUrl(filePath).SignedUrl
	return publicURL, nil
}

func (s *Supabase) UploadFileFromIOReader(bucket string, filePath string, contentType string, file io.Reader) (string, error) {
	wrapper := safeWrapper(func() error {
		s.Client.Storage.From(bucket).Upload(filePath, file, &storageGo.FileUploadOptions{ContentType: contentType})
		return nil
	})

	if wrapper != nil {
		return "", wrapper
	}

	publicURL := s.Client.Storage.From(bucket).GetPublicUrl(filePath).SignedUrl
	return publicURL, nil
}

func (s *Supabase) DeleteFile(bucket string, filePath string) error {
	wrapper := safeWrapper(func() error {
		s.Client.Storage.From(bucket).Remove([]string{filePath})
		return nil
	})

	return wrapper
}

func safeWrapper(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = response.ErrInternalServer(fmt.Sprintf("Unexpected panic: %v", r))
		}
	}()
	return f()
}
