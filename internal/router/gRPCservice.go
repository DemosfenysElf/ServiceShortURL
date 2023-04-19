package router

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"ServiceShortURL/internal/router/gRPC"
)

//GetShortURL(context.Context, *Long) (*Short, error)
//GetLongURL(context.Context, *Short) (*Long, error)
//GetBatchShort(context.Context, *Batch) (*Batch, error)

type ServerGRPC struct {
	*ServerShortener
}

func startServerGRPC() error {
	sGRPC := grpc.NewServer()
	srv := NewServerGRPC()
	gRPC.RegisterServiceShortUrlServer(sGRPC, srv)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	if err = sGRPC.Serve(l); err != nil {
		return err
	}
	return nil
}

func NewServerGRPC() (s *ServerGRPC) {
	return
}

func (s *ServerGRPC) GetShortURL(ctx context.Context, long *gRPC.Long) (*gRPC.Short, error) {
	short, err := s.StorageInterface.SetShortURL(ctx, long.Url)
	if err != nil {
		return nil, err
	}
	return &gRPC.Short{Url: short}, nil
}
func (s *ServerGRPC) GetLongURL(ctx context.Context, short *gRPC.Short) (*gRPC.Long, error) {
	long, err := s.StorageInterface.GetLongURL(ctx, short.Url)
	if err != nil {
		return nil, err
	}
	return &gRPC.Long{Url: long}, nil
}
func (s *ServerGRPC) GetBatchShort(ctx context.Context, b *gRPC.Batch) (*gRPC.Batch, error) {
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

	result := &gRPC.Batch{}
	result.Result = make([]*gRPC.Pack, 0, len(b.Result))

	for _, oneBatch := range shortURLBatch {
		result.Result = append(result.Result, &gRPC.Pack{
			CorrelationId: oneBatch.ID,
			Url:           oneBatch.ShortURL,
		})
	}
	return result, nil
}
