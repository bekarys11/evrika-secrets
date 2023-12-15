package users

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"os"
)

func (repo *Repository) activeDirSearch(email string) (*ldap.SearchResult, error) {
	filter := fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(email))

	searchReq := ldap.NewSearchRequest(
		os.Getenv("LDAP_BASE_DN"),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{},
		nil,
	)

	result, err := repo.LDAP.SearchWithPaging(searchReq, 1)
	if err != nil {
		return nil, fmt.Errorf("Ошибка поиска active directory: %s", err)
	}

	if len(result.Entries) > 0 {
		return result, nil
	} else {
		return nil, fmt.Errorf("Нет результатов")
	}
}
