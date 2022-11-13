package requests

import "github.com/verifa/coastline/server/oapi"

func (e *Engine) GetServices() *oapi.ServicesResp {
	svcMap := make(map[string]ServiceDef, 0)
	for _, tmpl := range e.templates {
		for _, svc := range tmpl.Def.Services {
			svcMap[svc.Name] = svc
		}
	}

	var resp oapi.ServicesResp
	resp.Services = make([]oapi.Service, 0, len(svcMap))
	for _, svc := range svcMap {
		resp.Services = append(resp.Services, oapi.Service{
			Name: svc.Name,
		})
	}
	return &resp
}
