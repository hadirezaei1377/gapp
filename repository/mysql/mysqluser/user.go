package mysqluser

import (
	"context"
	"database/sql"
	"fmt"
	"gapp/entity"
	"gapp/pkg/errmsg"
	"gapp/pkg/richerror"
	"gapp/repository/mysql"
	"time"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (d *DB) Register(u entity.User) (entity.User, error) {
	res, err := d.conn.Conn().Exec(`insert into users(name, phone_number, password, role) values(?, ?, ?, ?)`,
		u.Name, u.PhoneNumber, u.Password, u.Role.String())
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	// error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.conn.Conn().QueryRow(`select * from users where phone_number = ?`, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		// TODO - log unexpected error for better observability
		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserByID(ctx context.Context, userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"

	row := d.conn.Conn().QueryRowContext(ctx, `select * from users where id = ?`, userID)
	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(err).
				WithMessage(errmsg.ErrorMsgNotFound).WithKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgCantScanQueryResult).WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var createdAt time.Time
	var user entity.User

	var roleStr string

	err := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password, &roleStr)

	user.Role = entity.MapToRoleEntity(roleStr)

	return user, err
}
