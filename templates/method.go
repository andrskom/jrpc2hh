package templates

var Method string = `case "{{.Method}}":
		{{.ArgsBlock}}
		{{.ResultBlock}}
		err := s.{{.Method}}(args, &res)
		if err != nil {
			return nil, jModels.NewError(jModels.ErrorCodeInternalError, "Internal error", err.Error())
		}
		return res, nil`
