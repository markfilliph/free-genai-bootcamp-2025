require 'rspec'
require 'httparty'
require 'json'

def api_url
  'http://localhost:8080/api'
end

RSpec.describe 'Study API' do
  let(:base_url) { "#{api_url}/study" }
  let(:valid_session) do
    {
      group_id: nil, # Will be set in before block
      duration_minutes: 30,
      target_words: 10
    }
  end

  before(:all) do
    # Create a group
    group_response = HTTParty.post(
      "#{api_url}/groups",
      body: {
        name: 'Test Study Group',
        description: 'Group for testing study sessions'
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
          japanese: "単語#{i}",
          romaji: "tango#{i}",
          english: "word#{i}",
          parts: ['noun']
        }.to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @words << word_response.parsed_response

      # Add word to group
      HTTParty.post("#{api_url}/groups/#{@group_id}/words/#{word_response.parsed_response['id']}")
    end
  end

  describe 'POST /study/start' do
    context 'with valid parameters' do
      before do
        @response = HTTParty.post(
          "#{base_url}/start",
          body: valid_session.merge(group_id: @group_id).to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 201 status code' do
        expect(@response.code).to eq(201)
      end

      it 'returns the created study session' do
        expect(@response.parsed_response).to include(
          'id',
          'group_id' => @group_id,
          'duration_minutes' => valid_session[:duration_minutes],
          'target_words' => valid_session[:target_words],
          'status' => 'in_progress'
        )
      end
    end

    context 'with invalid parameters' do
      before do
        @response = HTTParty.post(
          "#{base_url}/start",
          body: { group_id: 999999 }.to_json,
          headers: { 'Content-Type' => 'application/json' }
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

  describe 'POST /study/:id/answer' do
    before do
      # Start a study session
      session_response = HTTParty.post(
        "#{base_url}/start",
        body: valid_session.merge(group_id: @group_id).to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @session_id = session_response.parsed_response['id']
      @word = @words.first
    end

    context 'with correct answer' do
      before do
        @response = HTTParty.post(
          "#{base_url}/#{@session_id}/answer",
          body: {
            word_id: @word['id'],
            answer: @word['english'],
            time_taken_ms: 1500
          }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns success result' do
        expect(@response.parsed_response).to include(
          'correct' => true,
          'word_id' => @word['id']
        )
      end
    end

    context 'with incorrect answer' do
      before do
        @response = HTTParty.post(
          "#{base_url}/#{@session_id}/answer",
          body: {
            word_id: @word['id'],
            answer: 'wrong answer',
            time_taken_ms: 1500
          }.to_json,
          headers: { 'Content-Type' => 'application/json' }
        )
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns failure result' do
        expect(@response.parsed_response).to include(
          'correct' => false,
          'word_id' => @word['id']
        )
      end
    end
  end

  describe 'POST /study/:id/end' do
    before do
      # Start a study session
      session_response = HTTParty.post(
        "#{base_url}/start",
        body: valid_session.merge(group_id: @group_id).to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @session_id = session_response.parsed_response['id']
    end

    context 'with valid session' do
      before do
        @response = HTTParty.post("#{base_url}/#{@session_id}/end")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns session summary' do
        expect(@response.parsed_response).to include(
          'session_id' => @session_id,
          'total_words' => be_a(Integer),
          'correct_answers' => be_a(Integer),
          'average_time_ms' => be_a(Integer)
        )
      end
    end

    context 'with invalid session' do
      before do
        @response = HTTParty.post("#{base_url}/999999/end")
      end

      it 'returns 404 status code' do
        expect(@response.code).to eq(404)
      end

      it 'returns an error message' do
        expect(@response.parsed_response).to include('error')
      end
    end
  end

  describe 'GET /study/history' do
    before do
      @response = HTTParty.get("#{base_url}/history")
    end

    it 'returns 200 status code' do
      expect(@response.code).to eq(200)
    end

    it 'returns a paginated list of study sessions' do
      expect(@response.parsed_response).to include('items', 'current_page', 'total_pages', 'total_items', 'items_per_page')
    end

    it 'returns items as an array' do
      expect(@response.parsed_response['items']).to be_an(Array)
    end
  end

  describe 'GET /study/:id' do
    before do
      # Start a study session
      session_response = HTTParty.post(
        "#{base_url}/start",
        body: valid_session.merge(group_id: @group_id).to_json,
        headers: { 'Content-Type' => 'application/json' }
      )
      @session_id = session_response.parsed_response['id']
    end

    context 'when session exists' do
      before do
        @response = HTTParty.get("#{base_url}/#{@session_id}")
      end

      it 'returns 200 status code' do
        expect(@response.code).to eq(200)
      end

      it 'returns session details' do
        expect(@response.parsed_response).to include(
          'id' => @session_id,
          'group_id' => @group_id,
          'status'
        )
      end
    end

    context 'when session does not exist' do
      before do
        @response = HTTParty.get("#{base_url}/999999")
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
