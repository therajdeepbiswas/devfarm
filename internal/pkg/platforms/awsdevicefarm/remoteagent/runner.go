package remoteagent

import (
	"fmt"
	"github.com/dena/devfarm/internal/pkg/executor"
	"github.com/dena/devfarm/internal/pkg/executor/adb"
	"github.com/dena/devfarm/internal/pkg/executor/iosdeploy"
	"github.com/dena/devfarm/internal/pkg/logging"
	"github.com/dena/devfarm/internal/pkg/platforms"
)

type Runner func() error

type RunnerBag interface {
	GetLogger() logging.SeverityLogger
	GetEnvGetter() executor.EnvGetter
	GetFinder() executor.ExecutableFinder
	GetExecutor() executor.Executor
	GetInteractiveExecutor() executor.InteractiveExecutor
}

func NewRunner(bag RunnerBag) Runner {
	return func() error {
		getEnv := bag.GetEnvGetter()
		getEnvVars := newEnvVarsGetter(getEnv)

		envVars, envVarsErr := getEnvVars()
		if envVarsErr != nil {
			return envVarsErr
		}

		switch envVars.OSName {
		case platforms.OSIsIOS:
			getIOSEnvVars := newIOSSpecificEnvVarsGetter(getEnv)

			iosEnvVars, iosEnvVarsErr := getIOSEnvVars()
			if iosEnvVarsErr != nil {
				return iosEnvVarsErr
			}

			udid := iosdeploy.UDID(envVars.DeviceUDID)
			unarchivedAppPath := iosdeploy.UnarchivedAppPath(iosEnvVars.UnarchivedAppPath)

			iosDeployCmd := iosdeploy.NewExecutor(
				bag.GetLogger(),
				iosEnvVars.IOSDeployBin,
			)
			runIOSApp := newIOSAppRunner(
				bag.GetLogger(),
				iosdeploy.NewAppLauncher(iosDeployCmd),
			)
			return runIOSApp(udid, unarchivedAppPath, platforms.IOSArgs(envVars.AppArgs), envVars.Lifetime)

		case platforms.OSIsAndroid:
			getAndroidEnvVars := newAndroidSpecificEnvVarsGetter(getEnv)
			androidEnvVars, androidEnvVarsErr := getAndroidEnvVars()
			if androidEnvVarsErr != nil {
				return androidEnvVarsErr
			}

			packageName := adb.PackageName(androidEnvVars.AppID)
			adbCmd := adb.NewExecutor(bag.GetFinder(), bag.GetExecutor())
			interactiveAdbCmd := adb.NewInteractiveExecutor(bag.GetFinder(), bag.GetInteractiveExecutor())
			getProp := adb.NewPropGetter(adbCmd)
			runAndroidApp := newAndroidRunner(
				bag.GetLogger(),
				packageName,
				adb.NewSerialNumberGetter(adbCmd),
				adb.NewWaitUntilBecomeReady(adb.NewReadyDetector(getProp), executor.NewWaiter()),
				adb.NewMainIntentFinder(adbCmd),
				adb.NewActivityStarter(adbCmd),
				newAndroidWatcher(bag.GetLogger(), adb.NewActivityMonitor(interactiveAdbCmd)),
			)
			return runAndroidApp(platforms.AndroidIntentExtras(envVars.AppArgs), envVars.Lifetime)

		default:
			return fmt.Errorf("unsupported os: %q", envVars.OSName)
		}
	}
}