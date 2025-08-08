package mock

const (
	ComputePutRestrictionTemplate = `$[?(@.spec && @.spec.skuRef && @.spec.zone && @.spec.bootVolume)]`
	StoragePutRestrictionTemplate = `$[?(@.spec )]`
)
