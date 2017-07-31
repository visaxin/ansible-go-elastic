package metadata

type DataSource interface {
	Save() error
	Update() error
	Delete() error
}


