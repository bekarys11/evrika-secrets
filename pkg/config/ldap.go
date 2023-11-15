package config

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"os"
)

func (app *Config) ConnectToLDAP() (err error) {
	app.LDAP, err = ldap.DialURL(fmt.Sprintf("ldap://%s:389", os.Getenv("LDAP_HOST")))
	if err != nil {
		return fmt.Errorf("dial error to LDAP: %v", err)
	}

	if err = app.bind(); err != nil {
		return err
	}
	return nil
}

func (app *Config) bind() error {
	if err := app.LDAP.Bind(os.Getenv("LDAP_USERNAME"), os.Getenv("LDAP_PASSWORD")); err != nil {
		return fmt.Errorf("ldap bind error: %s", err)
	}
	return nil
}
