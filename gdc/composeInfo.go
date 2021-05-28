package gdc

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//ComposeInfo containing information about the Docker Compose command
type ComposeInfo struct {
	MainFile []byte
  AdditionalFile string
	RequestedServices []string
	cachedConfiguredServices []string
}

func servicesFromCompose(data []byte) []string {
	cf := make(map[interface{}]interface{})
  err := yaml.Unmarshal(data, &cf)
	if (err != nil) {
		Exit("Error parsing compose file %s", err)
	}

	// get services from map because Go has no `keys` method...
	serviceMap := cf["services"].(map[interface{}]interface{})
	services := []string{}
  for k := range(serviceMap) {
		services = append(services, k.(string))
	}
	return services
}

func (compose ComposeInfo) configuredServices() []string {
	if (len(compose.cachedConfiguredServices) > 0) {
		return compose.cachedConfiguredServices
	}
	services := servicesFromCompose((compose.MainFile))
	if (len(compose.AdditionalFile) > 0) {
		file, err := ioutil.ReadFile(compose.AdditionalFile)
		if (err != nil) {
			Exit("Error reading additional Compose file: %s", err)
		}
		newServices := servicesFromCompose(file)
		for _, s := range(newServices) {
			services = append(services, s)
		}
	}
	compose.cachedConfiguredServices = services
	return services
}

// IsServiceConfigured in the compose files or not
func (compose ComposeInfo) IsServiceConfigured(service string) bool {
	found := false
	for _, s := range(compose.configuredServices()) {
			if (s == service) {
				found = true
				break
			}
	}
	return found
}

// IsServiceRequested in the command line or not
func (compose ComposeInfo) IsServiceRequested(service string) bool {
	found := false
	for _, s := range(compose.RequestedServices) {
			if (s == service) {
				found = true
				break
			}
	}
	return found
}
