package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func pull(repo string) bool {
	cmd := exec.Command("docker", "pull", repo)
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println(exitError.ExitCode())
			return false
		}
	}
	fmt.Println("image - " + repo + "found; pulled")
	return true
}

func getVars(Args []string) ( /*string, string, string, */ string, string, string) {
	// gitrepo := ""
	// username := ""
	// password := ""
	androidVersion := "30"
	androidVersionBuild := "30.0.2"
	gradleVersion := "6.1.1"

	// reGIT := regexp.MustCompile(`(?miU)(http).*git$`)
	// reUserName := regexp.MustCompile(`(?miU)^(-u)`)
	// rePassword := regexp.MustCompile(`(?miU)^(-p)`)
	reAndroidVersion := regexp.MustCompile(`(?miU)^(-g)`)
	reAndroidBuildVersion := regexp.MustCompile(`(?miU)^(-v)`)
	reGradleVersion := regexp.MustCompile(`(?miU)^(-b)`)
	for _, arg := range os.Args {
		switch {
		// case reGIT.MatchString(arg):
		// 	gitrepo = arg
		// case reUserName.MatchString(arg):
		// 	username = arg[2:]
		// case rePassword.MatchString(arg):
		// 	password = arg[2:]
		case reAndroidVersion.MatchString(arg):
			androidVersion = arg[2:]
		case reAndroidBuildVersion.MatchString(arg):
			androidVersionBuild = arg[2:]
		case reGradleVersion.MatchString(arg):
			gradleVersion = arg[2:]
		}
	}
	// for {
	// 	if gitrepo == "" {
	// 		fmt.Print("give github repository path (mendatory): ")
	// 		fmt.Scanln(&gitrepo)
	// 	}
	// 	if gitrepo != "" {
	// 		break
	// 	}
	// }

	// if gitrepo == "" {
	// 	fmt.Print("give github repository path (mendatory): ")
	// 	fmt.Scanln(&gitrepo)
	// }

	// if username == "" {
	// 	fmt.Print("username (if any) : ")
	// 	fmt.Scanln(&username)
	// }

	// if password == "" {
	// 	fmt.Print("password (if any) : ")
	// 	fmt.Scanln(&password)
	// }

	if androidVersion == "" {
		fmt.Print("Change Android Version (", androidVersion, ") : ")
		fmt.Scanln(&androidVersion)
	}

	if androidVersionBuild == "" {
		fmt.Print("change Android Build Tools Version (", androidVersionBuild, ") : ")
		fmt.Scanln(&androidVersionBuild)
	}

	if gradleVersion == "" {
		fmt.Print("change gradle Version (", gradleVersion, ") : ")
		fmt.Scanln(&gradleVersion)
	}

	return /*gitrepo, username, password, */ androidVersion, androidVersionBuild, gradleVersion
}

func createdockerImage(gradleVersion, androidVersion, androidVersionBuild, imgname string) {
	dockerFileString := fmt.Sprintf(`
	FROM gradle:%s-jdk8
	
	USER root
	
	ENV SDK_URL="https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip" \
		ANDROID_HOME="/usr/local/android-sdk" \
		ANDROID_VERSION=%s \
		ANDROID_BUILD_TOOLS_VERSION=%s
		
	RUN mkdir "$ANDROID_HOME" .android \
		&& cd "$ANDROID_HOME" \
		&& curl -o sdk.zip $SDK_URL \
		&& unzip sdk.zip \
		&& rm sdk.zip \
		&& mkdir "$ANDROID_HOME/licenses" || true \
		&& echo "24333f8a63b6825ea9c5514f83c2829b004d1fee" > "$ANDROID_HOME/licenses/android-sdk-license"
	
	RUN $ANDROID_HOME/tools/bin/sdkmanager --update
	
	RUN $ANDROID_HOME/tools/bin/sdkmanager "build-tools;${ANDROID_BUILD_TOOLS_VERSION}" \
		"platforms;android-${ANDROID_VERSION}" \
		"platform-tools"
	
	RUN apt-get update && apt-get install build-essential -y && apt-get install file -y && apt-get install apt-utils -y
		
		`, gradleVersion, androidVersion, androidVersionBuild)
	fmt.Println("opening file")
	f, err := os.OpenFile("dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		check(err)
	}
	fmt.Println("writting file")
	if _, err := f.Write([]byte(dockerFileString)); err != nil {
		f.Close()
		check(err)
	}
	if err := f.Close(); err != nil {
		check(err)
	}
	fmt.Println("building docker image")
	cmd := exec.Command("docker", "build", "-t", imgname, "--file", "./dockerfile")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println("building docker image : ")
			fmt.Println(exitError.ExitCode())
		}
	}
}

func pushImage(repo string) {
	fmt.Println("logging in")
	exec.Command("docker", "login", "-u", "shraddheya", "-p", "2totootwO")
	fmt.Println("tagging")
	exec.Command("docker", "tag", "imgname", repo)
	fmt.Println("pushing")
	exec.Command("docker", "push", repo)
	// exec.Command("docker", "logout")
}

func main() {
	/*_, _, _, */ androidVersion, androidVersionBuild, gradleVersion := getVars(os.Args)
	imgname := fmt.Sprintf("android-build:android%s-androidbuild%s-gradle%s", strings.ReplaceAll(androidVersion, ".", "-"), strings.ReplaceAll(androidVersionBuild, ".", "-"), strings.ReplaceAll(gradleVersion, ".", "-"))
	repo := "sraddheya/" + imgname
	if !pull(imgname) {
		fmt.Println("image - " + repo + "not found; creating...")
		createdockerImage(gradleVersion, androidVersion, androidVersionBuild, imgname)
		pushImage(repo)
	}

	// fmt.Println(gitrepo, username, password, androidVersion, androidVersionBuild, gradleVersion, imgname)
}
