module mqtt

require (
	github.com/afoninsky/noolite-go v0.0.0-20180908163927-5a8479451f27
	github.com/afoninsky/noolite-go/noolite v0.0.0-20180908163927-5a8479451f27
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/spf13/viper v1.2.1
	github.com/yosssi/gmq v0.0.1
	golang.org/x/sys v0.0.0-20181019084534-8f1d3d21f81b // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/afoninsky/noolite-go/noolite => ../noolite
