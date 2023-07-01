package main

import (
	"database/sql"
	"fmt"
	env "github.com/logotipiwe/dc_go_env_lib"
	"log"
)
import _ "github.com/go-sql-driver/mysql"

type PropResult struct {
	Id          string
	Name        string
	Value       string
	NamespaceId sql.NullString
	ServiceId   sql.NullString
	Active      bool
}

func ConnectDb() *sql.DB {
	println("Database connected!")
	connectionStr := fmt.Sprintf("%v:%v@tcp(%v)/%v", env.GetDbLogin(), env.GetDbPassword(),
		env.GetDbHost(), env.GetDbName())
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func nullable(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (service *Service) save() error {
	db := ConnectDb()
	_, err := db.Exec("insert into services (id, name) values (?, ?);", service.Id, service.Name)
	if err == nil {
		fmt.Printf("Servie with name %s saved!", service.Name)
	}
	return err
}

func GetAllServices() ([]Service, error) {
	db := ConnectDb()
	var services []Service
	rows, err := db.Query("SELECT * FROM services")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		s := Service{}
		err = rows.Scan(&s.Id, &s.Name)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func GetAllNamespaces() ([]Namespace, error) {
	db := ConnectDb()
	var namespaces []Namespace
	rows, err := db.Query("SELECT * FROM namespaces")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		n := Namespace{}
		err = rows.Scan(&n.Id, &n.Name)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, n)
	}
	return namespaces, nil
}
func GetAllProps() ([]Property, error) {
	db := ConnectDb()
	var result []PropResult
	rows, err := db.Query("SELECT id,service,namespace,is_active,name,`value` FROM config_entries")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		n := PropResult{}
		err = rows.Scan(&n.Id, &n.ServiceId, &n.NamespaceId, &n.Active, &n.Name, &n.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}

	properties := Map(result, toModel)

	return properties, nil
}

func GetProp(id string) Property {
	db := ConnectDb()
	var r PropResult
	err := db.QueryRow("SELECT id,service,namespace,is_active,name,`value` FROM config_entries WHERE id = ?",
		id).Scan(&r.Id, &r.ServiceId, &r.NamespaceId, &r.Active, &r.Name, &r.Value)
	if err != nil {
		log.Fatalln(err)
		return Property{}
	}
	return toModel(r)
}

func toModel(r PropResult) Property {
	return Property{
		r.Id,
		r.Name,
		r.Value,
		r.NamespaceId.String,
		r.ServiceId.String,
		r.Active,
	}
}

func (p Property) save() error {
	db := ConnectDb()
	var exists bool
	err := db.QueryRow("SELECT count(*) > 0 FROM config_entries WHERE id = ?", p.Id).Scan(&exists)
	if err != nil {
		println(err)
		return err
	}
	if exists {
		_, err := db.Exec("UPDATE config_entries SET name = ?, value = ?, is_active = ? WHERE id = ?",
			p.Name, p.Value, p.Active, p.Id)
		if err != nil {
			println(err)
			return err
		}
	} else {
		_, err := db.Exec("insert into config_entries (id, service, namespace, name, value) "+
			"VALUES (?,?,?,?,?)", p.Id, nullable(p.ServiceId), nullable(p.NamespaceId), p.Name, p.Value)
		if err != nil {
			println(err)
			return err
		}
	}
	return nil
}

func DeleteProperty(id string) error {
	db := ConnectDb()
	_, err := db.Exec("DELETE FROM config_entries WHERE id = ?", id)
	return err
}
