package syncer_test

import (
	"context"

	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
)

var _ = Describe("Translate Proxy", func() {

	var (
		xdsCache    *mockXdsCache
		sanitizer   *mockXdsSanitizer
		syncer      v1.ApiSyncer
		snap        *v1.ApiSnapshot
		settings    *v1.Settings
		proxyClient v1.ProxyClient
		ctx         context.Context
		cancel      context.CancelFunc
		proxyName   = "proxy-name"
		ref         = "syncer-test"
		ns          = "any-ns"
	)

	BeforeEach(func() {
		xdsCache = &mockXdsCache{}
		sanitizer = &mockXdsSanitizer{}
		ctx, cancel = context.WithCancel(context.Background())

		resourceClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}

		proxyClient, _ = v1.NewProxyClient(ctx, resourceClientFactory)

		upstreamClient, err := resourceClientFactory.NewResourceClient(ctx, factory.NewResourceClientParams{ResourceType: &v1.Upstream{}})
		Expect(err).NotTo(HaveOccurred())

		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
				Namespace: ns,
				Name:      proxyName,
			},
		}

		settings = &v1.Settings{}

		rep := reporter.NewReporter(ref, proxyClient.BaseClient(), upstreamClient)

		xdsHasher := &xds.ProxyKeyHasher{}
		syncer = NewTranslatorSyncer(&mockTranslator{true, false}, xdsCache, xdsHasher, sanitizer, rep, false, nil, settings)
		snap = &v1.ApiSnapshot{
			Proxies: v1.ProxyList{
				proxy,
			},
		}
		_, err = proxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		err = syncer.Sync(context.Background(), snap)
		Expect(err).NotTo(HaveOccurred())

		proxies, err := proxyClient.List(proxy.GetMetadata().Namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies).To(HaveLen(1))
		Expect(proxies[0]).To(BeAssignableToTypeOf(&v1.Proxy{}))
		Expect(proxies[0].Status).To(Equal(&core.Status{
			State:      2,
			Reason:     "1 error occurred:\n\t* hi, how ya doin'?\n\n",
			ReportedBy: ref,
		}))

		// NilSnapshot is always consistent, so snapshot will always be set as part of endpoints update
		Expect(xdsCache.called).To(BeTrue())

		// update rv for proxy
		p1, err := proxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		snap.Proxies[0] = p1

		syncer = NewTranslatorSyncer(&mockTranslator{false, false}, xdsCache, xdsHasher, sanitizer, rep, false, nil, settings)

		err = syncer.Sync(context.Background(), snap)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() { cancel() })

	It("writes the reports the translator spits out and calls SetSnapshot on the cache", func() {
		proxies, err := proxyClient.List(ns, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies).To(HaveLen(1))
		Expect(proxies[0]).To(BeAssignableToTypeOf(&v1.Proxy{}))
		Expect(proxies[0].Status).To(Equal(&core.Status{
			State:      1,
			ReportedBy: ref,
		}))

		Expect(xdsCache.called).To(BeTrue())
	})

	It("updates the cache with the sanitized snapshot", func() {
		sanitizer.snap = envoycache.NewEasyGenericSnapshot("easy")
		err := syncer.Sync(context.Background(), snap)
		Expect(err).NotTo(HaveOccurred())

		Expect(sanitizer.called).To(BeTrue())
		Expect(xdsCache.setSnap).To(BeEquivalentTo(sanitizer.snap))
	})

	It("uses listeners and routes from the previous snapshot when sanitization fails", func() {
		sanitizer.err = errors.Errorf("we ran out of coffee")

		oldXdsSnap := xds.NewSnapshotFromResources(
			envoycache.NewResources("", nil),
			envoycache.NewResources("", nil),
			envoycache.NewResources("", nil),
			envoycache.NewResources("old listeners from before the war", []envoycache.Resource{
				resource.NewEnvoyResource(&envoy_config_listener_v3.Listener{}),
			}),
		)

		// return this old snapshot when the syncer asks for it
		xdsCache.getSnap = oldXdsSnap
		err := syncer.Sync(context.Background(), snap)
		Expect(err).NotTo(HaveOccurred())

		Expect(sanitizer.called).To(BeTrue())
		Expect(xdsCache.called).To(BeTrue())

		oldListeners := oldXdsSnap.GetResources(resource.ListenerTypeV3)
		newListeners := xdsCache.setSnap.GetResources(resource.ListenerTypeV3)

		Expect(oldListeners).To(Equal(newListeners))

		oldRoutes := oldXdsSnap.GetResources(resource.RouteTypeV3)
		newRoutes := xdsCache.setSnap.GetResources(resource.RouteTypeV3)

		Expect(oldRoutes).To(Equal(newRoutes))
	})

})

var _ = Describe("Translate mulitple proxies with errors", func() {

	var (
		xdsCache       *mockXdsCache
		sanitizer      *mockXdsSanitizer
		syncer         v1.ApiSyncer
		snap           *v1.ApiSnapshot
		settings       *v1.Settings
		proxyClient    v1.ProxyClient
		upstreamClient v1.UpstreamClient
		proxyName      = "proxy-name"
		upstreamName   = "upstream-name"
		ref            = "syncer-test"
		ns             = "any-ns"
	)

	proxiesShouldHaveErrors := func(proxies v1.ProxyList, numProxies int) {
		Expect(proxies).To(HaveLen(numProxies))
		for _, proxy := range proxies {
			Expect(proxy).To(BeAssignableToTypeOf(&v1.Proxy{}))
			Expect(proxy.Status).To(Equal(&core.Status{
				State:      2,
				Reason:     "1 error occurred:\n\t* hi, how ya doin'?\n\n",
				ReportedBy: ref,
			}))

		}

	}
	writeUniqueErrsToUpstreams := func() {
		// Re-writes existing upstream to have an annotation
		// which triggers a unique error to be written from each proxy's mockTranslator
		upstreams, err := upstreamClient.List(ns, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(upstreams).To(HaveLen(1))

		us := upstreams[0]
		// This annotation causes the translator mock to generate a unique error per proxy on each upstream
		us.Metadata.Annotations = map[string]string{"uniqueErrPerProxy": "true"}
		_, err = upstreamClient.Write(us, clients.WriteOpts{OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())
		snap.Upstreams = upstreams
		err = syncer.Sync(context.Background(), snap)
		Expect(err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		var err error
		xdsCache = &mockXdsCache{}
		sanitizer = &mockXdsSanitizer{}

		resourceClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}

		proxyClient, _ = v1.NewProxyClient(context.Background(), resourceClientFactory)

		usClient, err := resourceClientFactory.NewResourceClient(context.Background(), factory.NewResourceClientParams{ResourceType: &v1.Upstream{}})
		Expect(err).NotTo(HaveOccurred())

		proxy1 := &v1.Proxy{
			Metadata: &core.Metadata{
				Namespace: ns,
				Name:      proxyName + "1",
			},
		}
		proxy2 := &v1.Proxy{
			Metadata: &core.Metadata{
				Namespace: ns,
				Name:      proxyName + "2",
			},
		}

		us := &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      upstreamName,
				Namespace: ns,
			},
		}

		settings = &v1.Settings{}

		rep := reporter.NewReporter(ref, proxyClient.BaseClient(), usClient)

		xdsHasher := &xds.ProxyKeyHasher{}
		syncer = NewTranslatorSyncer(&mockTranslator{true, true}, xdsCache, xdsHasher, sanitizer, rep, false, nil, settings)
		snap = &v1.ApiSnapshot{
			Proxies: v1.ProxyList{
				proxy1,
				proxy2,
			},
			Upstreams: v1.UpstreamList{
				us,
			},
		}

		_, err = usClient.Write(us, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = proxyClient.Write(proxy1, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = proxyClient.Write(proxy2, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		err = syncer.Sync(context.Background(), snap)
		Expect(err).NotTo(HaveOccurred())

		proxies, err := proxyClient.List(proxy1.GetMetadata().Namespace, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies).To(HaveLen(2))
		Expect(proxies[0]).To(BeAssignableToTypeOf(&v1.Proxy{}))
		Expect(proxies[0].Status).To(Equal(&core.Status{
			State:      2,
			Reason:     "1 error occurred:\n\t* hi, how ya doin'?\n\n",
			ReportedBy: ref,
		}))

		// NilSnapshot is always consistent, so snapshot will always be set as part of endpoints update
		Expect(xdsCache.called).To(BeTrue())

		upstreamClient, err = v1.NewUpstreamClient(context.Background(), resourceClientFactory)
		Expect(err).NotTo(HaveOccurred())
	})

	It("handles reporting errors on multiple proxies sharing an upstream reporting 2 different errors", func() {
		// Testing the scenario where we have multiple proxies,
		// each of which should report a different unique error on an upstream.
		proxies, err := proxyClient.List(ns, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		proxiesShouldHaveErrors(proxies, 2)

		writeUniqueErrsToUpstreams()

		upstreams, err := upstreamClient.List(ns, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())

		Expect(upstreams[0].Status).To(Equal(&core.Status{
			State:      2,
			Reason:     "2 errors occurred:\n\t* upstream is bad - determined by proxy-name1\n\t* upstream is bad - determined by proxy-name2\n\n",
			ReportedBy: ref,
		}))

		Expect(xdsCache.called).To(BeTrue())
	})

	It("handles reporting errors on multiple proxies sharing an upstream, each reporting the same upstream error", func() {
		// Testing the scenario where we have multiple proxies,
		// each of which should report the same error on an upstream.
		proxies, err := proxyClient.List(ns, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		proxiesShouldHaveErrors(proxies, 2)

		upstreams, err := upstreamClient.List(ns, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(upstreams).To(HaveLen(1))
		Expect(upstreams[0].Status).To(Equal(&core.Status{
			State:      2,
			Reason:     "1 error occurred:\n\t* generic upstream error\n\n",
			ReportedBy: ref,
		}))

		Expect(xdsCache.called).To(BeTrue())
	})
})

type mockTranslator struct {
	reportErrs         bool
	reportUpstreamErrs bool // Adds an error to every upstream in the snapshot
}

func (t *mockTranslator) Translate(params plugins.Params, proxy *v1.Proxy) (envoycache.Snapshot, reporter.ResourceReports, *validation.ProxyReport, error) {
	if t.reportErrs {
		rpts := reporter.ResourceReports{}
		rpts.AddError(proxy, errors.Errorf("hi, how ya doin'?"))
		if t.reportUpstreamErrs {
			for _, upstream := range params.Snapshot.Upstreams {
				if upstream.Metadata.Annotations["uniqueErrPerProxy"] == "true" {
					rpts.AddError(upstream, errors.Errorf("upstream is bad - determined by %s", proxy.Metadata.Name))
				} else {
					rpts.AddError(upstream, errors.Errorf("generic upstream error"))
				}
			}
		}
		return envoycache.NilSnapshot{}, rpts, &validation.ProxyReport{}, nil
	}
	return envoycache.NilSnapshot{}, nil, &validation.ProxyReport{}, nil
}

var _ envoycache.SnapshotCache = &mockXdsCache{}

type mockXdsCache struct {
	called bool
	// snap that is set
	setSnap envoycache.Snapshot
	// snap that is returned
	getSnap envoycache.Snapshot
}

func (*mockXdsCache) CreateWatch(envoycache.Request) (value chan envoycache.Response, cancel func()) {
	panic("implement me")
}

func (*mockXdsCache) Fetch(context.Context, envoycache.Request) (*envoycache.Response, error) {
	panic("implement me")
}

func (*mockXdsCache) GetStatusInfo(string) envoycache.StatusInfo {
	panic("implement me")
}

func (c *mockXdsCache) GetStatusKeys() []string {
	return []string{}
}

func (c *mockXdsCache) SetSnapshot(node string, snapshot envoycache.Snapshot) error {
	c.called = true
	c.setSnap = snapshot
	return nil
}

func (c *mockXdsCache) GetSnapshot(node string) (envoycache.Snapshot, error) {
	if c.getSnap != nil {
		return c.getSnap, nil
	}
	return &envoycache.NilSnapshot{}, nil
}

func (*mockXdsCache) ClearSnapshot(node string) {
	panic("implement me")
}

type mockXdsSanitizer struct {
	called bool
	snap   envoycache.Snapshot
	err    error
}

func (s *mockXdsSanitizer) SanitizeSnapshot(ctx context.Context, glooSnapshot *v1.ApiSnapshot, xdsSnapshot envoycache.Snapshot, reports reporter.ResourceReports) (envoycache.Snapshot, error) {
	s.called = true
	if s.snap != nil {
		return s.snap, nil
	}
	if s.err != nil {
		return nil, s.err
	}
	return xdsSnapshot, nil
}
