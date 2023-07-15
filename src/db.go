package main

import (
	"database/sql"
	"fmt"
	env "github.com/logotipiwe/dc_go_env_lib"
	. "github.com/logotipiwe/dc_go_utils/src"
	"strings"
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

var (
	db *sql.DB
)

func InitDb() error {
	connectionStr := fmt.Sprintf("%v:%v@tcp(%v)/%v", env.GetDbLogin(), env.GetDbPassword(),
		env.GetDbHost(), env.GetDbName())
	conn, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return err
	}
	if err := conn.Ping(); err != nil {
		println(fmt.Sprintf("Error connecting database: %s", err))
		return err
	}
	db = conn
	println("Database connected!")
	return nil
}

func toNullable(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (service *Service) save() error {
	_, err := db.Exec("insert into services (id, name) values (?, ?);", service.Id, service.Name)
	if err == nil {
		fmt.Printf("Servie with name %s saved!", service.Name)
	}
	return err
}

func GetAllServices() ([]Service, error) {
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
	var result []PropResult
	rows, err := db.Query("SELECT id,service,namespace,is_active,name,`value` FROM config_entries")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		n, err := scanPropResult(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}

	properties := Map(result, toModel)

	return properties, nil
}

func scanPropResult(rows *sql.Rows) (PropResult, error) {
	n := PropResult{}
	err := rows.Scan(&n.Id, &n.ServiceId, &n.NamespaceId, &n.Active, &n.Name, &n.Value)
	if err != nil {
		return PropResult{}, err
	}
	return n, nil
}

func GetPropsByNamespaceAndService(namespaceName, serviceName string) ([]Property, error) {
	isAllServices := serviceName == "*"
	isDefaultNamespace := namespaceName == "default"
	if isAllServices {
		serviceName = ""
	}
	if isDefaultNamespace {
		namespaceName = ""
	}
	rows, err := db.Query("select * from config_entries c "+
		"where is_active AND ("+
		"	(? = '' AND namespace is null)"+
		"	OR (namespace = (select id from namespaces where namespaces.name = ?)) "+
		") "+
		"AND ("+
		"	service is null"+
		"	OR (service = (select id from services where services.name = ?))"+
		") and ("+
		"	("+
		"	(select count(*)"+
		"		from config_entries c2"+
		"		where c.namespace = c2.namespace"+
		"		and c.name = c2.name"+
		"		and c2.service is not null"+
		"		and c2.is_active) > 0"+
		"		AND c.service is not null"+
		"	)"+
		"	OR ("+
		"		(select count(*)"+
		"			from config_entries c2"+
		"			where c.namespace = c2.namespace"+
		"			and c.name = c2.name"+
		"			and c2.service is not null"+
		"			and c2.is_active) = 0"+
		"		AND c.service is null"+
		"	)"+
		")", namespaceName, namespaceName, serviceName)
	if err != nil {
		return nil, err
	}
	var res []PropResult
	for rows.Next() {
		p, err := scanPropResult(rows)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	properties := Map(res, func(p PropResult) Property {
		return toModel(p)
	})
	return properties, nil
}

func GetProp(id string) (Property, error) {
	var r PropResult
	err := db.QueryRow("SELECT id,service,namespace,is_active,name,`value` FROM config_entries WHERE id = ?",
		id).Scan(&r.Id, &r.ServiceId, &r.NamespaceId, &r.Active, &r.Name, &r.Value)
	if err != nil {
		println(err.Error())
		return Property{}, err
	}
	return toModel(r), nil
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
			"VALUES (?,?,?,?,?)", p.Id, toNullable(p.ServiceId), toNullable(p.NamespaceId), p.Name, p.Value)
		if err != nil {
			println(err)
			return err
		}
	}
	return nil
}

func DeleteProperty(id string) error {
	_, err := db.Exec("DELETE FROM config_entries WHERE id = ?", id)
	return err
}

func importConfig(namespaces []Namespace, services []Service, props []Property) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = clearProps(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = clearNamespaces(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = clearServices(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = importNamespaces(tx, namespaces)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = importServices(tx, services)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = importProps(tx, props)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("Import done successfully!")
	return nil
}

func clearProps(tx *sql.Tx) error {
	_, err := tx.Exec("delete from config_entries")
	if err != nil {
		return err
	}
	return nil
}

func clearNamespaces(tx *sql.Tx) error {
	_, err := tx.Exec("delete from namespaces")
	if err != nil {
		return err
	}
	return nil
}

func clearServices(tx *sql.Tx) error {
	_, err := tx.Exec("delete from services")
	if err != nil {
		return err
	}
	return nil
}

func importProps(tx *sql.Tx, props []Property) error {
	var valuesStr []string
	var values []interface{}

	for _, prop := range props {
		valuesStr = append(valuesStr, "(?,?,?,?,?,?)")
		values = append(values,
			prop.Id,
			prop.Name,
			prop.Value,
			toNullable(prop.NamespaceId),
			toNullable(prop.ServiceId),
			prop.Active,
		)
	}

	query := fmt.Sprintf("INSERT INTO config_entries (id, name, value, namespace, service, is_active)"+
		" VALUES %s", strings.Join(valuesStr, ","))

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	println("Props imported successfully")
	return nil
}

func importNamespaces(tx *sql.Tx, namespaces []Namespace) error {
	var valuesStr []string
	var values []interface{}

	for _, namespace := range namespaces {
		valuesStr = append(valuesStr, "(?,?)")
		values = append(values,
			namespace.Id,
			namespace.Name,
		)
	}

	_, err := tx.Exec("delete from namespaces")
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO namespaces (id, name) VALUES %s", strings.Join(valuesStr, ","))

	_, err = tx.Exec(query, values...)
	if err != nil {
		return err
	}

	println("Namespaces imported successfully")
	return nil
}

func importServices(tx *sql.Tx, services []Service) error {
	var valuesStr []string
	var values []interface{}

	for _, service := range services {
		valuesStr = append(valuesStr, "(?,?)")
		values = append(values,
			service.Id,
			service.Name,
		)
	}

	query := fmt.Sprintf("INSERT INTO services (id, name) VALUES %s", strings.Join(valuesStr, ","))

	_, err := tx.Exec(query, values...)
	if err != nil {
		return err
	}

	println("Services imported successfully")
	return nil
}
