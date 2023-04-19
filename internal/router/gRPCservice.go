package router

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

// ServerGRPC я не мастер писать
type ServerGRPC struct {
	*ServerShortener
}

// startServerGRPC() итак понятно
func startServerGRPC() error {
	sGRPC := grpc.NewServer()
	srv := NewServerGRPC()
	RegisterServiceShortUrlServer(sGRPC, srv)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	go sGRPC.Serve(l)
	return nil
}

// NewServerGRPC создаём экземпляр сервера
func NewServerGRPC() (s *ServerGRPC) {
	return
}

// Не делал больше, так как не вывожу уже,
// и совсем не уверен правильное ли это всё

// как просто проверить работу пока не знаю,
// но зная себя мне потребуется дня 3 на осознание

// GetShortURL даём ссылку, получаем короткую
func (s *ServerGRPC) GetShortURL(ctx context.Context, long *Long) (*Short, error) {
	short, err := s.StorageInterface.SetShortURL(ctx, long.Url)
	if err != nil {
		return nil, err
	}
	return &Short{Url: short}, nil
}

// GetShortURL даём короткую, получаем оригинальную
func (s *ServerGRPC) GetLongURL(ctx context.Context, short *Short) (*Long, error) {
	long, err := s.StorageInterface.GetLongURL(ctx, short.Url)
	if err != nil {
		return nil, err
	}
	return &Long{Url: long}, nil
}

// GetBatchShort даём пачку ссылок, получаем пачку коротких ссылок
func (s *ServerGRPC) GetBatchShort(ctx context.Context, b *Batch) (*Batch, error) {
	shortURLOne := shortURLApiShortenBatch{}
	longURLBatch := make([]urlAPIShortenBatch, 0, len(b.Result))
	shortURLBatch := make([]shortURLApiShortenBatch, 0, len(longURLBatch))

	for _, url := range b.Result {
		longURLBatch = append(longURLBatch, urlAPIShortenBatch{
			ID:          url.CorrelationId,
			OriginalURL: url.Url,
		})
	}

	for i := range longURLBatch {
		short, err := s.SetShortURL(ctx, longURLBatch[i].OriginalURL)
		if err != nil {
			return nil, err
		}
		shortURLOne.ShortURL = s.Cfg.BaseURL + "/" + short
		shortURLOne.ID = longURLBatch[i].ID
		shortURLBatch = append(shortURLBatch, shortURLOne)
	}

	result := &Batch{}
	result.Result = make([]*Pack, 0, len(b.Result))

	for _, oneBatch := range shortURLBatch {
		result.Result = append(result.Result, &Pack{
			CorrelationId: oneBatch.ID,
			Url:           oneBatch.ShortURL,
		})
	}
	return result, nil
}
