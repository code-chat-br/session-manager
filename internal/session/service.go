package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"worker-session/pkg/db"
	"worker-session/pkg/db/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const INSTANCE_PATH string = "instances"

type Service struct {
	logger *logrus.Entry
}

var db_connections = make(db.ManagerDB)
var mutex sync.Locker

func connect_db(path, db_name string) (*gorm.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if db, exists := db_connections[db_name]; exists {
		return db, nil
	}

	conn, err := db.InitializerDB(filepath.Join(path, db_name))
	if err != nil {
		return nil, err
	}

	conn.AutoMigrate(&models.Creds{})

	db_connections[db_name] = conn

	return conn, nil
}

func New() *Service {
	logger := logrus.New()
	mutex = &sync.Mutex{}
	return &Service{
		logger: logger.WithField("desc", "session.Service"),
	}
}

func (s *Service) Create(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		// s.logger.WithField("func", "CreateFolder").Log(logrus.ErrorLevel, err)
		return os.MkdirAll(folder, os.ModePerm)
	}

	return nil
}

func (s *Service) CreateGroup(group string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group)
	if err := s.Create(path_name); err != nil {
		s.logger.WithField("func", "CreateInstanceFolder").Log(logrus.ErrorLevel, err)
		return http.StatusInternalServerError, err
	}

	s.logger.WithField("status", "criado").Log(logrus.WarnLevel, fmt.Sprintf("Folder %s", path_name))

	return http.StatusNoContent, nil
}

func (s *Service) RemoveGroup(group string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group)
	if err := os.RemoveAll(path_name); err != nil {
		s.logger.WithField("func", "RemoveFolder").Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, err
	}

	s.logger.Log(logrus.WarnLevel, fmt.Sprintf("Folder %s -> removed", path_name))

	return http.StatusOK, nil
}

func (s *Service) CreateInstanceDb(group, instance_name string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group)
	db, err := connect_db(path_name, instance_name)
	if err != nil {
		s.logger.WithField("func", "CreateInstanceDb").Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, err
	}

	db.AutoMigrate(&models.Creds{})
	db_connections[instance_name] = db

	return http.StatusOK, nil
}

func (s *Service) RemoveInstanceDb(group, instance_name string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group, instance_name+".db")
	if err := os.RemoveAll(path_name); err != nil {
		s.logger.WithField("func", "RemoveInstanceDb").Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *Service) WriterCredentials(group, instance_name, key string, value map[string]any) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group)
	db, err := connect_db(path_name, instance_name)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "WriterCredentials",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, err
	}

	binary, err := json.Marshal(value["data"])
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "WriterCredentials",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)

		return http.StatusBadRequest, err
	}

	creds := models.NewCreds(key, binary)

	tx := db.Create(&creds)
	if tx.Error != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "WriterCredentials",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, tx.Error)

		return http.StatusInternalServerError, tx.Error
	}

	return http.StatusNoContent, nil
}

func (s *Service) ReadCredentials(group, instance_name, key string) (int, []byte, error) {
	path_name := filepath.Join(INSTANCE_PATH, group)
	db, err := connect_db(path_name, instance_name)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "ReadCredentials",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, nil, err
	}

	var creds models.Creds
	tx := db.Where("name = ?", key).First(&creds)
	if tx.Error != nil {
		return http.StatusBadRequest, nil, tx.Error
	}

	return http.StatusOK, creds.Content, nil
}

func (s *Service) RemoveCredential(group, instance_name, key string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group)
	db, err := connect_db(path_name, instance_name)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "RemoveCredential.connect_db",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, err
	}

	tx := db.Where("name = ?", key).Delete(&models.Creds{})
	if tx.Error != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "RemoveCredential.Delete",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, tx.Error)

		return http.StatusInternalServerError, tx.Error
	}

	return http.StatusOK, nil
}

func (s *Service) ListInstances(group string) (int, *[]string, error) {
	source := filepath.Join(INSTANCE_PATH, group)

	if _, err := os.Stat(source); os.IsNotExist(err) {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "ListInstances",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, nil, err
	}

	files, err := os.ReadDir(source)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "ListInstances",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)

		return http.StatusInternalServerError, nil, err
	}

	list := make([]string, len(files))

	for i, file := range files {
		if !file.IsDir() {
			list[i] = strings.Replace(file.Name(), ".db", "", 1)
		}
	}

	return http.StatusOK, &list, nil
}
