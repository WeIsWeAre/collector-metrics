package storage

type Repository interface {
	Store(string, string) error
	Get(string) (string, error)
}

type MetricRepository struct{}

func (r *MetricRepository) Store(temp1, temp2 string) error {
	//_, err := r.db.Exec("INSERT INTO urls (short_code, long_url) VALUES (?, ?)", shortCode, longURL)
	return nil
}

func (r *MetricRepository) Get(temp string) (string, error) {
	//_, err := r.db.Exec("INSERT INTO urls (short_code, long_url) VALUES (?, ?)", shortCode, longURL)
	return "", nil
}
