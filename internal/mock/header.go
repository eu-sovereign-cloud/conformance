package mock

func headerParamsGeneric(authToken string) map[string]string {
	return map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + authToken,
	}
}

func headerParamsLimit(authToken string, limit string) map[string]string {
	return map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + authToken,
		limitHeaderKey:             limit,
	}
}

func headerParamsLabel(authToken string, labelKey string, labelValue string) map[string]string {
	return map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + authToken,
		labelsHeaderKey:            labelKey + "=" + labelValue,
	}
}

func headerParamsLimitAndLabel(authToken string, limit string, labelKey string, labelValue string) map[string]string {
	return map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + authToken,
		labelsHeaderKey:            labelKey + "=" + labelValue,
		limitHeaderKey:             limit,
	}
}
