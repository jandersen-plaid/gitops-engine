package kubetest

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/managedfields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/util/openapi"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
)

type KubectlOutput struct {
	Output string
	Err    error
}

type MockKubectlCmd struct {
	APIResources  []kube.APIResourceInfo
	Commands      map[string]KubectlOutput
	Events        chan watch.Event
	Version       string
	DynamicClient dynamic.Interface

	getResourceFunc *func(ctx context.Context, config *rest.Config, gvk schema.GroupVersionKind, name string, namespace string) (*unstructured.Unstructured, error)
}

// WithGetResourceFunc overrides the default ConvertToVersion behavior.
func (k *MockKubectlCmd) WithGetResourceFunc(getResourcefunc func(context.Context, *rest.Config, schema.GroupVersionKind, string, string) (*unstructured.Unstructured, error)) *MockKubectlCmd {
	k.getResourceFunc = &getResourcefunc
	return k
}

func (k *MockKubectlCmd) NewDynamicClient(config *rest.Config) (dynamic.Interface, error) {
	return k.DynamicClient, nil
}

func (k *MockKubectlCmd) GetAPIResources(config *rest.Config, preferred bool, resourceFilter kube.ResourceFilter) ([]kube.APIResourceInfo, error) {
	return k.APIResources, nil
}

func (k *MockKubectlCmd) GetResource(ctx context.Context, config *rest.Config, gvk schema.GroupVersionKind, name string, namespace string) (*unstructured.Unstructured, error) {
	if k.getResourceFunc != nil {
		return (*k.getResourceFunc)(ctx, config, gvk, name, namespace)
	}

	return nil, nil
}

func (k *MockKubectlCmd) PatchResource(ctx context.Context, config *rest.Config, gvk schema.GroupVersionKind, name string, namespace string, patchType types.PatchType, patchBytes []byte, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (k *MockKubectlCmd) DeleteResource(ctx context.Context, config *rest.Config, gvk schema.GroupVersionKind, name string, namespace string, deleteOptions metav1.DeleteOptions) error {
	command, ok := k.Commands[name]
	if !ok {
		return nil
	}
	return command.Err
}

func (k *MockKubectlCmd) CreateResource(ctx context.Context, config *rest.Config, gvk schema.GroupVersionKind, name string, namespace string, obj *unstructured.Unstructured, createOptions metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (k *MockKubectlCmd) GetServerVersion(config *rest.Config) (string, error) {
	return k.Version, nil
}

func (k *MockKubectlCmd) LoadOpenAPISchema(config *rest.Config) (openapi.Resources, *managedfields.GvkParser, error) {
	return nil, nil, nil
}

func (k *MockKubectlCmd) SetOnKubectlRun(onKubectlRun kube.OnKubectlRunFunc) {
}

func (k *MockKubectlCmd) ManageResources(config *rest.Config, openAPISchema openapi.Resources) (kube.ResourceOperations, func(), error) {
	return &MockResourceOps{}, func() {

	}, nil
}
