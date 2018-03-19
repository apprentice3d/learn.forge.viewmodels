package server

import (
	"log"
	"net/http"

	"github.com/apprentice3d/forge-api-go-client/dm"
	"github.com/apprentice3d/forge-api-go-client/md"
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"runtime"
	"os/exec"
)

// ForgeServices holds reference to all services required in this server
type ForgeServices struct {
	oauth.TwoLeggedAuth
	dm.BucketAPI
	md.ModelDerivativeAPI
}

// StartServer setups the endpoints and starts the http server with endpoints expected by the frontend
func StartServer(port, clientID, clientSecret string) {

	service := ForgeServices{
		oauth.NewTwoLeggedClient(clientID, clientSecret),
		dm.NewBucketAPIWithCredentials(clientID, clientSecret),
		md.NewAPIWithCredentials(clientID, clientSecret),
	}

	// serving static files
	static := http.FileServer(http.Dir("www"))
	http.Handle("/", static)

	// defining other endpoints
	http.HandleFunc("/api/forge/oauth/token", service.getAccessToken)
	http.HandleFunc("/api/forge/oss/buckets", service.manageBuckets)
	http.HandleFunc("/api/forge/oss/objects", service.manageObjects)
	http.HandleFunc("/api/forge/modelderivative/jobs", service.translateObject)

	go openBrowser("http://localhost" + port)
	if err := http.ListenAndServe("localhost" + port, nil); err != nil {
		log.Fatal(err.Error())
	}

}


// Idea taken from https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}