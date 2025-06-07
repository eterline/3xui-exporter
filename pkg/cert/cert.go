package cert

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CertErr string

func (t CertErr) Error() string {
	return "cert provider error: " + string(t)
}

const (
	ErrCertStoreDirNotSet CertErr = "certificate store directory is not set"
	ErrCertNotFound       CertErr = "certificate not found"
	ErrFileNotFound       CertErr = "file not found in certificate directory"
	ErrInvalidCertPair    CertErr = "invalid certificate/key pair"
)

type X509Type string

func (t X509Type) String() string {
	return string(t)
}

func (t X509Type) JoinName(name string) string {
	return name + t.String()
}

const (
	Pem  X509Type = ".pem"
	Cert X509Type = ".crt"
	Key  X509Type = ".key"
)

type CertData struct {
	name     string
	data     []byte
	fileType X509Type
}

func (t CertData) FullName() string {
	return t.name + t.fileType.String()
}

type CertStore struct {
	certMap []CertData
}

var (
	certStoreDir = ""
	store        = CertStore{
		certMap: []CertData{},
	}
)

func init() {
	certStoreDir = os.Getenv("ETERLINE_CERT_DIR")
	if certStoreDir != "" {
		loadCertificates()
	}
}

// loadCertificates - loads all certs from 'certStoreDir'
func loadCertificates() error {
	if certStoreDir == "" {
		return ErrCertStoreDirNotSet
	}

	files, err := ioutil.ReadDir(certStoreDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		filePath := filepath.Join(certStoreDir, fileName)

		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		// define file type by extension
		var fileType X509Type

		switch strings.ToLower(filepath.Ext(fileName)) {
		case Pem.String():
			fileType = Pem
		case Cert.String():
			fileType = Cert
		case Key.String():
			fileType = Key
		default:
			continue // skip unsupported
		}

		name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		store.certMap = append(store.certMap, CertData{
			name:     name,
			data:     data,
			fileType: fileType,
		})
	}

	return nil
}

// GetCertificate - return cert bytes by name and type
func GetCertificate(name string, fileType X509Type) ([]byte, error) {
	if certStoreDir == "" {
		return nil, ErrCertStoreDirNotSet
	}

	for _, cert := range store.certMap {
		if cert.FullName() == fileType.JoinName(name) {
			return cert.data, nil
		}
	}

	return nil, ErrCertNotFound
}

// LoadCertificate - loads new cert to storage
func LoadCertificate(name string, fileType X509Type) ([]byte, error) {
	if certStoreDir == "" {
		return nil, ErrCertStoreDirNotSet
	}

	if data, err := GetCertificate(name, fileType); err == nil {
		return data, nil
	}

	fileName := name + string(fileType)
	filePath := filepath.Join(certStoreDir, fileName)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	store.certMap = append(store.certMap, CertData{
		name:     name,
		data:     data,
		fileType: fileType,
	})

	return data, nil
}

func ForceDir(dir string) error {
	if dir == "" {
		return errors.New("directory path cannot be empty")
	}

	fileInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist", dir)
		}
		return fmt.Errorf("failed to access directory: %v", err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	file, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("cannot open directory: %v", err)
	}
	file.Close()

	certStoreDir = dir

	if err := loadCertificates(); err != nil {
		return fmt.Errorf("failed to load certificates from new directory: %v", err)
	}

	return nil
}

func Count() int {
	return len(store.certMap)
}

func X509Pair(name string) (*tls.Certificate, error) {
	certData, err := GetCertificate(name, Cert)
	if err != nil {
		return nil, err
	}

	keyData, err := GetCertificate(name, Key)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCertPair, err)
	}

	return &cert, nil
}
