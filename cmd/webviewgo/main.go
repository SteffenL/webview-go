package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func fetch_webview(libs_dir string, set_env bool, version string) {
	lib_dir := filepath.Join(libs_dir, "webview")
	_, err := os.Stat(filepath.Join(lib_dir, "webview.h"))
	if err != nil || version == "master" {
		if err == nil {
			err = os.RemoveAll(lib_dir)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("Fetching webview %s...\n", version)
		os.MkdirAll(filepath.Dir(lib_dir), os.ModePerm) // TODO: reduce permissions
		repo_owner := "SteffenL"
		repo_name := "webview-nogo"
		repo_slug := strings.Join([]string{repo_owner, repo_name}, "/")
		repo_base_url := "https://github.com/" + repo_slug
		archive_url := repo_base_url + "/archive"
		url := fmt.Sprintf("%s/%s.tar.gz", archive_url, version)
		curl_cmd := exec.Command("curl", "-sSL", url)
		curl_stdout, err := curl_cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		defer curl_stdout.Close()
		tar_cmd := exec.Command("tar", "-xf", "-", "-C", filepath.Dir(lib_dir))
		tar_cmd.Stdin = curl_stdout
		curl_cmd.Start()
		_, err = tar_cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		rename_from_path := filepath.Join(filepath.Dir(lib_dir), strings.Join([]string{repo_name, version}, "-"))
		rename_to_path := lib_dir
		err = os.Rename(rename_from_path, rename_to_path)
		if err != nil {
			log.Fatal(err)
		}
	}
	if set_env {
		get_env_cmd := exec.Command("go", "env", "-json")
		env_output, err := get_env_cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		env := map[string]string{}
		err = json.Unmarshal(env_output, &env)
		if err != nil {
			log.Fatal(err)
		}
		include_dir := lib_dir
		if strings.Contains(env["CGO_CXXFLAGS"], include_dir) {
			return
		}
		fmt.Printf("Updating the environment for webview %s...\n", version)
		cxxflags := fmt.Sprintf("CGO_CXXFLAGS=%s \"-I%s\"", env["CGO_CXXFLAGS"], include_dir)
		set_env_cmd := exec.Command("go", "env", "-w", cxxflags)
		err = set_env_cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func fetch_mswebview2(libs_dir string, set_env bool, version string) {
	lib_dir := filepath.Join(libs_dir, "mswebview2")
	_, err := os.Stat(filepath.Join(lib_dir, "Microsoft.Web.WebView2.nuspec"))
	if err != nil || version == "latest" {
		if err == nil {
			err = os.RemoveAll(lib_dir)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("Fetching WebView2 %s...\n", version)
		os.MkdirAll(lib_dir, os.ModePerm) // TODO: reduce permissions
		url := "https://www.nuget.org/api/v2/package/Microsoft.Web.WebView2/"
		if version != "latest" {
			url += version
		}
		curl_cmd := exec.Command("curl", "-sSL", url)
		curl_stdout, err := curl_cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		defer curl_stdout.Close()
		tar_cmd := exec.Command("tar", "-xf", "-", "-C", lib_dir)
		tar_cmd.Stdin = curl_stdout
		curl_cmd.Start()
		_, err = tar_cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
	}
	if set_env {
		get_env_cmd := exec.Command("go", "env", "-json")
		env_output, err := get_env_cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		env := map[string]string{}
		err = json.Unmarshal(env_output, &env)
		if err != nil {
			log.Fatal(err)
		}
		include_subdir := filepath.Join("build", "native", "include")
		include_dir := filepath.Join(lib_dir, include_subdir)
		if strings.Contains(env["CGO_CXXFLAGS"], include_dir) {
			return
		}
		fmt.Printf("Updating the environment for WebView2 %s...\n", version)
		cxxflags := fmt.Sprintf("CGO_CXXFLAGS=%s \"-I%s\"", env["CGO_CXXFLAGS"], include_dir)
		set_env_cmd := exec.Command("go", "env", "-w", cxxflags)
		err = set_env_cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func cmd_fetch_libs(cmd string, args []string) {
	var libs_dir string
	var set_env bool
	var webview_version string
	var mswebview2_version string
	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.StringVar(&libs_dir, "libs-dir", "libs", "Libraries output directory.")
	flags.BoolVar(&set_env, "set-env", false, "Set environment variables for go.")
	flags.StringVar(&webview_version, "webview-version", "master", "webview version to use (branch, tag, commit).")
	flags.StringVar(&mswebview2_version, "mswebview2-version", "1.0.1150.38", "Microsoft WebView2 version to use or \"latest\".")
	flags.Parse(args)

	if !filepath.IsAbs(libs_dir) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		libs_dir = filepath.Join(wd, libs_dir)
	}
	os := runtime.GOOS
	fetch_webview(libs_dir, set_env, webview_version)
	if os == "windows" {
		fetch_mswebview2(libs_dir, set_env, mswebview2_version)
	}
}

func main() {
	commands := map[string]func(string, []string){"fetch-libs": cmd_fetch_libs}
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Commands:")
		for k, _ := range commands {
			fmt.Printf("  %s\n", k)
		}
		os.Exit(1)
	}
	cmd := args[0]
	args = args[1:]
	for k, v := range commands {
		if (k == cmd) {
			v(cmd, args)
			return
		}
	}
	fmt.Printf("Invalid command: %s\n", cmd)
	os.Exit(1)
}
