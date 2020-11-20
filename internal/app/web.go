package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var stop = make(chan bool)

// Serve website and API
func Serve() {
	var s state
	var cs currentState
	cs.isReady = false
	// delete temp files with substituted env vars when the program terminates
	defer os.RemoveAll(tempFilesDir)
	if !flags.noCleanup {
		defer s.cleanup()
	}

	setupEndpoints(&cs)

	flags.readState(&s)
	if len(s.GroupMap) > 0 {
		s.TargetMap = s.getAppsInGroupsAsTargetMap()
		if len(s.TargetMap) == 0 {
			log.Info("No apps defined with -group flag were found, exiting...")
			os.Exit(0)
		}
	}
	if len(s.TargetMap) > 0 {
		s.TargetApps = s.getAppsInTargetsOnly()
		s.TargetNamespaces = s.getNamespacesInTargetsOnly()
		if len(s.TargetApps) == 0 {
			log.Info("No apps defined with -target flag were found, exiting...")
			os.Exit(0)
		}
	}
	settings = s.Settings
	curContext = s.Context

	// set the kubecontext to be used Or create it if it does not exist
	log.Info("Setting up kubectl...")
	if !setKubeContext(settings.KubeContext) {
		if err := createContext(&s); err != nil {
			log.Fatal(err.Error())
		}
	}

	// add repos -- fails if they are not valid
	log.Info("Setting up helm...")
	if err := addHelmRepos(s.HelmRepos); err != nil && !flags.destroy {
		log.Fatal(err.Error())
	}

	if flags.apply || flags.dryRun || flags.destroy {
		// add/validate namespaces
		if !flags.noNs {
			log.Info("Setting up namespaces...")
			if flags.nsOverride == "" {
				addNamespaces(&s)
			} else {
				createNamespace(flags.nsOverride)
				s.overrideAppsNamespace(flags.nsOverride)
			}
		}
	}

	// if !flags.skipValidation {
	// 	log.Info("Validating charts...")
	// 	// validate charts-versions exist in defined repos
	// 	if err := validateReleaseCharts(&s); err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// } else {
	// 	log.Info("Skipping charts' validation.")
	// }

	log.Info("Loading Image history...")
	LoadImageVersions("default")

	log.Info("Preparing plan...")
	buildState(&s, &cs)

	// p := cs.makePlan(&s)
	// if !flags.keepUntrackedReleases {
	// 	cs.cleanUntrackedReleases(&s, p)
	// }

	// p.sort()
	// p.print()
	// if flags.debug {
	// 	p.printCmds()
	// }
	// p.sendToSlack()

	// if flags.apply || flags.dryRun || flags.destroy {
	// 	p.exec()
	// }
	<-stop

}

func setupEndpoints(cs *currentState) {
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"isReady":  cs.isReady,
			"releases": cs.releases,
		})
	})
	r.StaticFile("/", "./public/index.html")
	r.Static("/js", "./public/js")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that it wont block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		stop <- true
	}()
}
