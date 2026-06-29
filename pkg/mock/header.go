package mock

func PathParamsLimit(limit string) map[string]string {
	return map[string]string{
		limitHeaderKey: limit,
	}
}

func PathParamsLabel(labelKey string, labelValue string) map[string]string {
	return map[string]string{
		labelsHeaderKey: labelKey + "=" + labelValue,
	}
}

func PathParamsLimitAndLabel(limit string, labelKey string, labelValue string) map[string]string {
	return map[string]string{
		labelsHeaderKey: labelKey + "=" + labelValue,
		limitHeaderKey:  limit,
	}
}
