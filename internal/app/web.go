package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Serve website and API
func Serve() {
	var s state

	// delete temp files with substituted env vars when the program terminates
	defer os.RemoveAll(tempFilesDir)
	if !flags.noCleanup {
		defer s.cleanup()
	}

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
	cs := buildState(&s)

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
	setupEndpoints(cs)
}

func setupEndpoints(cs *currentState) {
	r := gin.Default()

	r.GET("/releases", func(c *gin.Context) {
		var managedReleases []helmRelease
		for _, r := range cs.releases {
			managedReleases = append(managedReleases, r)
		}
		c.JSON(http.StatusOK, managedReleases)
	})
	r.StaticFile("/", "./public/index.html")
	r.Static("/js", "./public/js")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
