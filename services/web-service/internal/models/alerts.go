package models

type AlertType string

const (
	AlertSuccess AlertType = "success"
	AlertError   AlertType = "danger"
	AlertWarning AlertType = "warning"
	AlertInfo    AlertType = "info"
)

type Alert struct {
	Type AlertType
	Msg  string
}
