package validate

import "github.com/karst-pulse/karst-pulse/internal/model"

func SpringInput(input model.SpringInput) error {
	var errors Errors
	Required(&errors, "code", input.Code)
	Length(&errors, "code", input.Code, 2, 32)
	Required(&errors, "name", input.Name)
	Length(&errors, "name", input.Name, 2, 100)
	Required(&errors, "region", input.Region)
	Required(&errors, "aquifer", input.Aquifer)
	Finite(&errors, "latitude", input.Latitude)
	Finite(&errors, "longitude", input.Longitude)
	Range(&errors, "latitude", input.Latitude, -90, 90)
	Range(&errors, "longitude", input.Longitude, -180, 180)
	if !validSpringStatus(input.Status) {
		errors.Add("status", "must be active, inactive, or watch")
	}
	return errors.Err()
}

func validSpringStatus(value string) bool {
	return value == string(model.SpringActive) || value == string(model.SpringInactive) || value == string(model.SpringWatch)
}
