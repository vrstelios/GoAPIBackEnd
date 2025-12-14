package database

import (
	"GoAPIBackEnd/internal/models"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"strings"
)

func GetUser(conn *pgx.Conn, id string, name string) ([]*models.Users, error) {
	var results = make([]*models.Users, 0)
	var r pgx.Rows
	var err error
	var masterSql string

	masterSql = `SELECT id, name, password, email, role, coach_id, created_at FROM users`
	var counter int = 0
	var params []interface{}
	if len(id) > 0 {
		counter++
		params = append(params, id)
		masterSql = masterSql + fmt.Sprintf(" and id = $%d", counter)
	}
	if len(name) > 0 {
		counter++
		params = append(params, name)
		masterSql = masterSql + fmt.Sprintf(" and name = $%d", counter)
	}
	masterSql = strings.Replace(masterSql, " and ", " where ", 1)

	r, err = conn.Query(context.Background(), masterSql, params...)
	if err != nil {
		return nil, err
	}
	for r.Next() {
		e := &models.Users{}
		err = r.Scan(&e.Id, &e.Name, &e.Password, &e.Email, &e.Role, &e.CoachId, &e.CreatedAt)
		if err != nil {
			return results, err
		}
		results = append(results, e)
	}
	r.Close()

	return results, nil
}

func PostUser(conn *pgx.Conn, user models.Users) error {
	masterSql := `INSERT INTO users (id, name, password, email, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := conn.Exec(context.Background(), masterSql, user.Id, user.Name, user.Password, user.Email, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func PostExercise(conn *pgx.Conn, ex models.Exercises) error {
	masterSql := `INSERT INTO exercises (id, name, description, category) VALUES ($1, $2, $3, $4)`
	_, err := conn.Exec(context.Background(), masterSql, ex.Id, ex.Name, ex.Description, ex.Category)
	if err != nil {
		return err
	}
	return nil
}

func GetExercise(conn *pgx.Conn, id string, name string) ([]*models.Exercises, error) {
	var results = make([]*models.Exercises, 0)
	var r pgx.Rows
	var err error
	var masterSql string

	masterSql = `SELECT id, name, description, category, created_at FROM exercises`
	var counter int = 0
	var params []interface{}
	if len(id) > 0 {
		counter++
		params = append(params, id)
		masterSql = masterSql + fmt.Sprintf(" and id = $%d", counter)
	}
	if len(name) > 0 {
		counter++
		params = append(params, name)
		masterSql = masterSql + fmt.Sprintf(" and name = $%d", counter)
	}
	masterSql = strings.Replace(masterSql, " and ", " where ", 1)

	r, err = conn.Query(context.Background(), masterSql, params...)
	if err != nil {
		return nil, err
	}
	for r.Next() {
		e := &models.Exercises{}
		tmpDesc := sql.NullString{}
		tmpCat := sql.NullString{}
		err = r.Scan(&e.Id, &e.Name, &tmpDesc, &tmpCat, &e.CreatedAt)
		if err != nil {
			return results, err
		}
		if tmpDesc.Valid {
			e.Description = &tmpDesc.String
		}
		if tmpCat.Valid {
			e.Category = &tmpCat.String
		}
		results = append(results, e)
	}
	r.Close()

	return results, nil
}

func GetCoach(conn *pgx.Conn, name string, addRelations bool) ([]*models.Coach, error) {
	var results = make([]*models.Coach, 0)
	var r pgx.Rows
	var err error
	var masterSql string

	masterSql = `SELECT id, name, user_id FROM coach`
	var counter int = 0
	var params []interface{}
	if len(name) > 0 {
		counter++
		params = append(params, name)
		masterSql = masterSql + fmt.Sprintf(" and name = $%d", counter)
	}

	masterSql = strings.Replace(masterSql, " and ", " where ", 1)

	r, err = conn.Query(context.Background(), masterSql, params...)
	if err != nil {
		return nil, err
	}
	for r.Next() {
		e := &models.Coach{}
		err = r.Scan(&e.Id, &e.Name, &e.UserId)
		if err != nil {
			return results, err
		}

		results = append(results, e)
	}
	r.Close()

	if addRelations {
		for i := range results {
			relations := models.CoachRelations{}
			user, err := GetUser(conn, results[i].UserId, "") // Αλλαγή εδώ
			if err != nil {
				return nil, err
			}

			if len(user) > 0 {
				relations.User = user[0]
			} else {
				relations.User = nil
			}

			results[i].Relations = &relations
		}
	}

	return results, nil
}

func PostCoach(conn *pgx.Conn, coach models.Coach) error {
	masterSql := `INSERT INTO coach (id, name, user_id) VALUES ($1, $2, $3)`
	_, err := conn.Exec(context.Background(), masterSql, coach.Id, coach.Name, coach.UserId)
	if err != nil {
		return err
	}
	return nil
}
