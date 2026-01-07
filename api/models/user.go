
package models

type User struct {
	IdUser int    `json:"id_user" db:"id_user"`
	Nama   string `json:"nama" db:"nama"`
	Role   string `json:"role" db:"role"`
}