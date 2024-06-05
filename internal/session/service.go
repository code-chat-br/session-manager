package session

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

const INSTANCE_PATH string = "instances"

var mutex sync.Locker

type Service struct {
	logger *logrus.Entry
}

func New() *Service {
	logger := logrus.New()
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

func (s *Service) CreateFolder(group, instance_name string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group, instance_name)
	if err := s.Create(path_name); err != nil {
		s.logger.WithField("func", "CreateInstanceFolder").Log(logrus.ErrorLevel, err)
		return http.StatusInternalServerError, err
	}

	s.logger.WithField("status", "criado").Log(logrus.WarnLevel, fmt.Sprintf("Folder %s", path_name))

	return http.StatusNoContent, nil
}

func (s *Service) RemoveFolder(group, instance_name string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group, instance_name)
	if err := os.RemoveAll(path_name); err != nil {
		s.logger.WithField("func", "RemoveFolder").Log(logrus.ErrorLevel, err)
		return http.StatusInternalServerError, err
	}

	s.logger.Log(logrus.WarnLevel, fmt.Sprintf("Folder %s -> removed", path_name))

	return http.StatusOK, nil
}

func (s *Service) WriterCredentials(group, instance_name, key string, json map[string]string) (int, error) {
	mutex.Lock()
	defer mutex.Unlock()

	path_name := filepath.Join(INSTANCE_PATH, group, instance_name, key)

	file, err := os.Create(path_name)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "WriterCredentials",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)
		return http.StatusInternalServerError, err
	}
	defer file.Close()

	binary := []byte(json["data"])

	if _, err := file.Write(binary); err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "WriterCredentials",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

func (s *Service) ReadCredentials(group, instance_name, key string) (int, []byte, error) {
	path_name := filepath.Join(INSTANCE_PATH, group, instance_name, key)

	file, err := os.Open(path_name)
	if err != nil {
		// Este erro será o mais exibido no log.
		//
		// Sempre que a API buscar alguma chave, e ela não existir,
		// haverá uma falha ao abrir o arquivo.
		//
		// Um erro será retornado à API e uma nova chave como o mesmo
		// nome e valor será criada pela própria API.
		// s.logger.
		// 	WithFields(logrus.Fields{
		// 		"func": "ReadCredentials",
		// 		"status": "controlado",
		// 	}).
		// 	Log(logrus.ErrorLevel, err)
		return http.StatusBadRequest, nil, err
	}
	defer file.Close()

	binary, err := io.ReadAll(file)

	return http.StatusOK, binary, err
}

func (s *Service) RemoveCredential(group, instance_name, key string) (int, error) {
	path_name := filepath.Join(INSTANCE_PATH, group, instance_name, key)

	err := os.RemoveAll(path_name)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "RemoveCredential",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)
		return http.StatusBadRequest, err
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
		return http.StatusBadRequest, nil, err
	}

	files, err := os.ReadDir(source)
	if err != nil {
		s.logger.
			WithFields(logrus.Fields{
				"func":   "ListInstances",
				"status": "controlado",
			}).
			Log(logrus.ErrorLevel, err)
		return http.StatusBadRequest, nil, err
	}

	list := make([]string, len(files))

	for i, file := range files {
		if file.IsDir() {
			list[i] = file.Name()
		}
	}

	return http.StatusOK, &list, nil
}
