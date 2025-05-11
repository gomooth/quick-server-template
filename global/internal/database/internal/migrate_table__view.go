package internal

var views = map[string]string{
	"vw_users": userViewSql,
}

var userViewSql = `
CREATE VIEW vw_users AS
    SELECT u.*
    FROM users AS u
`
