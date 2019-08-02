/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package versioned

import (
	buildv1alpha1 "github.com/projectriff/system/pkg/client/clientset/versioned/typed/build/v1alpha1"
	corev1alpha1 "github.com/projectriff/system/pkg/client/clientset/versioned/typed/core/v1alpha1"
	knativev1alpha1 "github.com/projectriff/system/pkg/client/clientset/versioned/typed/knative/v1alpha1"
	streamingv1alpha1 "github.com/projectriff/system/pkg/client/clientset/versioned/typed/streaming/v1alpha1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	BuildV1alpha1() buildv1alpha1.BuildV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Build() buildv1alpha1.BuildV1alpha1Interface
	CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Core() corev1alpha1.CoreV1alpha1Interface
	KnativeV1alpha1() knativev1alpha1.KnativeV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Knative() knativev1alpha1.KnativeV1alpha1Interface
	StreamingV1alpha1() streamingv1alpha1.StreamingV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Streaming() streamingv1alpha1.StreamingV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	buildV1alpha1     *buildv1alpha1.BuildV1alpha1Client
	coreV1alpha1      *corev1alpha1.CoreV1alpha1Client
	knativeV1alpha1   *knativev1alpha1.KnativeV1alpha1Client
	streamingV1alpha1 *streamingv1alpha1.StreamingV1alpha1Client
}

// BuildV1alpha1 retrieves the BuildV1alpha1Client
func (c *Clientset) BuildV1alpha1() buildv1alpha1.BuildV1alpha1Interface {
	return c.buildV1alpha1
}

// Deprecated: Build retrieves the default version of BuildClient.
// Please explicitly pick a version.
func (c *Clientset) Build() buildv1alpha1.BuildV1alpha1Interface {
	return c.buildV1alpha1
}

// CoreV1alpha1 retrieves the CoreV1alpha1Client
func (c *Clientset) CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface {
	return c.coreV1alpha1
}

// Deprecated: Core retrieves the default version of CoreClient.
// Please explicitly pick a version.
func (c *Clientset) Core() corev1alpha1.CoreV1alpha1Interface {
	return c.coreV1alpha1
}

// KnativeV1alpha1 retrieves the KnativeV1alpha1Client
func (c *Clientset) KnativeV1alpha1() knativev1alpha1.KnativeV1alpha1Interface {
	return c.knativeV1alpha1
}

// Deprecated: Knative retrieves the default version of KnativeClient.
// Please explicitly pick a version.
func (c *Clientset) Knative() knativev1alpha1.KnativeV1alpha1Interface {
	return c.knativeV1alpha1
}

// StreamingV1alpha1 retrieves the StreamingV1alpha1Client
func (c *Clientset) StreamingV1alpha1() streamingv1alpha1.StreamingV1alpha1Interface {
	return c.streamingV1alpha1
}

// Deprecated: Streaming retrieves the default version of StreamingClient.
// Please explicitly pick a version.
func (c *Clientset) Streaming() streamingv1alpha1.StreamingV1alpha1Interface {
	return c.streamingV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.buildV1alpha1, err = buildv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.coreV1alpha1, err = corev1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.knativeV1alpha1, err = knativev1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.streamingV1alpha1, err = streamingv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.buildV1alpha1 = buildv1alpha1.NewForConfigOrDie(c)
	cs.coreV1alpha1 = corev1alpha1.NewForConfigOrDie(c)
	cs.knativeV1alpha1 = knativev1alpha1.NewForConfigOrDie(c)
	cs.streamingV1alpha1 = streamingv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.buildV1alpha1 = buildv1alpha1.New(c)
	cs.coreV1alpha1 = corev1alpha1.New(c)
	cs.knativeV1alpha1 = knativev1alpha1.New(c)
	cs.streamingV1alpha1 = streamingv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
