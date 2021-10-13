package backend

type Client interface {
	CodePlatform
	Storage
}

var cli Client

func GetClient() Client {
	return cli
}

func RegisterClient(c Client) {
	cli = c
}

type CodePlatform interface {
	ListRepos(org string) ([]string, error)

	ListBranchesOfRepo(org, repo string) ([]BranchInfo, error)

	ListAllFilesOfRepo(Branch) ([]RepoFile, error)

	GetFileConent(b Branch, path string) (string, string, error)
}

type Storage interface {
	SaveFiles(b Branch, branchSHA string, files []File) error
	GetFileSummary(b Branch, fileName string) ([]RepoFile, error)
}

type Branch struct {
	Org    string `json:"org" required:"true"`
	Repo   string `json:"repo" required:"true"`
	Branch string `json:"branch" required:"true"`
}

type BranchInfo struct {
	Name string `json:"name" required:"true"`
	SHA  string `json:"sha" required:"true"`
}

type RepoFile struct {
	Path string `json:"path" required:"true"`
	SHA  string `json:"sha" required:"true"`
}

type File struct {
	RepoFile

	// Allow empty file
	Content string `json:"content,omitempty"`
}
