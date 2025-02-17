require 'rspec'
require 'httparty'
require 'json'
require 'json-schema'
require 'faker'
require 'active_support/all'

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups
end

# Helper module for API testing
module ApiHelper
  def api_url
    'http://localhost:8080/api'
  end

  def json_response
    JSON.parse(response.body)
  end

  def response_schema_file(schema_name)
    File.join(File.dirname(__FILE__), 'schemas', "#{schema_name}.json")
  end

  def validate_schema(response, schema_name)
    schema = JSON.parse(File.read(response_schema_file(schema_name)))
    JSON::Validator.validate!(schema, response)
  end
end

RSpec.configure do |config|
  config.include ApiHelper
end
