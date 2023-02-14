package data

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Models struct{
	User User
}

type User struct{
	ID int `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"_"`
	Active    int `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

const dbTimeOut = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models{
	db = dbPool
	return Models{
		User:User{},
	}
}
// Exp ::::::: returning Pointer to a struct From Function !!!!
//In many cases this is preferable
//
//func (s Store) GetCar() *Car
//because it is more convenient and readable, but has consequences. All variables such as Car are created inside the function which means they are placed onto stack. After function return this memory of stack is marked as invalid and can be used again. It a bit differs for pointer values such as *Car. Because pointer is virtually means you want to share the value with other scope and return an address, the value has to be stored somewhere in order to be available for calling function. It is copied onto heap memory and stays there until garbage collection finds no references to it.
//
//It implies overheads:
//
//copying values from stack to heap
//additional work for garbage collection
//The overheads is not significant if the value is relatively small. This is a reason why we have to pass an argument in io.Reader and io.Writer rather than have the value in return.
//
//If you'd like to dive yourself into guts follow the links: Language Mechanics On Stacks And Pointers and Bad Go: pointer returns

func (u *User)GetByEmail(email string)(*User, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := "select id,email,first_name,last_name from users where email=$1"

	row, err := db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}

	var user *User
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("Error scanning", err)
		return nil, err
	}
	return user, nil
}

func (u *User)GetAll()([]*User, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := "select * from users order by last_name"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func(u *User)Update()error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `update users set
	email = $1,
	first_name = $2,
	last_name = $3,
	user_active = $4,
	updated_at = $5,
    where id = $6
	`
	_ , err := db.ExecContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		u.UpdatedAt,
		u.ID)

	if err != nil {
		return err
	}
	return nil
}

func(u *User)Delete()error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `delete from users 
    where id = $1
	`
	_ , err := db.ExecContext(ctx, query,
		u.ID)

	if err != nil {
		return err
	}
	return nil
}

func DeleteById(id int)error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	query := `delete from users
    where id = $1
	`
	_ , err := db.ExecContext(ctx, query,
		id)

	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	//The RETURN statement is used to unconditionally and immediately end an SQL procedure by returning the flow of control to the caller of the stored procedure.
	//When the RETURN statement runs, it must return an integer value. If the return value is not provided, the default is 0.
	//The RETURNING keyword in PostgreSQL gives an opportunity to return from the insert or update statement the values of any columns after the insert or update was run.

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
