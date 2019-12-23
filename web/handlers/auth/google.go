package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/internal/model"
	"github.com/weiqiang333/devops/internal/authentication"
	"github.com/weiqiang333/devops/internal/database"
)


func createSecret(name string) (string, error) {
	secret := authentication.NewGoogleAuth().GetSecret()
	sql := fmt.Sprintf(`
		INSERT INTO google_auth (name, secret, updated_at)
		VALUES ('%s', '%s', now() at time zone 'utc')
		ON CONFLICT (name)
		DO UPDATE SET
		name = EXCLUDED.name,
		secret = EXCLUDED.secret,
		updated_at = EXCLUDED.updated_at;`, name, secret)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		return "", fmt.Errorf("insert google_auth error for %s: %v", name, err)
	}
	return secret, nil
}


//CreateQRcode
func CreateQRcode(c *gin.Context)  {
	name := Me(c)
	secret, err := createSecret(fmt.Sprint(name))
	if err != nil {
		log.Printf("CreateQRcode error for %s", name)
		c.JSON(http.StatusNotImplemented, gin.H{"QR code URL:": err.Error()})
		return
	}
	qrCodeUrl := authentication.NewGoogleAuth().GetQrcodeUrl(fmt.Sprint(name), secret)
	log.Printf("CreateQRcode Success for %s", name)
	c.JSON(http.StatusOK, gin.H{"QRcodeURL": qrCodeUrl})
	return
}

func SearchQRcodeSecret(name string) (string, error) {
	sql := fmt.Sprintf("SELECT secret FROM google_auth WHERE name = '%s';", name)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		return "", fmt.Errorf("SearchQRcodeSecret db.Query error for %s: %v", name, err)
	}

	var g model.TableGoogleAuth
	row.Next()
	if err := row.Scan(&g.Secret); err != nil {
		return "", fmt.Errorf("SearchQRcodeSecret db rows scan error for %s: %v", name, err)
	}
	return g.Secret, nil
}

//SearchQRcodeUrl handler
func SearchQRcodeUrl(name string) (string, error) {
	secret, err := SearchQRcodeSecret(name)
	if err != nil {
		return "", err
	}
	qrCodeUrl := authentication.NewGoogleAuth().GetQrcodeUrl(name, secret)
	return qrCodeUrl, nil
}
