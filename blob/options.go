package blob

import "context"

type ProviderType string

const (
	defaultProviderKey ProviderType = "Default"
)

type SaveResult struct {
	Url string `json:"url"`
}

type CreatedEventData struct {
	Path   string `json:"fileName"`
	Size   int64  `json:"size"`
	Result *SaveResult
}
type SavedHandleFunc func(ctx context.Context, b *CreatedEventData) error

type Options struct {
	providers map[ProviderType]Provider
	handlers  []SavedHandleFunc

	// UseDefaultProvider if true, then if provider not found, default provider will be used
	// default value is true
	UserDefaultProvider bool
	ProvidersConfig     map[ProviderType]map[string]any `mapstructure:"Providers"`
}

func NewOptions() *Options {
	return &Options{
		providers:           make(map[ProviderType]Provider),
		UserDefaultProvider: true,
		ProvidersConfig:     map[ProviderType]map[string]any{},
	}
}

func (opts *Options) AddProvider(providerType ProviderType, provider Provider) *Options {
	opts.providers[providerType] = provider
	return opts
}

func (opts *Options) SetDefaultProvider(provider Provider) *Options {
	opts.providers[defaultProviderKey] = provider
	return opts
}

func (opts *Options) DefaultProvider() Provider {
	return opts.providers[defaultProviderKey]
}

func (opts *Options) GetProvider(providerType ProviderType) Provider {
	if p, ok := opts.providers[providerType]; ok {
		return p
	}
	return nil
}

func (opts *Options) OnSaved(handler SavedHandleFunc) *Options {
	opts.handlers = append(opts.handlers, handler)
	return opts
}

func (opts *Options) Handlers() []SavedHandleFunc {
	return opts.handlers
}

func (opts *Options) GetConfig(t ProviderType, key string) any {
	if c, ok := opts.ProvidersConfig[t]; ok {
		return c[key]
	}
	return nil
}
