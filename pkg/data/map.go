package data

func PatchMapOverwrite(in, patch map[string]any) error {
	for k, v := range patch {
		in[k] = v
	}

	return nil
}

func PatchMap(in, patch map[string]any) map[string]any {
	for k, v := range patch {
		if _, ok := v.(map[string]any); ok {
			// value is map
			if _, ok := in[k]; ok {
				// in also contains this key already
				in[k] = PatchMap(in[k].(map[string]any), patch[k].(map[string]any))
			} else {
				in[k] = v
			}
		} else {
			// replace in key with value
			in[k] = v
		}
	}

	return in
}
