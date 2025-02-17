require 'rspec'
require 'httparty'

RSpec.describe 'Simple Test' do
  it 'can make a request' do
    response = HTTParty.get('http://localhost:8080/api/words')
    expect(response.code).to eq(200)
  end
end
