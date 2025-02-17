require 'rspec'
require 'httparty'
require 'json'

def api_url
  'http://localhost:8080/api'
end

RSpec.describe 'Statistics API' do
  let(:base_url) { "#{api_url}/stats" }

  before(:all) do
    # Create a group
    group_response = HTTParty.post(
      "#{api_url}/groups",
      body: {
        name: 'Test Stats Group',
        description: 'Group for testing statistics'
      }.to_json,
      headers: { 'Content-Type' => 'application/json' }
    )
    @group_id = group_response.parsed_response['id']

    # Create some words
    @words = []
    3.times do |i|
      word_response = HTTParty.post(
        "#{api_url}/words",
        body: {
          japanese: "統計#{i}",
          romaji: "toukei#{i}",
          english: "statistic#{i}",
          parts: ['noun']
        }.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @words << word_response.parsed_response

      # Add word to group
      HTTParty.post("#{api_url}/groups/#{@group_id}/words/#{word_response.parsed_response['id']}")
    end

    # Create and complete a study session
    session_response = HTTParty.post(
      "#{api_url}/study/start",
      body: {
        group_id: @group_id,
        duration_minutes: 30,
        target_words: 10
      }.to_json,
      headers: { 'Content-Type' => 'application/json' }
    )
    @session_id = session_response.parsed_response['id']

    # Submit some answers
    @words.each do |word|
      HTTParty.post(
        "#{api_url}/study/#{@session_id}/answer",
        body: {
          word_id: word['id'],
          answer: word['english'],
          time_taken_ms: 1500
        }.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
    end

    # End the session
    HTTParty.post("#{api_url}/study/#{@session_id}/end")
  end

  describe 'GET /stats/dashboard' do
    before do
      @response = HTTParty.get("#{base_url}/dashboard")
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns dashboard statistics' do
      expect(@response.parsed_response).to include(
        'total_words',
        'total_study_time_minutes',
        'words_studied',
        'accuracy_percentage'
      )
    end

    it 'returns numeric values' do
      expect(@response.parsed_response['total_words']).to be_a(Integer)
      expect(@response.parsed_response['total_study_time_minutes']).to be_a(Integer)
      expect(@response.parsed_response['words_studied']).to be_a(Integer)
      expect(@response.parsed_response['accuracy_percentage']).to be_a(Float)
    end
  end

  describe 'GET /stats/progress' do
    before do
      @response = HTTParty.get(
        "#{base_url}/progress",
        query: {
          period: 'week'
        }
      )
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns progress statistics' do
      expect(@response.parsed_response).to include(
        'study_sessions',
        'words_learned',
        'study_time_minutes'
      )
    end

    it 'returns data points as arrays' do
      expect(@response.parsed_response['study_sessions']).to be_an(Array)
      expect(@response.parsed_response['words_learned']).to be_an(Array)
      expect(@response.parsed_response['study_time_minutes']).to be_an(Array)
    end

    context 'with invalid period' do
      before do
        @response = HTTParty.get(
          "#{base_url}/progress",
          query: {
            period: 'invalid'
          }
        )
      end

      it 'returns 400 status code' do
        expect(@response.code).to eq(400)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end

  describe 'GET /stats/words/:id' do
    before do
      @word = @words.first
      @response = HTTParty.get("#{base_url}/words/#{@word['id']}")
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns word statistics' do
      expect(@response.parsed_response).to include(
        'word_id' => @word['id'],
        'times_studied',
        'correct_answers',
        'average_time_ms',
        'last_studied_at'
      )
    end

    context 'when word does not exist' do
      before do
        @response = HTTParty.get("#{base_url}/words/999999")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end

  describe 'GET /stats/groups/:id' do
    before do
      @response = HTTParty.get("#{base_url}/groups/#{@group_id}")
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns group statistics' do
      expect(@response.parsed_response).to include(
        'group_id' => @group_id,
        'total_words',
        'words_studied',
        'accuracy_percentage',
        'total_study_time_minutes'
      )
    end

    context 'when group does not exist' do
      before do
        @response = HTTParty.get("#{base_url}/groups/999999")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end
end
