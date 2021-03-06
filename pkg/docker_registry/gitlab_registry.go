package docker_registry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"

	"github.com/werf/logboek"

	"github.com/werf/werf/pkg/image"
)

const GitLabRegistryImplementationName = "gitlab"

var gitlabPatterns = []string{`^gitlab\.com`}

type gitLabRegistry struct {
	*defaultImplementation
	deleteRepoImageFunc func(repoImage *image.Info) error
}

type gitLabRegistryOptions struct {
	defaultImplementationOptions
}

func newGitLabRegistry(options gitLabRegistryOptions) (*gitLabRegistry, error) {
	d, err := newDefaultImplementation(options.defaultImplementationOptions)
	if err != nil {
		return nil, err
	}

	gitLab := &gitLabRegistry{defaultImplementation: d}

	return gitLab, nil
}

func (r *gitLabRegistry) DeleteRepoImage(repoImageList ...*image.Info) error {
	for _, repoImage := range repoImageList {
		if err := r.deleteRepoImage(repoImage); err != nil {
			return err
		}
	}

	return nil
}

func (r *gitLabRegistry) deleteRepoImage(repoImage *image.Info) error {
	deleteRepoImageFunc := r.deleteRepoImageFunc
	if deleteRepoImageFunc != nil {
		return deleteRepoImageFunc(repoImage)
	} else {
		var err error
		for _, deleteFunc := range []func(repoImage *image.Info) error{
			r.deleteRepoImageWithAllScopes,
			r.deleteRepoImageWithUniversalScope,
			r.defaultImplementation.deleteRepoImage,
		} {
			if err = deleteFunc(repoImage); err != nil {
				if strings.Contains(err.Error(), "UNAUTHORIZED") {
					reference := strings.Join([]string{repoImage.Repository, repoImage.Tag}, ":")
					logboek.Debug.LogF("DEBUG: Tag %s deletion failed: %s", reference, err)
					continue
				}

				return err
			}

			r.deleteRepoImageFunc = deleteFunc
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// TODO https://gitlab.com/gitlab-org/gitlab-ce/issues/48968
func (r *gitLabRegistry) deleteRepoImageWithUniversalScope(repoImage *image.Info) error {
	return r.deleteRepoImageWithSpecificScopes(repoImage, func(ref name.Reference) []string {
		return []string{ref.Scope("*")}
	})
}

func (r *gitLabRegistry) deleteRepoImageWithAllScopes(repoImage *image.Info) error {
	return r.deleteRepoImageWithSpecificScopes(repoImage, func(ref name.Reference) []string {
		return []string{ref.Scope("push"), ref.Scope("pull"), ref.Scope("delete")}
	})
}

func (r *gitLabRegistry) deleteRepoImageWithSpecificScopes(repoImage *image.Info, scopesFunc func(ref name.Reference) []string) error {
	reference := strings.Join([]string{repoImage.Repository, repoImage.RepoDigest}, "@")

	ref, err := name.ParseReference(reference, r.api.parseReferenceOptions()...)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %v", reference, err)
	}

	auth, authErr := authn.DefaultKeychain.Resolve(ref.Context().Registry)
	if authErr != nil {
		return fmt.Errorf("getting creds for %q: %v", ref, authErr)
	}

	scopes := scopesFunc(ref)
	tr, err := transport.New(ref.Context().Registry, auth, r.api.getHttpTransport(), scopes)
	if err != nil {
		return err
	}
	c := &http.Client{Transport: tr}

	u := url.URL{
		Scheme: ref.Context().Registry.Scheme(),
		Host:   ref.Context().RegistryStr(),
		Path:   fmt.Sprintf("/v2/%s/manifests/%s", ref.Context().RepositoryStr(), ref.Identifier()),
	}

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		return nil
	default:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("unrecognized status code during DELETE: %v; %v", resp.Status, string(b))
	}
}

func (r *gitLabRegistry) String() string {
	return GitLabRegistryImplementationName
}
