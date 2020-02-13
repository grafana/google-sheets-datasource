package googlesheets

func getDefaultValue(cellType string) interface{} {
	switch cellType {
	case "TIME":
		return nil
	case "NUMBER":
		return 0.0
	}

	return ""
}
