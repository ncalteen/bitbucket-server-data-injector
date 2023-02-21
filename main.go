package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	// command line flags
	projectsCount = flag.Int("projects-count", 1, "Number of projects to create, default 1")
	reposCount    = flag.Int("repos-count", 1, "Number of repos to create, default 1")
)

func init() {
	// Check environment variables have values
	if os.Getenv("BBS_API_USER") == "" {
		log.Fatal("$BBS_API_USER must be set")
	}
	if os.Getenv("BBS_API_PASS") == "" {
		log.Fatal("$BBS_API_PASS must be set")
	}
	if os.Getenv("BBS_API_URL") == "" {
		log.Fatal("$BBS_API_URL must be set")
	}
}

func main() {
	// parse command line flags
	flag.Parse()

	// get the state of the bbs cluster as a health check and verification the API is working
	healthResp, healthStatus, err := healthCheck(os.Getenv("BBS_API_USER"), os.Getenv("BBS_API_PASS"), os.Getenv("BBS_API_URL"))

	// check for errors
	if err != nil {
		log.Fatal(err)
	}

	// check the status code
	if healthStatus != 200 {
		log.Fatal("Health check failed with status: ", healthStatus)
	}

	// read the response body
	healthBody, err := ioutil.ReadAll(healthResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// print the response body
	fmt.Println("Health check response: ")
	fmt.Println(string(healthBody))
	fmt.Println("Status code: ", healthStatus)
	fmt.Println("Adding " + fmt.Sprint(*projectsCount) + " projects and " + fmt.Sprint(*reposCount) + " repos to each project")

	// create projects
	for i := 0; i < *projectsCount; i++ {
		// create a new project
		projectResp, projectStatus, err := createProject(os.Getenv("BBS_API_USER"), os.Getenv("BBS_API_PASS"), os.Getenv("BBS_API_URL"), "project"+fmt.Sprint(i))
		if err != nil {
			log.Fatal(err)
		}

		// check the status code
		if projectStatus != 201 {
			log.Fatal("Project creation failed with status: ", projectStatus)
		}

		// read the response body
		projectBody, err := ioutil.ReadAll(projectResp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// print the response body
		fmt.Println("Project creation response: ")
		fmt.Println(string(projectBody))
		fmt.Println("Status code: ", projectStatus)

		// create repos
		for j := 0; j < *reposCount; j++ {
			// create a new repo
			repoResp, repoStatus, err := createRepo(os.Getenv("BBS_API_USER"), os.Getenv("BBS_API_PASS"), os.Getenv("BBS_API_URL"), "project"+fmt.Sprint(i), "repo"+fmt.Sprint(j))
			if err != nil {
				log.Fatal(err)
			}

			// check the status code
			if repoStatus != 201 {
				log.Fatal("Repo creation failed with status: ", repoStatus)
			}

			// read the response body
			repoBody, err := ioutil.ReadAll(repoResp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// print the response body
			fmt.Println("Repo creation response: ")
			fmt.Println(string(repoBody))
			fmt.Println("Status code: ", repoStatus)
		}
	}

}

// healthCheck returns a response and status code from the Bitbucket Server API or an error
func healthCheck(user string, password string, url string) (response http.Response, status int, err error) {
	// create a new request
	req, err := http.NewRequest("GET", url+"/status", nil)
	if err != nil {
		return response, status, err
	}
	// create a new basic auth header
	req.SetBasicAuth(user, password)

	// set the content type
	req.Header.Set("Content-Type", "application/json")

	// make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, status, err
	}

	// return the response and status code
	return *resp, resp.StatusCode, nil
}

func createProject(user string, password string, url string, project string) (response http.Response, status int, err error) {
	// create body for the request
	body := []byte(`{"key": "Project-` + project + `", "links": ""}`)

	// create a new request
	req, err := http.NewRequest("POST", url+"/rest/api/latest/projects", bytes.NewBuffer(body))
	if err != nil {
		return response, status, err
	}
	// create a new basic auth header
	req.SetBasicAuth(user, password)

	// set the content type
	req.Header.Set("Content-Type", "application/json")

	// make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, status, err
	}

	// return the response and status code
	return *resp, resp.StatusCode, nil
}

func createRepo(user string, password string, url string, project string, repo string) (response http.Response, status int, err error) {
	// create body for the request
	body := []byte(`{"name": "Repo-` + repo + `", "scmId": "git", "forkable": true}`)

	// create a new request
	req, err := http.NewRequest("POST", url+"/rest/api/latest/projects/Project-"+project+"/repos", bytes.NewBuffer(body))
	if err != nil {
		return response, status, err
	}
	// create a new basic auth header
	req.SetBasicAuth(user, password)

	// set the content type
	req.Header.Set("Content-Type", "application/json")

	// make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, status, err
	}

	// return the response and status code
	return *resp, resp.StatusCode, nil
}
