/*
Copyright 2023 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resultcollector

import (
	"context"
	"sync"
	"time"

	kcpscopedclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	kcpinformers "github.com/kcp-dev/kcp/pkg/client/informers/externalversions"
	"github.com/kcp-dev/kcp/pkg/client/informers/externalversions/tenancy/v1alpha1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
)

type InformerManager struct {
	sync.Mutex
	mbPreInformer v1alpha1.WorkspaceInformer
}

var informerManager InformerManager = InformerManager{}

func (i *InformerManager) GetMbPreInformer(ctx context.Context, cfg *rest.Config) (v1alpha1.WorkspaceInformer, error) {
	i.Lock()
	defer i.Unlock()
	if i.mbPreInformer != nil {
		return i.mbPreInformer, nil
	}
	espwClient := kcpscopedclient.NewForConfigOrDie(cfg)
	espwInformerFactory := kcpinformers.NewSharedScopedInformerFactory(espwClient, 0, "")
	mbPreInformer := espwInformerFactory.Tenancy().V1alpha1().Workspaces()
	go func() {
		mbPreInformer.Informer().Run(ctx.Done())
	}()
	espwInformerFactory.Start(ctx.Done())
	if err := wait.PollImmediateInfinite(5*time.Second, func() (bool, error) {
		return mbPreInformer.Informer().HasSynced(), nil
	}); err != nil {
		return nil, err
	}
	i.mbPreInformer = mbPreInformer
	return mbPreInformer, nil
}
