# frozen_string_literal: true

lib = File.expand_path('lib', __dir__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'version'

puts "GDC::VERSION: #{GDC::VERSION}"

Gem::Specification.new do |spec|
	spec.name                  = 'global_docker_compose'
	spec.version               = GDC::VERSION
	spec.date                  = %q{2021-04-20}
	spec.authors               = ['Chandeep Singh']
	spec.email                 = ['chandeep.singh@flipp.com']
	spec.summary               = 'A wrapper for a global docker compose file, making it easy to setup dependencies for various applications through docker, while only maintaining a single docker-compose file.'
	spec.files                 = Dir['lib/**']
	spec.require_paths         = ["lib"]
	spec.required_ruby_version = Gem::Requirement.new(">= 2.3.0")
	spec.bindir                = "exe"
	spec.executables           = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths           = ["lib"]
end
