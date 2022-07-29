package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

func quit() {
	stopContainer()
	os.Exit(0)
}

func SelectSetup() error {
	items := []string{"Use Docker on your PC", "Help", "Quit"}

	index := 1
	if !cfg.DisableKubernetes {
		items = append(items[:index+1], items[index:]...)
		items[index] = "Use a Kubernetes Cluster"
		index++
	}

	if !cfg.DisableIG {
		items = append(items[:index+1], items[index:]...)
		items[index] = "Install FHIR package"
	}

	prompt := promptui.Select{
		Label: "Please choose how you want to run the setup. \nChoose Docker if you're running on your PC. \nIf you want to run Instant on Kubernetes, then you have should been provided credentials or have Kubernetes running on your PC.",
		Items: items,
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "SelectSetup() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {
	case "Use Docker on your PC":
		err = debugDocker()
		if err != nil {
			return err
		}
		err = selectDefaultOrCustom()

	case "Use a Kubernetes Cluster":
		err = debugKubernetes()
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Install FHIR package":
		err = selectUtil()

	case "Help":
		fmt.Println(getHelpText(true, ""))
		SelectSetup()

	case "Quit":
		quit()
	}

	return err
}

func selectUtil() error {
	fmt.Println("Enter URL for the published package")
	// prompt for url
	prompt := promptui.Prompt{
		Label: "URL",
	}

	ig_url, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "selectUtil() prompt failed")
	}

	fhir_server, params, err := selectFHIR()
	if err != nil {
		return err
	}
	fmt.Println("FHIR Server target:", fhir_server)
	err = LoadIGpackage(ig_url, fhir_server, params)
	if err != nil {
		return err
	}
	return SelectSetup()
}

func selectDefaultOrCustom() error {
	prompt := promptui.Select{
		Label: "Great, now choose an installation type",
		Items: []string{"Default Install Options", "Custom Install Options", "Quit", "Back"},
		Size:  12,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "selectDefaultOrCustom() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {
	case "Default Install Options":
		err = selectDefaultAction()
	case "Custom Install Options":
		err = selectCustomOptions()
	case "Quit":
		quit()
	case "Back":
		err = SelectSetup()
	}

	return err
}

func selectCustomOptions() error {
	index := 1
	items := []string{
		"Choose deploy action (default is init)",
		"Specify deploy packages",
		"Specify environment variable file location",
		"Specify environment variables",
		"Specify custom package locations",
		"Toggle only flag",
		"Specify Image Version",
		"Toggle dev mode (default mode is prod)",
		"Execute with current options",
		"View current options set",
		"Reset to default options",
		"Help",
		"Quit",
		"Back",
	}

	if !cfg.DisableCustomTargetSelection {
		items = append(items[:index+1], items[index:]...)
		items[index] = "Choose target launcher (default is " + cfg.DefaultTargetLauncher + ")"
	}

	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: items,
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "selectCustomOptions() prompt failed")
	}

	switch result {
	case "Choose deploy action (default is init)":
		err = setStartupAction()
	case "Choose target launcher (default is " + cfg.DefaultTargetLauncher + ")":
		err = setTargetLauncher()
	case "Specify deploy packages":
		err = setStartupPackages()
	case "Specify environment variable file location":
		err = setEnvVarFileLocation()
	case "Specify environment variables":
		err = setEnvVars()
	case "Specify custom package locations":
		err = setCustomPackages()
	case "Toggle only flag":
		err = toggleOnlyFlag()
	case "Toggle dev mode (default mode is prod)":
		err = toggleDevMode()
	case "Specify Image Version":
		err = setImageVersion()
	case "Execute with current options":
		err = printAll(false)
		if err != nil {
			return err
		}
		err = executeCommand()
	case "View current options set":
		err = printAll(true)
	case "Reset to default options":
		resetAll()
		err = printAll(true)
	case "Help":
		fmt.Println(getHelpText(true, "Custom Options"))
		return selectCustomOptions()
	case "Quit":
		quit()
	case "Back":
		err = selectDefaultOrCustom()
	}

	return err
}

func resetAll() {
	customOptions.startupAction = "init"
	customOptions.startupPackages = make([]string, 0)
	customOptions.envVarFileLocation = ""
	customOptions.envVars = make([]string, 0)
	customOptions.customPackageFileLocations = make([]string, 0)
	customOptions.onlyFlag = false
	customOptions.imageVersion = "latest"
	customOptions.targetLauncher = cfg.DefaultTargetLauncher
	customOptions.devMode = false
	fmt.Println("\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\nAll custom options have been reset to default.\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}

func setStartupAction() error {
	prompt := promptui.Select{
		Label: "Great, now choose a deploy action",
		Items: []string{"init", "destroy", "up", "down", "test", "Help", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "setStartupAction() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {
	case "init", "destroy", "up", "down", "test":
		customOptions.startupAction = result
		err = selectCustomOptions()
	case "Help":
		fmt.Println(getHelpText(true, "Deploy Commands"))
		return setStartupAction()
	case "Quit":
		quit()
	case "Back":
		err = selectCustomOptions()
	}

	return err
}

func setTargetLauncher() error {
	prompt := promptui.Select{
		Label: "Choose a target launcher",
		Items: []string{"docker", "swarm", "kubernetes", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "setTargetLauncher() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {
	case "docker", "swarm", "kubernetes":
		customOptions.targetLauncher = result
		err = selectCustomOptions()
	case "Quit":
		quit()
	case "Back":
		err = selectCustomOptions()
	}

	return err
}

var DeployCommands []string

func executeCommand() error {
	DeployCommands = []string{customOptions.startupAction}

	if len(customOptions.startupPackages) == 0 {
		fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n" +
			"Warning: No package IDs specified, all default packages will be included in your command.\n" +
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n\n")
	}

	DeployCommands = append(DeployCommands, customOptions.startupPackages...)

	if customOptions.envVarFileLocation != "" && len(customOptions.envVarFileLocation) > 0 {
		DeployCommands = append(DeployCommands, "--env-file="+customOptions.envVarFileLocation)
	}
	if customOptions.envVars != nil && len(customOptions.envVars) > 0 {
		for _, e := range customOptions.envVars {
			DeployCommands = append(DeployCommands, "-e="+e)
		}
	}
	if customOptions.customPackageFileLocations != nil && len(customOptions.customPackageFileLocations) > 0 {
		for _, c := range customOptions.customPackageFileLocations {
			DeployCommands = append(DeployCommands, "-c="+c)
		}
	}
	if customOptions.onlyFlag {
		DeployCommands = append(DeployCommands, "--only")
	}
	if customOptions.devMode {
		DeployCommands = append(DeployCommands, "--dev")
	}
	DeployCommands = append(DeployCommands, "--image-version="+customOptions.imageVersion)
	DeployCommands = append(DeployCommands, "-t="+customOptions.targetLauncher)
	return runDeployCommand(DeployCommands)
}

func printSlice(slice []string) {
	for _, s := range slice {
		fmt.Printf("-%q\n", s)
	}
	fmt.Println()
}

func printAll(loopback bool) error {
	fmt.Println("\nCurrent Custom Options Specified\n---------------------------------")
	fmt.Println("Target Launcher:")
	fmt.Printf("-%q\n", customOptions.targetLauncher)
	fmt.Println("Startup Action:")
	fmt.Printf("-%q\n", customOptions.startupAction)
	fmt.Println("Startup Packages:")
	if customOptions.startupPackages != nil && len(customOptions.startupPackages) > 0 {
		printSlice(customOptions.startupPackages)
	}
	fmt.Println("Environment Variable File Path:")
	if customOptions.envVarFileLocation != "" && len(customOptions.envVarFileLocation) > 0 {
		fmt.Printf("-%q\n", customOptions.envVarFileLocation)
	}
	fmt.Println("Environment Variables:")
	if customOptions.envVars != nil && len(customOptions.envVars) > 0 {
		printSlice(customOptions.envVars)
	}
	if customOptions.customPackageFileLocations != nil && len(customOptions.customPackageFileLocations) > 0 {
		fmt.Println("Custom Packages:")
		printSlice(customOptions.customPackageFileLocations)
	}
	fmt.Println("Image Version:")
	fmt.Printf("-%q\n", customOptions.imageVersion)

	fmt.Println("Only Flag Setting:")
	if customOptions.onlyFlag {
		fmt.Printf("-%q\n", "On")
	} else {
		fmt.Printf("-%q\n", "Off")
	}
	fmt.Println("Dev Mode Setting:")
	if customOptions.devMode {
		fmt.Printf("-%q\n\n", "On")
	} else {
		fmt.Printf("-%q\n\n", "Off")
	}

	var err error
	if loopback {
		err = selectCustomOptions()
	}

	return err
}

func setStartupPackages() error {
	if customOptions.startupPackages != nil && len(customOptions.startupPackages) > 0 {
		fmt.Println("\nCurrent Startup Packages Specified:")
		printSlice(customOptions.startupPackages)
	}
	prompt := promptui.Prompt{
		Label: "Startup Package List(Comma Delimited). e.g. core,cdr",
	}
	packageList, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "setStartupPackages() prompt failed")
	}

	startupPackages := strings.Split(packageList, ",")

	for _, p := range startupPackages {
		if !sliceContains(customOptions.startupPackages, p) {
			customOptions.startupPackages = append(customOptions.startupPackages, p)
		} else {
			fmt.Printf(p + " package already exists in the list.\n")
		}
	}

	return selectCustomOptions()
}

func setCustomPackages() error {
	if customOptions.customPackageFileLocations != nil && len(customOptions.customPackageFileLocations) > 0 {
		fmt.Println("Current Custom Packages Specified:")
		printSlice(customOptions.customPackageFileLocations)
	}
	prompt := promptui.Prompt{
		Label: "Custom Package List(Comma Delimited). e.g. " + filepath.FromSlash("../project/cdr") + "," + filepath.FromSlash("../project/demo"),
	}
	customPackageList, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "setCustomPackages() prompt failed")
	}

	newCustomPackages := strings.Split(customPackageList, ",")

	for _, cp := range newCustomPackages {
		if strings.HasPrefix(cp, "http") || strings.HasPrefix(cp, "git") {
			if !sliceContains(customOptions.customPackageFileLocations, cp) {
				customOptions.customPackageFileLocations = append(customOptions.customPackageFileLocations, cp)
			} else {
				fmt.Printf(cp + " URL already exists in the list.\n")
			}
		} else {
			exists, fileErr := fileExists(cp)
			if exists {
				if !sliceContains(customOptions.customPackageFileLocations, cp) {
					customOptions.customPackageFileLocations = append(customOptions.customPackageFileLocations, cp)
				} else {
					fmt.Printf(cp + " path already exists in the list.\n")
				}
			} else {
				fmt.Printf("\nFile at location %q could not be found due to error: %v\n", cp, fileErr)
				fmt.Println("\n-----------------\nPlease try again.\n-----------------")
			}
		}
	}

	return selectCustomOptions()
}

func setEnvVarFileLocation() error {
	if customOptions.envVarFileLocation != "" && len(customOptions.envVarFileLocation) > 0 {
		fmt.Println("Current Environment Variable File Location Specified:")
		fmt.Printf("-%q\n", customOptions.envVarFileLocation)
	}
	prompt := promptui.Prompt{
		Label: "Environment Variable file location e.g. " + filepath.FromSlash("../project/prod.env"),
	}
	envVarFileLocation, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "setEnvVarFileLocation() prompt failed")
	}
	exists, fileErr := fileExists(envVarFileLocation)
	if exists {
		customOptions.envVarFileLocation = envVarFileLocation
	} else {
		fmt.Printf("\nFile at location %q could not be found due to error: %v\n", envVarFileLocation, fileErr)
		fmt.Println("\n-----------------\nPlease try again.\n-----------------")
	}

	return selectCustomOptions()
}

func setImageVersion() error {
	if customOptions.imageVersion != "latest" && len(customOptions.imageVersion) > 0 {
		fmt.Println("Current Image Version Specified:")
		fmt.Printf("-%q\n", customOptions.imageVersion)
	}
	prompt := promptui.Prompt{
		Label: "Image Version e.g. 0.0.9",
	}
	imageVersion, err := prompt.Run()

	if err != nil {
		return errors.Wrap(err, "setImageVersion() prompt failed")
	}

	customOptions.imageVersion = imageVersion
	return selectCustomOptions()
}

func setEnvVars() error {
	if customOptions.envVars != nil && len(customOptions.envVars) > 0 {
		fmt.Println("Current Environment Variables Specified:")
		printSlice(customOptions.envVars)
	}
	prompt := promptui.Prompt{
		Label: "Environment Variable List(Comma Delimited). e.g. NODE_ENV=PROD,DOMAIN_NAME=instant.com",
	}
	envVarList, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "setEnvVars() prompt failed")
	}

	newEnvVars := strings.Split(envVarList, ",")

	for _, env := range newEnvVars {
		if !sliceContains(customOptions.envVars, env) {
			customOptions.envVars = append(customOptions.envVars, env)
		} else {
			fmt.Printf(env + " environment variable already exists in the list.\n")
		}
	}
	return selectCustomOptions()
}

func toggleOnlyFlag() error {
	customOptions.onlyFlag = !customOptions.onlyFlag
	if customOptions.onlyFlag {
		fmt.Println("Only flag is now on")
	} else {
		fmt.Println("Only flag is now off")
	}
	return selectCustomOptions()
}

func toggleDevMode() error {
	customOptions.devMode = !customOptions.devMode
	if customOptions.devMode {
		fmt.Println("Dev mode is now on")
	} else {
		fmt.Println("Dev mode is now off")
	}
	return selectCustomOptions()
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

func selectDefaultAction() error {
	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: []string{
			"init",
			"down",
			"destroy",
			"up",
			"Help",
			"Back",
			"Quit",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "selectDefaultAction() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {
	case "Help":
		fmt.Println(getHelpText(true, "Deploy Commands"))
		return selectDefaultAction()
	case "Back":
		return selectDefaultOrCustom()
	case "Quit":
		quit()
		return nil
	}

	return selectDefaultPackage(result)
}

func selectDefaultPackage(action string) error {
	var optionItems []string
	for _, p := range cfg.Packages {
		optionItems = append(optionItems, p.Name)
	}
	optionItems = append(optionItems, "All", "Back", "Quit")

	prompt := promptui.Select{
		Label: "Which package would you like to perform the action on (Packages will also invoke their dependencies automatically)",
		Items: optionItems,
		Size:  12,
	}

	i, result, err := prompt.Run()
	if err != nil {
		return errors.Wrap(err, "selectDefaultPackage() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {
	case "All":
		fmt.Println("...Setting up All Packages")
		err = RunDeployCommand([]string{action, "-t=" + cfg.DefaultTargetLauncher})
		if err != nil {
			return err
		}
		return selectDefaultAction()
	case "Quit":
		quit()
		return nil
	case "Back":
		return selectDefaultAction()
	default:
		err = RunDeployCommand([]string{cfg.Packages[i].ID, action, "-t=" + cfg.DefaultTargetLauncher})
		if err != nil {
			return err
		}
		return selectDefaultAction()
	}
}

func selectPackageCluster() error {
	prompt := promptui.Select{
		Label: "Great, now choose an action",
		Items: []string{"Initialise Core (Required, Start Here)", "Launch Facility Registry", "Launch Workforce", "Stop and Cleanup Core", "Stop and Cleanup Facility Registry", "Stop and Cleanup Workforce", "Stop All Services and Cleanup Kubernetes", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return errors.Wrap(err, "selectPackageCluster() prompt failed")
	}

	fmt.Printf("\nYou chose %q\n========================================\n", result)

	switch result {
	case "Launch Core (Required, Start Here)":
		fmt.Println("...Setting up Core Package")
		err = RunDeployCommand([]string{"core", "init", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Launch Facility Registry":
		fmt.Println("...Setting up Facility Registry Package")
		err = RunDeployCommand([]string{"facility", "up", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Launch Workforce":
		fmt.Println("...Setting up Workforce Package")
		err = RunDeployCommand([]string{"healthworker", "up", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Stop and Cleanup Core":
		fmt.Println("Stopping and Cleaning Up Core...")
		err = RunDeployCommand([]string{"core", "destroy", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Stop and Cleanup Facility Registry":
		fmt.Println("Stopping and Cleaning Up Facility Registry...")
		err = RunDeployCommand([]string{"facility", "destroy", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Stop and Cleanup Workforce":
		fmt.Println("Stopping and Cleaning Up Workforce...")
		err = RunDeployCommand([]string{"healthworker", "destroy", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Stop All Services and Cleanup Kubernetes":
		fmt.Println("Stopping and Cleaning Up Everything...")
		err = RunDeployCommand([]string{"core", "destroy", "-t=k8s"})
		if err != nil {
			return err
		}
		err = RunDeployCommand([]string{"facility", "destroy", "-t=k8s"})
		if err != nil {
			return err
		}
		err = RunDeployCommand([]string{"healthworker", "destroy", "-t=k8s"})
		if err != nil {
			return err
		}
		err = selectPackageCluster()

	case "Quit":
		quit()

	case "Back":
		err = SelectSetup()
	}

	return err
}

func selectFHIR() (result_url string, params *Params, err error) {
	prompt := promptui.Select{
		Label: "Select or enter URL for a FHIR Server",
		Items: []string{"Docker Default", "Kubernetes Default", "Use Public HAPI Server", "Enter a Server URL", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", nil, errors.Wrap(err, "selectFHIR() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)

	switch result {

	case "Docker Default":
		result_url := "http://localhost:8080/fhir"
		params := &Params{}
		params.TypeAuth = "Custom"
		params.Token = "test"
		return result_url, params, nil

	case "Kubernetes Default":
		result_url := "http://localhost:8080/fhir"
		params := &Params{}
		params.TypeAuth = "Custom"
		params.Token = "test"
		return result_url, params, nil

	case "Use Public HAPI Server":
		result_url := "http://hapi.fhir.org/baseR4"
		params := &Params{}
		params.TypeAuth = "None"
		return result_url, params, nil

	case "Enter a Server URL":
		prompt := promptui.Prompt{
			Label: "URL",
		}
		result_url, err := prompt.Run()
		if err != nil {
			return "", nil, errors.Wrap(err, "Server URL in selectFHIR() prompt failed")
		}

		params, err := selectParams()
		return result_url, params, err

	case "Quit":
		quit()
		return "", &Params{}, nil

	case "Back":
		return "", &Params{}, selectUtil()

	}
	return result_url, params, nil

}

type Params struct {
	// none, token, basic, custom
	TypeAuth  string
	Token     string
	BasicUser string
	BasicPass string
}

func selectParams() (*Params, error) {
	params := &Params{}

	prompt := promptui.Select{
		Label: "Choose authentication type",
		Items: []string{"None", "Basic", "Token", "Custom", "Quit", "Back"},
		Size:  12,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, errors.Wrap(err, "selectParams() prompt failed")
	}

	fmt.Printf("You chose %q\n========================================\n", result)
	switch result {

	case "None":
		params.TypeAuth = "None"
		return params, nil

	case "Basic":
		params.TypeAuth = "Basic"

		prompt_basic_user := promptui.Prompt{
			Label: "Basic User",
		}
		result_basic_user, err := prompt_basic_user.Run()
		if err != nil {
			return nil, errors.Wrap(err, "Case 'Basic' in selectParams() prompt failed")
		}
		params.BasicUser = result_basic_user

		prompt_basic_pass := promptui.Prompt{
			Label: "Basic Password",
		}
		result_basic_pass, err := prompt_basic_pass.Run()
		if err != nil {
			return nil, errors.Wrap(err, "Case 'Basic' in selectParams() prompt failed")
		}
		params.BasicPass = result_basic_pass

		return params, nil

	case "Token":
		params.TypeAuth = "Token"

		prompt_token := promptui.Prompt{
			Label: "Bearer Token",
		}
		result_token, err := prompt_token.Run()
		if err != nil {
			return nil, errors.Wrap(err, "Case 'Token' in selectParams() prompt failed")
		}
		params.Token = result_token
		return params, nil

	case "Custom":
		params.TypeAuth = "Custom"

		prompt_ctoken := promptui.Prompt{
			Label: "Custom Token",
		}
		result_ctoken, err := prompt_ctoken.Run()
		if err != nil {
			return nil, errors.Wrap(err, "Case 'Custom' in selectParams() prompt failed")
		}
		params.Token = result_ctoken
		return params, nil

	case "Quit":
		quit()
		return params, nil

	case "Back":
		return params, selectUtil()
	}

	return params, err
}
