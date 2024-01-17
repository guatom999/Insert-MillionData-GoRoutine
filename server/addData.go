package server

import "onemildata/modules"

func (s *server) addDataService() {

	addData := modules.NewAddData(s.db)

	s.app.GET("/test", addData.InsertData)
	s.app.GET("/test2", addData.InsertDataTwo)

	s.app.GET("/migrate", addData.MigrateData)

}
