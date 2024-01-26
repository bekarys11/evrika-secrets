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
	repo.Logger.Debug("search request", searchReq)

	result, err := repo.LDAP.SearchWithPaging(searchReq, 1)
	if err != nil {
		repo.Logger.Error("Ошибка поиска active directory", err)
		return nil, fmt.Errorf("Ошибка поиска active directory: %s", err)
	}
	repo.Logger.Debug("active directory search result", result)

	if len(result.Entries) > 0 {
		repo.Logger.Debug("result entries", result.Entries)
		return result, nil
	} else {
		repo.Logger.Error("no result")
		return nil, fmt.Errorf("Нет результатов")
	}
}
