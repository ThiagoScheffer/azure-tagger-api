package models

type Resource struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Tags        map[string]string `json:"tags"`
	AzureID     string            `json:"azure_id"`
	CreatedUnix int64             `json:"create_unix"`
}
