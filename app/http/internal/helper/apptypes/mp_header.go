package apptypes

// MPHeader  通用请求 header
type MPHeader struct {
	headerI18N
	headerUTM

	MPPlatform string `header:"X-Mp-Platform" form:"h_mpp"`
	MPAppID    string `header:"X-Mp-App-Id" form:"h_mpi"`
}
