package repository

import (
	"database/sql"
	"fmt"

	"github.com/0xatanda/profileIntelligence/internal/model"
)

type Repo struct {
	DB *sql.DB
}

// --------------------
func (r *Repo) FindByName(name string) (*model.Profile, error) {

	row := r.DB.QueryRow(`
		SELECT id, name, gender, gender_probability,
		       sample_size, age, age_group,
		       country_id, country_probability, created_at
		FROM profiles
		WHERE LOWER(name)=LOWER($1)
	`, name)

	var p model.Profile

	err := row.Scan(
		&p.ID, &p.Name, &p.Gender, &p.GenderProbability,
		&p.SampleSize, &p.Age, &p.AgeGroup,
		&p.CountryID, &p.CountryProbability, &p.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// --------------------
func (r *Repo) Create(p *model.Profile) error {

	_, err := r.DB.Exec(`
		INSERT INTO profiles (
			id, name, gender, gender_probability,
			sample_size, age, age_group,
			country_id, country_probability, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`,
		p.ID, p.Name, p.Gender, p.GenderProbability,
		p.SampleSize, p.Age, p.AgeGroup,
		p.CountryID, p.CountryProbability, p.CreatedAt,
	)

	return err
}

// --------------------
func (r *Repo) FindByID(id string) (*model.Profile, error) {

	row := r.DB.QueryRow(`
		SELECT id, name, gender, gender_probability,
		       sample_size, age, age_group,
		       country_id, country_probability, created_at
		FROM profiles
		WHERE id=$1
	`, id)

	var p model.Profile

	err := row.Scan(
		&p.ID, &p.Name, &p.Gender, &p.GenderProbability,
		&p.SampleSize, &p.Age, &p.AgeGroup,
		&p.CountryID, &p.CountryProbability, &p.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// --------------------
func (r *Repo) FindAll(gender, countryID, ageGroup string) ([]model.Profile, error) {

	query := `
		SELECT id, name, gender, gender_probability,
		       sample_size, age, age_group,
		       country_id, country_probability, created_at
		FROM profiles
		WHERE 1=1
	`

	args := []interface{}{}
	i := 1

	if gender != "" {
		query += fmt.Sprintf(" AND LOWER(gender)=LOWER($%d)", i)
		args = append(args, gender)
		i++
	}

	if countryID != "" {
		query += fmt.Sprintf(" AND country_id=$%d", i)
		args = append(args, countryID)
		i++
	}

	if ageGroup != "" {
		query += fmt.Sprintf(" AND LOWER(age_group)=LOWER($%d)", i)
		args = append(args, ageGroup)
		i++
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Profile

	for rows.Next() {
		var p model.Profile

		err := rows.Scan(
			&p.ID, &p.Name, &p.Gender,
			&p.GenderProbability, &p.SampleSize,
			&p.Age, &p.AgeGroup,
			&p.CountryID, &p.CountryProbability,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

// --------------------
func (r *Repo) Delete(id string) error {

	res, err := r.DB.Exec("DELETE FROM profiles WHERE id=$1", id)
	if err != nil {
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}
