package skins

import "context"

type ListSourceStates struct {
	Storage SourceStateStorage
}

func (uc ListSourceStates) Execute(ctx context.Context) ([]SourceState, error) {
	if uc.Storage == nil {
		return nil, nil
	}
	return uc.Storage.ListSourceStates(ctx)
}

type GetDiagnostics struct {
	SourceStates SourceStateStorage
	Version      string
	DatabasePath string
	LogPath      string
}

func (uc GetDiagnostics) Execute(ctx context.Context) (Diagnostics, error) {
	states := []SourceState(nil)
	if uc.SourceStates != nil {
		loaded, err := uc.SourceStates.ListSourceStates(ctx)
		if err != nil {
			return Diagnostics{}, err
		}
		states = loaded
	}
	return Diagnostics{
		Version:      uc.Version,
		DatabasePath: uc.DatabasePath,
		LogPath:      uc.LogPath,
		Sources:      states,
	}, nil
}
