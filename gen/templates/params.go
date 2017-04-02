package templates

var ArgsEmpty string = `if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs`

var Args string = `var args %s
		if reqBody.HasParams() {
			err := json.Unmarshal(*reqBody.Params, &args)
			if err != nil {
				return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "Can't unmarshal params to args structure'", err.Error())
			}
		}`
