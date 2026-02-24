package resolve

import (
	"fmt"
	"io/ioutil"
	"os"

	"resodns/internal/pkg/console"
	"resodns/pkg/fileoperation"
	"resodns/pkg/progressbar"
	"resodns/pkg/wildcarder"
)

type WildcardFilterOptions struct {
	CacheFilename        string
	DomainOutputFilename string
	RootOutputFilename   string
	Resolvers            []string
	QueriesPerSecond     int
	ThreadCount          int
	ResolveTestCount     int
	BatchSize            int
}

type DefaultWildcardFilter struct {
	wc *wildcarder.Wildcarder
}

func NewDefaultWildcardFilter() *DefaultWildcardFilter {
	return &DefaultWildcardFilter{}
}

func (f *DefaultWildcardFilter) Filter(opt WildcardFilterOptions, totalCount int) (found int, roots []string, err error) {
	cacheReader, err := createCacheReader(opt.CacheFilename)
	if err != nil {
		return 0, nil, err
	}
	defer cacheReader.Close()
	f.wc = createWildcarder(opt)
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		return 0, nil, err
	}
	defer func() { tempFile.Close(); os.Remove(tempFile.Name()) }()
	tmpl := "[ETA {{ eta }}] {{ bar }} {{ current }}/{{ total }} queries: {{ queries }} (time: {{ time }})"
	bar := progressbar.New(f.updateProgressBar, int64(totalCount), progressbar.WithTemplate(tmpl), progressbar.WithWriter(console.Output))
	bar.Start()
	rootMap := make(map[string]struct{})
	for {
		precache, domainFile, count, err := prepareCache(cacheReader, tempFile.Name(), opt.BatchSize)
		if err != nil {
			return 0, nil, err
		}
		if count == 0 {
			break
		}
		f.wc.SetPreCache(precache)
		domains, rootsBatch := f.wc.Filter(domainFile)
		domainFile.Close()
		found += len(domains)
		if err := fileoperation.AppendLines(domains, opt.DomainOutputFilename); err != nil {
			return 0, nil, err
		}
		for _, root := range rootsBatch {
			rootMap[root] = struct{}{}
		}
	}
	var rootList []string
	for root := range rootMap {
		rootList = append(rootList, root)
	}
	if err := fileoperation.AppendLines(rootList, opt.RootOutputFilename); err != nil {
		return 0, nil, err
	}
	bar.Stop()
	return found, rootList, nil
}

func (f *DefaultWildcardFilter) updateProgressBar(bar *progressbar.ProgressBar) {
	bar.SetCurrent(int64(f.wc.Current()))
	bar.Set("queries", fmt.Sprintf("%d", f.wc.QueryCount()))
}

func createCacheReader(filename string) (*CacheReader, error) {
	cacheFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return NewCacheReader(cacheFile), nil
}

func createWildcarder(opt WildcardFilterOptions) *wildcarder.Wildcarder {
	qps := opt.QueriesPerSecond / len(opt.Resolvers)
	if qps == 0 && len(opt.Resolvers) > 0 {
		qps = 1
	}
	resolver := wildcarder.NewClientDNS(opt.Resolvers, 10, qps, 100)
	return wildcarder.New(opt.ThreadCount, opt.ResolveTestCount, wildcarder.WithResolver(resolver))
}

func prepareCache(cacheReader *CacheReader, tempFilename string, batchSize int) (*wildcarder.DNSCache, *os.File, int, error) {
	domainFile, err := os.Create(tempFilename)
	if err != nil {
		return nil, nil, 0, err
	}
	precache := wildcarder.NewDNSCache()
	totalCount, err := cacheReader.Read(domainFile, precache, batchSize)
	if err := domainFile.Sync(); err != nil {
		return nil, nil, 0, err
	}
	if _, err := domainFile.Seek(0, 0); err != nil {
		return nil, nil, 0, err
	}
	return precache, domainFile, totalCount, err
}
