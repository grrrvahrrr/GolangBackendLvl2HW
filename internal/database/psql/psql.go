package psql

import (
	"GoBeLvl2/internal/database/dbport"
	"GoBeLvl2/internal/entities"
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib" //Postgres Driver
)

var _ dbport.DbPort = &PgStorage{}

type PgData struct {
	User         entities.User
	Project      entities.Project
	Organization entities.Organization
	CorpGroup    entities.CorpGroup
	Community    entities.Community
}

type PgStorage struct {
	db *sql.DB
}

func NewPgStorage(dsn string) (*PgStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	us := &PgStorage{
		db: db,
	}

	return us, nil
}

func (pg *PgStorage) AddUser(ctx context.Context, user entities.User) error {
	dbd := &PgData{
		User: user,
	}

	_, err := pg.db.ExecContext(ctx, "INSERT", dbd.User)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgStorage) SearchEnv(ctx context.Context, params []string) (interface{}, error) {
	var i interface{}
	rows, err := pg.db.QueryContext(ctx, `SELECT`, params)
	if err != nil {
		return nil, err
	}
	rows.Close()
	return i, nil
}

func (pg *PgStorage) SearchUser(ctx context.Context, params []string) (*entities.User, error) {
	var user *entities.User
	rows, err := pg.db.QueryContext(ctx, `SELECT`, params)
	if err != nil {
		return nil, err
	}
	rows.Close()
	return user, nil
}

func (pg *PgStorage) DeleteUserFromCommunity(ctx context.Context, user entities.User, com entities.Community) error {
	dbd := &PgData{
		User:      user,
		Community: com,
	}
	_, err := pg.db.ExecContext(ctx, `DELETE`, dbd.User, dbd.Community)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) DeleteUserFromCorpGroup(ctx context.Context, user entities.User, cg entities.CorpGroup) error {
	dbd := &PgData{
		User:      user,
		CorpGroup: cg,
	}
	_, err := pg.db.ExecContext(ctx, `DELETE`, dbd.User, dbd.CorpGroup)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) DeleteUserFromOrganization(ctx context.Context, user entities.User, org entities.Organization) error {
	dbd := &PgData{
		User:         user,
		Organization: org,
	}
	_, err := pg.db.ExecContext(ctx, `DELETE`, dbd.User, dbd.Organization)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) DeleteUserFromProject(ctx context.Context, user entities.User, pr entities.Project) error {
	dbd := &PgData{
		User:    user,
		Project: pr,
	}
	_, err := pg.db.ExecContext(ctx, `DELETE`, dbd.User, dbd.Project)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) SetUserInCommunity(ctx context.Context, user entities.User, com entities.Community) error {
	dbd := &PgData{
		User:      user,
		Community: com,
	}
	_, err := pg.db.ExecContext(ctx, `INSERT`, dbd.User, dbd.Community)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) SetUserInCorpGroup(ctx context.Context, user entities.User, cg entities.CorpGroup) error {
	dbd := &PgData{
		User:      user,
		CorpGroup: cg,
	}
	_, err := pg.db.ExecContext(ctx, `INSERT`, dbd.User, dbd.CorpGroup)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) SetUserInOrganization(ctx context.Context, user entities.User, org entities.Organization) error {
	dbd := &PgData{
		User:         user,
		Organization: org,
	}
	_, err := pg.db.ExecContext(ctx, `INSERT`, dbd.User, dbd.Organization)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) SetUserInProject(ctx context.Context, user entities.User, pr entities.Project) error {
	dbd := &PgData{
		User:    user,
		Project: pr,
	}
	_, err := pg.db.ExecContext(ctx, `INSERT`, dbd.User, dbd.Project)
	if err != nil {
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *PgStorage) AddCommunity(ctx context.Context, com entities.Community) error {
	dbd := &PgData{
		Community: com,
	}
	_, err := pg.db.ExecContext(ctx, "INSERT", dbd.Community)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgStorage) AddCorpGroup(ctx context.Context, cg entities.CorpGroup) error {
	dbd := &PgData{
		CorpGroup: cg,
	}
	_, err := pg.db.ExecContext(ctx, "INSERT", dbd.CorpGroup)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgStorage) AddOrganization(ctx context.Context, org entities.Organization) error {
	dbd := &PgData{
		Organization: org,
	}
	_, err := pg.db.ExecContext(ctx, "INSERT", dbd.Organization)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgStorage) AddProject(ctx context.Context, pr entities.Project) error {
	dbd := &PgData{
		Project: pr,
	}
	_, err := pg.db.ExecContext(ctx, "INSERT", dbd.Project)
	if err != nil {
		return err
	}
	return nil
}
