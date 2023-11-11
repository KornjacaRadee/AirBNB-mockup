package services

type ConnectionService struct {
	address string
}

func NewConnectionService(address string) (ConnectionService, error) {
	return ConnectionService{
		address: address,
	}, nil
}

//
//func (c ConnectionService) getByUser(ctx context.Context, user domain.User) ([]domain.User, error) {
//	req, err := http.NewRequest(http.MethodGet, c.address+"/users/"+user.Id.String()+"/connections", nil)
//	if err != nil {
//		return nil, err
//	}
//	token := ctx.Value("token")
//	if token != nil {
//		req.Header["Auth-Token"] = []string{token.(string)}
//	}
//
//	httpResp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	resp := struct {
//		Connections []struct {
//			Id       string
//			Username string
//		}
//	}{
//		Connections: make([]struct{ Id, Username string }, 0),
//	}
//	err = json.NewDecoder(httpResp.Body).Decode(&resp)
//	if err != nil {
//		return nil, err
//	}
//
//	conns := make([]domain.User, 0)
//	for _, conn := range resp.Connections {
//		connId, err := uuid.Parse(conn.Id)
//		if err != nil {
//			continue
//		}
//		conns = append(conns, domain.User{Id: connId})
//	}
//	return conns, nil
//}
//
//func (c ConnectionService) connected(ctx context.Context, users [2]domain.User) bool {
//	conns, err := c.getByUser(ctx, users[0])
//	if err != nil {
//		conns, err = c.getByUser(ctx, users[1])
//		if err != nil {
//			return false
//		}
//	}
//	for _, conn := range conns {
//		if users[1].Equals(conn) {
//			return true
//		}
//	}
//	return false
//}
