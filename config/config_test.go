package config

import "testing"

func TestConfigIsProduction(t *testing.T) {
	if !(Config{AppEnv: "production"}).IsProduction() {
		t.Fatal("production environment should be recognized")
	}

	if (Config{AppEnv: "development"}).IsProduction() {
		t.Fatal("development environment should not be recognized as production")
	}
}
