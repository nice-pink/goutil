package data

func PatchMapOverwrite(in, patch map[string]any) map[string]any {
	if patch == nil {
		return in
	}
	if in == nil {
		return patch
	}

	for k, v := range patch {
		in[k] = v
	}
	return in
}

func PatchMap(in, patch map[string]any) map[string]any {
	if patch == nil {
		return in
	}
	if in == nil {
		return patch
	}

	for k, v := range patch {
		if _, ok := v.(map[string]any); ok {
			// value is map
			if _, ok := in[k]; ok {
				// in also contains this key -> go deeper into map
				in[k] = PatchMap(in[k].(map[string]any), patch[k].(map[string]any))
			} else {
				// in does not contain key -> set value for key from patch
				in[k] = v
			}
		} else {
			// value is not a map -> replace in value with patch value
			in[k] = v
		}
	}

	return in
}
