package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"sync"
)

//Service .  Это сервис для управления баннерами
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

//NewService . функция для создания нового сервиса
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

//Banner ..Структура нашего баннера
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

//это стартовый ID но для каждого создание баннера его изменяем
var sID int64 = 0

//All ...
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.items, nil
}

//ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.items {

		if v.ID == id {

			return v, nil
		}
	}

	return nil, errors.New("item not found")
}

//Save ...
func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {

		sID++

		item.ID = sID

		if item.Image != "" {

			item.Image = fmt.Sprint(item.ID) + "." + item.Image

			err := uploadFile(file, "./web/banners/"+item.Image)

			if err != nil {
				return nil, err
			}
		}

		//и после этих действий мы добавляем item в слайс
		s.items = append(s.items, item)

		return item, nil
	}

	for k, v := range s.items {

		if v.ID == item.ID {

			if item.Image != "" {

				item.Image = fmt.Sprint(item.ID) + "." + item.Image

				err := uploadFile(file, "./web/banners/"+item.Image)

				if err != nil {
					return nil, err
				}
			} else {

				item.Image = s.items[k].Image
			}

			s.items[k] = item

			return item, nil
		}
	}

	return nil, errors.New("item not found")
}

//RemoveByID ... Метод для удаления
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, v := range s.items {

		if v.ID == id {

			s.items = append(s.items[:k], s.items[k+1:]...)

			return v, nil
		}
	}

	return nil, errors.New("item not found")
}

//это функция сохраняет файл в сервере в заданной папке path и возврашает nil если все успешно
func uploadFile(file multipart.File, path string) error {

	var data, err = ioutil.ReadAll(file)

	if err != nil {
		return errors.New("not readble data")
	}

	err = ioutil.WriteFile(path, data, 0666)

	if err != nil {
		return errors.New("not saved from folder ")
	}

	return nil
}
