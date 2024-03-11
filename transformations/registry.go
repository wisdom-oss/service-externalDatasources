package transformations

func GetTransformation(transformationKey string) interface{} {
	switch transformationKey {
	case AddQueryParametersKey:
		return AddQueryParameters{}
	case SortQueryParametersKey:
		return SortQueryParameters{}
	case SetQueryParametersKey:
		return SetQueryParameters{}
	case ExternalScriptTransformationKey:
		return ExternalScript{}
	}
	return nil
}
