package database

import (
	"GoAPIBackEnd/internal/models"
	"context"
	_ "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
)

func GetUser(conn *pgx.Conn, name string) ([]*models.Users, error) {
	var results = make([]*models.Users, 0)
	var r pgx.Rows
	var err error
	var masterSql string

	masterSql = `SELECT id, name, password, email, role FROM users WHERE name = $1`
	r, err = conn.Query(context.Background(), masterSql, name)
	if err != nil {
		log.Printf("Error Querying the Table")
		return nil, err
	}
	for r.Next() {
		e := &models.Users{}
		err = r.Scan(&e.Id, &e.Name, &e.Password, &e.Email, &e.Role)
		if err != nil {
			log.Printf("Error Fetching Book Details")
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
		log.Println("Error Inserting Book Details")
		return err
	}
	return nil
}
