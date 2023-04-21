package router

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"ServiceShortURL/internal/shorturlservice"
)

func TestFlagParse(t *testing.T) {
	type wantFlags struct {
		ServerAddress string
		BaseURL       string
		Storage       string
		ConnectDB     string
		EnableHTTPS   bool
		TrustedSubnet string
		Config        string
	}
	tests := []struct {
		name string
		wantFlags
	}{
		{
			name: "TestFlagParse1",
			wantFlags: wantFlags{
				ServerAddress: ":8080",
				BaseURL:       "http://localhost:8080",
				Storage:       "shortsURl.log",
				ConnectDB:     "postgres://postgres:0000@localhost:5432/postgres",
				EnableHTTPS:   false,
				TrustedSubnet: "",
				Config:        "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rout := InitServer()
			rout.ConfigParse()

			if rout.Cfg.ServerAddress != tt.ServerAddress {
				t.Errorf("Expected flag %s, got %s", tt.ServerAddress, rout.Cfg.ServerAddress)
			}

			if rout.Cfg.BaseURL != tt.BaseURL {
				t.Errorf("Expected flag %s, got %s", tt.BaseURL, rout.Cfg.BaseURL)
			}

			if rout.Cfg.Storage != tt.Storage {
				t.Errorf("Expected flag %s, got %s", tt.Storage, rout.Cfg.Storage)
			}

			if rout.Cfg.ConnectDB != tt.ConnectDB {
				t.Errorf("Expected flag %s, got %s", tt.ConnectDB, rout.Cfg.ConnectDB)
			}

			if rout.Cfg.TrustedSubnet != tt.TrustedSubnet {
				t.Errorf("Expected flag %s, got %s", tt.TrustedSubnet, rout.Cfg.TrustedSubnet)
			}

			if rout.Cfg.EnableHTTPS != tt.EnableHTTPS {
				t.Errorf("Expected flag %t, got %t", tt.EnableHTTPS, rout.Cfg.EnableHTTPS)
			}
		})
	}
}

func TestDBOpen(t *testing.T) {
	tests := []struct {
		name   string
		result driver.Result
	}{
		{
			name:   "TestBDApiShortenBatch1",
			result: sqlmock.NewResult(1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))

			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			mockDB := &shorturlservice.Database{}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB

			mock.ExpectExec(regexp.QuoteMeta(`CREATE TABLE ShortenerURL`)).
				WillReturnResult(tt.result).WillReturnError(nil)
			mock.ExpectExec(regexp.QuoteMeta("CREATE UNIQUE INDEX URL_index ON ShortenerURL (url)")).
				WillReturnResult(tt.result).WillReturnError(nil)
			mock.ExpectPing()

			err = rout.DB.CreateTable()
			if err != nil {
				t.Errorf("DB connect error 1")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err = rout.DB.Ping(ctx)

			if err != nil {
				t.Errorf("DB connect error 2")
			}
		})
	}
}
