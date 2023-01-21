package shorturlservice

import (
	"fmt"
	"io"
)

type StorageInterface interface {
	SetURL(url string) (short string, err error)
	GetURL(short string) (url string, err error)
	Delete(user string, listURL []string)
}

// Хранение в памяти:
type MemoryStorage struct {
	data []URLInfo
}

func InitMem() *MemoryStorage {
	return &MemoryStorage{data: make([]URLInfo, 0)}
}

func (ms *MemoryStorage) GetURL(short string) (url string, err error) {
	for i, data := range ms.data {
		if short == data.ShortURL {
			url = data.URL
			break
		}
		if i >= len(ms.data) {
			return "", fmt.Errorf("no found url")
		}
	}
	return url, err
}

func (ms *MemoryStorage) SetURL(url string) (short string, err error) {
	short = shortURL()
	if len(ms.data) != 0 {
		for i, data := range ms.data {
			if short == data.ShortURL {
				break
			}
			if i+1 >= len(ms.data) {
				break
			}

			short = shortURL()
		}
	}

	SetStructURL(url, short)
	ms.data = append(ms.data, *GetStructURL())
	return short, nil
}

func (ms *MemoryStorage) Delete(user string, listURL []string) {
	fmt.Println(">>>Storage_Delete_list<<<  ", listURL, "User: ", user)
	fmt.Println(">>>>хранение и удаление кук в памяти не реализованно")

}

// Хранение в файле:
type FileStorage struct {
	Writer   io.Writer
	FilePath string
}

func (fs *FileStorage) GetURL(short string) (url string, err error) {

	consumerURL, err := NewConsumer(fs.FilePath)
	if err != nil {
		return "", err
	}
	defer consumerURL.Close()

	for {
		readURL, err := consumerURL.ReadURLInfo()
		if err != nil {
			break
		}
		if readURL.Deleted == "delet" {
			return "", fmt.Errorf("deleted")
		}
		if readURL.ShortURL == short {
			return readURL.URL, nil
		}
	}
	return "", fmt.Errorf("no found url")
}

func (fs *FileStorage) SetURL(url string) (short string, err error) {
	short = shortURL()
	for {
		_, err := fs.GetURL(short)
		if err != nil {
			break
		}
		short = shortURL()
	}

	urli := SetStructURL(url, short)
	producerURL, err := NewProducer(fs.FilePath)
	if err != nil {
		return "", err
	}
	defer producerURL.Close()
	if err := producerURL.WriteURL(urli); err != nil {
		return "", err
	}

	return short, nil
}

//////////////////////////////test

func (fs *FileStorage) Delete(user string, listURL []string) {
	fmt.Println(">>>Storage_Delete_list<<<  ", listURL, "User: ", user)

	fileRW, _ := fs.newRW(fs.FilePath)
	defer fileRW.Close()

	var iSlice = make([]int64, len(listURL), cap(listURL))
	var position int64

	for {
		readURL, len1string, err := fileRW.ReadURLInfo()
		if err != nil {
			break
		}
		fmt.Println("readURL,	", readURL.ShortURL)
		position = position + int64(len1string) + 1

		for i, u := range listURL {
			if (readURL.ShortURL == u) && (readURL.CookiesAuthentication.ValueUser == user) && (readURL.Deleted == "false") {
				//if (readURL.ShortURL == u) && (readURL.Deleted == "false") {
				iSlice[i] = position - 8
			}
		}
	}
	for i, _ := range iSlice {
		fileRW.WriteDelet(iSlice[i])
	}

}

// Хранение в БД>: databaseInterface
