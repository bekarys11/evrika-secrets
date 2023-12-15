package config

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"os"
)

func connectToLDAP() (ldapConn *ldap.Conn, err error) {
	ldapConn, err = ldap.DialURL(fmt.Sprintf("ldap://%s:389", os.Getenv("LDAP_HOST")))
	if err != nil {
		return nil, fmt.Errorf("dial error to LDAP: %v", err)
	}

	if err = bind(ldapConn); err != nil {
		return nil, err
	}

	return ldapConn, nil
}

func bind(ldapConn *ldap.Conn) error {
	if err := ldapConn.Bind(os.Getenv("LDAP_USERNAME"), os.Getenv("LDAP_PASSWORD")); err != nil {
		return fmt.Errorf("ldap bind error: %s", err)
	}
	return nil
}
