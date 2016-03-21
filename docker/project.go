package docker

import (
	"github.com/Sirupsen/logrus"

	"github.com/mefellows/libcompose/lookup"
	"github.com/mefellows/libcompose/project"
)

// NewProject creates a Project with the specified context.
func NewProject(context *Context) (*project.Project, error) {
	if context.ResourceLookup == nil {
		context.ResourceLookup = &lookup.FileConfigLookup{}
	}

	if context.EnvironmentLookup == nil {
		context.EnvironmentLookup = &lookup.OsEnvLookup{}
	}

	if context.ServiceFactory == nil {
		context.ServiceFactory = &ServiceFactory{
			context: context,
		}
	}

	if context.Builder == nil {
		context.Builder = NewDaemonBuilder(context)
	}

	if context.ClientFactory == nil {
		factory, err := NewDefaultClientFactory(ClientOpts{})
		if err != nil {
			return nil, err
		}
		context.ClientFactory = factory
	}

	p := project.NewProject(&context.Context)

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.open(); err != nil {
		logrus.Errorf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	return p, err
}
