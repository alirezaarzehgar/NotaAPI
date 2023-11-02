package validations

import (
	"os"

	"github.com/Asrez/NotaAPI/config"
)

func IsValidAsset(pathname string) bool {
	_, err := os.Stat(config.Assets() + "/" + pathname)
	return err == nil
}
