package router

import "testing"

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
			rout.FlagParse()

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
