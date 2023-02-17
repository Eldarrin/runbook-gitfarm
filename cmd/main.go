package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTargetWorkflowQueueLength = 1
)

type githubRunnerMetadata struct {
	githubApiURL                        string
	owner                               string
	personalAccessToken                 string
	repo                                *string
	targetWorkflowQueueLength           int64
	activationTargetWorkflowQueueLength int64
	scalerIndex                         int
}

type WorkflowRuns struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	Id               int64   `json:"id"`
	Name             string  `json:"name"`
	NodeId           string  `json:"node_id"`
	HeadBranch       string  `json:"head_branch"`
	HeadSha          string  `json:"head_sha"`
	Path             string  `json:"path"`
	DisplayTitle     string  `json:"display_title"`
	RunNumber        int     `json:"run_number"`
	Event            string  `json:"event"`
	Status           string  `json:"status"`
	Conclusion       *string `json:"conclusion"`
	WorkflowId       int     `json:"workflow_id"`
	CheckSuiteId     int64   `json:"check_suite_id"`
	CheckSuiteNodeId string  `json:"check_suite_node_id"`
	Url              string  `json:"url"`
	HtmlUrl          string  `json:"html_url"`
	PullRequests     []struct {
		Url    string `json:"url"`
		Id     int    `json:"id"`
		Number int    `json:"number"`
		Head   struct {
			Ref  string `json:"ref"`
			Sha  string `json:"sha"`
			Repo struct {
				Id   int    `json:"id"`
				Url  string `json:"url"`
				Name string `json:"name"`
			} `json:"repo"`
		} `json:"head"`
		Base struct {
			Ref  string `json:"ref"`
			Sha  string `json:"sha"`
			Repo struct {
				Id   int    `json:"id"`
				Url  string `json:"url"`
				Name string `json:"name"`
			} `json:"repo"`
		} `json:"base"`
	} `json:"pull_requests"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Actor     struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"actor"`
	RunAttempt          int           `json:"run_attempt"`
	ReferencedWorkflows []interface{} `json:"referenced_workflows"`
	RunStartedAt        time.Time     `json:"run_started_at"`
	TriggeringActor     struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"triggering_actor"`
	JobsUrl            string  `json:"jobs_url"`
	LogsUrl            string  `json:"logs_url"`
	CheckSuiteUrl      string  `json:"check_suite_url"`
	ArtifactsUrl       string  `json:"artifacts_url"`
	CancelUrl          string  `json:"cancel_url"`
	RerunUrl           string  `json:"rerun_url"`
	PreviousAttemptUrl *string `json:"previous_attempt_url"`
	WorkflowUrl        string  `json:"workflow_url"`
	HeadCommit         struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"committer"`
	} `json:"head_commit"`
	Repository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string `json:"html_url"`
		Description      string `json:"description"`
		Fork             bool   `json:"fork"`
		Url              string `json:"url"`
		ForksUrl         string `json:"forks_url"`
		KeysUrl          string `json:"keys_url"`
		CollaboratorsUrl string `json:"collaborators_url"`
		TeamsUrl         string `json:"teams_url"`
		HooksUrl         string `json:"hooks_url"`
		IssueEventsUrl   string `json:"issue_events_url"`
		EventsUrl        string `json:"events_url"`
		AssigneesUrl     string `json:"assignees_url"`
		BranchesUrl      string `json:"branches_url"`
		TagsUrl          string `json:"tags_url"`
		BlobsUrl         string `json:"blobs_url"`
		GitTagsUrl       string `json:"git_tags_url"`
		GitRefsUrl       string `json:"git_refs_url"`
		TreesUrl         string `json:"trees_url"`
		StatusesUrl      string `json:"statuses_url"`
		LanguagesUrl     string `json:"languages_url"`
		StargazersUrl    string `json:"stargazers_url"`
		ContributorsUrl  string `json:"contributors_url"`
		SubscribersUrl   string `json:"subscribers_url"`
		SubscriptionUrl  string `json:"subscription_url"`
		CommitsUrl       string `json:"commits_url"`
		GitCommitsUrl    string `json:"git_commits_url"`
		CommentsUrl      string `json:"comments_url"`
		IssueCommentUrl  string `json:"issue_comment_url"`
		ContentsUrl      string `json:"contents_url"`
		CompareUrl       string `json:"compare_url"`
		MergesUrl        string `json:"merges_url"`
		ArchiveUrl       string `json:"archive_url"`
		DownloadsUrl     string `json:"downloads_url"`
		IssuesUrl        string `json:"issues_url"`
		PullsUrl         string `json:"pulls_url"`
		MilestonesUrl    string `json:"milestones_url"`
		NotificationsUrl string `json:"notifications_url"`
		LabelsUrl        string `json:"labels_url"`
		ReleasesUrl      string `json:"releases_url"`
		DeploymentsUrl   string `json:"deployments_url"`
	} `json:"repository"`
	HeadRepository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string `json:"html_url"`
		Description      string `json:"description"`
		Fork             bool   `json:"fork"`
		Url              string `json:"url"`
		ForksUrl         string `json:"forks_url"`
		KeysUrl          string `json:"keys_url"`
		CollaboratorsUrl string `json:"collaborators_url"`
		TeamsUrl         string `json:"teams_url"`
		HooksUrl         string `json:"hooks_url"`
		IssueEventsUrl   string `json:"issue_events_url"`
		EventsUrl        string `json:"events_url"`
		AssigneesUrl     string `json:"assignees_url"`
		BranchesUrl      string `json:"branches_url"`
		TagsUrl          string `json:"tags_url"`
		BlobsUrl         string `json:"blobs_url"`
		GitTagsUrl       string `json:"git_tags_url"`
		GitRefsUrl       string `json:"git_refs_url"`
		TreesUrl         string `json:"trees_url"`
		StatusesUrl      string `json:"statuses_url"`
		LanguagesUrl     string `json:"languages_url"`
		StargazersUrl    string `json:"stargazers_url"`
		ContributorsUrl  string `json:"contributors_url"`
		SubscribersUrl   string `json:"subscribers_url"`
		SubscriptionUrl  string `json:"subscription_url"`
		CommitsUrl       string `json:"commits_url"`
		GitCommitsUrl    string `json:"git_commits_url"`
		CommentsUrl      string `json:"comments_url"`
		IssueCommentUrl  string `json:"issue_comment_url"`
		ContentsUrl      string `json:"contents_url"`
		CompareUrl       string `json:"compare_url"`
		MergesUrl        string `json:"merges_url"`
		ArchiveUrl       string `json:"archive_url"`
		DownloadsUrl     string `json:"downloads_url"`
		IssuesUrl        string `json:"issues_url"`
		PullsUrl         string `json:"pulls_url"`
		MilestonesUrl    string `json:"milestones_url"`
		NotificationsUrl string `json:"notifications_url"`
		LabelsUrl        string `json:"labels_url"`
		ReleasesUrl      string `json:"releases_url"`
		DeploymentsUrl   string `json:"deployments_url"`
	} `json:"head_repository"`
}

type Repos struct {
	Repo []Repo
}

type Repo struct {
	Id       int    `json:"id"`
	NodeId   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private          bool        `json:"private"`
	HtmlUrl          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	Url              string      `json:"url"`
	ArchiveUrl       string      `json:"archive_url"`
	AssigneesUrl     string      `json:"assignees_url"`
	BlobsUrl         string      `json:"blobs_url"`
	BranchesUrl      string      `json:"branches_url"`
	CollaboratorsUrl string      `json:"collaborators_url"`
	CommentsUrl      string      `json:"comments_url"`
	CommitsUrl       string      `json:"commits_url"`
	CompareUrl       string      `json:"compare_url"`
	ContentsUrl      string      `json:"contents_url"`
	ContributorsUrl  string      `json:"contributors_url"`
	DeploymentsUrl   string      `json:"deployments_url"`
	DownloadsUrl     string      `json:"downloads_url"`
	EventsUrl        string      `json:"events_url"`
	ForksUrl         string      `json:"forks_url"`
	GitCommitsUrl    string      `json:"git_commits_url"`
	GitRefsUrl       string      `json:"git_refs_url"`
	GitTagsUrl       string      `json:"git_tags_url"`
	GitUrl           string      `json:"git_url"`
	IssueCommentUrl  string      `json:"issue_comment_url"`
	IssueEventsUrl   string      `json:"issue_events_url"`
	IssuesUrl        string      `json:"issues_url"`
	KeysUrl          string      `json:"keys_url"`
	LabelsUrl        string      `json:"labels_url"`
	LanguagesUrl     string      `json:"languages_url"`
	MergesUrl        string      `json:"merges_url"`
	MilestonesUrl    string      `json:"milestones_url"`
	NotificationsUrl string      `json:"notifications_url"`
	PullsUrl         string      `json:"pulls_url"`
	ReleasesUrl      string      `json:"releases_url"`
	SshUrl           string      `json:"ssh_url"`
	StargazersUrl    string      `json:"stargazers_url"`
	StatusesUrl      string      `json:"statuses_url"`
	SubscribersUrl   string      `json:"subscribers_url"`
	SubscriptionUrl  string      `json:"subscription_url"`
	TagsUrl          string      `json:"tags_url"`
	TeamsUrl         string      `json:"teams_url"`
	TreesUrl         string      `json:"trees_url"`
	CloneUrl         string      `json:"clone_url"`
	MirrorUrl        string      `json:"mirror_url"`
	HooksUrl         string      `json:"hooks_url"`
	SvnUrl           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Language         interface{} `json:"language"`
	ForksCount       int         `json:"forks_count"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Size             int         `json:"size"`
	DefaultBranch    string      `json:"default_branch"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	IsTemplate       bool        `json:"is_template"`
	Topics           []string    `json:"topics"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	HasDownloads     bool        `json:"has_downloads"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	Visibility       string      `json:"visibility"`
	PushedAt         time.Time   `json:"pushed_at"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Permissions      struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	AllowRebaseMerge    bool        `json:"allow_rebase_merge"`
	TemplateRepository  interface{} `json:"template_repository"`
	TempCloneToken      string      `json:"temp_clone_token"`
	AllowSquashMerge    bool        `json:"allow_squash_merge"`
	AllowAutoMerge      bool        `json:"allow_auto_merge"`
	DeleteBranchOnMerge bool        `json:"delete_branch_on_merge"`
	AllowMergeCommit    bool        `json:"allow_merge_commit"`
	SubscribersCount    int         `json:"subscribers_count"`
	NetworkCount        int         `json:"network_count"`
	License             struct {
		Key     string `json:"key"`
		Name    string `json:"name"`
		Url     string `json:"url"`
		SpdxId  string `json:"spdx_id"`
		NodeId  string `json:"node_id"`
		HtmlUrl string `json:"html_url"`
	} `json:"license"`
	Forks      int `json:"forks"`
	OpenIssues int `json:"open_issues"`
	Watchers   int `json:"watchers"`
}

var githubApiURL = "https://api.github.com"
var httpClient = &http.Client{}
var personalAccessToken = ""
var owner = "Eldarrin"
var repo *string

// getUserRepositories returns a list of repositories for a given organization
func getUserRepositories(ctx context.Context) (*[]string, error) {
	url := fmt.Sprintf("%s/users/Eldarrin/repos", githubApiURL)
	body, err := getGithubRequest(ctx, url, httpClient)
	if err != nil {
		return nil, err
	}

	var repos []Repo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		fmt.Println(err, "Cannot unmarshal GitHub Workflow Repos API response")
		return nil, err
	}

	fmt.Println(string(body))

	var repoList []string
	for _, repo := range repos {
		fmt.Println(repo.Name)

		repoList = append(repoList, repo.Name)
	}

	return &repoList, nil
}

func getGithubRequest(ctx context.Context, url string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+personalAccessToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	r, err := httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}
	r.Body.Close()

	if r.StatusCode != 200 {
		return []byte{}, fmt.Errorf("the GitHub REST API returned error. url: %s status: %d response: %s", url, r.StatusCode, string(b))
	}

	return b, nil
}

func stripDeadRuns(allWfrs []WorkflowRuns) []WorkflowRun {
	var filtered []WorkflowRun
	for _, wfrs := range allWfrs {
		for _, wfr := range wfrs.WorkflowRuns {
			if wfr.Status == "queued" {
				filtered = append(filtered, wfr)
			}
		}
	}
	return filtered
}

// GetWorkflowQueueLength returns the number of workflows in the queue
func GetWorkflowQueueLength(ctx context.Context) (int64, error) {

	/*
		curl \
		-H "Accept: application/vnd.github+json" \
		-H "Authorization: Bearer <YOUR-TOKEN>"\
		-H "X-GitHub-Api-Version: 2022-11-28" \
		https://api.github.com/repos/OWNER/REPO/actions/runs
	*/

	var repos *[]string
	var err error

	if repo == nil {
		repos, err = getUserRepositories(ctx)
		if err != nil {
			return -1, err
		}
		//s.logger.Info("Found repos")
	} else {
		repos = &[]string{*repo}
	}

	var allWfrs []WorkflowRuns

	for _, repo := range *repos {
		url := fmt.Sprintf("%s/repos/%s/%s/actions/runs", githubApiURL, owner, repo)
		body, err := getGithubRequest(ctx, url, httpClient)
		if err != nil {
			return -1, err
		}

		var wfrs WorkflowRuns
		err = json.Unmarshal(body, &wfrs)
		if err != nil {
			fmt.Println(err, "Cannot unmarshal GitHub Workflow Runs API response")
			return -1, err
		}
		allWfrs = append(allWfrs, wfrs)
	}

	return int64(len(stripDeadRuns(allWfrs))), nil
}

func main() {
	//repos := "runbook-gitfarm"
	//repo = &repos
	ctx := context.Background()
	a, err := GetWorkflowQueueLength(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)

}
