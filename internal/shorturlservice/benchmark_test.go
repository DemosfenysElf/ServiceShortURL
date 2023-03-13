package shorturlservice

import (
	"context"
	"math/rand"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func BenchmarkXxx(b *testing.B) {
	cfgDB := "postgres://postgres:0000@localhost:5432/postgres"
	var userList []string
	DB := &Database{}
	DB.Connect(cfgDB)

	b.ResetTimer()

	b.Run("URl to Short to URL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			url := "https://" + randomString(5) + ".ru/" + randomString(10)
			user := randomString(7)
			SetStructCookies("Authentication", user)
			b.StartTimer()
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
			defer cancel()
			short, _ := DB.SetURL(ctx, url)
			if short != "" {
				DB.GetURL(ctx, short)
			}
			b.StopTimer()
			if i%2 == 0 {
				userList = append(userList, short)

			}
			b.StartTimer()
			if i%5 == 0 {
				DB.Delete(ctx, user, userList)
				b.StopTimer()
				userList = nil
				b.StartTimer()
			}
		}
	})

}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)

}
