package mock

func pathParamsLimit(limit string) map[string]string {
	return map[string]string{
		limitHeaderKey: limit,
	}
}

func pathParamsLabel(labelKey string, labelValue string) map[string]string {
	return map[string]string{
		labelsHeaderKey: labelKey + "=" + labelValue,
	}
}

func pathParamsLimitAndLabel(limit string, labelKey string, labelValue string) map[string]string {
	return map[string]string{
		labelsHeaderKey: labelKey + "=" + labelValue,
		limitHeaderKey:  limit,
	}
}
